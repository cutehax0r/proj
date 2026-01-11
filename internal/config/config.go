package config

import (
	"log/slog"
	"path/filepath"
	"proj/internal/paths"

	"github.com/spf13/viper"
)

func InitGlobal(file string) error {
	if file != "" {
		viper.SetConfigFile(file)
	} else {
		configPaths := paths.GlobalConfigPaths()
		slog.Debug("Adding global config paths", slog.Any("Paths", configPaths))
		for _, path := range configPaths { 
			viper.AddConfigPath(path)
		}
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
