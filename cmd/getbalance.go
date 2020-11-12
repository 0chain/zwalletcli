package cmd

import (
	"fmt"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var getbalancecmd = &cobra.Command{
	Use:   "getbalance",
	Short: "Get balance from sharders",
	Long:  `Get balance from sharders`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetBalance(statusBar)
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}
		if statusBar.success {
			token := zcncore.ConvertToToken(statusBar.balance)
			tokenUSD, err := zcncore.ConvertTokenToUSD(token)
			if err != nil {
				ExitWithError("\nGet balance failed. " + err.Error() + "\n")
			} else {
				fmt.Printf("\nBalance: %v (%v USD)\n", token, tokenUSD)
			}
		} else {
			ExitWithError("\nGet balance failed. " + statusBar.errMsg + "\n")
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(getbalancecmd)
}
