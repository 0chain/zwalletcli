package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/transaction"
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

		txn, err := transaction.VerifyTransaction(hash)
		if err != nil {
			ExitWithError(err)
		}

		fmt.Printf("\nTransaction verification success\nTransactionStatus: %v\nTransactionOutput: %v",
			txn.Status, txn.TransactionOutput)
	},
}

func init() {
	rootCmd.AddCommand(verifycmd)
	verifycmd.PersistentFlags().String("hash", "", "hash of the transaction")
	verifycmd.MarkFlagRequired("hash")
}
