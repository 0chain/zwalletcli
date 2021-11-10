package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

var minerscUpdateNodeSettings = &cobra.Command{
	Use:   "mn-update-node-settings",
	Short: "Change miner/sharder settings in Miner SC.",
	Long:  "Change miner/sharder settings in Miner SC by delegate wallet.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags   = cmd.Flags()
			id      string
			sharder bool
			err     error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if sharder, err = flags.GetBool("sharder"); err != nil {
			log.Fatal(err)
		}

		var (
			miner     *zcncore.MinerSCMinerInfo
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		if err = zcncore.GetMinerSCNodeInfo(id, statusBar); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		miner = new(zcncore.MinerSCMinerInfo)
		err = json.Unmarshal([]byte(statusBar.errMsg), miner)
		if err != nil {
			log.Fatal(err)
		}

		// remove not settings fields
		miner = &zcncore.MinerSCMinerInfo{SimpleMinerSCMinerInfo: &zcncore.SimpleMinerSCMinerInfo{
			NumberOfDelegates: miner.NumberOfDelegates,
			MinStake:          miner.MinStake,
			MaxStake:          miner.MaxStake,
			ID:                id,
		},
		}

		if flags.Changed("num_delegates") {
			miner.NumberOfDelegates, err = flags.GetInt("num_delegates")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("min_stake") {
			var min float64
			if min, err = flags.GetFloat64("min_stake"); err != nil {
				log.Fatal(err)
			}
			miner.MinStake = common.Balance(zcncore.ConvertToValue(min))
		}

		if flags.Changed("max_stake") {
			var max float64
			if max, err = flags.GetFloat64("max_stake"); err != nil {
				log.Fatal(err)
			}
			miner.MaxStake = common.Balance(zcncore.ConvertToValue(max))
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if sharder {
			err = txn.MinerSCSharderSettings(miner)
		} else {
			err = txn.MinerSCMinerSettings(miner)
		}
		if err != nil {
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

		fmt.Println("settings updated")
	},
}

func init() {
	rootCmd.AddCommand(minerscUpdateNodeSettings)
	minerscUpdateNodeSettings.PersistentFlags().String("id", "", "miner/sharder ID to update")
	minerscUpdateNodeSettings.PersistentFlags().Bool("sharder", false, "set true for sharder node")
	minerscUpdateNodeSettings.PersistentFlags().Int("num_delegates", 0, "max number of delegate pools")
	minerscUpdateNodeSettings.PersistentFlags().Float64("min_stake", 0.0, "min stake allowed")
	minerscUpdateNodeSettings.PersistentFlags().Float64("max_stake", 0.0, "max stake allowed")
	minerscUpdateNodeSettings.MarkFlagRequired("id")
}