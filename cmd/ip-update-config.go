package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var updateInterestPoolConfigCmd = &cobra.Command{
	Use:   "ip-update-config",
	Short: "Update the interest pool configurations.",
	Long:  `Update the interest pool configurations.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		input := new(zcncore.InputMap)
		input.Fields, err = setupInputMap(cmd.Flags())
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		statusBar := &ZCNStatus{wg: &wg}
		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		if err = txn.InterestPoolUpdateConfig(input); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		statusBar.success = false
		wg.Add(1)
		if err = txn.Verify(); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Printf("interest pool smart contract settings updated\nHash: %v\n", txn.GetTransactionHash())
	},
}

func init() {
	rootCmd.AddCommand(updateInterestPoolConfigCmd)
	updateInterestPoolConfigCmd.PersistentFlags().Float64("min_lock", 0, "minimum tokens that can be locked")
	updateInterestPoolConfigCmd.PersistentFlags().Float64("apr", 0, "apr, apr")
	updateInterestPoolConfigCmd.PersistentFlags().Duration("min_lock_period", 0.0, "minimum lock period")
	updateInterestPoolConfigCmd.PersistentFlags().Float64("max_mint", 0.0, "minimum tokes interest sc can mint")
}
