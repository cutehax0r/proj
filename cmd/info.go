package cmd

import (
	"fmt"
	"log/slog"
	"os"
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

	viper.BindPFlags(infoCmd.Flags())
}

func runInfo(cmd *cobra.Command, args []string) {
	templateRoot, err := paths.Resolve(viper.GetString("template-root"))
	if err != nil {
		slog.Error("Failed to resolve template root", slog.Any("error", err))
		os.Exit(1)
	}

	catalog := info.NewCatalog()
	projectCtx, err := catalog.FindProjectContext(".")
	if err != nil {
		slog.Error("Failed to inspect project context", slog.Any("error", err))
		os.Exit(1)
	}

	switch len(args) {
	case 0:
		runInfoNoArgs(catalog, templateRoot, projectCtx)
	case 1:
		runInfoTemplate(catalog, templateRoot, args[0], projectCtx)
	case 2:
		runInfoDefinition(catalog, templateRoot, args[0], args[1], projectCtx)
	}
}

func runInfoNoArgs(catalog *info.Catalog, templateRoot string, projectCtx *info.ProjectContext) {
	if projectCtx == nil {
		templates, err := catalog.ListTemplates(templateRoot)
		if err != nil {
			slog.Error("Failed to list templates", slog.String("template-root", templateRoot), slog.Any("error", err))
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

func runInfoTemplate(catalog *info.Catalog, templateRoot, templateName string, projectCtx *info.ProjectContext) {
	summaries, err := catalog.TemplateSummaries(templateRoot, templateName, projectCtx)
	if err != nil {
		slog.Error("Failed to read template details", slog.String("template", templateName), slog.Any("error", err))
		os.Exit(1)
	}

	pageData := info.BuildTemplatePageData(templateName, summaries)
	page, err := info.RenderTemplatePage(pageData)
	if err != nil {
		slog.Error("Failed to render template page", slog.Any("error", err))
		os.Exit(1)
	}
	fmt.Print(page)
}

func runInfoDefinition(catalog *info.Catalog, templateRoot, templateName, definitionName string, projectCtx *info.ProjectContext) {
	detail, err := catalog.DefinitionDetails(templateRoot, templateName, definitionName, projectCtx)
	if err != nil {
		slog.Error("Failed to inspect definition", slog.String("template", templateName), slog.String("definition", definitionName), slog.Any("error", err))
		os.Exit(1)
	}

	pageData := info.BuildDefinitionPageData(templateName, detail)
	page, err := info.RenderDefinitionPage(pageData)
	if err != nil {
		slog.Error("Failed to render definition page", slog.Any("error", err))
		os.Exit(1)
	}
	fmt.Print(page)
}
