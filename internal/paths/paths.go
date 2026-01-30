package paths

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Paths struct {
	TargetRoot         string
	TargetPath         string
	TargetConfigFile   string
	TemplateRoot       string
	TemplatePath       string
	TemplateConfigFile string
	GlobalConfigPath   string
	GlobalConfigRoot   string
}

func NewPaths(targetRoot string, targetPath string, templateRoot string, templatePath string, targetConfigFile string, globalConfigFile string, globalConfigRoot string) (*Paths, error) {
	resolvedTargetRoot, err := resolve(targetRoot)
	if err != nil {
		slog.Error("failed to resolve target root", slog.Any("error", err))
		return nil, err
	}

	resolvedTargetPath, err := resolve(targetPath)
	if err != nil {
		slog.Error("failed to resolve target path", slog.Any("error", err))
		return nil, err
	}

	resolvedTargetConfigFile, err := resolve(targetConfigFile)
	if err != nil {
		slog.Error("failed to resolve target configuration file", slog.Any("error", err))
		return nil, err
	}

	resolvedTemplateRoot, err := resolve(templateRoot)
	if err != nil {
		slog.Error("failed to resolve template root", slog.Any("error", err))
		return nil, err
	}

	resolvedTemplatePath, err := resolve(templatePath)
	if err != nil {
		slog.Error("failed to resolve template path", slog.Any("error", err))
		return nil, err
	}

	resolvedTemplateConfigFile, err := resolve(templatePath, TemplateConfigFile)
	if err != nil {
		slog.Error("failed to resolve template configuration file", slog.Any("error", err))
		return nil, err
	}

	resolvedGlobalConfigFile, err := resolve(globalConfigFile)
	if err != nil {
		slog.Error("failed to resolve global configuration file", slog.Any("error", err))
		return nil, err
	}

	resolvedGlobalConfigDir, err := resolve(globalConfigRoot)
	if err != nil {
		slog.Error("failed to resolve global configuration dir", slog.Any("error", err))
		return nil, err
	}
	return &Paths{
		TargetRoot:         resolvedTargetRoot,
		TargetPath:         resolvedTargetPath,
		TargetConfigFile:   resolvedTargetConfigFile,
		TemplateRoot:       resolvedTemplateRoot,
		TemplatePath:       resolvedTemplatePath,
		TemplateConfigFile: resolvedTemplateConfigFile,
		GlobalConfigPath:   resolvedGlobalConfigFile,
		GlobalConfigRoot:   resolvedGlobalConfigDir,
	}, nil
}

func NewPathsFromConfig(config map[string]any) (*Paths, error) {
	targetConfigFile := config["target-config-file"].(string)
	targetPath := config["target-path"].(string)
	templatePath := config["template-path"].(string)
	if templatePath == "" {
		templatePath = filepath.Join(config["template-root"].(string), config["template-name"].(string))
	}

	// If we have a target configuration file then we can work back from there to figure out the target
	// root. If not then we can build the target path - either from being set explicitly
	// or from target root. Once we have target path we can build the 'expected' target config
	// file path.
	if targetConfigFile != "" {
		targetPath = filepath.Dir(targetConfigFile)
	} else {
		if targetPath == "" {
			targetPath = filepath.Join(config["target-root"].(string), config["target-name"].(string))
		}
		targetConfigFile = filepath.Join(targetPath, TargetConfigFileDir, TargetConfigFile)
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
		"targetRoot":         p.TargetRoot,
		"targetPath":         p.TargetPath,
		"targetConfigFile":   p.TargetConfigFile,
		"templateRoot":       p.TemplateRoot,
		"templatePath":       p.TemplatePath,
		"templateConfigFile": p.TemplateConfigFile,
		"globalConfigRoot":   p.GlobalConfigRoot,
		"globalConfigPath":   p.GlobalConfigPath,
	}
}

func resolve(components ...string) (string, error) {
	path := filepath.Join(components...)

	expanded, err := expandPath(path)
	if err != nil {
		return "", err
	}
	absPath, err := filepath.Abs(expanded)
	if err != nil {
		return "", err
	}
	return absPath, nil
}

func expandPath(path string) (string, error) {
	expanded := os.ExpandEnv(path)
	if strings.HasPrefix(expanded, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if expanded == "~" {
			return home, nil
		}
		if strings.HasPrefix(expanded, "~/") {
			expanded = filepath.Join(home, expanded[2:])
		}
	}
	return expanded, nil
}
