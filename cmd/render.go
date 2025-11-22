package cmd

import (
	"github.com/JamesClonk/plato/pkg/render"
	"github.com/spf13/cobra"
)

var (
	removeTerraformFiles bool
	removeAllDirectories bool
)

var renderCmd = &cobra.Command{
	Use:   "render [OPTIONS]",
	Short: "Render all template files and inject secrets",
	Long: `Render all template files from 'plato.source' into 'plato.target',
and also inject all configuration data and secrets from plato.yaml, secrets.yaml and further files.
The secrets are automatically decrypted on the fly via SOPS.`,
	Run: func(cmd *cobra.Command, args []string) {
		render.RenderTemplates(removeTerraformFiles, removeAllDirectories)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().BoolVarP(&removeTerraformFiles, "cleanup-terraform", "t", false, "Cleanup all .terraform directories in target path before rendering")
	renderCmd.Flags().BoolVarP(&removeAllDirectories, "remove-directories", "d", false, "Clean entire target path before rendering")
}
