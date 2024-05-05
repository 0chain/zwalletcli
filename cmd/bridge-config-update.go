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

var updateBridgeGlobalConfigCmd = &cobra.Command{
	Use:   "bridge-config-update",
	Short: "Update ZCNSC bridge global settings",
	Long:  `Update ZCNSC bridge global settings.`,
	Args:  cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		input := new(zcncore.InputMap)
		input.Fields = setupInputMap(cmd.Flags(), "keys", "values")
		if err != nil {
			log.Fatal(err)
		}

		var wg sync.WaitGroup
		statusBar := &ZCNStatus{wg: &wg}
		txn, err := zcncore.NewTransaction(statusBar, getTxnFee(), nonce)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		if err = txn.ZCNSCUpdateGlobalConfig(input); err != nil {
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
	cmd := updateBridgeGlobalConfigCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	cmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
