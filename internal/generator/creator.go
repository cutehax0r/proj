package generator

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"proj/internal/config"
	"proj/internal/luabridge"
	"proj/internal/paths"
	"strings"

	"github.com/spf13/viper"
	"go.yaml.in/yaml/v3"
)

type Creator struct {
	cfg     *Config
	paths   *paths.Paths
	reqs    *config.RequirementSpec
	scripts config.ScriptSpec
	files   []config.FileSpec
	vars    map[string]any
	luaenv  *luabridge.Runtime
}

func NewCreator(cfg *Config) (*Creator, error) {
	return &Creator{
		cfg:   cfg,
		paths: cfg.Paths,
	}, nil
}

func (c *Creator) Create() error {
	if err := c.setupConfig(); err != nil {
		return err
	}

	if err := c.runBeforeScripts(); err != nil {
		return err
	}

	if err := c.processFiles(); err != nil {
		return err
	}

	if err := c.runAfterScripts(); err != nil {
		return err
	}

	c.vars = c.luaenv.GetVariables()
	slog.Debug("Variables after after-scripts", slog.Any("vars", c.vars))

	if err := c.writeConfig(); err != nil {
		return err
	}
	c.luaenv.CloseState()

	return nil
}

func (c *Creator) writeConfig() error {
	projDir := filepath.Join(c.paths.TargetPath, ".proj")
	if err := os.MkdirAll(projDir, 0755); err != nil {
		slog.Error("Failed to create .proj directory", slog.String("directory", projDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created .proj directory", slog.String("directory", projDir))

	targetConfig := &config.TargetConfig{
		TemplateName: c.cfg.TemplateName,
		Variables:    c.vars,
		Scripts:      make(map[string]any),
		Definitions:  make(map[string]any),
	}

	yamlData, err := yaml.Marshal(targetConfig)
	if err != nil {
		slog.Error("Failed to marshal target config to YAML", slog.Any("error", err))
		return err
	}

	configPath := filepath.Join(projDir, "proj.yml")
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		slog.Error("Failed to write target config file", slog.String("path", configPath), slog.Any("error", err))
		return err
	}
	slog.Debug("Wrote target config", slog.String("path", configPath))

	return nil
}

func (c *Creator) setupConfig() error {
	if err := config.InitTemplate(c.paths.TemplateConfigPath); err != nil {
		slog.Error("Template configuration load failure", slog.Any("error", err))
		return err
	}
	slog.Debug("template config loaded", "file", viper.ConfigFileUsed())

	if _, err := os.Stat(c.paths.TargetPath); err == nil {
		slog.Error("Target path exists", slog.String("path", c.paths.TargetPath))
		return fmt.Errorf("target path already exists: %s", c.paths.TargetPath)
	}

	defPath := strings.Join([]string{"definitions", c.cfg.DefinitionName}, ".")
	if !viper.IsSet(defPath) {
		slog.Error("Definition does not exist in template", slog.String("path", defPath), slog.String("definition-name", c.cfg.DefinitionName), slog.String("template name", c.cfg.TemplateName), slog.String("template config file", viper.GetString("template-config-file")))
		return fmt.Errorf("definition '%s' does not exist in template '%s'", c.cfg.DefinitionName, c.cfg.TemplateName)
	}

	reqs, err := loadRequirements()
	if err != nil {
		return err
	}
	c.reqs = reqs

	vars, scripts, files, luaenv, err := loadCreatorState(c.cfg, c.paths, c.reqs)
	if err != nil {
		return err
	}
	c.vars = vars
	c.scripts = scripts
	c.files = files
	c.luaenv = luaenv

	return nil
}

func (c *Creator) runBeforeScripts() error {
	vars, err := runBeforeScriptsAndValidateVars(c.luaenv, c.scripts.BeforeScripts(), c.reqs)
	if err != nil {
		return err
	}
	c.vars = vars

	return nil
}

func (c *Creator) runAfterScripts() error {
	return runAfterScripts(c.luaenv, c.scripts.AfterScripts())
}

func (c *Creator) processFiles() error {
	return processFiles(c.cfg, c.paths, c.files, c.vars)
}
