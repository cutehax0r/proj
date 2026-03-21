package config

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"proj/internal/paths"
	"strings"

	"github.com/spf13/viper"
)

type FileSpec struct {
	Source     string      `yaml:"source"`
	Target     string      `yaml:"target"`
	Parse      bool        `yaml:"parse" default:"true"`
	Raw        string      `yaml:"raw"`
	Rendered   string      `yaml:"rendered"`
	SourceMode fs.FileMode `yaml:"source_mode"`
}

func (f *FileSpec) ToMap() map[string]any {
	return map[string]any{
		"source":      f.Source,
		"target":      f.Target,
		"parse":       f.Parse,
		"source_mode": f.SourceMode,
	}
}

func NewFileSpecs(p *paths.Paths) ([]FileSpec, error) {
	result := make([]FileSpec, 0)

	datapath := strings.Join([]string{"definitions", viper.GetString("definition-name"), "files"}, ".")
	files := viper.Get(datapath)

	fileSlice, ok := files.([]any)
	if !ok {
		slog.Error("FileSpec parse failure", slog.Any("Files", files))
		return nil, errors.New("FileSpec parse failure: entire files section is busted")
	}

	// Get the source directory for this definition
	defName := viper.GetString("definition-name")
	sourceDir := GetDefinitionSource(defName)
	templateSourceDir := viper.GetString("template-path")
	slog.Debug("Resolving file sources relative to", slog.String("definition", defName), slog.String("source-dir", sourceDir))

	for i, fileDef := range fileSlice {
		spec, ok := fileDef.(map[string]any)
		if !ok {
			slog.Error("Invalid file declaration", slog.Int("Index", i), slog.Any("definition", fileDef))
			return nil, errors.New("FileSpec parse failure: bad file declaration")
		}

		// Resolve source relative to the definition's source directory
		source, err := p.ResolveFrom(sourceDir, spec["source"].(string))
		if err != nil {
			slog.Error("Failed to resolve source path", slog.String("source", spec["source"].(string)), slog.String("sourceDir", sourceDir), slog.Any("error", err))
			return nil, err
		}

		// If source file doesn't exist and we have a different template source dir, try fallback
		info, err := os.Stat(source)
		if err != nil && sourceDir != templateSourceDir {
			slog.Debug("Source file not found in primary location, checking fallback", slog.String("source", source), slog.String("fallback-dir", templateSourceDir))
			fallbackSource, fallbackErr := p.ResolveFrom(templateSourceDir, spec["source"].(string))
			if fallbackErr == nil {
				if fallbackInfo, fallbackStatErr := os.Stat(fallbackSource); fallbackStatErr == nil {
					slog.Debug("Found source file in fallback location", slog.String("fallback-source", fallbackSource))
					source = fallbackSource
					info = fallbackInfo
					err = nil
				}
			}
		}

		if err != nil {
			slog.Error("Failed to stat source file", slog.String("source", source), slog.Any("spec", spec), slog.Any("error", err))
			return nil, err
		}
		sourceMode := info.Mode()

		target, _ := p.ResolveTarget(spec["target"].(string))

		parse, ok := spec["parse"].(bool)
		if !ok {
			parse = true
		}
		result = append(result, FileSpec{Source: source, Target: target, Parse: parse, SourceMode: sourceMode})
	}

	return result, nil
}
