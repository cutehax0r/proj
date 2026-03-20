package info

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"
)

//go:embed templates/templates.tmpl.md
var templatesPageTemplate string

//go:embed templates/template.tmpl.md
var templatePageTemplate string

//go:embed templates/definition.tmpl.md
var definitionPageTemplate string

type TemplatesPageData struct {
	InProject        bool
	TemplateNames    []string
	CurrentTemplate  string
	LocalDefinitions []string
}

type TemplatePageData struct {
	TemplateName string
	NewDefs      []DefinitionSummary
	AddDefs      []DefinitionSummary
}

type DefinitionVariableRow struct {
	Name     string
	Type     string
	Required bool
	Value    any
	Default  any
}

type DefinitionPageData struct {
	TemplateName   string
	DefinitionName string
	Source         string
	Targets        []string
	Variables      []DefinitionVariableRow
}

func RenderTemplatesPage(data TemplatesPageData) (string, error) {
	return renderPage("templates", templatesPageTemplate, data)
}

func RenderTemplatePage(data TemplatePageData) (string, error) {
	return renderPage("template", templatePageTemplate, data)
}

func RenderDefinitionPage(data DefinitionPageData) (string, error) {
	return renderPage("definition", definitionPageTemplate, data)
}

func BuildTemplatePageData(templateName string, summaries []DefinitionSummary) TemplatePageData {
	newDefs := make([]DefinitionSummary, 0)
	addDefs := make([]DefinitionSummary, 0)
	for _, s := range summaries {
		if s.Local {
			addDefs = append(addDefs, s)
			continue
		}
		newDefs = append(newDefs, s)
	}

	return TemplatePageData{
		TemplateName: templateName,
		NewDefs:      newDefs,
		AddDefs:      addDefs,
	}
}

func BuildDefinitionPageData(templateName string, detail *DefinitionDetail) DefinitionPageData {
	vars := make([]DefinitionVariableRow, 0, len(detail.Variables))
	for _, v := range detail.Variables {
		vars = append(vars, DefinitionVariableRow{
			Name:     v.Name,
			Type:     "",
			Required: true,
			Value:    nil,
			Default:  v.Default,
		})
	}

	return DefinitionPageData{
		TemplateName:   templateName,
		DefinitionName: detail.Name,
		Source:         detail.Source,
		Targets:        detail.Targets,
		Variables:      vars,
	}
}

func renderPage(name, tmpl string, data any) (string, error) {
	t, err := template.New(name).Funcs(template.FuncMap{
		"valueCell": formatValueCell,
		"boolCell": func(v bool) string {
			if v {
				return "yes"
			}
			return "no"
		},
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return "", err
	}
	return out.String(), nil
}

func formatValueCell(value any, fallback any) string {
	if value != nil {
		return fmt.Sprintf("%v", value)
	}
	if fallback != nil {
		return fmt.Sprintf("%v (default)", fallback)
	}
	return ""
}
