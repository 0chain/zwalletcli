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
			var accountAddrIndex zcnbridge.AccountAddressIndex

			if c.Flags().Changed(OptionAccountIndex) {
				var err error
				accountAddrIndex.AccountIndex, err = c.Flags().GetInt(OptionAccountIndex)
				if err != nil {
					ExitWithError(err)
					return
				}
			}

			if c.Flags().Changed(OptionAddressIndex) {
				var err error
				accountAddrIndex.AddressIndex, err = c.Flags().GetInt(OptionAddressIndex)
				if err != nil {
					ExitWithError(err)
					return
				}
			}

			_, err := zcnbridge.ImportAccount(path, mnemonic, password, accountAddrIndex)
			if err != nil {
				ExitWithError(err)
				return
			}
		},
	}

	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String(OptionMnemonic, "", "Ethereum mnemonic")
	cmd.PersistentFlags().String(OptionKeyPassword, "", "Password to lock and unlock account to sign transaction")
	cmd.PersistentFlags().Int(OptionAccountIndex, 0, "Index of the account to use, default 0")
	cmd.PersistentFlags().Int(OptionAddressIndex, 0, "Index of the address to use, default 0")
	cmd.PersistentFlags().String(OptionConfigFolder, GetConfigDir(), "Home config directory")

	cmd.MarkFlagRequired(OptionMnemonic)
	cmd.MarkFlagRequired(OptionKeyPassword)
}
