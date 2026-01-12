package config

import (
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type VariableSpec struct {
	Name string
	Default any
}

type RequirementSpec struct {
	Local bool `yaml:"local" default:"false"`
	Variables []VariableSpec `yaml:"variables"`
}

func BuildRequirements() (RequirementSpec, error) {
	var result RequirementSpec
	path := strings.Join([]string{"definitions", viper.GetString("definition"), "requirements"}, ".")
	reqs := viper.Get(path)
	slog.Debug("Loaded requirements", slog.String("path", path), slog.Any("Requirements", reqs))
	return result, nil

	
	// fileSlice, ok := files.([]any)
	// if !ok {
	// 	slog.Error("FileSpec parse failure", slog.Any("Files", files))
	// 	return nil, errors.New("FileSpec parse failure: entire files section is busted")
	// }
	// for i, fileDef := range fileSlice {
	// 	spec, ok := fileDef.(map[string]any)
	// 	if !ok {
	// 		slog.Error("Invalid file declaration", slog.Int("Index", i), slog.Any("definition", fileDef))
	// 		return nil, errors.New("FileSpec parse failure: bad file declaration")
	// 	}
	// 	source, _ := spec["source"].(string)
	// 	target, _ := spec["target"].(string)
	// 	parse, ok := spec["parse"].(bool)
	// 	if !ok {
	// 		parse = true
	// 	}
	// 	result = append(result, FileSpec{Source: source, Target: target, Parse: parse})
	// }
	//
	// return result, nil
}

func BuildVariables() (map[string]any, error) {
	slog.Debug("Prepare Variables", slog.String("Placeholder", "resolving variables"), slog.String("Desc", "Generate all combinations of variables and values by combining global, template, and definition descriptions of variables"))
	result := make(map[string]any)
	setvars := buildMapFromSetVariables()
	slog.Debug("Explicitly set varaibles", slog.Any("from set-variable", setvars))
	// now we walk through requirements.variables
	// ans for each do result[name] if setvars[name] is present
	return result, nil
	// step through requirements.variables
	// result[name] = setvars[name] (have to parse that into a map[string]any)
	// result[name] = global[name]
	// result[name] = default
}

func buildMapFromSetVariables() map[string]any {
	vars := make(map[string]any)
	slog.Debug("vars", "len", len(vars))
	for i, rawset := range viper.GetStringSlice("set-variables") {
		key, value, found := strings.Cut(rawset, "=")
		if !found {
			slog.Warn("Invalid argument for set-variable. skipping", slog.Int("Index", i), slog.String("raw", rawset))
			continue
		}
		slog.Debug("Parsing set-variable", slog.String("raw", rawset), slog.String("key", key), slog.String("value", value))
		vars[key] = value // maybe consider casting? json parsing?
	}
	return vars
}
