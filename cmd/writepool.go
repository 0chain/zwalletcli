package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var getWritePoolStatCmd = &cobra.Command{
	Use:   "getwritelockedtokens",
	Short: "Get locked tokens of write pool",
	Long:  `Get info about locked tokens of write pool`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var flags = cmd.Flags()
		if flags.Changed("allocation_id") == false {
			log.Fatal("error: allocation_id flag is missing")
		}

		var allocID, err = flags.GetString("allocation_id")
		if err != nil {
			log.Fatal("error: invalid allocation id:", err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		if err := zcncore.GetWritePoolStat(statusBar, allocID); err != nil {
			log.Fatal(err)
		}
		wg.Wait()
		if statusBar.success {
			log.Printf("\nWrite pool locked tokens:\n %s\n", statusBar.errMsg)
			return
		}
		log.Fatalf("\nFailed to get locked tokens. %s\n", statusBar.errMsg)
	},
}

var writePoolLockCmd = &cobra.Command{
	Use:   "writelock",
	Short: "Lock tokens in write pool",
	Long:  `Lock tokens in write pool for a duration.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var flags = cmd.Flags()
		if flags.Changed("tokens") == false {
			log.Fatal("error: tokens flag is missing")
		}

		if flags.Changed("allocation_id") == false {
			log.Fatal("error: allocation_id flag is missing")
		}

		var toks, err = flags.GetFloat64("tokens")
		if err != nil {
			log.Fatal("error: invalid number of tokens:", err)
		}

		var allocID string
		if allocID, err = flags.GetString("allocation_id"); err != nil {
			log.Fatal("error: invalid allocation id:", err)
		}

		var fee float64
		if fee, err = flags.GetFloat64("fee"); err != nil {
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
		err = txn.WritePoolLock(allocID, zcncore.ConvertToValue(toks))
		if err != nil {
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
				log.Printf("\nTokens (%f) locked successfully\n", toks)
				return
			}
		}

		log.Fatalf("\nFailed to lock tokens. %s\n", statusBar.errMsg)
	},
}

var writePoolUnlockCmd = &cobra.Command{
	Use:   "writeunlock",
	Short: "Unlock tokens in write pool",
	Long:  `Unlock expired tokens in write pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var flags = cmd.Flags()
		if flags.Changed("allocation_id") == false {
			log.Fatal("error: allocation_id flag is missing")
		}

		var allocID, err = flags.GetString("allocation_id")
		if err != nil {
			log.Fatal("error: invalid allocation id:", err)
		}

		var fee float64
		if fee, err = flags.GetFloat64("fee"); err != nil {
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
		if err = txn.WritePoolUnlock(allocID); err != nil {
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
				log.Printf("\nTokens of %s unlocked successfully\n", allocID)
				return
			}
		}

		log.Fatalf("\nFailed to unlock tokens. %s\n", statusBar.errMsg)
	},
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	rootCmd.AddCommand(getWritePoolStatCmd)
	rootCmd.AddCommand(writePoolLockCmd)
	rootCmd.AddCommand(writePoolUnlockCmd)

	getWritePoolStatCmd.PersistentFlags().String("allocation_id", "",
		"allocation identifier")
	getWritePoolStatCmd.MarkFlagRequired("allocation_id")

	writePoolLockCmd.PersistentFlags().String("allocation_id", "",
		"allocation identifier")
	writePoolLockCmd.PersistentFlags().Float64("tokens", 0,
		"number of tokens to lock")
	writePoolLockCmd.PersistentFlags().Float64("fee", 0, "transaction fee")
	writePoolLockCmd.MarkFlagRequired("tokens")
	writePoolLockCmd.MarkFlagRequired("allocation_id")

	writePoolUnlockCmd.PersistentFlags().String("allocation_id", "",
		"allocation identifier")
	writePoolUnlockCmd.PersistentFlags().Float64("fee", 0, "transaction fee")
	writePoolUnlockCmd.MarkFlagRequired("allocation_id")
}
