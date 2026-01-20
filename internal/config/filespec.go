package config

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type FileSpec struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Parse  bool   `yaml:"parse" default:"true"`
}

func (f *FileSpec) ToMap() map[string]any {
	return map[string]any{
		"source": f.Source,
		"target": f.Target,
		"parse":  f.Parse,
	}
}

func ParseFileSpecs() (*[]FileSpec, error) {
	var result []FileSpec

	path := strings.Join([]string{"definitions", viper.GetString("definition-name"), "files"}, ".")
	files := viper.Get(path)

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
		source, _ := spec["source"].(string)
		target, _ := spec["target"].(string)
		parse, ok := spec["parse"].(bool)
		if !ok {
			parse = true
		}
		result = append(result, FileSpec{Source: source, Target: target, Parse: parse})
	}

	return &result, nil
}
