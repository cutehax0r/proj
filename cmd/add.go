package cmd

import (
	"log/slog"
	"os"
	"proj/internal/generator"
	"proj/internal/paths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addCmd = &cobra.Command{
	Use:   "add <kind> <name>",
	Args:  cobra.ExactArgs(2),
	Short: "Add file(s) called NAME based on the KIND definition from the template used to create this project",
	Long:  `Add file(s) called NAME based on the KIND definition from the template used to create this project`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	Run: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	cwd, err := os.Getwd()
	if err != nil {
		cwd = ""
	}

	addCmd.Flags().StringP("target-path", "r", cwd, "Directory containing .proj/proj.yml")
	addCmd.Flags().StringP("template-root", "s", paths.TemplateRootDir(), "Path containing project templates")
	addCmd.Flags().StringP("template-path", "t", "", "Path to read files from")
	addCmd.Flags().StringArrayP("set-variable", "v", []string{}, "Set a variable using key=value")

	viper.BindPFlags(addCmd.Flags())
}

func runAdd(cmd *cobra.Command, args []string) {
	slog.Debug("Execute Add Command", slog.String("DefinitionName", args[0]), slog.String("TargetName", args[1]))

	cfg, err := generator.AddConfig(args[0], args[1])
	if err != nil {
		slog.Error("Failed to setup configuration", slog.Any("error", err))
		os.Exit(1)
	}

	adder, err := generator.NewAdder(cfg)
	if err != nil {
		slog.Error("Failed to create adder", slog.Any("error", err))
		os.Exit(1)
	}

	if err := adder.Add(); err != nil {
		slog.Error("Adding files failed", slog.Any("error", err))
		os.Exit(1)
	}
}
