package cmd

import (
	"os"
	"proj/internal/config"
	"proj/internal/logger"
	"proj/internal/paths"

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
	viper.SetConfigName(paths.GlobalConfigFile)

	rootCmd.PersistentFlags().BoolP("no-write", "w", false, "Print the plan and don't write anything")
	rootCmd.PersistentFlags().IntP("log-level", "l", 0, "How much to log [0-3], bigger = more")
	rootCmd.PersistentFlags().StringVarP(&globalConfigFile, "global-config-file", "g", "", "Use specific global configuration file")
	rootCmd.PersistentFlags().String("template-git", "https://github.com/", "Default git source for templates")
	viper.BindPFlags(rootCmd.PersistentFlags())
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	logLevel := 0
	if cmd.Flags().Changed("log-level") {
		logLevel = viper.GetInt("log-level")
	}

	logger.Init(logLevel)

	config.InitGlobal(globalConfigFile)

	if !cmd.Flags().Changed("log-level") && viper.IsSet("log-level") {
		configLogLevel := viper.GetInt("log-level")
		if configLogLevel != logLevel {
			logger.Init(configLogLevel)
		}
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	cmd.Help()
}
