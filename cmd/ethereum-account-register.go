package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var registerAccount = &cobra.Command{
	Use:   "eth-register-account",
	Short: "register ethereum account in local key storage",
	Long:  `register ethereum account using mnemonic and protected with password`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		check(cmd,
			"mnemonic",
			"password")

		mnemonic := cmd.Flag("mnemonic").Value.String()
		password := cmd.Flag("password").Value.String()

		address, err := zcnbridge.ImportAccount(mnemonic, password)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("Imported account, address: %s", address)
	},
}

func init() {
	rootCmd.AddCommand(registerAccount)

	registerAccount.PersistentFlags().String("mnemonic", "", "mnemonic")
	registerAccount.PersistentFlags().String("password", "", "password")

	_ = registerAccount.MarkFlagRequired("mnemonic")
	_ = registerAccount.MarkFlagRequired("password")
}
