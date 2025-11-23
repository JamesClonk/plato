package cmd

import (
	"github.com/JamesClonk/plato/pkg/config"
	"github.com/JamesClonk/plato/pkg/render"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
)

var templateCmd = &cobra.Command{
	Use:   "template [input] [output]",
	Short: "Renders given template and injects secrets",
	Long: `Renders given template via STDIN or file to either STDOUT or an output file,',
and injects all configuration data and secrets from plato.yaml and (optional) secrets.yaml.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		inputFile = "/dev/stdin"   // use STDIN as default
		outputFile = "/dev/stdout" // use STDOUT as default

		if len(args) > 0 {
			inputFile = args[0]
		}
		if len(args) > 1 {
			outputFile = args[1]
		}

		if inputFile == "/dev/stdin" {
			log.Disable() // disable all log output if we read from STDIN
		} else if outputFile == "/dev/stdout" {
			log.Disable() // disable all log output if we render to STDOUT
		}

		config.InitConfig()
		render.RenderFile(inputFile, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
