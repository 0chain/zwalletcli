package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var getStakePoolStatCmd = &cobra.Command{
	Use:   "getstakelockedtokens",
	Short: "Get locked tokens of stake pool",
	Long:  `Get info about locked tokens of stake pool`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var flags = cmd.Flags()
		if flags.Changed("blobber_id") == false {
			log.Fatal("error: blobber_id flag is missing")
		}

		var blobberID, err = flags.GetString("blobber_id")
		if err != nil {
			log.Fatal("error: invalid allocation id:", err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		if err := zcncore.GetStakePoolStat(statusBar, blobberID); err != nil {
			log.Fatal(err)
		}
		wg.Wait()
		if statusBar.success {
			log.Printf("\nStake pool locked tokens:\n %s\n", statusBar.errMsg)
			return
		}
		log.Fatalf("\nFailed to get locked tokens. %s\n", statusBar.errMsg)
	},
}

var stakePoolUnlockCmd = &cobra.Command{
	Use:   "stakeunlock",
	Short: "Unlock tokens in stake pool",
	Long:  `Unlock expired tokens in stake pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags    = cmd.Flags()
			fee, err = flags.GetFloat64("fee")
		)
		if err != nil {
			log.Fatal("error: invalid fee value:", err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)

		var txn zcncore.TransactionScheme
		txn, err = zcncore.NewTransaction(statusBar,
			zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		if err = txn.StakePoolUnlock(); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if statusBar.success {
			statusBar.success = false

			wg.Add(1)
			if err = txn.Verify(); err != nil {
				log.Fatal(err)
			}
			wg.Wait()

			if statusBar.success {
				log.Printf("\nTokens of stake pool unlocked successfully\n")
				return
			}
		}

		log.Fatalf("\nFailed to unlock tokens. %s\n", statusBar.errMsg)
	},
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	rootCmd.AddCommand(getStakePoolStatCmd)
	rootCmd.AddCommand(stakePoolUnlockCmd)

	getStakePoolStatCmd.PersistentFlags().String("blobber_id", "",
		"allocation identifier")
	getStakePoolStatCmd.MarkFlagRequired("blobber_id")

	stakePoolUnlockCmd.PersistentFlags().Float64("fee", 0, "transaction fee")
}
