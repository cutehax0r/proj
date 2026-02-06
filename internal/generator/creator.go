package generator

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"proj/internal/config"
	"proj/internal/luabridge"
	"proj/internal/paths"
	"strings"
	"text/template"

	"github.com/spf13/viper"
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

func (g *Creator) Create() error {
	if err := g.setupConfig(); err != nil {
		return err
	}

	if err := g.runBeforeScripts(); err != nil {
		return err
	}

	if err := g.processFiles(); err != nil {
		return err
	}

	if err := g.runAfterScripts(); err != nil {
		return err
	}

	return nil
}

func (g *Creator) setupConfig() error {
	if err := config.InitTemplate(g.paths.TemplateConfigPath); err != nil {
		slog.Error("Template configuration load failure", slog.Any("error", err))
		return err
	}
	slog.Debug("template config loaded", "file", viper.ConfigFileUsed())

	if _, err := os.Stat(g.paths.TargetPath); err == nil {
		slog.Error("Target path exists", slog.String("path", g.paths.TargetPath))
		return err
	}

	defPath := strings.Join([]string{"definitions", g.cfg.DefinitionName}, ".")
	if !viper.IsSet(defPath) {
		slog.Error("Definition does not exist in template", slog.String("path", defPath), slog.String("definition-name", g.cfg.DefinitionName), slog.String("template name", g.cfg.TemplateName), slog.String("template config file", viper.GetString("template-config-file")))
		return os.ErrNotExist
	}

	reqs, err := config.NewRequirements()
	if err != nil {
		slog.Error("Failed to load requirements", slog.Any("error", err))
		return err
	}
	g.reqs = reqs
	slog.Debug("Final Requirements", slog.Any("reqs", reqs))

	vars, err := config.BuildVariables(reqs.Variables)
	if err != nil {
		slog.Error("Failed to build variables", slog.Any("error", err))
		return err
	}
	g.vars = vars
	slog.Debug("Final Variables", slog.Any("vars", vars))

	scripts, err := config.NewScriptSpec(g.paths)
	if err != nil {
		slog.Error("Couldn't build scripts", slog.Any("error", err))
		return err
	}
	g.scripts = scripts
	slog.Debug("Final scripts", slog.Any("scripts", scripts))

	files, err := config.NewFileSpecs(g.paths)
	if err != nil {
		slog.Error("Failed to load files from template definition", slog.Any("error", err))
		return err
	}
	g.files = files
	slog.Debug("Final files", slog.Any("files", files))

	g.luaenv = luabridge.NewRuntime(g.vars, g.paths, g.reqs, &g.files, g.cfg.NoWrite)

	return nil
}

func (g *Creator) runBeforeScripts() error {
	for _, script := range g.scripts.BeforeScripts() {
		if err := g.luaenv.Run(script); err != nil {
			slog.Error("Error in lua script. Aborting", slog.Any("error", err), slog.String("script", script))
			return err
		}
	}

	g.vars = g.luaenv.GetVariables()
	slog.Debug("Variables after before-scripts", slog.Any("vars", g.vars))

	for _, varspec := range g.reqs.Variables {
		if g.vars[varspec.Name] == nil {
			slog.Error("Required variable is not set. Use --set-variable. Aborting.", slog.Any("Name", varspec.Name))
			slog.Info("All variables", slog.Any("vars", g.vars))
			return os.ErrInvalid
		}
	}
	slog.Debug("All the variables are ready so we can do the work")

	return nil
}

func (g *Creator) runAfterScripts() error {
	for _, script := range g.scripts.AfterScripts() {
		if err := g.luaenv.Run(script); err != nil {
			slog.Error("Error in lua script. Aborting", slog.Any("error", err), slog.String("script", script))
			return err
		}
	}
	return nil
}

