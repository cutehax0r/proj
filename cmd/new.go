package cmd

import (
	"log/slog"
	"os"
	"proj/internal/generator"
	"proj/internal/paths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newCmd = &cobra.Command{
	Use:   "new <kind> <name>",
	Args:  cobra.ExactArgs(2),
	Short: "Create a project called NAME based on the KIND template",
	Long:  `Create a new project based on a project template in the current directory.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	Run: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("target-root", "r", ".", "Path to create the project in (pwd)")
	newCmd.Flags().StringP("target-path", "p", "", "Path to write files at (pwd/foobar)")
	newCmd.Flags().StringP("template-root", "s", paths.TemplateRootDir(), "Path containing project templates")
	newCmd.Flags().StringP("template-path", "t", "", "Path to read files from")
	newCmd.Flags().StringArrayP("set-variable", "v", []string{}, "Set a variable using key=value")
	newCmd.Flags().StringP("definition-name", "d", "new", "Definition in template to use")

	viper.BindPFlags(newCmd.Flags())
}

func runNew(cmd *cobra.Command, args []string) {
	slog.Debug("Execute New Command", slog.String("TemplateName", args[0]), slog.String("TargetName", args[1]))

	cfg, err := generator.NewConfig(args[0], args[1])
	if err != nil {
		slog.Error("Failed to setup configuration", slog.Any("error", err))
		os.Exit(1)
	}

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
