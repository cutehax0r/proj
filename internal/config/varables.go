package config

import (
	"fmt"
	"log/slog"
	"maps"
	"strings"

	"github.com/spf13/viper"
)

type VariableSpec struct {
	Name    string
	Default any
}

// normalizeVariableName converts a variable name to lowercase.
// All variable names are normalized to lowercase for consistency.
// For example: "UserName" becomes "username", "user_name" stays "user_name"
func normalizeVariableName(name string) string {
	return strings.ToLower(name)
}

// warnIfMixedCase logs a warning if the variable name contains mixed case.
// This helps users understand that their variable names will be normalized to lowercase.
func warnIfMixedCase(originalName, normalizedName string) {
	if originalName != normalizedName {
		slog.Warn("Variable name contains mixed case and has been normalized to lowercase",
			slog.String("original", originalName),
			slog.String("normalized", normalizedName))
	}
}

func BuildVariables(reqvars []VariableSpec) (map[string]any, error) {
	result := make(map[string]any)

	globalvars := viper.GetStringMap("variables.global")
	// Normalize global variable names
	normalizedGlobalVars := make(map[string]any)
	for k, v := range globalvars {
		normalized := normalizeVariableName(k)
		warnIfMixedCase(k, normalized)
		normalizedGlobalVars[normalized] = v
	}
	maps.Copy(result, normalizedGlobalVars)

	setvars := buildMapFromSetVariables()
	maps.Copy(result, setvars)

	reqd := make(map[string]any)
	for _, v := range reqvars {
		normalizedName := normalizeVariableName(v.Name)
		warnIfMixedCase(v.Name, normalizedName)

		var finalval any
		// Priority: CLI set-variable > viper variables key > global variables > default
		if _, ok := setvars[normalizedName]; ok {
			finalval = setvars[normalizedName]
		} else if viper.IsSet(fmt.Sprintf("variables.%s", normalizedName)) {
			finalval = viper.Get(fmt.Sprintf("variables.%s", normalizedName))
		} else if _, ok := normalizedGlobalVars[normalizedName]; ok {
			finalval = normalizedGlobalVars[normalizedName]
		} else {
			finalval = v.Default
		}
		reqd[normalizedName] = finalval
	}
	maps.Copy(result, reqd)

	// System variables are always lowercase
	result["targetname"] = viper.GetString("target-name")
	result["templatename"] = viper.GetString("template-name")
	result["definitionname"] = viper.GetString("definition-name")

	return result, nil
}

func buildMapFromSetVariables() map[string]any {
	vars := make(map[string]any)
	slog.Debug("vars", "len", len(vars))
	for i, rawset := range viper.GetStringSlice("set-variable") {
		key, value, found := strings.Cut(rawset, "=")
		if !found {
			slog.Warn("Invalid argument for set-variable. skipping", slog.Int("Index", i), slog.String("raw", rawset))
			continue
		}
		normalizedKey := normalizeVariableName(key)
		warnIfMixedCase(key, normalizedKey)
		slog.Debug("Parsing set-variable", slog.String("raw", rawset), slog.String("key", key), slog.String("normalizedKey", normalizedKey), slog.String("value", value))
		vars[normalizedKey] = value // maybe consider casting? json parsing?
	}
	return vars
}
