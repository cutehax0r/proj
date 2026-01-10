package cmd

import (
	"os"
	"proj/internal/config"
	"proj/internal/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var globalConfigFile string

var rootCmd = &cobra.Command{
	Use:              "proj",
	Short:            "Reusable templates for creating new projects and add files to existing ones.",
	Long:             `This makes projects. Detailed description coming later`,
	PersistentPreRun: persistentPreRun,
	Run:              runRoot,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	viper.SetConfigName("proj")

	rootCmd.PersistentFlags().BoolP("no-write", "w", false, "Print the plan and don't write anything")
	rootCmd.PersistentFlags().IntP("log-level", "l", 0, "How much to log [0-3], bigger = more")
	rootCmd.PersistentFlags().StringVarP(&globalConfigFile, "global-config-file", "g", "", "Use specific global configuration file")
	viper.BindPFlags(rootCmd.PersistentFlags())

	// TODO: make better defaults
	viper.SetDefault("requirements", make(map[string]any)) // { variables: array[string] } - get better later
	viper.SetDefault("variables", make(map[string]any))    // map[string]any
	viper.SetDefault("scripts", make(map[string]string))   // { before: nil, after: nil)
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	logger.Init(viper.GetInt("log-level"))
	config.InitGlobal(globalConfigFile)
}

func runRoot(cmd *cobra.Command, args []string) {
	cmd.Help()
}
