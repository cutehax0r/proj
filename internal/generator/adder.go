package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"proj/internal/config"
	"proj/internal/luabridge"
	"proj/internal/paths"
	"strings"
	"text/template"

	"github.com/spf13/afero"
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
	fs := afero.NewOsFs()
	projData, err := afero.ReadFile(fs, projYmlPath)
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

	if _, err := fs.Stat(a.paths.TargetPath); err != nil {
		slog.Error("Target path does not exist", slog.String("path", a.paths.TargetPath))
		return err
	}
	slog.Debug("Target path exists", slog.String("path", a.paths.TargetPath))

	defPath := strings.Join([]string{"definitions", a.cfg.DefinitionName}, ".")
	if !viper.IsSet(defPath) {
		slog.Error("Definition does not exist", slog.String("path", defPath), slog.String("definition-name", a.cfg.DefinitionName), slog.String("template name", a.cfg.TemplateName))
		return afero.ErrFileNotFound
	}

	reqs, err := config.NewRequirements()
	if err != nil {
		slog.Error("Failed to load requirements", slog.Any("error", err))
		return err
	}
	a.reqs = reqs
	slog.Debug("Final Requirements", slog.Any("reqs", reqs))

	vars, err := config.BuildVariables(reqs.Variables)
	if err != nil {
		slog.Error("Failed to build variables", slog.Any("error", err))
		return err
	}
	a.vars = vars
	slog.Debug("Final Variables", slog.Any("vars", vars))

	scripts, err := config.NewScriptSpecWithFS(a.cfg.Fs, a.paths)
	if err != nil {
		slog.Error("Couldn't build scripts", slog.Any("error", err))
		return err
	}
	a.scripts = scripts
	slog.Debug("Final scripts", slog.Any("scripts", scripts))

	files, err := config.NewFileSpecsWithFS(a.cfg.Fs, a.paths)
	if err != nil {
		slog.Error("Failed to load files from template definition", slog.Any("error", err))
		return err
	}
	a.files = files
	slog.Debug("Final files", slog.Any("files", files))

	a.luaenv = luabridge.NewRuntime(a.vars, a.paths, a.reqs, &a.files, a.cfg.NoWrite)

	return nil
}

func (a *Adder) runBeforeScripts() error {
	for _, script := range a.scripts.BeforeScriptsWithFS(a.cfg.Fs) {
		if err := a.luaenv.Run(script); err != nil {
			slog.Error("Error in lua script. Aborting", slog.Any("error", err), slog.String("script", script))
			return err
		}
	}

	a.vars = a.luaenv.GetVariables()
	slog.Debug("Variables after before-scripts", slog.Any("vars", a.vars))

	for _, varspec := range a.reqs.Variables {
		if a.vars[varspec.Name] == nil {
			slog.Error("Required variable is not set. Use --set-variable. Aborting.", slog.Any("Name", varspec.Name))
			slog.Info("All variables", slog.Any("vars", a.vars))
			return errors.New("required variable not set")
		}
	}
	slog.Debug("All the variables are ready so we can do the work")

	return nil
}

