package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

var updateInterestPoolConfigCmd = &cobra.Command{
	Use:   "ip-update-config",
	Short: "Update the interest pool configurations.",
	Long:  `Update the interest pool configurations.`,
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
		if flags.Changed("min_lock") {
			if minLock, err := flags.GetFloat64("min_lock"); err != nil {
				log.Fatal(err)
			} else {
				conf.Fields["min_lock"] = common.Balance(zcncore.ConvertToValue(minLock))
			}

		}
		if flags.Changed("apr") {
			if conf.Fields["apr"], err = flags.GetFloat64("apr"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_mint") {
			if maxMint, err := flags.GetFloat64("max_mint"); err != nil {
				log.Fatal(err)
			} else {
				conf.Fields["max_mint"] = common.Balance(zcncore.ConvertToValue(maxMint))
			}
		}

		if flags.Changed("min_lock_period") {
			if conf.Fields["min_lock_period"], err = flags.GetDuration("min_lock_period"); err != nil {
				log.Fatal(err)
			}
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.InterestPoolUpdateConfig(conf); err != nil {
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

		fmt.Println("interest pool smart contract settings updated")
	},
}

func init() {
	rootCmd.AddCommand(updateInterestPoolConfigCmd)
	updateInterestPoolConfigCmd.PersistentFlags().Float64("min_lock", 0, "minimum tokens that can be locked")
	updateInterestPoolConfigCmd.PersistentFlags().Float64("apr", 0, "apr, apr")
	updateInterestPoolConfigCmd.PersistentFlags().Duration("min_lock_period", 0.0, "minimum lock period")
	updateInterestPoolConfigCmd.PersistentFlags().Float64("max_mint", 0.0, "minimum tokes interest sc can mint")
}
