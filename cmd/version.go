package cmd

import (
	"github.com/JamesClonk/plato/pkg/config"
	"github.com/JamesClonk/plato/pkg/util/color"
	"github.com/JamesClonk/plato/pkg/util/log"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays PLATO version",
	Long:  `Displays PLATO version`,
	Run: func(cmd *cobra.Command, args []string) {
		config.InitConfig()
		log.Infof("version: \t%s", color.Green(version))
		log.Infof("date: \t%s", color.Green(date))
		log.Infof("commit: \t%s", commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
