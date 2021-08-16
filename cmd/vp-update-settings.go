package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

var updateVestingPoolConfigCmd = &cobra.Command{
	Use:   "vp-update-config",
	Short: "Update the vesting pool configurations.",
	Long:  `Update the vesting pool configurations.`,
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
			var minLock float64
			if minLock, err = flags.GetFloat64("min_lock"); err != nil {
				log.Fatal(err)
			}
			conf.Fields["min_lock"] = common.Balance(zcncore.ConvertToValue(minLock))
		}
		if flags.Changed("max_destinations") {
			if conf.Fields["max_destinations"], err = flags.GetInt("max_destinations"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_description_length") {
			if conf.Fields["max_description_length"], err = flags.GetInt("max_description_length"); err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("min_duration") {
			if conf.Fields["min_duration"], err = flags.GetDuration("min_duration"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_duration") {
			if conf.Fields["max_duration"], err = flags.GetDuration("max_duration"); err != nil {
				log.Fatal(err)
			}
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.VestingUpdateConfig(conf); err != nil {
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

		fmt.Println("vesting smart contract settings updated")
	},
}

func init() {
	rootCmd.AddCommand(updateVestingPoolConfigCmd)
	updateVestingPoolConfigCmd.PersistentFlags().Int("max_description_length", 0, "max length for descriptions")
	updateVestingPoolConfigCmd.PersistentFlags().Int("max_destinations", 0, "max destinations allowed")
	updateVestingPoolConfigCmd.PersistentFlags().Float64("min_lock", 0.0, "minimum lock for vesting")
	updateVestingPoolConfigCmd.PersistentFlags().Duration("min_duration", 0.0, "minimum duration for vesting")
	updateVestingPoolConfigCmd.PersistentFlags().Duration("max_duration", 0.0, "max duration for vesting")
}
