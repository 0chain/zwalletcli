package cmd

import (
	"fmt"
	"strconv"
	"strings"
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
		txn, err := zcncore.NewTransaction(statusBar, 0, nonce)
		if err != nil {
			ExitWithError(err)
		}

		token := float64(0)
		token, err = cmd.Flags().GetFloat64("tokens")
		wg.Add(1)
		_, err = txn.ExecuteSmartContract(zcncore.FaucetSmartContractAddress,
			methodName, input, zcncore.ConvertToValue(token), zcncore.WithNoEstimateFee())
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
				switch txn.GetVerifyConfirmationStatus() {
				case zcncore.ChargeableError:
					ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
				case zcncore.Success:
					fmt.Println("Execute faucet smart contract success with txn : ", txn.GetTransactionHash())
				default:
					ExitWithError("\nExecute faucet smart contract failed. Unknown status code: " +
						strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
				}
				return
			}
		}
		ExitWithError("\nExecute faucet smart contract failed. " + statusBar.errMsg + "\n")
	},
}

func init() {
	rootCmd.AddCommand(faucetcmd)
	faucetcmd.PersistentFlags().String("methodName string", "", "methodName")
	faucetcmd.PersistentFlags().String("input string", "", "input")
	faucetcmd.PersistentFlags().Float64("tokens float", 0, "Token request")
	faucetcmd.MarkFlagRequired("methodName")
	faucetcmd.MarkFlagRequired("input")
}
