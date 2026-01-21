package config

import (
	"errors"
	"log/slog"
	"proj/internal/paths"
	"strings"

	"github.com/spf13/viper"
)

type FileSpec struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
	Parse  bool   `yaml:"parse" default:"true"`
	Raw  string `yaml:"raw"`
	Rendered string `yaml:"rendered"`

}

func (f *FileSpec) ToMap() map[string]any {
	return map[string]any{
		"source": f.Source,
		"target": f.Target,
		"parse":  f.Parse,
	}
}

func ParseFileSpecs(paths *paths.Paths) (*[]FileSpec, error) {
	var result []FileSpec

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
		source, _ := resolve(paths.TemplatePath, spec["source"].(string))
		target, _ := resolve(paths.TargetPath, spec["target"].(string))
		parse, ok := spec["parse"].(bool)
		if !ok {
			parse = true
		}
		result = append(result, FileSpec{Source: source, Target: target, Parse: parse})
	}

	return &result, nil
}
