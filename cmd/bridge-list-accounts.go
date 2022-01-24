package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var listEthAccounts = &cobra.Command{
	Use:   "bridge-list-accounts",
	Short: "List Ethereum account registered in local key chain",
	Long:  `List available ethereum accounts`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(*cobra.Command, []string) {
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
