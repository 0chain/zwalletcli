package cmd

import (
	"fmt"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var verifycmd = &cobra.Command{
	Use:   "verify",
	Short: "verify transaction",
	Long: `verify transaction.
	        <hash>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("hash") == false {
			ExitWithError("Error: hash flag is missing")
		}
		hash := cmd.Flag("hash").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, gTxnFee, nonce)
		if err != nil {
			ExitWithError(err)
		}

		txn.SetTransactionHash(hash)
		wg.Add(1)
		err = txn.Verify()
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}
		if statusBar.success {
			statusBar.success = false
			fmt.Printf("\nTransaction verification success\nTransactionStatus: %v\nTransactionOutput: %v",
				txn.GetVerifyConfirmationStatus(), txn.GetVerifyOutput())
			return
		}
		ExitWithError("\nVerification failed." + statusBar.errMsg + "\n")
	},
}

func init() {
	rootCmd.AddCommand(verifycmd)
	verifycmd.PersistentFlags().String("hash", "", "hash of the transaction")
	verifycmd.MarkFlagRequired("hash")
}
