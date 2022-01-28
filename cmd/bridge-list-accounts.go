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
	Run: func(cmd *cobra.Command, _ []string) {
		fflags := cmd.Flags()
		path, err := fflags.GetString(OptionConfigFolder)
		if err != nil {
			fmt.Printf("Flag '%s' not found, defaulting to %s\n", OptionConfigFolder, GetConfigDir())
		}

		accounts := zcnbridge.ListStorageAccounts(path)
		if len(accounts) == 0 {
			fmt.Println("Accounts not found")
		}

		fmt.Println("Ethereum available account:")
		for _, acc := range accounts {
			fmt.Println(acc.Hex())
		}
	},
}

func init() {
	f := listEthAccounts
	rootCmd.AddCommand(listEthAccounts)

	f.PersistentFlags().String(OptionConfigFolder, GetConfigDir(), "Configuration dir")
}
