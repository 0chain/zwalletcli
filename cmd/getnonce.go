package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/client"
	"github.com/spf13/cobra"
)

var getnoncecmd = &cobra.Command{
	Use:   "getnonce",
	Short: "Get nonce from sharders",
	Long:  `Get nonce from sharders`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		bal, err := client.GetBalance()
		if err != nil {
			ExitWithError(err)
			return
		}
		fmt.Printf("\nNonce: %v\n", bal.Nonce)
	},
}

func init() {
	rootCmd.AddCommand(getnoncecmd)
}
