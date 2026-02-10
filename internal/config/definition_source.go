package config

import (
	"log/slog"
	"path/filepath"

	"github.com/spf13/viper"
)

// DefinitionSource tracks where a definition came from
type DefinitionSource struct {
	Name string // Name of the definition (e.g., "page", "post")
	Path string // Directory where the definition sources are relative to
}

// SetDefinitionSource stores metadata about where a definition comes from
func SetDefinitionSource(defName, sourceDir string) {
	key := "definition-sources." + defName
	viper.Set(key, sourceDir)
	slog.Debug("Set definition source", slog.String("definition", defName), slog.String("path", sourceDir))
}

// GetDefinitionSource retrieves the source directory for a definition
// Priority: project .proj/proj.yml > template > global ~/.local/share/proj
func GetDefinitionSource(defName string) string {
	key := "definition-sources." + defName
	if viper.IsSet(key) {
		return viper.GetString(key)
	}
	// Default to template path if not explicitly set
	return viper.GetString("template-path")
}

// GetGlobalDefinitionsPath returns the path to global definitions
func GetGlobalDefinitionsPath() (string, error) {
	globalConfigRoot := viper.GetString("global-config-root")
	if globalConfigRoot == "" {
		return "", nil
	}
	return globalConfigRoot, nil
}

// SetProjectDefinitionSources marks all definitions in projectCfg as coming from the project root
func SetProjectDefinitionSources(projectRoot string, projectCfg *TargetConfig) {
	if projectCfg.Definitions != nil {
		projectDefDir := filepath.Join(projectRoot, ".proj")
		for defName := range projectCfg.Definitions {
			SetDefinitionSource(defName, projectDefDir)
		}
	}
}

// SetGlobalDefinitionSources marks definitions from global config as coming from global config root
func SetGlobalDefinitionSources(globalConfigRoot string) {
	// For global definitions, use the global config root as the base
	viper.Set("global-definitions-path", globalConfigRoot)
	slog.Debug("Set global definitions path", slog.String("path", globalConfigRoot))
}
