package cmd

import (
	"github.com/JamesClonk/plato/pkg/config"
	"github.com/JamesClonk/plato/pkg/render"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/spf13/cobra"
)

var (
	templateFile string
	outputFile   string
)

var fileCmd = &cobra.Command{
	Use:   "file [template-file] [output-file]",
	Short: "Renders given template file and inject secrets",
	Long: `Renders given template file to either STDOUT or to an output file (optional parameter)',
and injects all configuration data and secrets from plato.yaml and secrets.yaml.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		templateFile = "/dev/stdin" // STDIN
		outputFile = "/dev/stdout"  // STDOUT

		if len(args) > 0 {
			templateFile = args[0]
			// -- TODO: make templateFile optional too. if not given, automatically assume /dev/stdin for input, easy peasy lemon squeezy!
		}
		if len(args) > 1 {
			outputFile = args[1]
		}

		if templateFile == "/dev/stdin" {
			log.Disable() // disable all log output if we read from STDIN
		} else if outputFile == "/dev/stdout" {
			log.Disable() // disable all log output if we render to STDOUT
		}

		config.InitConfig()
		render.RenderFile(templateFile, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
}
