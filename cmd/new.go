package cmd

import (
	"log/slog"
	"proj/internal/config"
	"proj/internal/paths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
	// "github.com/yuin/gopher-lua"
)

var newCmd = &cobra.Command{
	Use:   "new <kind> <name>",
	Args:  cobra.ExactArgs(2),
	Short: "Create a project called NAME based on the KIND template",
	Long: `Create a new project based on a project template in the current directory.`,
	Run: runNew,
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

	viper.Set("template-name", args[0])
	viper.Set("target-name", args[1])
	viper.Set("target-config-file", "") // this should not exist yet if we're running new

	dumpViper("Global config loaded", viper.AllSettings())

	paths, err := paths.NewPathsFromConfig(viper.AllSettings())
	if err != nil {
		slog.Error("Couldn't build paths", slog.Any("Error", err))
	}
	slog.Debug("Paths", slog.Any("paths", paths))

	config.InitTemplate(paths.TemplateConfigFile)
	dumpViper("template config loaded", viper.AllSettings())

	slog.Info("Running new - checking paths")
	// check that path doesn't exist || path is empty
	// check that if nothing in teh target root heiarchy has a .proj/ in it.
}

func dumpViper(desc string, settings map[string]any) {
	yamlBytes, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		slog.Error("Marshal merged failed", slog.Any("error", err))
		return
	}
	slog.Debug(desc, "all settings", string(yamlBytes))
}

// func runNew(cmd *cobra.Command, args []string) {
// 	// quick hack -- should be target-root, then target-path
// 	lua_path := filepath.Join(viper.GetString("target_root"), "test.lua")
//
// 	script_runner := luabridge.NewRuntime(template_config.AllSettings(), lua_path)
// 	defer script_runner.CloseState()
// 	script_runner.Run()
// }
