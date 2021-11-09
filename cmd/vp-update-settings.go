package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var updateVestingPoolConfigCmd = &cobra.Command{
	Use:   "vp-update-config",
	Short: "Update the vesting pool configurations.",
	Long:  `Update the vesting pool configurations.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		input := new(zcncore.InputMap)
		input.Fields = setupInputMap(cmd.Flags())
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
		if err = txn.VestingUpdateConfig(input); err != nil {
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

		fmt.Printf("vesting smart contract settings updated\nHash: %v\n", txn.GetTransactionHash())
	},
}

func init() {
	rootCmd.AddCommand(updateVestingPoolConfigCmd)
	updateVestingPoolConfigCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateVestingPoolConfigCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
