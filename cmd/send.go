package cmd

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

func checkBalanceBeforeSend(tokens, fee uint64) {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	wg.Add(1)
	err := zcncore.GetBalance(statusBar)
	if err != nil {
		return // continue sending txn even if getBalance fails
	}
	wg.Wait()
	if !statusBar.success {
		return // continue sending txn even if getBalance fails
	}
	b := statusBar.balance

	if uint64(b) < tokens+fee {
		ExitWithError("Insufficient balance for this transaction.")
	}
	return
}

var sendcmd = &cobra.Command{
	Use:   "send",
	Short: "Send ZCN tokens to another wallet",
	Long: `Send ZCN tokens to another wallet.
	        <to_client_id> <tokens> <description> [transaction fee]`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("to_client_id") == false {
			ExitWithError("Error: to_client_id flag is missing")
		}
		if fflags.Changed("tokens") == false {
			ExitWithError("Error: tokens flag is missing")
		}
		if fflags.Changed("desc") == false {
			ExitWithError("Error: Description flag is missing")
		}
		tokenZCN, err := cmd.Flags().GetFloat64("tokens")
		if err != nil {
			ExitWithError("Error: invalid 'tokens' flag", err)
		}
		if tokenZCN < 0 {
			ExitWithError("invalid tokens amount: negative")
		}

		toClientID := cmd.Flag("to_client_id").Value.String()
		doJSON, _ := cmd.Flags().GetBool("json")
		desc := cmd.Flag("desc").Value.String()

		tokens := zcncore.ConvertToValue(tokenZCN)
		fee := getTxnFee()
		if fee > 0 {
			checkBalanceBeforeSend(tokens, fee)
		}

		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, fee, nonce)
		if err != nil {
			ExitWithError(err)
		}

		wg.Add(1)
		err = txn.Send(toClientID, tokens, desc)
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
				if doJSON {
					j := map[string]string{
						"status": "success",
						"tx":     txn.Hash(),
						"nonce":  strconv.FormatInt(txn.GetTransactionNonce(), 10)}
					util.PrintJSON(j)
					return
				}
				fmt.Println("Send tokens success: ", txn.Hash())
				return
			}
		}
		ExitWithError("Send tokens failed. " + statusBar.errMsg)
	},
}

func init() {
	rootCmd.AddCommand(sendcmd)
	sendcmd.PersistentFlags().String("to_client_id", "", "to_client_id")
	sendcmd.PersistentFlags().Float64("tokens", 0, "Token to send")
	sendcmd.PersistentFlags().String("desc", "", "Description")
	sendcmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	sendcmd.MarkFlagRequired("to_client_id")
	sendcmd.MarkFlagRequired("tokens")
	sendcmd.MarkFlagRequired("desc")
	sendcmd.Flags().Bool("json", false, "pass this option to print response as json data")
}
