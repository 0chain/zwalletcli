package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var minerScPayReward = &cobra.Command{
	Use:   "collect-reward",
	Short: "Pay accrued rewards for a stake pool.",
	Long:  "Pay accrued rewards for a stake pool.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()

		var providerId string
		var err error
		var hash string

		if flags.Changed("provider_id") {
			providerId, err = flags.GetString("provider_id")
			if err != nil {
				log.Fatal(err)
			}
		}

		providerName, err := flags.GetString("provider_type")
		if err != nil {
			log.Fatal(err)
		}

		switch providerName {
		case "miner":
			hash, _, _, _, err = zcncore.MinerSCCollectReward(providerId, zcncore.ProviderMiner)
		case "sharder":
			hash, _, _, _, err = zcncore.MinerSCCollectReward(providerId, zcncore.ProviderSharder)
		case "authorizer":
			hash, _, _, _, err = zcncore.ZCNSCCollectReward(providerId, zcncore.ProviderAuthorizer)
		default:
			log.Fatal("unknown provider type")
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("locked with:", hash)
	},
}

func init() {
	rootCmd.AddCommand(minerScPayReward)

	minerScPayReward.PersistentFlags().String("provider_id", "", "miner or sharder id")
	minerScPayReward.MarkFlagRequired("provider_id")
	minerScPayReward.PersistentFlags().String("provider_type", "miner", "provider type, miner or sharder or authorizer")
	minerScPayReward.MarkFlagRequired("provider_type")
}
