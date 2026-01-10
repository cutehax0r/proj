package config

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitGlobalConfig(file string) error {
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
			configPath := filepath.Join(xdgConfigHome, "proj")
			viper.AddConfigPath(configPath)
		}
		viper.AddConfigPath("$HOME/.config/proj")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Unable to load global configuration", "Error", err)
		return err
	}

	slog.Debug("Loaded global configuration", "file", viper.ConfigFileUsed())
	return nil
}
