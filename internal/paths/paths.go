package paths

import (
	"fmt"
	"log/slog"
	"path/filepath"
)

const configFileName = "proj.yml"
const configFileFolder = ".proj"

type Paths struct {
	TargetRoot string
	TargetPath string
	TargetConfigDir string
	TargetConfigFile string
	TemplateRoot string
	TemplatePath string
	TemplateConfigFile string
}

func NewPaths(targetRoot string, targetPath string, templateRoot string, templatePath string, targetConfigFile string) (*Paths, error) {
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

	resolvedTemplateConfigFile, err := resolve("Template configuration file", templatePath, configFileName)
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
		targetConfigFile = definedOrDefault("Target configuration file", "", targetPath, configFileFolder, configFileName)
	}

	return NewPaths(
		config["target-root"].(string),
		targetPath,
		config["template-root"].(string),
		templatePath,
		targetConfigFile,
	)
}

func definedOrDefault(desc string, configVal string, components ...string) string {
	if configVal != "" {
		slog.Debug("Using defined value", slog.String("Description", desc), slog.String("Value", configVal))
		return configVal
	}
	slog.Debug("Value not provided, constructing it", slog.String("Description", desc), slog.Any("Componets", components))
	return filepath.Join(components...)
}

func resolve(desc string, components ...string) (string, error) {
	path := filepath.Join(components...)
	absPath, err := filepath.Abs(path)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to resolve %s", desc), slog.String("Path", path), slog.Any("Error", err))
	}
	return absPath, err
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

func (p *Paths) DefinitionSourcePath(source string) (string, error) {
	// for resources in the template's definition's files source section: 
	return "eh", nil
}

func (p *Paths) DefinitionTargetPath(target string, variables map[string]any) (string, error) {
	// for resources in the template's definition's files target section: 
	return "eh", nil
}
