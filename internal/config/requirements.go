// TODO: This should be replaced by some proper YAML parsing library.
package config

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/spf13/viper"
)

type RequirementSpec struct {
	Local     bool           `yaml:"local" default:"false"`
	Variables []VariableSpec `yaml:"variables"`
}

func NewRequirements() (*RequirementSpec, error) {
	var result RequirementSpec
	path := strings.Join([]string{"definitions", viper.GetString("definition-name"), "requirements"}, ".")
	reqs := viper.Get(path)
	slog.Debug("Loaded requirements", slog.String("path", path), slog.Any("Requirements", reqs))

	reqsMap, ok := reqs.(map[string]any)
	if !ok {
		slog.Error("Requirements parse failure", slog.Any("requirements", reqs))
		return &result, errors.New("RequirementsSpec parse failed: entire requirements section is busted")
	}

	local, ok := reqsMap["local"].(bool)
	if !ok {
		local = false
	}

	var resultVars []VariableSpec
	if variables, ok := reqsMap["variables"]; ok && variables != nil {
		varSlice := variables.([]any)
		for i, varDef := range varSlice {
			vdef, ok := varDef.(map[string]any)
			if !ok {
				slog.Error("Invalid variable declaration", slog.Int("index", i), slog.Any("definition-name", varDef))
				return &result, errors.New("RequirementSpec parse error: bad variable declaration")
			}
			name, _ := vdef["name"].(string)
			def, _ := vdef["default"]
			resultVars = append(resultVars, VariableSpec{Name: name, Default: def})
		}
	}

	result.Local = local
	result.Variables = resultVars
	slog.Debug("Parsed requirements", slog.Bool("local", result.Local), slog.Any("variables", result.Variables))
	return &result, nil
}
