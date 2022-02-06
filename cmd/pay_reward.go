package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
	"sync"
)

var minerScPayReward = &cobra.Command{
	Use:   "pay_reward",
	Short: "Pay accrued rewards for a stake pool.",
	Long:  "Pay accrued rewards for a stake pool.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		if !flags.Changed("pool_id") {
			log.Fatal("missing pool id flag")
		}

		pool_id, err := flags.GetString("pool_id")
		if err != nil {
			log.Fatal(err)
		}

		if !flags.Changed("provider_type") {
			log.Fatal("missing tokens flag")
		}

		providerType, err := flags.GetString("provider_type")
		if err != nil {
			log.Fatal(err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		err = txn.MinerSCPayReward(pool_id, providerType)
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

	minerScPayReward.PersistentFlags().String("pool_id", "", "stake pool id")
	minerScPayReward.MarkFlagRequired("pool_id")
	minerScPayReward.PersistentFlags().String("provider_type", "blobber", "provider type")
	minerScPayReward.MarkFlagRequired("provider_type")
}
