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

type Adder struct {
	cfg     *Config
	paths   *paths.Paths
	reqs    *config.RequirementSpec
	scripts config.ScriptSpec
	files   []config.FileSpec
	vars    map[string]any
	luaenv  *luabridge.Runtime
}

func NewAdder(cfg *Config) (*Adder, error) {
	return &Adder{
		cfg:   cfg,
		paths: cfg.Paths,
	}, nil
}

func (a *Adder) Add() error {
	if err := a.setupConfig(); err != nil {
		return err
	}

	if err := a.runBeforeScripts(); err != nil {
		return err
	}

	if err := a.validateFiles(); err != nil {
		return err
	}

	if err := a.processFiles(); err != nil {
		return err
	}

	if err := a.runAfterScripts(); err != nil {
		return err
	}

	a.luaenv.CloseState()
	return nil
}

func (a *Adder) setupConfig() error {
	if err := config.InitTemplate(a.paths.TemplateConfigPath); err != nil {
		slog.Error("Template configuration load failure", slog.Any("error", err))
		return err
	}
	slog.Debug("template config loaded", "file", viper.ConfigFileUsed())

	viper.Set("definition-name", a.cfg.DefinitionName)
	viper.Set("target-name", a.cfg.TargetName)
	viper.Set("template-name", a.cfg.TemplateName)

	projectRoot := viper.GetString("target-root")
	projYmlPath := filepath.Join(projectRoot, paths.TargetConfigFileDir, paths.TargetConfigFile)
	projData, err := os.ReadFile(projYmlPath)
	if err != nil {
		slog.Error("Failed to read project config", slog.String("path", projYmlPath), slog.Any("error", err))
		return err
	}

	projectCfg := &config.TargetConfig{}
	if err := yaml.Unmarshal(projData, projectCfg); err != nil {
		slog.Error("Failed to unmarshal project config", slog.String("path", projYmlPath), slog.Any("error", err))
		return err
	}

	if projectCfg.Variables != nil {
		for key, val := range projectCfg.Variables {
			viper.Set(fmt.Sprintf("variables.%s", key), val)
		}
	}

	if projectCfg.Scripts != nil {
		for key, val := range projectCfg.Scripts {
			viper.Set(fmt.Sprintf("scripts.%s", key), val)
		}
	}

	if projectCfg.Definitions != nil {
		for key, val := range projectCfg.Definitions {
			viper.Set(fmt.Sprintf("definitions.%s", key), val)
		}
		config.SetProjectDefinitionSources(projectRoot, projectCfg)
	}

	slog.Debug("Project config merged", slog.Any("variables", projectCfg.Variables), slog.Any("definitions", projectCfg.Definitions))

	if _, err := os.Stat(a.paths.TargetPath); err != nil {
		slog.Error("Target path does not exist", slog.String("path", a.paths.TargetPath))
		return err
	}
	slog.Debug("Target path exists", slog.String("path", a.paths.TargetPath))

	defPath := strings.Join([]string{"definitions", a.cfg.DefinitionName}, ".")
	if !viper.IsSet(defPath) {
		slog.Error("Definition does not exist", slog.String("path", defPath), slog.String("definition-name", a.cfg.DefinitionName), slog.String("template name", a.cfg.TemplateName))
		return fmt.Errorf("definition '%s' does not exist in template '%s'", a.cfg.DefinitionName, a.cfg.TemplateName)
	}

	reqs, err := loadRequirements()
	if err != nil {
		return err
	}
	a.reqs = reqs

	vars, scripts, files, luaenv, err := loadAdderState(a.cfg, a.paths, a.reqs)
	if err != nil {
		return err
	}
	a.vars = vars
	a.scripts = scripts
	a.files = files
	a.luaenv = luaenv

	return nil
}

func (a *Adder) runBeforeScripts() error {
	vars, err := runBeforeScriptsAndValidateVars(a.luaenv, a.scripts.BeforeScripts(), a.reqs)
	if err != nil {
		return err
	}
	a.vars = vars

	return nil
}

func (a *Adder) validateFiles() error {
	var conflicts []string

	for _, file := range a.files {
		// Render the destination filename
		destPath, err := renderTargetPath(file.Target, a.vars)
		if err != nil {
			slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", file.Target))
			return err
		}

		// Check if file already exists
		if _, err := os.Stat(destPath); err == nil {
			conflicts = append(conflicts, destPath)
			slog.Debug("Target file already exists", slog.String("target", destPath))
		}
	}

	if len(conflicts) > 0 {
		conflictMsg := strings.Join(conflicts, "\n  ")
		slog.Error("Target files already exist", slog.String("conflicts", conflictMsg))
		return fmt.Errorf("the following target files already exist:\n  %s", conflictMsg)
	}

	slog.Debug("File validation passed - no conflicts")
	return nil
}

func (a *Adder) processFiles() error {
	return processFiles(a.cfg, a.paths, a.files, a.vars)
}

func (a *Adder) runAfterScripts() error {
	return runAfterScripts(a.luaenv, a.scripts.AfterScripts())
}
