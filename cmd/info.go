package cmd

import (
	"log/slog"
	"os"
	"proj/internal/generator"
	"proj/internal/info"
	"proj/internal/paths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use:   "info [template] [definition]",
	Args:  cobra.RangeArgs(0, 2),
	Short: "Show templates and definitions",
	Long:  "Inspect available templates, their definitions, and definition details.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	Run: runInfo,
}

func init() {
	rootCmd.AddCommand(infoCmd)

	infoCmd.Flags().StringP("template-root", "s", paths.TemplateRootDir(), "Path containing project templates")
	infoCmd.Flags().StringP("target-path", "p", "", "Directory containing .proj/proj.yml (defaults to current directory)")
	infoCmd.Flags().BoolP("all", "a", false, "Show all templates (force global view)")

	viper.BindPFlags(infoCmd.Flags())
}

func runInfo(cmd *cobra.Command, args []string) {
	var templateName, definitionName string

	if len(args) > 0 && !viper.GetBool("all") {
		templateName = args[0]
	}
	if len(args) > 1 {
		definitionName = args[1]
	}

	slog.Debug("Execute Info Command", slog.String("TemplateName", templateName), slog.String("DefinitionName", definitionName))

	cfg, err := generator.InfoConfig(templateName, definitionName)
	if err != nil {
		slog.Error("Failed to setup configuration", slog.Any("error", err))
		os.Exit(1)
	}

	explainer := info.NewExplainer(cfg)

	if viper.GetBool("all") {
		if err := explainer.ExplainGlobal(); err != nil {
			slog.Error("Failed to explain", slog.Any("error", err))
			os.Exit(1)
		}
		return
	}

	if err := explainer.Explain(); err != nil {
		slog.Error("Failed to explain", slog.Any("error", err))
		os.Exit(1)
	}
}
