package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var unlockcmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock tokens",
	Long: `Unlock previously locked tokens .
	        <pool_id> [transaction fee]`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("pool_id") == false {
			ExitWithError("Error: pool_id flag is missing")
		}
		pool_id := cmd.Flag("pool_id").Value.String()
		fee := float64(0)
		fee, err := cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee), nonce)
		if err != nil {
			ExitWithError(err)
		}
		wg.Add(1)
		err = txn.UnlockTokens(pool_id)
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}
		if statusBar.success {
			statusBar.success = false
			wg.Add(1)
			err := txn.Verify()
			if err == nil {
				wg.Wait()
			} else {
				ExitWithError(err.Error())
			}
			if statusBar.success {
				switch txn.GetVerifyConfirmationStatus() {
				case zcncore.ChargeableError:
					ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
				case zcncore.Success:
					fmt.Printf("\nUnlock tokens success\nHash: %v\n", txn.GetTransactionHash())
				default:
					ExitWithError("\nFailed to unlock tokens. Unknown status code: " +
						strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
				}
				return
			}
		}
		ExitWithError("\nFailed to unlock tokens. " + statusBar.errMsg + "\n")
		return
	},
}

func init() {
	rootCmd.AddCommand(unlockcmd)
	unlockcmd.PersistentFlags().String("pool_id", "", "Pool ID - hash of the locked transaction")
	unlockcmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	unlockcmd.MarkFlagRequired("pool_id")
}
