package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var installCmd = &cobra.Command{
	Use: "install [source] [target]",
	Args: cobra.RangeArgs(1, 2),
	Short: "Install a template",
	Long: "Add a template to the local template-root from a git repository",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	Run: runInstaller,
}

func init() {
	rootCmd.AddCommand(installCmd)
	viper.BindPFlags(installCmd.Flags())
}

func runInstaller(cmd *cobra.Command, args []string) {
	var source, target string

	source = args[0]

	if len(args) == 2 {
		target = args[1]
	} else {
		target = source // todo, "trim" this to just the last path component
	}

	slog.Debug("Execute Install command", slog.String("source", source), slog.String("target", target))
}
