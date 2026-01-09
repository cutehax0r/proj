/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

var rootCmd = &cobra.Command{
	Use:   "proj",
	Short: "Tool for creating projects and files from templates",
	Long: `Create templates in a shared directory. Build new projects using those templates as
a scaffold. Files in the templates can include variables and define scripts to run.

General format is "proj [template] [name] OPTIONS. Examples:

  "proj new rails foobar": create a rails project called foobar
  "proj add model baz": create a model called baz in an existing rails project
  "proj new html foo --log-level 3": create a new HTML project with very verbose logging
	`,
	PersistentPreRun: persistentPreRun,
	Run:              runRoot,
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	// setup logging
	ll := slog.LevelDebug
	switch viper.GetInt("log_level") {
	case 0:
		ll = slog.LevelError
	case 1:
		ll = slog.LevelWarn
	case 2:
		ll = slog.LevelInfo
	case 3:
		ll = slog.LevelDebug
	}
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: ll})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)
	slog.Debug("Loaded configuration", "path", viper.ConfigFileUsed(), "settings", viper.AllSettings())

	// read in global config
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// TODO: use os.userconfigdir and os.userhomedir. Also, absoluteize the paths
		if xdg_config_home := os.Getenv("XDG_CONFIG_HOME"); xdg_config_home != "" {
			viper.AddConfigPath(filepath.Join(xdg_config_home, "proj"))
		}
		viper.AddConfigPath("$HOME/.config/proj")
		viper.AddConfigPath(".")

	}
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("Error reading configuration", "path", viper.ConfigFileUsed(), "err", err)
		os.Exit(1)
	} else {
		viper.Set("config", viper.ConfigFileUsed())
		slog.Debug("Read configuration", "path", viper.ConfigFileUsed(), "settings", viper.AllSettings())
	}

}

func runRoot(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func init() {
	viper.SetConfigName("proj")

	if xdg_data_home := os.Getenv("XDG_DATA_HOME"); xdg_data_home != "" {
		viper.SetDefault("template_root", filepath.Join(xdg_data_home, "proj"))
	} else {
		viper.SetDefault("template_root", "$HOME/.local/share/proj")
	}

	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "Print the plan, don't write anything")
	viper.BindPFlag("dry_run", rootCmd.PersistentFlags().Lookup("dry-run"))

	rootCmd.PersistentFlags().IntP("log-level", "l", 0, "How much to log [0-3], bigger = more")
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Use specific global config file")

	// TODO: make better defaults
	viper.SetDefault("requirements", make(map[string]any)) // { variables: array[string] } - get better later
	viper.SetDefault("variables", make(map[string]any))    // map[string]any
	viper.SetDefault("scripts", make(map[string]string))   // { before: nil, after: nil)

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
