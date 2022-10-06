package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var updateStoragScConfigCmd = &cobra.Command{
	Use:   "sc-update-config",
	Short: "Update the storage smart contract",
	Long:  `Update the storage smart contract.`,
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
		txn, err := zcncore.NewTransaction(statusBar, transactionFee(), nonce)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		if err = txn.StorageScUpdateConfig(input); err != nil {
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

		if statusBar.success {
			switch txn.GetVerifyConfirmationStatus() {
			case zcncore.ChargeableError:
				ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
			case zcncore.Success:
				fmt.Printf("storagesc smart contract settings updated\nHash: %v\n", txn.GetTransactionHash())
			default:
				ExitWithError("\nExecute storagesc smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateStoragScConfigCmd)
	updateStoragScConfigCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateStoragScConfigCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
