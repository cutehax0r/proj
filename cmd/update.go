package cmd

import (
	"errors"
	"log/slog"
	"os"

	"proj/internal/generator"
	"proj/internal/paths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update [target]",
	Args:  cobra.RangeArgs(0, 1),
	Short: "Update installed templates",
	Long:  "Update installed plugins automatically or all at once",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		return nil
	},
	Run: runUpdater,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().StringP("template-root", "s", paths.TemplateRootDir(), "Path containing project templates")
	updateCmd.Flags().BoolP("force", "f", false, "Delete and reinstall even if the repository has uncommitted changes")

	viper.BindPFlags(updateCmd.Flags())
}

func runUpdater(cmd *cobra.Command, args []string) {
	templateRoot := viper.GetString("template-root")

	var templateNames []string

	if len(args) == 1 {
		templateNames = []string{args[0]}
	} else {
		entries, err := os.ReadDir(templateRoot)
		if err != nil {
			slog.Error("Failed to read template root", slog.Any("error", err))
			os.Exit(1)
		}
		for _, e := range entries {
			if e.IsDir() {
				templateNames = append(templateNames, e.Name())
			}
		}
	}

	for _, name := range templateNames {
		cfg, err := generator.UpdaterConfig(name, templateRoot)
		if err != nil {
			slog.Error("Failed to setup updater config", slog.String("name", name), slog.Any("error", err))
			continue
		}

		updater, err := generator.NewUpdater(cfg)
		if err != nil {
			slog.Error("Failed to create updater", slog.String("name", name), slog.Any("error", err))
			continue
		}

		if err := updater.Update(); err != nil {
			if errors.Is(err, generator.ErrNotInstalled) {
				slog.Error("Template not found", slog.String("name", name))
			} else {
				slog.Error("Update failed", slog.String("name", name), slog.Any("error", err))
			}
		}
	}
}