func (a *Adder) validateFiles() error {
	var conflicts []string

	for _, file := range a.files {
		// Render the destination filename
		destPath, err := a.renderTargetPath(file.Target)
		if err != nil {
			slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", file.Target))
			return err
		}

		// Check if file already exists
		if _, err := a.cfg.Fs.Stat(destPath); err == nil {
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
	for i, file := range a.files {
		destPath, err := a.renderTargetPath(file.Target)
		if err != nil {
			slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", file.Target))
			return err
		}

		if file.Parse {
			if err := a.renderAndWriteFile(i, destPath); err != nil {
				return err
			}
		} else {
			if err := a.copyFile(i, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Adder) renderTargetPath(targetTemplate string) (string, error) {
	desttemp, err := template.New("filename").Parse(targetTemplate)
	if err != nil {
		slog.Error("Couldn't template the target filename", slog.String("target", targetTemplate), slog.Any("err", err))
		return "", err
	}

	var deststr bytes.Buffer
	err = desttemp.Execute(&deststr, a.vars)
	if err != nil {
		slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", targetTemplate))
		return "", err
	}

	return deststr.String(), nil
}

func (a *Adder) renderAndWriteFile(fileIdx int, destPath string) error {
	file := a.files[fileIdx]

	slog.Info("parsing content of file", slog.String("source", file.Source), slog.String("target", destPath))

	rawtemp, err := afero.ReadFile(a.cfg.Fs, file.Source)
	if err != nil {
		slog.Error("Couldn't read the raw template data", slog.String("Source", file.Source), slog.Any("err", err))
		return err
	}
	file.Raw = string(rawtemp)
	slog.Debug("rendering template", slog.String("raw data", file.Raw))

	conttemp, err := template.New("template").Parse(file.Raw)
	if err != nil {
		slog.Error("Error parsing template", slog.Any("err", err), slog.Any("file", file.Source), slog.Any("paths", a.paths))
		return err
	}

	var contbuff bytes.Buffer
	if err := conttemp.Execute(&contbuff, a.vars); err != nil {
		slog.Error("Error executing template", slog.Any("error", err))
		return err
	}
	file.Rendered = contbuff.String()
	slog.Info("result", slog.String("rendered", file.Rendered))

	if a.cfg.NoWrite {
		slog.Debug("No-write set: skipping write", slog.String("source", file.Source), slog.String("target", destPath))
		return nil
	}

	targetDir := filepath.Dir(destPath)
	if err := a.cfg.Fs.MkdirAll(targetDir, 0755); err != nil {
		slog.Error("Failed to create target directory", slog.String("directory", targetDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created target directory", slog.String("directory", targetDir))

	if err := afero.WriteFile(a.cfg.Fs, destPath, []byte(file.Rendered), file.SourceMode); err != nil {
		slog.Error("Failed to write rendered file", slog.String("target", destPath), slog.Any("error", err))
		return err
	}
	slog.Debug("Wrote rendered file", slog.String("target", destPath))

	return nil
}

func (a *Adder) copyFile(fileIdx int, destPath string) error {
	file := a.files[fileIdx]

	slog.Info("Nothing to render, skipping parse.", slog.String("source", file.Source), slog.String("target", destPath))

	if a.cfg.NoWrite {
		slog.Debug("No-write set: skipping copy", slog.String("source", file.Source), slog.String("target", destPath))
		return nil
	}

	targetDir := filepath.Dir(destPath)
	if err := a.cfg.Fs.MkdirAll(targetDir, 0755); err != nil {
		slog.Error("Failed to create target directory", slog.String("directory", targetDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created target directory", slog.String("directory", targetDir))

	sourceFile, err := a.cfg.Fs.Open(file.Source)
	if err != nil {
		slog.Error("Failed to open source file", slog.String("source", file.Source), slog.Any("error", err))
		return err
	}
	defer sourceFile.Close()

	targetFile, err := a.cfg.Fs.Create(destPath)
	if err != nil {
		slog.Error("Failed to create target file", slog.String("target", destPath), slog.Any("error", err))
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, sourceFile); err != nil {
		slog.Error("Failed to copy file", slog.String("source", file.Source), slog.String("target", destPath), slog.Any("error", err))
		return err
	}
	slog.Debug("Copied file", slog.String("source", file.Source), slog.String("target", destPath))

	if err := a.cfg.Fs.Chmod(destPath, file.SourceMode); err != nil {
		slog.Warn("Failed to set permissions on copied file", slog.String("target", destPath), slog.Any("error", err))
	}
	slog.Debug("Set permissions on copied file", slog.String("target", destPath), slog.Any("mode", file.SourceMode))

	return nil
}

func (a *Adder) runAfterScripts() error {
	for _, script := range a.scripts.AfterScriptsWithFS(a.cfg.Fs) {
		if err := a.luaenv.Run(script); err != nil {
			slog.Error("Error in lua script. Aborting", slog.Any("error", err), slog.String("script", script))
			return err
		}
	}
	return nil
}
