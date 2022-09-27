package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var createWalletCmd = &cobra.Command{
	Use:   "create-wallet",
	Short: "Create wallet and logs it into stdout",
	Long:  `Create wallet and logs it into standard output`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		statusBar, err := createWallet()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create wallet: %v", err)
		}
		fmt.Fprintf(os.Stdout, "\n\t======WALLET_START======\n\n")
		fmt.Fprintf(os.Stdout, statusBar.walletString)
		fmt.Fprintf(os.Stdout, "\n\n\t======WALLET_END======\n")
	},
}

func init() {
	rootCmd.AddCommand(WithoutWallet(createWalletCmd))
}
