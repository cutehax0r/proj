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

func NewFileSpecs(paths *paths.Paths) ([]FileSpec, error) {
	result := make([]FileSpec, 0)

	datapath := strings.Join([]string{"definitions", viper.GetString("definition-name"), "files"}, ".")
	files := viper.Get(datapath)

	fileSlice, ok := files.([]any)
	if !ok {
		slog.Error("FileSpec parse failure", slog.Any("Files", files))
		return nil, errors.New("FileSpec parse failure: entire files section is busted")
	}
	for i, fileDef := range fileSlice {
		spec, ok := fileDef.(map[string]any)
		if !ok {
			slog.Error("Invalid file declaration", slog.Int("Index", i), slog.Any("definition", fileDef))
			return nil, errors.New("FileSpec parse failure: bad file declaration")
		}
		source, _ := paths.ResolveTemplate(spec["source"].(string))
		target, _ := paths.ResolveTarget(spec["target"].(string))

		info, err := os.Stat(source)
		if err != nil {
			slog.Error("Failed to stat source file", slog.String("source", source), slog.Any("spec", spec), slog.Any("error", err))
			return nil, err
		}
		sourceMode := info.Mode()

		parse, ok := spec["parse"].(bool)
		if !ok {
			parse = true
		}
		result = append(result, FileSpec{Source: source, Target: target, Parse: parse, SourceMode: sourceMode})
	}

	return result, nil
}
