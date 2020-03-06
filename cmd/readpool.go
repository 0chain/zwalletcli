package cmd

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

func createReadPool() (err error) {
	var (
		txn       zcncore.TransactionScheme
		wg        sync.WaitGroup
		statusBar = &ZCNStatus{wg: &wg}
	)

	if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
		return
	}

	wg.Add(1)
	if err = txn.CreateReadPool(); err != nil {
		return
	}
	wg.Wait()

	if statusBar.success {
		statusBar.success = false

		wg.Add(1)
		if err = txn.Verify(); err != nil {
			return
		}
		wg.Wait()

		if statusBar.success {
			return // nil
		}
	}

	return errors.New(statusBar.errMsg)
}

var createReadPoolCmd = &cobra.Command{
	Use:   "createreadpool",
	Short: "Create read pool",
	Long:  "Create read pool for blobbers reading if the pool is missing",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var err = createReadPool()
		if err != nil {
			log.Fatalf("\nFailed to create read pool: %v\n", err)
		}
		log.Printf("\nRead pool created successfully\n")
	},
}

var getReadPoolsStatsCmd = &cobra.Command{
	Use:   "getreadlockedtokens",
	Short: "Get locked tokens of read pool",
	Long:  `Get info about locked tokens of read pool`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		if err := zcncore.GetReadPoolsStats(statusBar); err != nil {
			log.Fatal(err)
		}
		wg.Wait()
		if statusBar.success {
			log.Printf("\nRead pool locked tokens:\n %s\n", statusBar.errMsg)
			return
		}
		log.Fatalf("\nFailed to get locked tokens. %s\n", statusBar.errMsg)
	},
}

var readPoolLockCmd = &cobra.Command{
	Use:   "readlock",
	Short: "Lock tokens in read pool",
	Long:  `Lock tokens in read pool for a duration.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var flags = cmd.Flags()
		if flags.Changed("tokens") == false {
			log.Fatal("error: tokens flag is missing")
		}

		if flags.Changed("duration") == false {
			log.Fatal("error: duration is missing")
		}

		var toks, err = flags.GetFloat64("tokens")
		if err != nil {
			log.Fatal("error: invalid number of tokens:", err)
		}

		var dur time.Duration
		if dur, err = flags.GetDuration("duration"); err != nil {
			log.Fatal("error: invalid duration:", err)
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
		err = txn.ReadPoolLock(zcncore.ConvertToValue(toks), dur)
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

var readPoolUnlockCmd = &cobra.Command{
	Use:   "readunlock",
	Short: "Unlock tokens in read pool",
	Long:  `Unlock expired tokens in read pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var flags = cmd.Flags()
		if flags.Changed("pool_id") == false {
			log.Fatal("error: pool id flag is missing")
		}

		var poolID, err = flags.GetString("pool_id")
		if err != nil {
			log.Fatal("error: invalid pool id:", err)
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
		if err = txn.ReadPoolUnlock(poolID); err != nil {
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
				log.Printf("\nTokens of %s unlocked successfully\n", poolID)
				return
			}
		}

		log.Fatalf("\nFailed to unlock tokens. %s\n", statusBar.errMsg)
	},
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	rootCmd.AddCommand(createReadPoolCmd)
	rootCmd.AddCommand(getReadPoolsStatsCmd)
	rootCmd.AddCommand(readPoolLockCmd)
	rootCmd.AddCommand(readPoolUnlockCmd)

	readPoolLockCmd.PersistentFlags().Float64("tokens", 0, "number of tokens to lock")
	readPoolLockCmd.PersistentFlags().Duration("duration", 0, "duration")
	readPoolLockCmd.PersistentFlags().Float64("fee", 0, "transaction fee")
	readPoolLockCmd.MarkFlagRequired("tokens")
	readPoolLockCmd.MarkFlagRequired("duration")

	readPoolUnlockCmd.PersistentFlags().String("pool_id", "", "pool identifier")
	readPoolUnlockCmd.PersistentFlags().Float64("fee", 0, "transaction fee")
	readPoolUnlockCmd.MarkFlagRequired("pool_id")
}
