package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/client"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
	"strconv"
)

func checkBalanceBeforeSend(tokens, fee uint64) {
	b, err := client.GetBalance()
	if err != nil {
		return // continue sending txn even if getBalance fails
	}

	if uint64(b.Balance) < tokens+fee {
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

		hash, _, nonce, _, err := zcncore.Send(toClientID, tokens, desc)
		if err != nil {
			ExitWithError("Send tokens failed. " + err.Error())
		}
		if doJSON {
			j := map[string]string{
				"status": "success",
				"tx":     hash,
				"nonce":  strconv.FormatInt(nonce, 10)}
			util.PrintJSON(j)
			return
		}
		fmt.Println("Send tokens success: ", hash)
		return
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
