package cmd

import (
	"fmt"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var getnoncecmd = &cobra.Command{
	Use:   "getnonce",
	Short: "Get nonce from sharders",
	Long:  `Get nonce from sharders`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetNonce(statusBar)
		if err != nil {
			ExitWithError(err)
			return
		}
		wg.Wait()
		b := int64(0)
		if statusBar.success {
			b = statusBar.nonce
		}
		fmt.Printf("\nNonce: %v\n", b)
	},
}

func init() {
	rootCmd.AddCommand(getnoncecmd)
}
