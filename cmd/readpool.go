package cmd

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var createReadPoolCmd = &cobra.Command{
	Use:   "createreadpool",
	Short: "Create read pool",
	Long:  "Create read pool for blobbers reading if the pool is missing",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
			txn, err  = zcncore.NewTransaction(statusBar, 0)
		)

		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		if err = txn.CreateReadPool(); err != nil {
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
				log.Printf("\nRead pool created successfully\n")
				return
			}
		}

		log.Fatalf("\nFailed to create read pool: %s\n", statusBar.errMsg)
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
		if flags.Changed("t") == false {
			log.Fatal("error: token flag is missing")
		}

		if flags.Changed("d") == false {
			log.Fatal("error: duration is missing")
		}

		var toks, err = flags.GetFloat64("t")
		if err != nil {
			log.Fatal("error: invalid number of tokens:", err)
		}

		var dur time.Duration
		if dur, err = flags.GetDuration("d"); err != nil {
			log.Fatal("error: invalid duration:", err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)

		var txn zcncore.TransactionScheme
		if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
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
		if flags.Changed("p") == false {
			log.Fatal("error: pool id flag is missing")
		}

		var poolID, err = flags.GetString("p")
		if err != nil {
			log.Fatal("error: invalid pool id:", err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)

		var txn zcncore.TransactionScheme
		if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
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

	readPoolLockCmd.PersistentFlags().Float64("t", 0, "number of tokens to lock")
	readPoolLockCmd.PersistentFlags().Duration("d", 0, "duration")
	readPoolLockCmd.MarkFlagRequired("t")
	readPoolLockCmd.MarkFlagRequired("d")

	readPoolUnlockCmd.PersistentFlags().String("p", "", "pool identifier")
	readPoolUnlockCmd.MarkFlagRequired("p")
}
