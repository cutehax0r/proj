/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "proj",
	Short: "Tool for creating projects and files from templates",
	Long: `Create projects using templates stored in ~/.local/share/proj/FOO/new.
Add files to those projects using templates in ~/.local/share/proj/FOO/files/BAR.

"proj new rails foobar": create a rails project called foobar
"proj new model baz": create a model called baz in an existing rails project
`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
