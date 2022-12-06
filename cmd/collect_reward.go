package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
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

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		txn, err := zcncore.NewTransaction(statusBar, getTxnFee(), 0)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		switch providerName {
		case "miner":
			err = txn.MinerSCCollectReward(providerId, zcncore.ProviderMiner)
		case "sharder":
			err = txn.MinerSCCollectReward(providerId, zcncore.ProviderSharder)
		case "authorizer":
			log.Fatal("not implemented yet")
		default:
			log.Fatal("unknown provider type")
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

		if statusBar.success {
			switch txn.GetVerifyConfirmationStatus() {
			case zcncore.ChargeableError:
				ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
			case zcncore.Success:
				fmt.Println("locked with:", txn.GetTransactionHash())
			default:
				ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
			return
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}
	},
}

func init() {
	rootCmd.AddCommand(minerScPayReward)

	minerScPayReward.PersistentFlags().String("provider_id", "", "miner or sharder id")
	minerScPayReward.MarkFlagRequired("provider_id")
	minerScPayReward.PersistentFlags().String("provider_type", "miner", "provider type, miner or sharder")
	minerScPayReward.MarkFlagRequired("provider_type")
}
