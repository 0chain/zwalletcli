package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var updateMinerScConfigCmd = &cobra.Command{
	Use:   "mn-update-config",
	Short: "Update the miner smart contract",
	Long:  `Update the miner smart contract.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		input := new(zcncore.InputMap)
		input.Fields = setupInputMap(cmd.Flags(), "keys", "values")
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		statusBar := &ZCNStatus{wg: &wg}
		txn, err := zcncore.NewTransaction(statusBar, 0, nonce)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		if err = txn.MinerScUpdateConfig(input); err != nil {
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

		fmt.Printf("storagesc smart contract settings updated\nHash: %v\n", txn.GetTransactionHash())
	},
}

func init() {
	rootCmd.AddCommand(updateMinerScConfigCmd)
	updateMinerScConfigCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateMinerScConfigCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
