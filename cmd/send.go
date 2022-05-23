package cmd

import (
	"fmt"
	"sync"
	"strconv"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

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
		to_client_id := cmd.Flag("to_client_id").Value.String()
		token, err := cmd.Flags().GetFloat64("tokens")
		if err != nil {
			ExitWithError("Error: invalid 'tokens' flag", err)
		}
		doJSON, _ := cmd.Flags().GetBool("json")
		desc := cmd.Flag("desc").Value.String()
		fee := float64(0)
		fee, err = cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee), nonce)
		if err != nil {
			ExitWithError(err)
		}
		wg.Add(1)
		err = txn.Send(to_client_id, zcncore.ConvertToValue(token), desc)
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
					j := map[string]string {
		  			"status": "success",
		  			"tx": txn.Hash(),
		  			"nonce": strconv.FormatInt(txn.GetTransactionNonce(),10) }
					util.PrintJSON(j)
					return
				}
fmt.Println("Send tokens success:")
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
