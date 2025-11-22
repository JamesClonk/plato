package cmd

import (
	"github.com/JamesClonk/plato/pkg/config"
	"github.com/JamesClonk/plato/pkg/render"
	"github.com/spf13/cobra"
)

var (
	removeTerraformFiles bool
	removeAllDirectories bool
)

var renderCmd = &cobra.Command{
	Use:   "render [OPTIONS]",
	Short: "Renders all template files and inject secrets",
	Long: `Renders all template files from 'plato.source' into 'plato.target',
and injects all configuration data and secrets from plato.yaml and secrets.yaml.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		render.RenderTemplates(removeTerraformFiles, removeAllDirectories)
	},
}

func init() {
	rootCmd.AddCommand(renderCmd)
	renderCmd.Flags().BoolVarP(&removeTerraformFiles, "cleanup-terraform", "t", false, "Cleanup all .terraform directories in target path before rendering")
	renderCmd.Flags().BoolVarP(&removeAllDirectories, "remove-directories", "d", false, "Clean entire target path before rendering")
}
