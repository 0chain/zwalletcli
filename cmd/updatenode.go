package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var minerscUpdateNodeSettings = &cobra.Command{
	Use:   "mn-update-settings",
	Short: "Change miner/sharder settings in Miner SC.",
	Long:  "Change miner/sharder settings in Miner SC by delegate wallet.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags   = cmd.Flags()
			id      string
			sharder bool
			err     error
			hash    string
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

		miner := &zcncore.MinerSCMinerInfo{
			SimpleMiner: zcncore.SimpleMiner{
				ID: id,
			},
		}

		if flags.Changed("num_delegates") {
			numDelegates, err := flags.GetInt("num_delegates")
			if err != nil {
				log.Fatal(err)
			}
			miner.Settings.NumDelegates = &numDelegates
		}

		if flags.Changed("service_charge") {
			serviceCharge, err := flags.GetFloat64("service_charge")
			if err != nil {
				log.Fatal(err)
			}
			miner.Settings.ServiceCharge = &serviceCharge
		}

		if sharder {
			hash, _, _, _, err = zcncore.MinerSCSharderSettings(miner)
		} else {
			hash, _, _, _, err = zcncore.MinerSCMinerSettings(miner)
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("settings updated\nHash: %v", hash)
	},
}

func init() {
	rootCmd.AddCommand(minerscUpdateNodeSettings)
	minerscUpdateNodeSettings.PersistentFlags().String("id", "", "miner/sharder ID to update")
	minerscUpdateNodeSettings.PersistentFlags().Bool("sharder", false, "set true for sharder node")
	minerscUpdateNodeSettings.PersistentFlags().Int("num_delegates", 0, "max number of delegate pools")
	minerscUpdateNodeSettings.PersistentFlags().Float64("service_charge", 0, "service charge")
	minerscUpdateNodeSettings.MarkFlagRequired("id")
}
