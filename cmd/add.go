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
	Run:   runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringP("template-root", "s", paths.TemplateRootDir(), "Path containing project templates")
	viper.BindPFlag("template-root", addCmd.Flags().Lookup("template-root"))

	addCmd.Flags().StringP("template-path", "t", "", "Path to read files from")
	viper.BindPFlag("template-path", addCmd.Flags().Lookup("template-path"))

	addCmd.Flags().StringArrayP("set-variable", "v", []string{}, "Set a variable using key=value")
	viper.BindPFlag("set-variables", addCmd.Flags().Lookup("set-variable"))

}

func runAdd(cmd *cobra.Command, args []string) {
	slog.Debug("Execute Add Command", slog.String("DefinitionName", args[0]), slog.String("TargetName", args[1]))

	// maybe we should allow setting target-path so that you can 'apply' an add to an existing
	// project from outside the current working directory. Find project root would start with
	// PWD but if specified it would use target-path
	projectRoot, err := paths.FindProjectRoot()
	if err != nil {
		slog.Error("Failed to find project root", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Debug("Found project root", slog.String("path", projectRoot))

	cfg, err := generator.AddConfig(projectRoot, args[0], args[1])
	if err != nil {
		slog.Error("Failed to setup configuration", slog.Any("error", err))
		os.Exit(1)
	}

	// before writing files we should check that none of the target paths already exist.
	// bail if they're already taken.
	creator, err := generator.NewCreator(cfg)
	if err != nil {
		slog.Error("Failed to create generator", slog.Any("error", err))
		os.Exit(1)
	}

	if err := creator.Create(); err != nil {
		slog.Error("Project creation failed", slog.Any("error", err))
		os.Exit(1)
	}
}
