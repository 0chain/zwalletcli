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
		if err != nil {
			ExitWithError(err)
			return
		}
		wg.Wait()
		if !statusBar.success {
			ExitWithError(fmt.Sprintf("\nFailed to get balance: %s\n", statusBar.errMsg))
			return
		}
		b := statusBar.balance
		usd, err := zcncore.ConvertTokenToUSD(b.ToToken())
		if err != nil {
			ExitWithError(fmt.Sprintf("\nBalance: %v (Failed to get USD: %v)\n", b, err))
			return
		}
		fmt.Printf("\nBalance: %v (%.2f USD)\n", b, usd)
	},
}

func init() {
	rootCmd.AddCommand(getbalancecmd)
}
