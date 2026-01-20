package paths

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

type Paths struct {
	TargetRoot string
	TargetPath string
	TargetConfigFile string
	TemplateRoot string
	TemplatePath string
	TemplateConfigFile string
	GlobalConfigPath string
	GlobalConfigRoot string
}

func NewPaths(targetRoot string, targetPath string, templateRoot string, templatePath string, targetConfigFile string, globalConfigFile string, globalConfigRoot string) (*Paths, error) {
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

	resolvedGlobalConfigDir, err := resolve("Global configuration dir", globalConfigRoot)
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
		GlobalConfigPath: resolvedGlobalConfigFile,
		GlobalConfigRoot: resolvedGlobalConfigDir,
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
		targetConfigFile = definedOrDefault("Target configuration file", "", targetPath, TargetConfigFileDir, TargetConfigFile)
	}

	return NewPaths(
		config["target-root"].(string),
		targetPath,
		config["template-root"].(string),
		templatePath,
		targetConfigFile,
		config["global-config-file"].(string),
		filepath.Dir(config["global-config-file"].(string)),
	)
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
		slog.String("GlobalConfigRoot", p.GlobalConfigRoot),
		slog.String("GlobalConfigPath", p.GlobalConfigPath),
	)
}

func (p *Paths) ToMap() map[string]string {
	return map[string]string{
		"targetRoot": p.TargetRoot,
		"targetPath": p.TargetPath,
		"templateRoot": p.TemplateRoot,
		"templatePath": p.TemplatePath,
		"templateConfigFile": p.TemplateConfigFile,
		"globalConfigRoot": p.GlobalConfigRoot,
		"globalConfigPath": p.GlobalConfigPath,
	}
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
