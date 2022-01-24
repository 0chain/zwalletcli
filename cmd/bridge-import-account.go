package cmd

import (
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

//goland:noinspection GoUnhandledErrorResult
func init() {
	cmd := &cobra.Command{
		Use:   "bridge-import-account",
		Short: "Import Ethereum account to local key storage (default $HOME/.zcn/wallets)",
		Long:  "Import account to local key storage using mnemonic, protected with password (default $HOME/.zcn/wallets)",
		Args:  cobra.MinimumNArgs(0),
		Run: func(c *cobra.Command, _ []string) {
			check(c, OptionMnemonic, OptionKeyPassword)

			path := c.Flag(OptionConfigFolder).Value.String()
			mnemonic := c.Flag(OptionMnemonic).Value.String()
			password := c.Flag(OptionKeyPassword).Value.String()

			_, err := zcnbridge.ImportAccount(path, mnemonic, password)
			if err != nil {
				ExitWithError(err)
				return
			}
		},
	}

	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String(OptionMnemonic, "", "Ethereum mnemonic")
	cmd.PersistentFlags().String(OptionKeyPassword, "", "Password to lock and unlock account to sign transaction")
	cmd.PersistentFlags().String(OptionConfigFolder, GetConfigDir(), "Home config directory")

	cmd.MarkFlagRequired(OptionMnemonic)
	cmd.MarkFlagRequired(OptionKeyPassword)
}
