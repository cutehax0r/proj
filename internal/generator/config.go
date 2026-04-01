package generator

import (
	"log/slog"
	"os"
	"path/filepath"
	"proj/internal/config"
	"proj/internal/paths"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	TemplateName     string
	TargetName       string
	DefinitionName   string
	SetVariables     []string
	NoWrite          bool
	Force            bool
	GlobalConfigFile string

	Paths *paths.Paths
}

func NewConfig(templateName, targetName string) (*Config, error) {
	cfg := &Config{
		TemplateName:     templateName,
		TargetName:       targetName,
		DefinitionName:   viper.GetString("definition-name"),
		SetVariables:     viper.GetStringSlice("set-variable"),
		NoWrite:          viper.GetBool("no-write"),
		GlobalConfigFile: viper.GetString("global-config-file"),
	}

	slog.Debug("Configuration loaded", slog.Bool("no-write", cfg.NoWrite))

	if err := cfg.setupPaths(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func AddConfig(kind, name string) (*Config, error) {
	// Find project root from current directory
	targetPath := viper.GetString("target-path")
	projectRoot, err := paths.FindProjectRoot(targetPath)
	if err != nil {
		return nil, err
	}

	projYmlPath := filepath.Join(projectRoot, paths.TargetConfigFileDir, paths.TargetConfigFile)
	projData, err := os.ReadFile(projYmlPath)
	if err != nil {
		return nil, err
	}

	targetCfg := &config.TargetConfig{}
	if err := yaml.Unmarshal(projData, targetCfg); err != nil {
		return nil, err
	}

	cfg := &Config{
		TemplateName:     targetCfg.TemplateName,
		TargetName:       name,
		DefinitionName:   kind,
		SetVariables:     viper.GetStringSlice("set-variable"),
		NoWrite:          viper.GetBool("no-write"),
		GlobalConfigFile: viper.GetString("global-config-file"),
	}

	slog.Debug("Configuration loaded", slog.Bool("no-write", cfg.NoWrite))

	viper.Set("target-root", projectRoot)
	viper.Set("target-path", projectRoot) // For 'add', files should be written relative to projectRoot
	viper.Set("template-name", cfg.TemplateName)
	viper.Set("target-name", cfg.TargetName)
	viper.Set("target-config-file", projYmlPath)

	if err := cfg.setupPaths(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func InfoConfig(templateName, definitionName string) (*Config, error) {
	cfg := &Config{
		TemplateName:     templateName,
		DefinitionName:   definitionName,
		GlobalConfigFile: viper.GetString("global-config-file"),
	}

	// If template name not provided, try to find it from project context
	if cfg.TemplateName == "" {
		targetPath := viper.GetString("target-path")
		projectRoot, err := paths.FindProjectRoot(targetPath)
		if err == nil {
			projYmlPath := filepath.Join(projectRoot, paths.TargetConfigFileDir, paths.TargetConfigFile)
			projData, err := os.ReadFile(projYmlPath)
			if err == nil {
				targetCfg := &config.TargetConfig{}
				if err := yaml.Unmarshal(projData, targetCfg); err == nil {
					cfg.TemplateName = targetCfg.TemplateName
					viper.Set("target-root", projectRoot)
					viper.Set("target-path", projectRoot)
					viper.Set("target-config-file", projYmlPath)
				}
			}
		}
	}

	slog.Debug("Info configuration loaded", slog.String("template", cfg.TemplateName), slog.String("definition", cfg.DefinitionName))

	// Set template name in viper so setupPaths can use it
	viper.Set("template-name", cfg.TemplateName)
	viper.Set("target-name", "") // Info doesn't need a target name

	if err := cfg.setupPaths(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func InstallerConfig(templateName, targetName, templateRoot, templateGit string) (*Config, error) {
	cfg := &Config{
		TemplateName:     templateName,
		TargetName:       targetName,
		DefinitionName:   templateGit,
		NoWrite:          viper.GetBool("no-write"),
		Force:            viper.GetBool("force"),
		GlobalConfigFile: viper.GetString("global-config-file"),
	}

	return cfg, nil
}

func UninstallerConfig(targetName, templateRoot string) (*Config, error) {
	viper.Set("template-root", templateRoot)
	viper.Set("template-name", targetName)
	viper.Set("target-name", targetName)
	viper.Set("target-config-file", "")

	cfg := &Config{
		TargetName:       targetName,
		TemplateName:     targetName,
		NoWrite:          viper.GetBool("no-write"),
		Force:            viper.GetBool("force"),
		GlobalConfigFile: viper.GetString("global-config-file"),
	}

	if err := cfg.setupPaths(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) setupPaths() error {
	viper.Set("template-name", cfg.TemplateName)
	viper.Set("target-name", cfg.TargetName)
	viper.Set("target-config-file", "") // this should not exist yet if we're running new

	p, err := paths.NewPathsFromConfig(viper.AllSettings())
	if err != nil {
		slog.Error("Failed to setup paths", slog.Any("error", err))
		return err
	}
	cfg.Paths = p

	// Set the template path in viper so it can be used during file spec resolution
	viper.Set("template-path", p.TemplatePath)

	return nil
}
