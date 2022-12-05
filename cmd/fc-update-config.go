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

var updateFaucetCmd = &cobra.Command{
	Use:   "fc-update-config",
	Short: "Update the Faucet smart contract",
	Long:  `Update the Faucet smart contract.`,
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
		txn, err := zcncore.NewTransaction(statusBar, gTxnFee, nonce)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		if err = txn.FaucetUpdateConfig(input); err != nil {
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

		if statusBar.success {
			//fmt.Printf("Hash:%v\nNonce:%v\n", txn.GetTransactionHash(), txn.GetTransactionNonce())
			switch txn.GetVerifyConfirmationStatus() {
			case zcncore.ChargeableError:
				ExitWithError(strings.Trim(txn.GetVerifyOutput(), "\""))
			case zcncore.Success:
				fmt.Printf("faucet smart contract settings updated\nHash: %v\n", txn.GetTransactionHash())
			default:
				ExitWithError("Execute faucet smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
			return
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}

	},
}

func init() {
	rootCmd.AddCommand(updateFaucetCmd)
	updateFaucetCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateFaucetCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
