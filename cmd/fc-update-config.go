package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

var updateFaucetCmd = &cobra.Command{
	Use:   "fc-update-config",
	Short: "Update the Faucet smart contract",
	Long:  `Update the Faucet smart contract.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags     = cmd.Flags()
			err       error
			conf      = new(zcncore.InputMap)
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		conf.Fields = make(map[string]interface{})
		if flags.Changed("pour_amount") {
			var pourAmount float64
			pourAmount, err = flags.GetFloat64("pour_amount")
			if err != nil {
				log.Fatal(err)
			}
			conf.Fields["pour_amount"] = common.Balance(zcncore.ConvertToValue(pourAmount))
		}

		if flags.Changed("max_pour_amount") {
			var maxPourAmount float64
			maxPourAmount, err = flags.GetFloat64("max_pour_amount")
			if err != nil {
				log.Fatal(err)
			}
			conf.Fields["max_pour_amount"] = common.Balance(zcncore.ConvertToValue(maxPourAmount))
		}

		if flags.Changed("periodic_limit") {
			var periodicLimit float64
			periodicLimit, err = flags.GetFloat64("periodic_limit")
			if err != nil {
				log.Fatal(err)
			}
			conf.Fields["periodic_limit"] = common.Balance(zcncore.ConvertToValue(periodicLimit))
		}

		if flags.Changed("global_limit") {
			var globalLimit float64
			globalLimit, err = flags.GetFloat64("global_limit")
			if err != nil {
				log.Fatal(err)
			}
			conf.Fields["global_limit"] = common.Balance(zcncore.ConvertToValue(globalLimit))
		}

		if flags.Changed("individual_reset") {
			conf.Fields["individual_reset"], err = flags.GetDuration("individual_reset")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("global_rest") {
			conf.Fields["global_rest"], err = flags.GetDuration("global_rest")
			if err != nil {
				log.Fatal(err)
			}
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.FaucetUpdateConfig(conf); err != nil {
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
	rootCmd.AddCommand(updateFaucetCmd)
	updateFaucetCmd.PersistentFlags().Float64("pour_amount", 0, "pour amount")
	updateFaucetCmd.PersistentFlags().Float64("max_pour_amount", 0, "maximum pour amount")
	updateFaucetCmd.PersistentFlags().Float64("periodic_limit", 0.0, "periodic limit")
	updateFaucetCmd.PersistentFlags().Float64("global_limit", 0.0, "global limit")
	updateFaucetCmd.PersistentFlags().Duration("individual_reset", 0.0, "individual reset duration")
	updateFaucetCmd.PersistentFlags().Duration("global_rest", 0.0, "global reset duration")
}
