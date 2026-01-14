// This should be replaced by some proper YAML parsing library. We should split the regular proj
// config with the definition of templates. Maybe support splitting individual templates into their
// own config files.
package config

import (
	"errors"
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

	reqsMap, ok := reqs.(map[string]any)
	if !ok {
		slog.Error("Requirements parse failure", slog.Any("requirements", reqs))
		return result, errors.New("RequirementsSpec parse failed: entire requirements section is busted")
	}
	// pluck out Local
	local, ok := reqsMap["local"].(bool)
	if !ok {
		local = false
	}

	// pluck out variables
	var resultVars []VariableSpec
	varSlice := reqsMap["variables"].([]any)
	for i, varDef := range varSlice {
		vdef, ok := varDef.(map[string]any)
		if !ok {
			slog.Error("Invalid variable declaration", slog.Int("index", i), slog.Any("definition", varDef))
			return result, errors.New("RequirementSpec parse error: bad variable declaration")
		}
		name, _ := vdef["name"].(string)
		def, _ := vdef["default"]
		resultVars = append(resultVars, VariableSpec{Name: name, Default: def})
	}

	result.Local = local
	result.Variables = resultVars
	slog.Debug("Parsed requirements", slog.Bool("local", result.Local), slog.Any("variables", result.Variables))
	return result, nil
}

func BuildVariables(reqvars []VariableSpec) (map[string]any, error) {
	result := make(map[string]any)

	slog.Debug("Prepare Variables", slog.String("Placeholder", "resolving variables"), slog.String("Desc", "Generate all combinations of variables and values by combining global, template, and definition descriptions of variables"))
	setvars := buildMapFromSetVariables()
	globalvars := viper.GetStringMap("variables.global")
	slog.Debug("Building variables", slog.Any("globalvars", globalvars), slog.Any("from set-variable", setvars))

	for _, v := range reqvars {
		var finalval any
		if _, ok := setvars[v.Name]; ok {
			finalval = setvars[v.Name]
		} else if _, ok := globalvars[v.Name]; ok {
			finalval = globalvars[v.Name]
		} else {
			finalval = v.Default
		}
		result[v.Name] = finalval
	}

	return result, nil
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
