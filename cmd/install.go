package cmd

import (
	"log/slog"
	"os"
	"path/filepath"
	"proj/internal/generator"
	"proj/internal/paths"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var installCmd = &cobra.Command{
	Use:   "install [source] [target]",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Install a template",
	Long:  "Add a template to the local template-root from a git repository",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		return nil
	},
	Run: runInstaller,
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringP("template-root", "s", paths.TemplateRootDir(), "Path containing project templates")

	viper.BindPFlags(installCmd.Flags())
}

func runInstaller(cmd *cobra.Command, args []string) {
	templateName := args[0]

	targetName := filepath.Base(templateName)
	if len(args) == 2 {
		targetName = args[1]
	}

	slog.Debug("Execute Install command", slog.String("templateName", templateName), slog.String("targetName", targetName))

	templateRoot := viper.GetString("template-root")
	templateGit := viper.GetString("template-git")

	cfg, err := generator.InstallerConfig(templateName, targetName, templateRoot, templateGit)
	if err != nil {
		slog.Error("Failed to setup installer config", slog.Any("error", err))
		os.Exit(1)
	}

	installer, err := generator.NewInstaller(cfg, templateRoot, templateGit)
	if err != nil {
		slog.Error("Failed to create installer", slog.Any("error", err))
		os.Exit(1)
	}

	if err := installer.Install(); err != nil {
		slog.Error("Installation failed", slog.Any("error", err))
		os.Exit(1)
	}
}
