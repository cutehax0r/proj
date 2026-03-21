package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"proj/internal/generator"
	"proj/internal/info"
	"proj/internal/paths"
	"sort"

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

	viper.BindPFlags(infoCmd.Flags())
}

func runInfo(cmd *cobra.Command, args []string) {
	var templateName, definitionName string

	if len(args) > 0 {
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

	catalog := info.NewCatalog()

	switch {
	case cfg.TemplateName == "" && cfg.DefinitionName == "":
		runInfoNoArgs(catalog, cfg)
	case cfg.DefinitionName == "":
		runInfoTemplate(catalog, cfg)
	default:
		runInfoDefinition(catalog, cfg)
	}
}

func runInfoNoArgs(catalog *info.Catalog, cfg *generator.Config) {
	// Check if we're in a project by looking at paths
	inProject := cfg.Paths.TargetConfigPath != ""

	if !inProject {
		templates, err := catalog.ListTemplates(cfg.Paths.TemplateRoot)
		if err != nil {
			slog.Error("Failed to list templates", slog.String("template-root", cfg.Paths.TemplateRoot), slog.Any("error", err))
			os.Exit(1)
		}
		page, err := info.RenderTemplatesPage(info.TemplatesPageData{TemplateNames: templates})
		if err != nil {
			slog.Error("Failed to render templates page", slog.Any("error", err))
			os.Exit(1)
		}
		fmt.Print(page)
		return
	}

	// We're in a project, get project context
	projectCtx, err := catalog.FindProjectContext(cfg.Paths.TargetRoot)
	if err != nil {
		slog.Error("Failed to inspect project context", slog.Any("error", err))
		os.Exit(1)
	}

	defNames := make([]string, 0, len(projectCtx.Definitions))
	for name := range projectCtx.Definitions {
		defNames = append(defNames, name)
	}
	sort.Strings(defNames)

	page, err := info.RenderTemplatesPage(info.TemplatesPageData{
		InProject:        true,
		CurrentTemplate:  projectCtx.TemplateName,
		LocalDefinitions: defNames,
	})
	if err != nil {
		slog.Error("Failed to render project templates page", slog.Any("error", err))
		os.Exit(1)
	}
	fmt.Print(page)
}

func runInfoTemplate(catalog *info.Catalog, cfg *generator.Config) {
	// Get project context if in a project
	projectCtx, _ := catalog.FindProjectContext(cfg.Paths.TargetRoot)

	summaries, err := catalog.TemplateSummaries(cfg.Paths.TemplateRoot, cfg.TemplateName, projectCtx)
	if err != nil {
		slog.Error("Failed to read template details", slog.String("template", cfg.TemplateName), slog.Any("error", err))
		os.Exit(1)
	}

	pageData := info.BuildTemplatePageData(cfg.TemplateName, summaries)
	page, err := info.RenderTemplatePage(pageData)
	if err != nil {
		slog.Error("Failed to render template page", slog.Any("error", err))
		os.Exit(1)
	}
	fmt.Print(page)
}

func runInfoDefinition(catalog *info.Catalog, cfg *generator.Config) {
	// Get project context if in a project
	projectCtx, _ := catalog.FindProjectContext(cfg.Paths.TargetRoot)

	detail, err := catalog.DefinitionDetails(cfg.Paths.TemplateRoot, cfg.TemplateName, cfg.DefinitionName, projectCtx)
	if err != nil {
		slog.Error("Failed to inspect definition", slog.String("template", cfg.TemplateName), slog.String("definition", cfg.DefinitionName), slog.Any("error", err))
		os.Exit(1)
	}

	pageData := info.BuildDefinitionPageData(cfg.TemplateName, detail)
	page, err := info.RenderDefinitionPage(pageData)
	if err != nil {
		slog.Error("Failed to render definition page", slog.Any("error", err))
		os.Exit(1)
	}
	fmt.Print(page)
}
