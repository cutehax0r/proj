package cmd

import (
	"log/slog"
	"os"
	"proj/internal/config"
	"proj/internal/luabridge"
	"proj/internal/paths"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
	// "github.com/yuin/gopher-lua"
)

var newCmd = &cobra.Command{
	Use:   "new <kind> <name>",
	Args:  cobra.ExactArgs(2),
	Short: "Create a project called NAME based on the KIND template",
	Long:  `Create a new project based on a project template in the current directory.`,
	Run:   runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("target-root", "r", ".", "Path to create the project in")
	viper.BindPFlag("target-root", newCmd.Flags().Lookup("target-root"))

	newCmd.Flags().StringP("target-path", "p", "", "Path to write files at")
	viper.BindPFlag("target-path", newCmd.Flags().Lookup("target-path"))

	newCmd.Flags().StringP("template-root", "s", paths.TemplateRootDir(), "Path containing project templates")
	viper.BindPFlag("template-root", newCmd.Flags().Lookup("template-root"))

	newCmd.Flags().StringP("template-path", "t", "", "Path to read files from")
	viper.BindPFlag("template-path", newCmd.Flags().Lookup("template-path"))

	newCmd.Flags().StringArrayP("set-variable", "v", []string{}, "Set a variable using key=value")
	viper.BindPFlag("set-variables", newCmd.Flags().Lookup("set-variable"))

	newCmd.Flags().StringP("definition", "d", "new", "Definition in template to use")
	viper.BindPFlag("definition", newCmd.Flags().Lookup("definition"))
}

func runNew(cmd *cobra.Command, args []string) {
	slog.Debug("Execute New Command", slog.String("Definition", viper.GetString("definition")), slog.Group("Arguments", slog.String("Template Name", args[0]), slog.String("Target Name", args[1])))

	// mangling viper at this point feels kinda wrong. Maybe we update the 'new paths from config' to just take some arguments?
	viper.Set("template-name", args[0])
	viper.Set("target-name", args[1])
	viper.Set("target-config-file", "") // this should not exist yet if we're running new

	logViperDebug("Global config loaded", viper.AllSettings())

	paths, err := paths.NewPathsFromConfig(viper.AllSettings())
	if err != nil {
		os.Exit(1)
	}
	slog.Debug("Paths", slog.Any("paths", paths))

	if err = config.InitTemplate(paths.TemplateConfigFile); err != nil {
		os.Exit(1)
	}
	logViperDebug("template config loaded", viper.AllSettings())

	if _, err := os.Stat(paths.TargetPath); err == nil {
		slog.Error("Target path exists", slog.String("path", paths.TargetPath))
		os.Exit(1)
	}

	defPath := strings.Join([]string{"definitions", viper.GetString("definition")}, ".")
	if viper.IsSet(defPath) == false {
		slog.Error("Definition does not exist in template", slog.String("path", defPath), slog.String("definition", viper.GetString("definition")), slog.String("template name", viper.GetString("template-name")), slog.String("template config file", viper.GetString("template-config-file")))
	}

	// ioy
	reqs, _ := config.BuildRequirements()
	vars, _ := config.BuildVariables(reqs.Variables)

	scripts, err := config.ParseScriptSpecs(paths)
	if err != nil {
		slog.Error("Couldn't build scripts")
		os.Exit(1)
	}

	// will need files && requirements
	luaenv := luabridge.NewRuntime(vars, paths, reqs, viper.GetBool("no-write"))
	for _, script := range scripts.BeforeScripts() {
		luaenv.Run(script)
		// check for errors and if they exist then os.exit
	}

	slog.Debug("Final Requirements", slog.Any("reqs", reqs))
	slog.Debug("Final Variables", slog.Any("vars", vars))

	for key, value := range vars {
		if value == nil {
			slog.Error("Required variable is not set. Use --set-variable. Aborting.", slog.Any(key, value))
			// this is where you put a 'while nil { prompt }' loop in v2.
			os.Exit(1)
		}
	}

	if viper.GetBool("no-write") == true {
		slog.Info("No-write set: skipping copy")
	} else {
		copyFiles()
	}

	for _, script := range scripts.AfterScripts() {
		luaenv.Run(script)
	}
}

func copyFiles() error {
	files, err := config.ParseFileSpecs()
	if err != nil {
		slog.Error("Failed to load files from template definition", slog.Any("error", err))
		return err
	}
	for _, file := range files {
		slog.Info("Copying", slog.Bool("parse", file.Parse), slog.String("source", file.Source), slog.String("target", file.Target))
	}
	return nil
}

func logViperDebug(desc string, settings map[string]any) {
	yamlBytes, err := yaml.Marshal(settings)
	if err != nil {
		slog.Error("Marshal merged failed", slog.Any("error", err))
		return
	}
	slog.Debug(desc, "all settings", string(yamlBytes))
}
