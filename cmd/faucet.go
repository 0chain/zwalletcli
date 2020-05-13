package cmd

import (
	"fmt"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var faucetcmd = &cobra.Command{
	Use:   "faucet",
	Short: "Faucet smart contract",
	Long: `Faucet smart contract.
	        <methodName> <input>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("methodName") == false {
			ExitWithError("Error: Methodname flag is missing")
		}
		if fflags.Changed("input") == false {
			ExitWithError("Error: Input flag is missing")
		}

		methodName := cmd.Flag("methodName").Value.String()
		input := cmd.Flag("input").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			ExitWithError(err)
		}
		token := float64(0)
		token, err = cmd.Flags().GetFloat64("token")
		wg.Add(1)
		err = txn.ExecuteSmartContract(zcncore.FaucetSmartContractAddress, methodName, input, zcncore.ConvertToValue(token))
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
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
				fmt.Println("Execute faucet smart contract success with txn : ", txn.GetTransactionHash())
				return
			}
		}
		ExitWithError("\nExecute faucet smart contract failed. " + statusBar.errMsg + "\n")
	},
}

func init() {
	rootCmd.AddCommand(faucetcmd)
	faucetcmd.PersistentFlags().String("methodName", "", "methodName")
	faucetcmd.PersistentFlags().String("input", "", "input")
	faucetcmd.PersistentFlags().Float64("token", 0, "Token request")
	faucetcmd.MarkFlagRequired("methodName")
	faucetcmd.MarkFlagRequired("input")
}
