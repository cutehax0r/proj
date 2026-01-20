package paths

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	// https://github.com/adrg/xdg
)

const GlobalConfigDir = "proj"
const GlobalConfigFile = "proj"
const TargetConfigFile = "proj.yml"
const TargetConfigFileDir = ".proj"
const TemplateConfigFile = "proj.yml"

// os.UserConfDir() and os.UserDataDir() on mac OS return things that a CLI-user won't expect:
// paths in ~/Library. They'll expect UNIX-y paths so we construct those instead.

func GlobalConfigPaths() (paths []string) {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		paths = append(paths, filepath.Join(xdgConfigHome, GlobalConfigDir))
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Couldn't get home directory", slog.Any("Error", err))
	}
	paths = append(paths, filepath.Join(homeDir, ".config", GlobalConfigDir))

	paths = append(paths, ".")

	confDir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Couldn't get user config dir", slog.Any("Error", err))
	}
	paths = append(paths, confDir)

	return uniqPreserveOrder(paths)
}

func TemplateRootPaths() (paths []string) {
	if xdgDataHome := os.Getenv("XDG_DATA_HOME"); xdgDataHome != "" {
		paths = append(paths, filepath.Join(xdgDataHome, GlobalConfigDir))
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Couldn't get home directory", slog.Any("Error", err))
	}
	paths = append(paths, filepath.Join(homeDir, ".local", "share", GlobalConfigDir))

	dataDir, err := os.UserConfigDir()
	if err != nil {
		slog.Error("Couldn't get user config dir", slog.Any("Error", err))
	}
	paths = append(paths, filepath.Join(dataDir, GlobalConfigDir))
	return uniqPreserveOrder(paths)
}

func TemplateRootDir() string {
	paths := TemplateRootPaths()
	slog.Debug("Finding template root path from", slog.Any("paths", paths))
	for _, p := range  paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return TemplateRootPaths()[0]
}

func uniqPreserveOrder(paths []string) (uniq []string) {
	seen := make(map[string]struct{})
	for _, p := range paths {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			uniq = append(uniq, p)
		}

	}
	return uniq
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
