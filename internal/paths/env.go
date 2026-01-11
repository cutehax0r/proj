package paths

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const GlobalConfigDir = "proj"
const GlobalConfigFile = "proj"
const TargetConfigFile = "proj.yml"
const TargetConfigFileDir = ".proj"
const TemplateConfigFile = "proj.yml"
// templaterootdefault
// configfilepase

// we explicitly construct XDG_CONFIG_HOME ones here because OSX users probably expect that from a
// command line util even though go's userconfdir doesn't do that.
func GlobalConfigPaths() []string {
	// $XDG_CONFIG_HOME/proj/
	// os.UserHomeDir()/.config/proj
	// os.UserHomeDir()/proj
	// os.UserConfDir()

	// .
	return []string{ "." }
}

// go's userdatadir returns "weird" paths on OSX so we explicitly construct the xdgdatahome and
// 'normal nix path' versions to check first. that's what people will expect.
func TemplateRootPaths() []string {
	// This builds the full set
	// $XDG_DATA_HOME/proj/
	// os.UserHomeDir()/.local/share/proj/
	// os.UserDataDir()/proj
	return []string { "." }
}
func TemplateRootDir() string {
	// check if it exists - return the first one found
	// if all not present, return the last
	return TemplateRootPaths()[0]
}
func resolve(desc string, components ...string) (string, error) {
	path := filepath.Join(components...)

	expanded, err := expandPath(path)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to expand %s", desc), slog.String("Path", path), slog.Any("Error", err))
	}
	absPath, err := filepath.Abs(expanded)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to resolve %s", desc), slog.String("Path", path), slog.String("Expanded", expanded), slog.Any("Error", err))
	}
	slog.Debug(fmt.Sprintf("Resolved %s", desc), slog.Any("Components", components), slog.String("Expanded", expanded), slog.String("Absolute", absPath))
	return absPath, err
}

func expandPath(path string) (string, error) {
	expanded := os.ExpandEnv(path)
	if strings.HasPrefix(expanded, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		if expanded == "~" {
			return home, nil
		}
		if strings.HasPrefix(expanded, "~/") {
			expanded = filepath.Join(home, expanded[2:])
		}
	}
	return expanded, nil
}

