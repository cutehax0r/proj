package config

import (
	"log/slog"
	"maps"
	"strings"

	"github.com/spf13/viper"
)

type VariableSpec struct {
	Name    string
	Default any
}

func BuildVariables(reqvars []VariableSpec) (map[string]any, error) {
	result := make(map[string]any)

	globalvars := viper.GetStringMap("variables.global")
	maps.Copy(result, globalvars)

	setvars := buildMapFromSetVariables()
	maps.Copy(result, setvars)

	reqd := make(map[string]any)
	for _, v := range reqvars {
		var finalval any
		if _, ok := setvars[v.Name]; ok {
			finalval = setvars[v.Name]
		} else if _, ok := globalvars[v.Name]; ok {
			finalval = globalvars[v.Name]
		} else {
			finalval = v.Default
		}
		reqd[v.Name] = finalval
	}
	maps.Copy(result, reqd)

	result["targetName"] = viper.GetString("target-name")
	result["templateName"] = viper.GetString("template-name")
	result["definitionName"] = viper.GetString("definition-name")

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
