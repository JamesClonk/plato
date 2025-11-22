package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	version string
	commit  string
	date    string
)

var rootCmd = &cobra.Command{
	Use:   "plato",
	Short: "SOPS Template Renderer - CLI",
	Long:  `The plato CLI tool is used to render template files with automatic SOPS secret injection.`,
}

func Execute(v, c, d string) {
	// store build information
	version = v
	commit = c
	date = d

	// execute CLI
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