func (g *Creator) processFiles() error {
	// We're going to read all of the files that we're going to parse into memory. Two copies
	// (before and after parsing). That's not very ram efficient but it lets us do some nice
	// debugging with --no-write by forcing parsing to happen without writing. Computers have
	// lots of ram and text files are small so we're not sweating it.
	for i, file := range g.files {
		// Render the destination filename
		destPath, err := g.renderTargetPath(file.Target)
		if err != nil {
			slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", file.Target))
			return err
		}

		// Render file content if needed
		if file.Parse {
			if err := g.renderAndWriteFile(i, destPath); err != nil {
				return err
			}
		} else {
			if err := g.copyFile(i, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Creator) renderTargetPath(targetTemplate string) (string, error) {
	desttemp, err := template.New("filename").Parse(targetTemplate)
	if err != nil {
		slog.Error("Couldn't template the target filename", slog.String("target", targetTemplate), slog.Any("err", err))
		return "", err
	}

	var deststr bytes.Buffer
	err = desttemp.Execute(&deststr, g.vars)
	if err != nil {
		slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", targetTemplate))
		return "", err
	}

	return deststr.String(), nil
}

func (g *Creator) renderAndWriteFile(fileIdx int, destPath string) error {
	file := g.files[fileIdx]

	slog.Info("parsing content of file", slog.String("source", file.Source), slog.String("target", destPath))

	rawtemp, err := os.ReadFile(file.Source)
	if err != nil {
		slog.Error("Couldn't read the raw template data", slog.String("Source", file.Source), slog.Any("err", err))
		return err
	}
	file.Raw = string(rawtemp)
	slog.Debug("rendering template", slog.String("raw data", file.Raw))

	conttemp, err := template.ParseFiles(file.Source)
	if err != nil {
		slog.Error("Error parsing template", slog.Any("err", err), slog.Any("file", file.Source), slog.Any("paths", g.paths))
		return err
	}

	var contbuff bytes.Buffer
	if err := conttemp.Execute(&contbuff, g.vars); err != nil {
		slog.Error("Error executing template", slog.Any("error", err))
		return err
	}
	file.Rendered = contbuff.String()
	slog.Info("result", slog.String("rendered", file.Rendered))

	if g.cfg.NoWrite {
		slog.Debug("No-write set: skipping write", slog.String("source", file.Source), slog.String("target", destPath))
		return nil
	}

	targetDir := filepath.Dir(destPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		slog.Error("Failed to create target directory", slog.String("directory", targetDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created target directory", slog.String("directory", targetDir))

	if err := os.WriteFile(destPath, []byte(file.Rendered), file.SourceMode); err != nil {
		slog.Error("Failed to write rendered file", slog.String("target", destPath), slog.Any("error", err))
		return err
	}
	slog.Debug("Wrote rendered file", slog.String("target", destPath))

	return nil
}

func (g *Creator) copyFile(fileIdx int, destPath string) error {
	file := g.files[fileIdx]

	slog.Info("Nothing to render, skipping parse.", slog.String("source", file.Source), slog.String("target", destPath))

	if g.cfg.NoWrite {
		slog.Debug("No-write set: skipping copy", slog.String("source", file.Source), slog.String("target", destPath))
		return nil
	}

	targetDir := filepath.Dir(destPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		slog.Error("Failed to create target directory", slog.String("directory", targetDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created target directory", slog.String("directory", targetDir))

	sourceFile, err := os.Open(file.Source)
	if err != nil {
		slog.Error("Failed to open source file", slog.String("source", file.Source), slog.Any("error", err))
		return err
	}
	defer sourceFile.Close()

	targetFile, err := os.Create(destPath)
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

	if err := os.Chmod(destPath, file.SourceMode); err != nil {
		slog.Warn("Failed to set permissions on copied file", slog.String("target", destPath), slog.Any("error", err))
	}
	slog.Debug("Set permissions on copied file", slog.String("target", destPath), slog.Any("mode", file.SourceMode))

	return nil
}
