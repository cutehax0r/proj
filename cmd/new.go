package cmd

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yuin/gopher-lua"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [kind] [name]",
	Args: cobra.ExactArgs(2),
	Short: "Create a project called NAME based on the KIND template",
	Long: `Create a new project based on a project template in the current directory.

Uses the project template, name, and standard variables from config to create a new project ready to
work on. Projects typically include skeleton files as well as shell scripts to pre-configure things
like packaging, remote repositories, database creation, etc.
`,
	Run: newProject,
}

func init() {
	rootCmd.AddCommand(newCmd)
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// responsible for building the new project
func newProject(cmd *cobra.Command, args []string) {
	kind := args[0]
	name := args[1]
	slog.Debug("New called", "kind", kind, "name", name)

	// generally it's better to just try and then crash with errors but I'm learning about how
	// file management works in go so I'm being a bit cautious. This means I'm risking a race
	// where the time between check & use allows things to change. I'll eat that for now and
	// refactor later.
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Couldn't get home directory", "home", home)
		os.Exit(1)
	}

	src := filepath.Join(home, ".local", "share", "proj", kind)
	slog.Debug("Set compute src path", "src", src, "home", home)
	_, err = os.Stat(src)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Error("Template definition not found at expected path", "kind", kind, "src", src)
		} else {
			slog.Error("Template definition couldn't be read", "kind", kind, "src", src)
		}
		os.Exit(1)
	}

	slog.Debug("Parsing template config", "src", src)

	// check that config.yml exists
	// use viper to 'MergeInConfig' for the new config file.
	// verify config.yml has mandatory keys/correct structure

	// check destination
	dest, err := os.Getwd()
	if err != nil {
		slog.Error("Couldn't get working directory")
		os.Exit(1)
	}

	// check working dir to see if we can mess with it
	info, err := os.Stat(dest)
	if err != nil {
		slog.Error("Couldn't read working directory", "dest", dest)
		os.Exit(1)
	} else {
		if info.Mode()&0200 == 0 {
			slog.Error("No permission to write to working directory", "dest", dest)
			os.Exit(1)
		}
	}

	// then ensure our new target dir is okay
	dest = filepath.Join(dest, name)
	slog.Debug("Set dest path", "dest", dest)
	_, err = os.Stat(dest)
	if err == nil {
		slog.Error("Destination path already exists", "dest", dest)
		os.Exit(1)
	} else {
		if errors.Is(err, os.ErrNotExist) {
			slog.Debug("Destination path doesn't exist, all good", "dest", dest)
		} else {
			slog.Error("Error checking for destination", "dest", "dest")
			os.Exit(1)
		}
	}

	// that's a lot of os.exits but if you get here we should be able to start the real work
	config_vars := viper.GetStringMapString("variables")
	slog.Debug("building variables", "variables", config_vars)

	// run "before" script -- should we have separate environments for each script?
	before := filepath.Join(src, "scripts", "before.lua")
	slog.Debug("Attempt to run 'before.lua' script", "path", before)
	info, err = os.Stat(before)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Debug("No before script defined, skipping", "path", before)
		} else {
			slog.Error("Error accessing before script, aborting", "path", before)
			os.Exit(1)
		}
	} else {
		slog.Info("Running before script", "path", before)
		L := lua.NewState()
		defer L.Close()
		L.OpenLibs()

		// this needs to get more clever and recursively transform types
		variables := L.NewTable()
		for k, v := range config_vars {
			variables.RawSetString(k, lua.LString(v))
		}
		// also expose source, target, log level kind, and name
		L.SetGlobal("VARIABLES", variables)

		err = L.DoFile(before)
		if err != nil {
			slog.Error("Lua error in before script", "path", before, "error", err)
			os.Exit(1)
		}
		// We should copy back the changed variables to Go here.
		raw_lua_variables := L.GetGlobal("VARIABLES")
		lua_variables, ok := raw_lua_variables.(*lua.LTable)
		if !ok {
			slog.Error("Lua screwed up the variables")
		}
		go_variables := make(map[string]string)
		lua_variables.ForEach(func(k, v lua.LValue) {
			go_variables[k.String()] = v.String()
		})
		config_vars = go_variables

	}

	slog.Info("Creating target directory", "dest", dest)
	// really need a 'dry run' flag
	err = os.MkdirAll(dest, 0777) // umask will make this more restrictive
	if err != nil {
		slog.Error("Failed to create destination", "dest", dest)
		os.Exit(1)
	}

	slog.Info("Copy each file to destination")

	slog.Info("Run 'after' lua script")

	// run "after" script - should this just be a do_script using the before env?
	after := filepath.Join(src, "scripts", "after.lua")
	slog.Debug("Attempt to run 'after.lua' script", "path", after)
	info, err = os.Stat(after)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Debug("No after script defined, skipping", "path", after)
		} else {
			slog.Error("Error accessing after script, aborting", "path", after)
			os.Exit(1)
		}
	} else {
		slog.Info("Running after script", "path", after)
		L := lua.NewState()
		defer L.Close()
		L.OpenLibs()

		// this needs to get more clever and recursively transform types
		variables := L.NewTable()
		for k, v := range config_vars { // note this was modified by before
			variables.RawSetString(k, lua.LString(v))
		}
		// also expose source, target, log level kind, and name
		L.SetGlobal("VARIABLES", variables)

		err = L.DoFile(after)
		if err != nil {
			slog.Error("Lua error in after script", "path", after, "error", err)
			os.Exit(1)
		}
		// we should copy stuff back here so that it's available for 'global' scripts even
		// though this currently is the last step in the chain.
	}
	
}
