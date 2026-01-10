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
		slog.Error("Global configuration load failure", "Error", err)
		return err
	}

	absolutePath, err := filepath.Abs(viper.ConfigFileUsed())
	if err != nil {
		slog.Error("Global configuration file path resolution failed", slog.String("globalConfig", viper.ConfigFileUsed()))
		return err
	}
	slog.Debug("Loaded global configuration", "file", viper.ConfigFileUsed())

	viper.Set("global-config", absolutePath)
	return nil
}
