package cmd

import (
	"os"

	"github.com/JamesClonk/plato/pkg/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plato",
	Short: "SOPS Template Renderer - CLI",
	Long:  `The plato CLI tool is used to render template files with automatic SOPS secret injection.`,
}

func Execute(v, bs, gh string) {
	// setup shell completion
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use != "completion" {
			cmd.PreRun = func(cmd *cobra.Command, args []string) {
				config.InitConfig()
			}
		}
	}

	// execute CLI
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
