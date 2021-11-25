package cmd

import (
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:     "convert",
	Aliases: []string{"c, convert"},
	Short:   "Convert to ZCN tokens from WZCN (wrapped ZCN) tokens",
	Long:    `Convert to ZCN tokens from WZCN (wrapped ZCN) tokens`,
	Example: "-convert",
	Args:    cobra.MinimumNArgs(3),
	Version: "1.0.0",
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()

		if fflags.Changed("tokens") == false {
			ExitWithError("Error: tokens flag is missing")
		}
		if fflags.Changed("fee") == false {
			ExitWithError("Error: fee flag is missing")
		}

		//zcnbridge.ConfirmEthereumTransaction("", 50, time.Second)
		//
		//token, err := cmd.Flags().GetFloat64("tokens")
		//if err != nil {
		//	ExitWithError("Error: invalid 'tokens' flag", err)
		//}
		//fee := float64(0)
		//fee, err = cmd.Flags().GetFloat64("fee")

		// Steps:

		// 1.
		// Sender: the client who owns WZCN
		// Spender: the bridge contract
		// Before burning the client should approve required amount to spend by spender

		// 2.
		// Check the balance: sender should have enough amounts

		// 1. SDK: call burn in Ethereum (who is the private key owner to sign transaction:
		// SDK or Authorizer function)
		// 2. SDK: Check transaction status (SDK client function)
		// 3. SDK: Call mint sc

		//wg := &sync.WaitGroup{}
		//statusBar := &ZCNStatus{wg: wg}
		//wg.Add(1)
		//err = zcncore.GetLockConfig(statusBar)
		//if err != nil {
		//	ExitWithError(err)
		//}
		//wg.Wait()
		//if statusBar.success {
		//	fmt.Printf("\nConfiguration:\n %v\n", statusBar.errMsg)
		//	return
		//}
		//ExitWithError("\nFailed to get lock config." + statusBar.errMsg + "\n")
		//return
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	lockcmd.PersistentFlags().Float64("tokens", 0, "Number to tokens to exchange")
	lockcmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")

	_ = lockcmd.MarkFlagRequired("tokens")
	_ = lockcmd.MarkFlagRequired("fee")
}
