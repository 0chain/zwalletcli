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
	updateVestingPoolConfigCmd.PersistentFlags().Int("max_description_length", 0, "max length for descriptions")
	updateVestingPoolConfigCmd.PersistentFlags().Int("max_destinations", 0, "max destinations allowed")
	updateVestingPoolConfigCmd.PersistentFlags().Float64("min_lock", 0.0, "minimum lock for vesting")
	updateVestingPoolConfigCmd.PersistentFlags().Duration("min_duration", 0.0, "minimum duration for vesting")
	updateVestingPoolConfigCmd.PersistentFlags().Duration("max_duration", 0.0, "max duration for vesting")
}
