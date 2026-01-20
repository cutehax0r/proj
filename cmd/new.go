package cmd

import (
	"bytes"
	"log/slog"
	"os"
	"proj/internal/config"
	"proj/internal/luabridge"
	"proj/internal/paths"
	"strings"
	"text/template"

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

	newCmd.Flags().StringP("definition-name", "d", "new", "Definition in template to use")
	viper.BindPFlag("definition-name", newCmd.Flags().Lookup("definition-name"))
}

func runNew(cmd *cobra.Command, args []string) {
	slog.Debug("Execute New Command", slog.String("Definition", viper.GetString("definition-name")), slog.Group("Arguments", slog.String("Template Name", args[0]), slog.String("Target Name", args[1])))

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

	defPath := strings.Join([]string{"definitions", viper.GetString("definition-name")}, ".")
	if viper.IsSet(defPath) == false {
		slog.Error("Definition does not exist in template", slog.String("path", defPath), slog.String("definition-name", viper.GetString("definition-name")), slog.String("template name", viper.GetString("template-name")), slog.String("template config file", viper.GetString("template-config-file")))
	}

	reqs, _ := config.BuildRequirements()
	slog.Debug("Final Requirements", slog.Any("reqs", reqs))

	vars, _ := config.BuildVariables(reqs.Variables)
	slog.Debug("Final Variables", slog.Any("vars", vars))

	scripts, err := config.ParseScriptSpecs(paths)
	if err != nil {
		slog.Error("Couldn't build scripts")
		os.Exit(1)
	}
	slog.Debug("Final scripts", slog.Any("scripts", scripts))

	files, err := config.ParseFileSpecs()
	if err != nil {
		slog.Error("Failed to load files from template definition", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Debug("Final files", slog.Any("files", files))

	luaenv := luabridge.NewRuntime(vars, paths, reqs, files, viper.GetBool("no-write"))
	for _, script := range scripts.BeforeScripts() {
		luaenv.Run(script)
		// check for errors and if they exist then os.exit
	}

	// this should only happen for 'required' vars
	for _, varspec := range reqs.Variables {
		if vars[varspec.Name] == nil {
			slog.Error("Required variable is not set. Use --set-variable. Aborting.", slog.Any("Name", varspec.Name))
			// this is where you put a 'while nil { prompt }' loop in v2.
			slog.Info("All variables", slog.Any("vars", vars))
			os.Exit(1)
		}
	}
	slog.Debug("All the variables are ready so we can do the work")

	for _, file := range *files {
		desttemp, err := template.New("filename").Parse(file.Target)
		var deststr bytes.Buffer
		// consider adding 'funcs' here
		if err != nil {
			slog.Error("Couldn't template the target filename", slog.String("target", file.Target), slog.Any("err", err))
		}
		err = desttemp.Execute(&deststr, vars)
		if err != nil {
			slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", file.Target))
		}

		if viper.GetBool("no-write") == true {
			slog.Info("No-write set: skipping copy", slog.String("source", file.Source), slog.String("target", deststr.String()))
		} else {
			slog.Info("Copying", slog.Bool("parse", file.Parse), slog.String("source", file.Source), slog.String("target", deststr.String()))
		}
	}

	for _, script := range scripts.AfterScripts() {
		luaenv.Run(script)
	}
}

func logViperDebug(desc string, settings map[string]any) {
	yamlBytes, err := yaml.Marshal(settings)
	if err != nil {
		slog.Error("Marshal merged failed", slog.Any("error", err))
		return
	}
	slog.Debug(desc, "all settings", string(yamlBytes))
}
