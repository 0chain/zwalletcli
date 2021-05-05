package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var updateinterestpoolcmd = &cobra.Command{
	Use:   "update-interestpool",
	Short: "Update the Interest Pool Smart Contract",
	Long:  `Update the Interest Pool Smart Contract.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags      = cmd.Flags()
			err        error
			globalNode *zcncore.InterestPoolSCConfig
			wg         sync.WaitGroup
			statusBar  = &ZCNStatus{wg: &wg}
		)

		globalNode = new(zcncore.InterestPoolSCConfig)

		if flags.Changed("max_mint") {
			var maxMint float64
			maxMint, err = flags.GetFloat64("max_mint")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MaxMint = common.Balance(zcncore.ConvertToValue(maxMint))
		}

		if flags.Changed("min_lock") {
			var minLock float64
			minLock, err = flags.GetFloat64("min_lock")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MinLock = common.Balance(zcncore.ConvertToValue(minLock))
		}

		if flags.Changed("apr") {
			globalNode.APR, err = flags.GetFloat64("apr")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("min_lock_period") {
			globalNode.MinLockPeriod, err = flags.GetDuration("min_lock_period")
			if err != nil {
				log.Fatal(err)
			}
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.InterestPoolSCSettings(globalNode); err != nil {
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
	rootCmd.AddCommand(updateinterestpoolcmd)

	updateinterestpoolcmd.PersistentFlags().Float64("max_mint", 0, "the max amount allower to be mint")
	updateinterestpoolcmd.PersistentFlags().Float64("min_lock", 0, "the minimum amount to lock")
	updateinterestpoolcmd.PersistentFlags().Float64("apr", 0, "Annual percentage rate")
	updateinterestpoolcmd.PersistentFlags().Duration("min_lock_period", 0, "the minimum amount of time to lock")
}
