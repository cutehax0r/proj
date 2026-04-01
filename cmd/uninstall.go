package cmd

import (
	"log/slog"
	"os"
	"proj/internal/generator"
	"proj/internal/paths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [target]",
	Args:  cobra.ExactArgs(1),
	Short: "Uninstall a template",
	Long:  "Remove a template installed from github (not manually installed)",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		return nil
	},
	Run: runUninstaller,
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringP("template-root", "s", paths.TemplateRootDir(), "Path containing project templates")
	uninstallCmd.Flags().Bool("force", false, "Force removal even if not from git or has local changes")

	viper.BindPFlags(uninstallCmd.Flags())
}

func runUninstaller(cmd *cobra.Command, args []string) {
	targetName := args[0]

	slog.Debug("Execute Uninstall command", slog.String("targetName", targetName))

	templateRoot := viper.GetString("template-root")

	cfg, err := generator.UninstallerConfig(targetName, templateRoot)
	if err != nil {
		slog.Error("Failed to setup uninstaller config", slog.Any("error", err))
		os.Exit(1)
	}

	uninstaller, err := generator.NewUninstaller(cfg)
	if err != nil {
		slog.Error("Failed to create uninstaller", slog.Any("error", err))
		os.Exit(1)
	}

	if err := uninstaller.Uninstall(); err != nil {
		slog.Error("Uninstallation failed", slog.Any("error", err))
		os.Exit(1)
	}
}
