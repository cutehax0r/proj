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
	// `new` creates a new project
	// `add` adds a file to an existing project
	// `config.global` edits the global config
	// `config.local` edits the config for a project
	// `list templates` lists the local templates available
	// `list [template] files` lists the files you can add to a template
	// long term a way to search for install, and remove templates would be good.
	// want a `add` to add files to aproject
	// figure out how to do templating of file names. Just use go?
	// figure out how to support scripting: just one script that runs after the template step
	// want to use lua for that.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
