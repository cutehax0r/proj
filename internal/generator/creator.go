package generator

import (
	"bytes"
	"errors"
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

	if err := c.writeConfig(); err != nil {
		return err
	}
	c.luaenv.CloseState()

	return nil
}

func (c *Creator) writeConfig() error {
	projDir := filepath.Join(c.paths.TargetPath, ".proj")
	if err := c.cfg.Fs.MkdirAll(projDir, 0755); err != nil {
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
	if err := afero.WriteFile(c.cfg.Fs, configPath, yamlData, 0644); err != nil {
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

	if _, err := c.cfg.Fs.Stat(c.paths.TargetPath); err == nil {
		slog.Error("Target path exists", slog.String("path", c.paths.TargetPath))
		return afero.ErrFileExists
	}

	defPath := strings.Join([]string{"definitions", c.cfg.DefinitionName}, ".")
	if !viper.IsSet(defPath) {
		slog.Error("Definition does not exist in template", slog.String("path", defPath), slog.String("definition-name", c.cfg.DefinitionName), slog.String("template name", c.cfg.TemplateName), slog.String("template config file", viper.GetString("template-config-file")))
		return afero.ErrFileNotFound
	}

	reqs, err := config.NewRequirements()
	if err != nil {
		slog.Error("Failed to load requirements", slog.Any("error", err))
		return err
	}
	c.reqs = reqs
	slog.Debug("Final Requirements", slog.Any("reqs", reqs))

	vars, err := config.BuildVariables(reqs.Variables)
	if err != nil {
		slog.Error("Failed to build variables", slog.Any("error", err))
		return err
	}
	c.vars = vars
	slog.Debug("Final Variables", slog.Any("vars", vars))

	scripts, err := config.NewScriptSpec(c.paths)
	if err != nil {
		slog.Error("Couldn't build scripts", slog.Any("error", err))
		return err
	}
	c.scripts = scripts
	slog.Debug("Final scripts", slog.Any("scripts", scripts))

	files, err := config.NewFileSpecs(c.paths)
	if err != nil {
		slog.Error("Failed to load files from template definition", slog.Any("error", err))
		return err
	}
	c.files = files
	slog.Debug("Final files", slog.Any("files", files))

	c.luaenv = luabridge.NewRuntime(c.vars, c.paths, c.reqs, &c.files, c.cfg.NoWrite)

	return nil
}

func (c *Creator) runBeforeScripts() error {
	for _, script := range c.scripts.BeforeScripts() {
		if err := c.luaenv.Run(script); err != nil {
			slog.Error("Error in lua script. Aborting", slog.Any("error", err), slog.String("script", script))
			return err
		}
	}

	c.vars = c.luaenv.GetVariables()
	slog.Debug("Variables after before-scripts", slog.Any("vars", c.vars))

	for _, varspec := range c.reqs.Variables {
		if c.vars[varspec.Name] == nil {
			slog.Error("Required variable is not set. Use --set-variable. Aborting.", slog.Any("Name", varspec.Name))
			slog.Info("All variables", slog.Any("vars", c.vars))
			return errors.New("required variable not set")
		}
	}
	slog.Debug("All the variables are ready so we can do the work")

	return nil
}

func (c *Creator) runAfterScripts() error {
	for _, script := range c.scripts.AfterScripts() {
		if err := c.luaenv.Run(script); err != nil {
			slog.Error("Error in lua script. Aborting", slog.Any("error", err), slog.String("script", script))
			return err
		}
	}
	return nil
}

func (c *Creator) processFiles() error {
	// We're going to read all of the files that we're going to parse into memory. Two copies
	// (before and after parsing). That's not very ram efficient but it lets us do some nice
	// debugging with --no-write by forcing parsing to happen without writing. Computers have
	// lots of ram and text files are small so we're not sweating it.
	for i, file := range c.files {
		// Render the destination filename
		destPath, err := c.renderTargetPath(file.Target)
		if err != nil {
			slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", file.Target))
			return err
		}

		// Render file content if needed
		if file.Parse {
			if err := c.renderAndWriteFile(i, destPath); err != nil {
				return err
			}
		} else {
			if err := c.copyFile(i, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Creator) renderTargetPath(targetTemplate string) (string, error) {
	desttemp, err := template.New("filename").Parse(targetTemplate)
	if err != nil {
		slog.Error("Couldn't template the target filename", slog.String("target", targetTemplate), slog.Any("err", err))
		return "", err
	}

	var deststr bytes.Buffer
	err = desttemp.Execute(&deststr, c.vars)
	if err != nil {
		slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", targetTemplate))
		return "", err
	}

	return deststr.String(), nil
}

func (c *Creator) renderAndWriteFile(fileIdx int, destPath string) error {
	file := c.files[fileIdx]

	slog.Info("parsing content of file", slog.String("source", file.Source), slog.String("target", destPath))

	rawtemp, err := afero.ReadFile(c.cfg.Fs, file.Source)
	if err != nil {
		slog.Error("Couldn't read the raw template data", slog.String("Source", file.Source), slog.Any("err", err))
		return err
	}
	file.Raw = string(rawtemp)
	slog.Debug("rendering template", slog.String("raw data", file.Raw))

	conttemp, err := template.New("template").Parse(file.Raw)
	if err != nil {
		slog.Error("Error parsing template", slog.Any("err", err), slog.Any("file", file.Source), slog.Any("paths", c.paths))
		return err
	}

	var contbuff bytes.Buffer
	if err := conttemp.Execute(&contbuff, c.vars); err != nil {
		slog.Error("Error executing template", slog.Any("error", err))
		return err
	}
	file.Rendered = contbuff.String()
	slog.Info("result", slog.String("rendered", file.Rendered))

	if c.cfg.NoWrite {
		slog.Debug("No-write set: skipping write", slog.String("source", file.Source), slog.String("target", destPath))
		return nil
	}

	targetDir := filepath.Dir(destPath)
	if err := c.cfg.Fs.MkdirAll(targetDir, 0755); err != nil {
		slog.Error("Failed to create target directory", slog.String("directory", targetDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created target directory", slog.String("directory", targetDir))

	if err := afero.WriteFile(c.cfg.Fs, destPath, []byte(file.Rendered), file.SourceMode); err != nil {
		slog.Error("Failed to write rendered file", slog.String("target", destPath), slog.Any("error", err))
		return err
	}
	slog.Debug("Wrote rendered file", slog.String("target", destPath))

	return nil
}

func (c *Creator) copyFile(fileIdx int, destPath string) error {
	file := c.files[fileIdx]

	slog.Info("Nothing to render, skipping parse.", slog.String("source", file.Source), slog.String("target", destPath))

	if c.cfg.NoWrite {
		slog.Debug("No-write set: skipping copy", slog.String("source", file.Source), slog.String("target", destPath))
		return nil
	}

	targetDir := filepath.Dir(destPath)
	if err := c.cfg.Fs.MkdirAll(targetDir, 0755); err != nil {
		slog.Error("Failed to create target directory", slog.String("directory", targetDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created target directory", slog.String("directory", targetDir))

	sourceFile, err := c.cfg.Fs.Open(file.Source)
	if err != nil {
		slog.Error("Failed to open source file", slog.String("source", file.Source), slog.Any("error", err))
		return err
	}
	defer sourceFile.Close()

	targetFile, err := c.cfg.Fs.Create(destPath)
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

	if err := c.cfg.Fs.Chmod(destPath, file.SourceMode); err != nil {
		slog.Warn("Failed to set permissions on copied file", slog.String("target", destPath), slog.Any("error", err))
	}
	slog.Debug("Set permissions on copied file", slog.String("target", destPath), slog.Any("mode", file.SourceMode))

	return nil
}
