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
	TargetConfigPath   string
	TemplateRoot       string
	TemplatePath       string
	TemplateConfigPath string
	GlobalConfigPath   string
	GlobalConfigRoot   string
}

func (p *Paths) ResolveGlobal(path string) (string, error) {
	return resolve(p.GlobalConfigRoot, path)
}

func (p *Paths) ResolveTemplate(path string) (string, error) {
	return resolve(p.TemplatePath, path)
}

func (p *Paths) ResolveTarget(path string) (string, error) {
	return resolve(p.TargetPath, path)
}

func NewPaths(targetRoot string, targetPath string, templateRoot string, templatePath string, targetConfigPath string, globalConfigPath string) (*Paths, error) {
	p := &Paths{}
	var err error

	if p.TargetRoot, err = resolve(targetRoot); err != nil {
		slog.Error("Resolve target root failed", slog.Any("error", err))
		return nil, err
	}

	if p.TargetPath, err = resolve(targetPath); err != nil {
		slog.Error("Resolve target path failed", slog.Any("error", err))
		return nil, err
	}

	if p.TargetConfigPath, err = resolve(targetConfigPath); err != nil {
		slog.Error("Resolve target configuration path failed", slog.Any("error", err))
		return nil, err
	}

	if p.TemplateRoot, err = resolve(templateRoot); err != nil {
		slog.Error("Resolve template root failed", slog.Any("error", err))
		return nil, err
	}

	if p.TemplatePath, err = resolve(templatePath); err != nil {
		slog.Error("Resolve template path failed", slog.Any("error", err))
		return nil, err
	}

	if p.TemplateConfigPath, err = resolve(templatePath, TemplateConfigFile); err != nil {
		slog.Error("Resolve template configuration file failed", slog.Any("error", err))
		return nil, err
	}

	if p.GlobalConfigPath, err = resolve(globalConfigPath); err != nil {
		slog.Error("Resolve global configuration path failed", slog.Any("error", err))
		return nil, err
	}

	p.GlobalConfigRoot = filepath.Dir(p.GlobalConfigPath)

	return p, nil
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
	)
}

func (p *Paths) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("TargetRoot", p.TargetRoot),
		slog.String("TargetPath", p.TargetPath),
		slog.String("TargetConfigPath", p.TargetConfigPath),
		slog.String("TemplateRoot", p.TemplateRoot),
		slog.String("TemplatePath", p.TemplatePath),
		slog.String("TemplateConfigPath", p.TemplateConfigPath),
		slog.String("GlobalConfigRoot", p.GlobalConfigRoot),
		slog.String("GlobalConfigPath", p.GlobalConfigPath),
	)
}

func (p *Paths) ToMap() map[string]string {
	return map[string]string{
		"targetRoot":         p.TargetRoot,
		"targetPath":         p.TargetPath,
		"targetConfigPath":   p.TargetConfigPath,
		"templateRoot":       p.TemplateRoot,
		"templatePath":       p.TemplatePath,
		"templateConfigPath": p.TemplateConfigPath,
		"globalConfigRoot":   p.GlobalConfigRoot,
		"globalConfigPath":   p.GlobalConfigPath,
	}
}

func resolve(components ...string) (string, error) {
	path := filepath.Join(components...)

	expanded, err := expandPath(path)
	if err != nil {
		slog.Error("Resolve failed to expand path", slog.String("Path", path), slog.Any("Error", err))
		return "", err
	}
	absPath, err := filepath.Abs(expanded)
	if err != nil {
		slog.Error("Resolve failed to make path absolute", slog.String("Path", path), slog.Any("Error", err))
		return "", err
	}
	return absPath, nil
}

func expandPath(path string) (string, error) {
	expanded := os.ExpandEnv(path)
	if strings.HasPrefix(expanded, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			slog.Error("Expand path couldn't get home directory", slog.String("Path", path), slog.Any("Error", err))
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
