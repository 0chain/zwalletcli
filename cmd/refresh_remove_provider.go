package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var refreshRemoveProvider = &cobra.Command{
	Use:   "rrp",
	Short: "refresh remove provider",
	Long:  "refresh remove provider",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		hash, _, _, _, err := zcncore.RefreshRemoveProviders()
		if err != nil {
			ExitWithError("Refresh failed : ", err.Error())
		}

		fmt.Println("Refresh success with transaction hash : ", hash)
	},
}

func init() {
	rootCmd.AddCommand(refreshRemoveProvider)
}
