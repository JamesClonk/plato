package cmd

import (
	"github.com/JamesClonk/plato/pkg/store"
	"github.com/spf13/cobra"
)

var storeCmd = &cobra.Command{
	Aliases: []string{"save", "store"},
	Use:     "store-secrets",
	Short:   "Stores generated secrets in encrypted SOPS file",
	Long: `Stores all dynamically generated secrets found under 'plato.secrets' back into the encrypted SOPS secret file. Property paths are determined by naming convention.

By convention the filenames under 'plato.secrets'/* will translate into YAML paths, for example:
./rendered/secrets/tls.key -> "tls.key:" in secrets.yaml

It will also re-encrypt all formerly named *.sops_enc files back to their original location under 'plato.source'.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		store.StoreGeneratedSecrets()
	},
}

func init() {
	rootCmd.AddCommand(storeCmd)
}
