/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "proj",
	Short: "Tool for creating projects and files from templates",
	Long: `Create projects using templates stored in ~/.local/share/proj/FOO/new.
Add files to those projects using templates in ~/.local/share/proj/FOO/files/BAR.

"proj new rails foobar": create a rails project called foobar
"proj new model baz": create a model called baz in an existing rails project
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	logLevel := slog.LevelDebug
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel })
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// lets load the default global config
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/proj")
	err := viper.ReadInConfig()
	if err != nil {
		slog.Error("Error reading config file", "err", err, "path", viper.ConfigFileUsed())
		// maybe tried to read a file that doesn't exist - can create defaults
	} else {
		slog.Info("Read Configuration configuration", "path", viper.ConfigFileUsed())
	}

	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//
	// flags
	// - config: load a new default config (~/.config/proj/config.yml default)
	// - dir: set the target directory for a new project (cwd default)
	// - projects: path to where templates can be found
}
