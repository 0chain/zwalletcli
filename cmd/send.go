package cmd

import (
	"fmt"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var sendcmd = &cobra.Command{
	Use:   "send",
	Short: "Send ZCN token to another wallet",
	Long: `Send ZCN token to another wallet.
	        <to_client_id> <token> <description> [transaction fee]`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("to_client_id") == false {
			ExitWithError("Error: to_client_id flag is missing")
		}
		if fflags.Changed("token") == false {
			ExitWithError("Error: token flag is missing")
		}
		if fflags.Changed("desc") == false {
			ExitWithError("Error: Description flag is missing")
		}
		to_client_id := cmd.Flag("to_client_id").Value.String()
		token, err := cmd.Flags().GetFloat64("token")
		if err != nil {
			ExitWithError("Error: invalid token.", err)
		}
		desc := cmd.Flag("desc").Value.String()
		fee := float64(0)
		fee, err = cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee))
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
				fmt.Println("Send token success")
				return
			}
		}
		ExitWithError("Send token failed. " + statusBar.errMsg)
	},
}

func init() {
	rootCmd.AddCommand(sendcmd)
	sendcmd.PersistentFlags().String("to_client_id", "", "to_client_id")
	sendcmd.PersistentFlags().Float64("token", 0, "Token to send")
	sendcmd.PersistentFlags().String("desc", "", "Description")
	sendcmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	sendcmd.MarkFlagRequired("to_client_id")
	sendcmd.MarkFlagRequired("token")
	sendcmd.MarkFlagRequired("desc")
}
