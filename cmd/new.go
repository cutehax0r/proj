/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [kind] [name]",
	Args: cobra.ExactArgs(2),
	Short: "Create a project called NAME based on the KIND template",
	Long: `Create a new project based on a project template.

Uses the project template, name, and standard variables from config to create a new project ready to
work on. Projects typically include skeleton files as well as shell scripts to pre-configure things
like packaging, remote repositories, database creation, etc.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// kind must be one of the root folders in ~/.local/share/proj
		kind := args[0]
		name := args[1]
		fmt.Printf("new called to make a new '%s' called '%s'\n", kind, name)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
