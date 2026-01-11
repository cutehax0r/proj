package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"proj/internal/paths"

	"github.com/spf13/viper"
)

func InitGlobal(file string) error {
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		// I'm choosing these to avoid a library dependency on https://github.com/adrg/xdg
		// on Macos you get "surprising paths" like ~/Library/Application Support/ for
		// USER_DATA_HOME which nobody expects.

		// do something similar for XDG_DATA_HOME
		// use os.UserHomeDir()
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			configPath := filepath.Join(xdgConfigHome, paths.GlobalConfigDir)
			viper.AddConfigPath(configPath)
		}
		viper.AddConfigPath("$HOME/.config/proj")
		viper.AddConfigPath(".")
	}

	// needs to set a default for template root  XDG_DATA_HOME


	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Global configuration load failure", "Error", err)
		return err
	}

	absolutePath, err := filepath.Abs(viper.ConfigFileUsed())
	if err != nil {
		slog.Error("Global configuration file path resolution failed", slog.String("globalConfig", viper.ConfigFileUsed()))
		return err
	}
	slog.Debug("Loaded global configuration", "file", viper.ConfigFileUsed())

	viper.Set("global-config-file", absolutePath)
	return nil
}
