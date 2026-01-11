package paths

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

type Paths struct {
	TargetRoot string
	TargetPath string
	TargetConfigDir string
	TargetConfigFile string
	TemplateRoot string
	TemplatePath string
	TemplateConfigFile string
	GlobalConfigFile string
}

func NewPaths(targetRoot string, targetPath string, templateRoot string, templatePath string, targetConfigFile string, globalConfigFile string) (*Paths, error) {
	resolvedTargetRoot, err := resolve("target root", targetRoot)
	if err != nil {
		return nil, err
	}

	resolvedTargetPath, err := resolve("target path", targetPath)
	if err != nil {
		return nil, err
	}

	resolvedTargetConfigFile, err := resolve("target configuration file", targetConfigFile)
	if err != nil {
		return nil, err
	}

	resolvedTemplateRoot, err := resolve("template root", templateRoot)
	if err != nil {
		return nil, err
	}

	resolvedTemplatePath, err := resolve("template path", templatePath)
	if err != nil {
		return nil, err
	}

	resolvedTemplateConfigFile, err := resolve("Template configuration file", templatePath, TemplateConfigFile)
	if err != nil {
		return nil, err
	}

	resolvedGlobalConfigFile, err := resolve("Global configuration file", globalConfigFile)
	if err != nil {
		return nil, err
	}

	return &Paths {
		TargetRoot: resolvedTargetRoot,
		TargetPath: resolvedTargetPath,
		TargetConfigFile: resolvedTargetConfigFile,
		TemplateRoot: resolvedTemplateRoot,
		TemplatePath: resolvedTemplatePath,
		TemplateConfigFile: resolvedTemplateConfigFile,
		GlobalConfigFile: resolvedGlobalConfigFile,
	}, nil
}

func NewPathsFromConfig(config map[string]any) (*Paths, error) {
	targetConfigFile := config["target-config-file"].(string)
	targetPath := config["target-path"].(string)
	templatePath := definedOrDefault("Template path", config["template-path"].(string), config["template-root"].(string), config["template-name"].(string))

	// If we have a target configuration file then we can work back from there to figure out the target
	// root. If not then we can build the target path - either from being set explicitly
	// or from target root. Once we have target path we can build the 'expected' target config
	// file path.
	if targetConfigFile != "" {
		targetPath = filepath.Dir(targetConfigFile)
	} else {
		targetPath = definedOrDefault("Target path", targetPath, config["target-root"].(string), config["target-name"].(string))
		targetConfigFile = definedOrDefault("Target configuration file", "", targetPath, TargetConfigFileDir, targetConfigFile)
	}

	return NewPaths(
		config["target-root"].(string),
		targetPath,
		config["template-root"].(string),
		templatePath,
		targetConfigFile,
		config["global-config-file"].(string),
	)
}

func (p *Paths) DefinitionSourcePath(source string) (string, error) {
	// for resources in the template's definition's files source section: 
	return "templatepath+source", nil
}

func (p *Paths) DefinitionTargetPath(target string, variables map[string]any) (string, error) {
	// for resources in the template's definition's files target section: 
	return "targetpath+templateify(target,vars)", nil
}

func (p *Paths) LogValue() slog.Value {
	if p == nil {
		return slog.Value{}
	}
	return slog.GroupValue(
		slog.String("TargetRoot", p.TargetRoot),
		slog.String("TargetPath", p.TargetPath),
		slog.String("TargetConfigFile", p.TargetConfigFile),
		slog.String("TemplateRoot", p.TemplateRoot),
		slog.String("TemplatePath", p.TemplatePath),
		slog.String("TemplateConfigFile", p.TemplateConfigFile),
	)
}


func definedOrDefault(desc string, configVal string, components ...string) string {
	if configVal != "" {
		slog.Debug(fmt.Sprintf("%s provided, using it",  desc), slog.String("Path", configVal))
		return configVal
	}
	path := filepath.Join(components...)
	slog.Debug(fmt.Sprintf("%s not provided, constructing it", desc), slog.Any("Componets", components), slog.String("Path", path))
	return path
}

