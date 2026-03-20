package generator

import (
	"bytes"
	"io"
	"log/slog"
	"path/filepath"
	"proj/internal/config"
	"proj/internal/paths"
	"text/template"

	"github.com/spf13/afero"
)

func processFiles(cfg *Config, p *paths.Paths, files []config.FileSpec, vars map[string]any) error {
	for i, file := range files {
		destPath, err := renderTargetPath(file.Target, vars)
		if err != nil {
			slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", file.Target))
			return err
		}

		if file.Parse {
			if err := renderAndWriteFile(cfg, p, files, i, destPath, vars); err != nil {
				return err
			}
		} else {
			if err := copyFile(cfg, files, i, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func renderTargetPath(targetTemplate string, vars map[string]any) (string, error) {
	desttemp, err := template.New("filename").Parse(targetTemplate)
	if err != nil {
		slog.Error("Couldn't template the target filename", slog.String("target", targetTemplate), slog.Any("err", err))
		return "", err
	}

	var deststr bytes.Buffer
	err = desttemp.Execute(&deststr, vars)
	if err != nil {
		slog.Error("Error templating target filename", slog.Any("error", err), slog.String("target", targetTemplate))
		return "", err
	}

	return deststr.String(), nil
}

func renderAndWriteFile(cfg *Config, p *paths.Paths, files []config.FileSpec, fileIdx int, destPath string, vars map[string]any) error {
	file := files[fileIdx]

	slog.Info("parsing content of file", slog.String("source", file.Source), slog.String("target", destPath))

	rawtemp, err := afero.ReadFile(cfg.Fs, file.Source)
	if err != nil {
		slog.Error("Couldn't read the raw template data", slog.String("Source", file.Source), slog.Any("err", err))
		return err
	}
	file.Raw = string(rawtemp)
	slog.Debug("rendering template", slog.String("raw data", file.Raw))

	conttemp, err := template.New("template").Parse(file.Raw)
	if err != nil {
		slog.Error("Error parsing template", slog.Any("err", err), slog.Any("file", file.Source), slog.Any("paths", p))
		return err
	}

	var contbuff bytes.Buffer
	if err := conttemp.Execute(&contbuff, vars); err != nil {
		slog.Error("Error executing template", slog.Any("error", err))
		return err
	}
	file.Rendered = contbuff.String()
	slog.Info("result", slog.String("rendered", file.Rendered))

	if cfg.NoWrite {
		slog.Debug("No-write set: skipping write", slog.String("source", file.Source), slog.String("target", destPath))
		return nil
	}

	targetDir := filepath.Dir(destPath)
	if err := cfg.Fs.MkdirAll(targetDir, 0o755); err != nil {
		slog.Error("Failed to create target directory", slog.String("directory", targetDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created target directory", slog.String("directory", targetDir))

	if err := afero.WriteFile(cfg.Fs, destPath, []byte(file.Rendered), file.SourceMode); err != nil {
		slog.Error("Failed to write rendered file", slog.String("target", destPath), slog.Any("error", err))
		return err
	}
	slog.Debug("Wrote rendered file", slog.String("target", destPath))

	return nil
}

func copyFile(cfg *Config, files []config.FileSpec, fileIdx int, destPath string) error {
	file := files[fileIdx]

	slog.Info("Nothing to render, skipping parse.", slog.String("source", file.Source), slog.String("target", destPath))

	if cfg.NoWrite {
		slog.Debug("No-write set: skipping copy", slog.String("source", file.Source), slog.String("target", destPath))
		return nil
	}

	targetDir := filepath.Dir(destPath)
	if err := cfg.Fs.MkdirAll(targetDir, 0o755); err != nil {
		slog.Error("Failed to create target directory", slog.String("directory", targetDir), slog.Any("error", err))
		return err
	}
	slog.Debug("Created target directory", slog.String("directory", targetDir))

	sourceFile, err := cfg.Fs.Open(file.Source)
	if err != nil {
		slog.Error("Failed to open source file", slog.String("source", file.Source), slog.Any("error", err))
		return err
	}
	defer sourceFile.Close()

	targetFile, err := cfg.Fs.Create(destPath)
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

	if err := cfg.Fs.Chmod(destPath, file.SourceMode); err != nil {
		slog.Warn("Failed to set permissions on copied file", slog.String("target", destPath), slog.Any("error", err))
	}
	slog.Debug("Set permissions on copied file", slog.String("target", destPath), slog.Any("mode", file.SourceMode))

	return nil
}
