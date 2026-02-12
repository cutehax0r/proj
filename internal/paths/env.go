package paths

// TODO: The XDG situation is a bit of a mess on macOS. You can't count on a person to have
// explicitly set the environment variables to make it work. Despite that, most users will expect
// those paths to work. Go provides an os.UserConfigDir() and os.UserDataDir() but on macOS they
// will point to a path in ~/Library which most users will not expect. For this reason I'm choosing
// not to use those functions. Unfortunately this choice means that Windows is not going to work.
//
// https://github.com/adrg/xdg
// Provides an interface to XDG directories along with sensible defaults that might be worth adding.

// TemplateRootPaths() returns candidate paths, but TemplateRootDir() returns the first existing
// one. The naming could be clearer (TemplateRootPathCandidates() and FindTemplateRootDir()).
// maybe merge them?
//
// this whole thing should probably merge with paths.go
//
// There's a lot of filepath.join when we should really be using the Resolve() function

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

const GlobalConfigDir = "proj"
const GlobalConfigFile = "proj"
const TargetConfigFile = "proj.yml"
const TargetConfigFileDir = ".proj"
const TemplateConfigFile = "proj.yml"

func GlobalConfigPaths() (paths []string) {
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		paths = append(paths, filepath.Join(xdgConfigHome, GlobalConfigDir))
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Couldn't get home directory", slog.Any("Error", err))
	}
	paths = append(paths, filepath.Join(homeDir, ".config", GlobalConfigDir))

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
		slog.Debug("Couldn't get home directory", slog.Any("Error", err))
	}
	paths = append(paths, filepath.Join(homeDir, ".local", "share", GlobalConfigDir))

	dataDir, err := os.UserConfigDir()
	if err != nil {
		slog.Debug("Couldn't get user config dir", slog.Any("Error", err))
	}
	paths = append(paths, filepath.Join(dataDir, GlobalConfigDir))
	return uniqPreserveOrder(paths)
}

// what do we do if none of these paths exist?
func TemplateRootDir() string {
	return TemplateRootDirWithFS(afero.NewOsFs())
}

func TemplateRootDirWithFS(fs afero.Fs) string {
	paths := TemplateRootPaths()
	slog.Debug("Finding template root path from", slog.Any("paths", paths))
	for _, p := range paths {
		if _, err := fs.Stat(p); err == nil {
			return p
		}
	}
	return paths[0]
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
