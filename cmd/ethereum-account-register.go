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
		fflags := cmd.Flags()
		if fflags.Changed("mnemonic") == false {
			ExitWithError("Error: 'mnemonic' flag is missing")
		}
		if fflags.Changed("password") == false {
			ExitWithError("Error: 'password' flag is missing")
		}

		mnemonic := cmd.Flag("mnemonic").Value.String()
		password := cmd.Flag("password").Value.String()

		err := zcnbridge.ImportAccount(mnemonic, password)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(registerAccount)

	registerAccount.PersistentFlags().String("mnemonic", "", "mnemonic")
	registerAccount.PersistentFlags().String("password", "", "password")

	_ = registerAccount.MarkFlagRequired("mnemonic")
	_ = registerAccount.MarkFlagRequired("password")
}
