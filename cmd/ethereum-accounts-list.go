package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var listEthAccounts = &cobra.Command{
	Use:   "eth-list-accounts",
	Short: "list Ethereum account registered in local key chain",
	Long:  `list available ethereum accounts`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(_ *cobra.Command, _ []string) {
		accounts := zcnbridge.ListStorageAccounts()
		if len(accounts) == 0 {
			fmt.Printf("Accounts not found")
		}

		fmt.Println("Ethereum available account:")
		for _, acc := range accounts {
			fmt.Println(acc.Hex())
		}
	},
}

func init() {
	rootCmd.AddCommand(listEthAccounts)
}
