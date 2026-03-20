package info

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"proj/internal/config"
	"proj/internal/paths"
	"sort"
	"strings"

	"go.yaml.in/yaml/v3"
)

type Catalog struct{}

type ProjectContext struct {
	Root         string
	TemplateName string
	Definitions  map[string]any
}

type DefinitionSummary struct {
	Name   string
	Local  bool
	Source string
}

type VariableDetail struct {
	Name    string
	Default any
}

type DefinitionDetail struct {
	Name      string
	Local     bool
	Source    string
	Targets   []string
	Variables []VariableDetail
}

type templateConfig struct {
	Definitions map[string]templateDefinition `yaml:"definitions"`
}

type templateDefinition struct {
	Requirements templateRequirements `yaml:"requirements"`
	Files        []templateFile       `yaml:"files"`
}

type templateRequirements struct {
	Local     bool               `yaml:"local"`
	Variables []templateVariable `yaml:"variables"`
}

type templateVariable struct {
	Name    string `yaml:"name"`
	Default any    `yaml:"default"`
}

type templateFile struct {
	Target string `yaml:"target"`
}

func NewCatalog() *Catalog {
	return &Catalog{}
}

func (c *Catalog) FindProjectContext(startPath string) (*ProjectContext, error) {
	startAbs, err := paths.Resolve(startPath)
	if err != nil {
		return nil, err
	}

	projectRoot, err := paths.FindProjectRoot(startAbs)
	if err != nil {
		if errors.Is(err, paths.ErrNotInProjDirectory) {
			return nil, nil
		}
		return nil, err
	}

	cfgPath := filepath.Join(projectRoot, paths.TargetConfigFileDir, paths.TargetConfigFile)
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	targetCfg := &config.TargetConfig{}
	if err := yaml.Unmarshal(data, targetCfg); err != nil {
		return nil, err
	}

	return &ProjectContext{
		Root:         projectRoot,
		TemplateName: targetCfg.TemplateName,
		Definitions:  targetCfg.Definitions,
	}, nil
}

func (c *Catalog) ListTemplates(templateRoot string) ([]string, error) {
	entries, err := os.ReadDir(templateRoot)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		templatePath := filepath.Join(templateRoot, entry.Name(), paths.TemplateConfigFile)
		if _, err := os.Stat(templatePath); err == nil {
			result = append(result, entry.Name())
		}
	}

	sort.Strings(result)
	return result, nil
}

func (c *Catalog) TemplateSummaries(templateRoot, templateName string, projectCtx *ProjectContext) ([]DefinitionSummary, error) {
	defs, err := c.loadTemplateDefinitions(templateRoot, templateName)
	if err != nil {
		return nil, err
	}

	result := make([]DefinitionSummary, 0, len(defs))
	for name, def := range defs {
		result = append(result, DefinitionSummary{
			Name:   name,
			Local:  def.Requirements.Local,
			Source: "template",
		})
	}

	if projectCtx != nil && projectCtx.TemplateName == templateName {
		for name, rawDef := range projectCtx.Definitions {
			localDef := parseLocalDefinition(rawDef)
			result = append(result, DefinitionSummary{
				Name:   name,
				Local:  localDef.Requirements.Local,
				Source: "local",
			})
		}
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].Local != result[j].Local {
			return !result[i].Local
		}
		if result[i].Name != result[j].Name {
			return result[i].Name < result[j].Name
		}
		return result[i].Source < result[j].Source
	})

	return result, nil
}

func (c *Catalog) DefinitionDetails(templateRoot, templateName, definitionName string, projectCtx *ProjectContext) (*DefinitionDetail, error) {
	if projectCtx != nil && projectCtx.TemplateName == templateName {
		if rawDef, ok := projectCtx.Definitions[definitionName]; ok {
			localDef := parseLocalDefinition(rawDef)
			return buildDefinitionDetail(definitionName, localDef.Requirements.Local, "local", localDef.Files, localDef.Requirements.Variables), nil
		}
	}

	defs, err := c.loadTemplateDefinitions(templateRoot, templateName)
	if err != nil {
		return nil, err
	}

	def, ok := defs[definitionName]
	if !ok {
		return nil, fmt.Errorf("definition %q not found in template %q", definitionName, templateName)
	}

	targets := make([]string, 0, len(def.Files))
	for _, f := range def.Files {
		targets = append(targets, f.Target)
	}
	variables := make([]VariableDetail, 0, len(def.Requirements.Variables))
	for _, v := range def.Requirements.Variables {
		variables = append(variables, VariableDetail{Name: v.Name, Default: v.Default})
	}

	return buildDefinitionDetail(definitionName, def.Requirements.Local, "template", targets, variables), nil
}

func buildDefinitionDetail(name string, local bool, source string, targets []string, variables []VariableDetail) *DefinitionDetail {
	sort.Strings(targets)
	sort.Slice(variables, func(i, j int) bool { return variables[i].Name < variables[j].Name })

	return &DefinitionDetail{
		Name:      name,
		Local:     local,
		Source:    source,
		Targets:   targets,
		Variables: variables,
	}
}

func (c *Catalog) loadTemplateDefinitions(templateRoot, templateName string) (map[string]templateDefinition, error) {
	templateConfigPath, err := paths.Resolve(templateRoot, templateName, paths.TemplateConfigFile)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(templateConfigPath)
	if err != nil {
		return nil, err
	}

	tpl := &templateConfig{}
	if err := yaml.Unmarshal(data, tpl); err != nil {
		return nil, err
	}

	if tpl.Definitions == nil {
		return map[string]templateDefinition{}, nil
	}

	return tpl.Definitions, nil
}

type localDefinition struct {
	Requirements localRequirements
	Files        []string
}

type localRequirements struct {
	Local     bool
	Variables []VariableDetail
}

func parseLocalDefinition(raw any) localDefinition {
	result := localDefinition{
		Requirements: localRequirements{Local: true},
	}

	defMap, ok := raw.(map[string]any)
	if !ok {
		return result
	}

	if reqRaw, ok := defMap["requirements"].(map[string]any); ok {
		if local, ok := reqRaw["local"].(bool); ok {
			result.Requirements.Local = local
		}
		if varsRaw, ok := reqRaw["variables"].([]any); ok {
			for _, item := range varsRaw {
				if varMap, ok := item.(map[string]any); ok {
					name, _ := varMap["name"].(string)
					if strings.TrimSpace(name) == "" {
						continue
					}
					result.Requirements.Variables = append(result.Requirements.Variables, VariableDetail{Name: name, Default: varMap["default"]})
				}
			}
		}
	}

	if filesRaw, ok := defMap["files"].([]any); ok {
		for _, item := range filesRaw {
			if fileMap, ok := item.(map[string]any); ok {
				target, _ := fileMap["target"].(string)
				if strings.TrimSpace(target) == "" {
					continue
				}
				result.Files = append(result.Files, target)
			}
		}
	}

	return result
}
