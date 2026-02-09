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
	GlobalConfigFile string

	Paths *paths.Paths
}

func NewConfig(templateName, targetName string) (*Config, error) {
	cfg := &Config{
		TemplateName:     templateName,
		TargetName:       targetName,
		DefinitionName:   viper.GetString("definition-name"),
		SetVariables:     viper.GetStringSlice("set-variables"),
		NoWrite:          viper.GetBool("no-write"),
		GlobalConfigFile: viper.GetString("global-config-file"),
	}

	// Compute paths once
	if err := cfg.setupPaths(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) setupPaths() error {
	// Set values in viper for path resolution
	viper.Set("template-name", cfg.TemplateName)
	viper.Set("target-name", cfg.TargetName)
	viper.Set("target-config-file", "") // this should not exist yet if we're running new

	p, err := paths.NewPathsFromConfig(viper.AllSettings())
	if err != nil {
		slog.Error("Failed to setup paths", slog.Any("error", err))
		return err
	}
	cfg.Paths = p
	return nil
}

// AddConfig creates a Config for the add command by loading the existing project's
// proj.yml and using it to determine the template name and other settings.
func AddConfig(projectRoot, kind, name string) (*Config, error) {
	// Load the existing proj.yml
	projYmlPath := filepath.Join(projectRoot, paths.TargetConfigFileDir, paths.TargetConfigFile)
	projData, err := os.ReadFile(projYmlPath)
	if err != nil {
		slog.Error("Failed to read proj.yml", slog.String("path", projYmlPath), slog.Any("error", err))
		return nil, err
	}

	// Unmarshal into TargetConfig to extract template name
	targetCfg := &config.TargetConfig{}
	if err := yaml.Unmarshal(projData, targetCfg); err != nil {
		slog.Error("Failed to unmarshal proj.yml", slog.String("path", projYmlPath), slog.Any("error", err))
		return nil, err
	}

	// Create Config with the template name from proj.yml
	cfg := &Config{
		TemplateName:     targetCfg.TemplateName,
		TargetName:       name,
		DefinitionName:   kind,
		SetVariables:     viper.GetStringSlice("set-variables"),
		NoWrite:          viper.GetBool("no-write"),
		GlobalConfigFile: viper.GetString("global-config-file"),
	}

	// Set up paths with the project root as target-root
	viper.Set("target-root", projectRoot)
	viper.Set("template-name", cfg.TemplateName)
	viper.Set("target-name", cfg.TargetName)
	viper.Set("target-config-file", projYmlPath)

	if err := cfg.setupPaths(); err != nil {
		return nil, err
	}

	return cfg, nil
}
