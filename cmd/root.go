/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/Marlliton/slogpretty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var globalConfig string

var rootCmd = &cobra.Command{
	Use:              "proj",
	Short:            "Reusable templates for creating new projects and add files to existing ones.",
	Long:             `This makes projects. Detailed description coming later`,
	PersistentPreRun: persistentPreRun,
	Run:              runRoot,
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	logOpts := &slogpretty.Options{
		Level:      slog.LevelDebug,
		AddSource:  true,
		Colorful:   true,
		Multiline:  true,
		TimeFormat: "15:04:05",
	}
	logHandler := slogpretty.New(os.Stdout, logOpts)
	slog.SetDefault(slog.New(logHandler))
}

func runRoot(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func init() {
	viper.SetConfigName("proj")

	rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "Print the plan, don't write anything")
	viper.BindPFlag("dry_run", rootCmd.PersistentFlags().Lookup("dry-run"))

	rootCmd.PersistentFlags().IntP("log-level", "l", 0, "How much to log [0-3], bigger = more")
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().StringVarP(&globalConfig, "global-config", "c", "", "Use specific global configuration file")

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
