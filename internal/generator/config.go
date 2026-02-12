package generator

import (
	"log/slog"
	"path/filepath"
	"proj/internal/config"
	"proj/internal/paths"

	"github.com/spf13/afero"
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
	Fs    afero.Fs
}

func NewConfig(templateName, targetName string) (*Config, error) {
	cfg := &Config{
		TemplateName:     templateName,
		TargetName:       targetName,
		DefinitionName:   viper.GetString("definition-name"),
		SetVariables:     viper.GetStringSlice("set-variable"),
		NoWrite:          viper.GetBool("no-write"),
		GlobalConfigFile: viper.GetString("global-config-file"),
		Fs:               afero.NewOsFs(),
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

func AddConfig(kind, name string) (*Config, error) {
	// Find project root from current directory
	targetPath := viper.GetString("target-path")
	fs := afero.NewOsFs()
	projectRoot, err := paths.FindProjectRootWithFS(fs, targetPath)
	if err != nil {
		return nil, err
	}

	projYmlPath := filepath.Join(projectRoot, paths.TargetConfigFileDir, paths.TargetConfigFile)
	projData, err := afero.ReadFile(fs, projYmlPath)
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
		Fs:               afero.NewOsFs(),
	}

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
