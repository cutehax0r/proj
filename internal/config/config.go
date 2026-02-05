package config

import (
	"log/slog"
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

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Global configuration load failure", slog.String("file", file), "Error", err)
		return err
	}
	slog.Debug("global configuration read", "file", viper.ConfigFileUsed())
	viper.Set("global-config-file", viper.ConfigFileUsed())
	return nil
}

func InitTemplate(file string) error {
	conf := viper.New()
	conf.SetConfigFile(file)
	if err := conf.ReadInConfig(); err != nil {
		slog.Error("Template configuration load failure", slog.String("file", file), slog.Any("error", err))
		return err
	}
	conf.Set("template-config-file", conf.ConfigFileUsed())
	viper.MergeConfigMap(conf.AllSettings())
	return nil
}
