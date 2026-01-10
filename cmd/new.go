package cmd

import (
	// "errors"
	"log/slog"
	// "os"
	// "path/filepath"

	// "proj/internal/config"
	// "proj/internal/luabridge"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	// "github.com/yuin/gopher-lua"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new <kind> <name>",
	Args:  cobra.ExactArgs(2),
	Short: "Create a project called NAME based on the KIND template",
	Long: `Create a new project based on a project template in the current directory.

Uses the project template, name, and standard variables from config to create a new project ready to
work on. Projects typically include skeleton files as well as shell scripts to pre-configure things
like packaging, remote repositories, database creation, etc.
`,
	Run: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().String("target-root", ".", "Path to create the project in")
	viper.BindPFlag("target_root", newCmd.Flags().Lookup("target-root"))

	newCmd.Flags().String("target-path", "", "Path to write files at")
	viper.BindPFlag("target_path", newCmd.Flags().Lookup("target-path"))

	newCmd.Flags().String("template-root", ".", "Path containing project templates")
	viper.BindPFlag("template_root", newCmd.Flags().Lookup("template-root"))

	newCmd.Flags().String("template-path", "", "Path to read files from")
	viper.BindPFlag("template_path", newCmd.Flags().Lookup("template-path"))

	newCmd.Flags().StringArray("set", []string{}, "Set a variable using key=value")
	viper.BindPFlag("set", newCmd.Flags().Lookup("set"))
}

func runNew(cmd *cobra.Command, args []string) {
	slog.Info("Hello from new")
}

// func buildAbsolutePath(source string, base string, ext string) {
// 	value := viper.GetString(source)
// 	if value != "" {
// 		path, err := filepath.Abs(value)
// 		if err != nil {
// 			slog.Error("Failed to absolute-ize path", "value", value, "error", err)
// 		}
// 		slog.Debug("key provided", "source", source, "value", viper.GetString(source), "path", path)
// 		viper.Set(source, path)
// 	} else {
// 		slog.Debug("BUILDING", "base", base, "val", viper.GetString(base))
// 		path, err := filepath.Abs(filepath.Join(viper.GetString(base), ext))
// 		if err != nil {
// 			slog.Error("failed to absolute-ize path", "source", source, "value", viper.GetString(source), "ext", ext)
// 		}
// 		viper.Set(source, path)
// 		slog.Debug("path not set, using default", "path", path)
// 	}
// 	slog.Debug("Added path", source, viper.GetString(source))
// }
//
// // TODO: maybe this should take a config to merge in with?
// func loadTemplateConfig(template_path string) (*viper.Viper, error) {
// 	template_config := viper.New()
// 	template_config.SetConfigName("config")
// 	template_config.AddConfigPath(template_path)
// 	if err := template_config.ReadInConfig(); err != nil {
// 		slog.Error("Error reading template config", "template_path", template_path, "file", template_config.ConfigFileUsed(), "error", err)
// 		return nil, err
// 	}
// 	slog.Debug("Read project config", "file", template_config.ConfigFileUsed(), "settings", template_config.AllSettings())
// 	return template_config, nil
// }

// func runNew(cmd *cobra.Command, args []string) {
// 	slog.Debug("New called", "args", args)
//
// 	viper.Set("kind", args[0])
// 	buildAbsolutePath("template_root", "template_root", "")
// 	buildAbsolutePath("template_path", "template_root", args[0])
//
// 	viper.Set("name", args[1])
// 	buildAbsolutePath("target_root", "target_root", "")
// 	buildAbsolutePath("target_path", "target_root", args[1])
//
// 	slog.Debug("loading template config")
// 	template_config, err := loadTemplateConfig(viper.GetString("template_path"))
// 	if err != nil {
// 		os.Exit(1)
// 	}
// 	if err := template_config.MergeConfigMap(viper.AllSettings()); err != nil {
// 		slog.Debug("Failed to merge global config with project config", "error", err)
// 		os.Exit(1)
// 	}
//
//
// 	// TODO:
// 	// run before scripts
// 	// resolve variables
// 	// check if target is okay (bail unless --force if it exists)
// 	// do the copy / template dance
// 	// run after scripts
//
// 	// quick hack -- should be target-root, then target-path
// 	lua_path := filepath.Join(viper.GetString("target_root"), "test.lua")
//
// 	script_runner := luabridge.NewRuntime(template_config.AllSettings(), lua_path)
// 	defer script_runner.CloseState()
// 	script_runner.Run()
// }
