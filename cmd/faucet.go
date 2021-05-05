package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/0chain/gosdk/core/common"
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
		token, err = cmd.Flags().GetFloat64("tokens")
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

var updatefaucetcmd = &cobra.Command{
	Use:   "update-faucet",
	Short: "Update the Faucet smart contract",
	Long:  `Update the Faucet smart contract.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags      = cmd.Flags()
			err        error
			globalNode = new(zcncore.FaucetSCConfig)
			wg         sync.WaitGroup
			statusBar  = &ZCNStatus{wg: &wg}
		)

		if flags.Changed("pour_amount") {
			var pourAmount float64
			pourAmount, err = flags.GetFloat64("pour_amount")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.PourAmount = common.Balance(zcncore.ConvertToValue(pourAmount))
		}

		if flags.Changed("max_pour_amount") {
			var maxPourAmount float64
			maxPourAmount, err = flags.GetFloat64("max_pour_amount")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MaxPourAmount = common.Balance(zcncore.ConvertToValue(maxPourAmount))
		}

		if flags.Changed("periodic_limit") {
			var periodicLimit float64
			periodicLimit, err = flags.GetFloat64("periodic_limit")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.PeriodicLimit = common.Balance(zcncore.ConvertToValue(periodicLimit))
		}

		if flags.Changed("global_limit") {
			var globalLimit float64
			globalLimit, err = flags.GetFloat64("global_limit")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.GlobalLimit = common.Balance(zcncore.ConvertToValue(globalLimit))
		}

		if flags.Changed("individual_reset") {
			globalNode.IndividualReset, err = flags.GetDuration("individual_reset")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("global_rest") {
			globalNode.GlobalReset, err = flags.GetDuration("global_rest")
			if err != nil {
				log.Fatal(err)
			}
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.FaucetSCSettings(globalNode); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		statusBar.success = false
		wg.Add(1)
		if err = txn.Verify(); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Printf("faucet smart contract settings updated\nHash: %v\n", txn.GetTransactionHash())
	},
}

func init() {
	rootCmd.AddCommand(faucetcmd)
	rootCmd.AddCommand(updatefaucetcmd)

	faucetcmd.PersistentFlags().String("methodName", "", "methodName")
	faucetcmd.PersistentFlags().String("input", "", "input")
	faucetcmd.PersistentFlags().Float64("tokens", 0, "Token request")
	faucetcmd.MarkFlagRequired("methodName")
	faucetcmd.MarkFlagRequired("input")

	updatefaucetcmd.PersistentFlags().Float64("pour_amount", 0, "the amount to pour per request")
	updatefaucetcmd.PersistentFlags().Float64("max_pour_amount", 0, "the max amount poured allower per request")
	updatefaucetcmd.PersistentFlags().Float64("periodic_limit", 0, "the max amount poured to an individual in a time limit")
	updatefaucetcmd.PersistentFlags().Float64("global_limit", 0, "the max the smart contract can pour total in a time limit")
	updatefaucetcmd.PersistentFlags().Duration("individual_reset", 0, "the time to reset the periodic limit for user")
	updatefaucetcmd.PersistentFlags().Duration("global_rest", 0, "the time to reset the global limit for the smart contract")
}
