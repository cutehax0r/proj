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

	slog.Debug("Configuration loaded", slog.Bool("no-write", cfg.NoWrite))

	if err := cfg.setupPaths(); err != nil {
		return nil, err
	}

	return cfg, nil
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
	fs := afero.NewOsFs()
	cfg := &Config{
		TemplateName:     templateName,
		DefinitionName:   definitionName,
		GlobalConfigFile: viper.GetString("global-config-file"),
		Fs:               fs,
	}

	// If template name not provided, try to find it from project context
	// (unless --all flag is set, which forces global view)
	if cfg.TemplateName == "" && !viper.GetBool("all") {
		targetPath := viper.GetString("target-path")
		projectRoot, err := paths.FindProjectRootWithFS(fs, targetPath)
		if err == nil {
			projYmlPath := filepath.Join(projectRoot, paths.TargetConfigFileDir, paths.TargetConfigFile)
			projData, err := afero.ReadFile(fs, projYmlPath)
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
