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

var updateGlobalConfigCmd = &cobra.Command{
	Use:   "global-update-config",
	Short: "Update global settings",
	Long:  `Update global settings.`,
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
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(txFee), nonce)
		if err != nil {
			log.Fatal(err)
		}

		if err := txn.AdjustTransactionFee(txVelocity.toZCNFeeType()); err != nil {
			log.Fatal("failed to adjust transaction fee: ", err)
		}

		wg.Add(1)
		if err = txn.MinerScUpdateGlobals(input); err != nil {
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
				fmt.Printf("global settings updated\nHash: %v\n", txn.GetTransactionHash())
			default:
				ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
			return
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}

	},
}

func init() {
	rootCmd.AddCommand(updateGlobalConfigCmd)
	updateGlobalConfigCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateGlobalConfigCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
