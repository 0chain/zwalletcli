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

var addHardForkCmd = &cobra.Command{
	Use:   "add-hardfork",
	Short: "Add hardfork",
	Long:  `Add hardfork`,
	Args:  cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		input := new(zcncore.InputMap)
		input.Fields = setupInputMap(cmd.Flags(), "names", "rounds")
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
		if err = txn.AddHardfork(input); err != nil {
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
	addHardForkCmd.PersistentFlags().StringSliceP("names", "n", nil, "list of names")
	addHardForkCmd.PersistentFlags().StringSliceP("rounds", "r", nil, "list of rounds")

	rootCmd.AddCommand(addHardForkCmd)

}
