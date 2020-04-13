package cmd

import (
	"log"
	"os"
	"time"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var getVestingPoolConfigCmd = &cobra.Command{
	Use:   "vp-config",
	Short: "Check out vesting pool configurations.",
	Long:  `Check out vesting pool configurations.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := zcncore.GetVestingSCConfig()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("CONFIG:", conf)
	},
}

var getVestingPoolInfoCmd = &cobra.Command{
	Use:   "vp-info",
	Short: "Check out vesting pool information.",
	Long:  `Check out vesting pool information.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var flags = cmd.Flags()
		if !flags.Changed("pool_id") {
			log.Fatal("missing required 'pool_id' flag")
		}
		poolID, err := flags.GetString("pool_id")
		if err != nil {
			log.Fatalf("can't get 'pool_id' flag: %v", err)
		}
		info, err := zcncore.GetVestingPoolInfo(common.Key(poolID))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("INFO:", info)
	},
}

var getVestingClientPoolsCmd = &cobra.Command{
	Use:   "vp-list",
	Short: "Check out vesting pools list.",
	Long:  `Check out vesting pools list.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags    = cmd.Flags()
			clientID string
			list     []common.Key
			err      error
		)
		if flags.Changed("client_id") {
			if clientID, err = flags.GetString("client_id"); err != nil {
				log.Fatalf("error in 'client_id' flag: %v", err)
			}
		}
		list, err = zcncore.GetVestingClientList(common.Key(clientID))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("LIST:", list)
	},
}

func toKeys(ss []string) (keys []common.Key) {
	keys = make([]common.Key, 0, len(ss))
	for _, s := range ss {
		keys = append(keys, common.Key(s))
	}
	return
}

var vestingPoolUpdateConfigCmd = &cobra.Command{
	Use:   "vp-update-config",
	Short: "Update vesting pool config",
	Long:  "Update vesting pool config by SC owner.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var conf, err = zcncore.GetVestingSCConfig()
		if err != nil {
			log.Fatal(err)
		}
		var (
			flags   = cmd.Flags()
			changed bool
		)
		if flags.Changed("t") {
			var triggers []string
			if triggers, err = flags.GetStringSlice("t"); err != nil {
				log.Fatalf("parsing 't' flag: %v", err)
			}
			if len(triggers) == 0 {
				log.Fatal("empty triggers list")
			}
			conf.Triggers, changed = toKeys(triggers), true
		}
		if flags.Changed("min_lock") {
			var minLock float64
			if minLock, err = flags.GetFloat64("min_lock"); err != nil {
				log.Fatalf("parsing 'min_lock' flag: %v", err)
			}
			if minLock < 0 {
				log.Fatal("negative 'min_lock' value")
			}
			conf.MinLock = common.Balance(zcncore.ConvertToValue(minLock))
			changed = true
		}
		if flags.Changed("min_duration") {
			var minDur time.Duration
			if minDur, err = flags.GetDuration("min_duration"); err != nil {
				log.Fatalf("parsing 'min_duration' flag: %v", err)
			}
			if minDur < 0 {
				log.Fatal("negative 'min_duration'")
			}
			conf.MinDuration, changed = minDur, true
		}
		if flags.Changed("max_duration") {
			var maxDur time.Duration
			if maxDur, err = flags.GetDuration("max_duration"); err != nil {
				log.Fatalf("parsing 'max_duration' flag: %v", err)
			}
			if maxDur < conf.MinDuration {
				log.Fatal("max_duration less then min_duaration")
			}
			conf.MaxDuration, changed = maxDur, true
		}
		if flags.Changed("min_friquency") {
			var minFriq time.Duration
			if minFriq, err = flags.GetDuration("min_friquency"); err != nil {
				log.Fatalf("parsing 'min_friquency' flag: %v", err)
			}
			if minFriq < 0 {
				log.Fatal("negative 'min_friquency'")
			}
			conf.MinFriquency, changed = minFriq, true
		}
		if flags.Changed("max_friquency") {
			var maxFriq time.Duration
			if maxFriq, err = flags.GetDuration("max_friquency"); err != nil {
				log.Fatalf("parsing 'max_friquency' flag: %v", err)
			}
			if maxFriq < conf.MinFriquency {
				log.Fatal("max_friquency less then min_friquency")
			}
			conf.MaxFriquency, changed = maxFriq, true
		}
		if flags.Changed("max_dests") {
			var maxDests int
			if maxDests, err = flags.GetInt("max_dests"); err != nil {
				log.Fatalf("parsing 'max_dests' flag: %v", err)
			}
			if maxDests <= 0 {
				log.Fatal("max_dests is negative or zero")
			}
			conf.MaxDestinations, changed = maxDests, true
		}
		if flags.Changed("max_descr") {
			var maxDescr int
			if maxDescr, err = flags.GetInt("max_descr"); err != nil {
				log.Fatalf("parsing 'max_descr' flag: %v", err)
			}
			if maxDescr < 0 {
				log.Fatal("max_descr is negative")
			}
			conf.MaxDescriptionLength, changed = maxDescr, true
		}
		if !changed {
			log.Fatal("no changes")
		}

		var (
			statusBar = NewZCNStatus()
			txn       zcncore.TransactionScheme
		)
		if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		if err = txn.VestingUpdateConfig(conf); err != nil {
			log.Fatal(err)
		}
		statusBar.Wait()

		if statusBar.success {
			statusBar.success = false

			statusBar.Begin()
			if err = txn.Verify(); err != nil {
				log.Fatal(err)
			}
			statusBar.Wait()

			if statusBar.success {
				log.Println("\nConfigurations updated successfully")
				return
			}
		}

		log.Fatalf("\nFailed to update configurations. %s\n", statusBar.errMsg)
	},
}

var vestingPoolAddCmd = &cobra.Command{
	Use:   "vp-add",
	Short: "Add a vesting pool",
	Long:  "Add a vesting pool.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags = cmd.Flags()
			add   zcncore.VestingAddRequest
			err   error
		)

		// description, optional

		if flags.Changed("description") {
			var description string
			if description, err = flags.GetString("description"); err != nil {
				log.Fatalf("parsing 'description' flag: %v", err)
			}
			add.Description = description
		}

		// start time, optional, default is now

		if flags.Changed("start_time") {
			var startTime int64
			if startTime, err = flags.GetInt64("start_time"); err != nil {
				log.Fatalf("parsing 'start_time' flag: %v", err)
			}
			add.StartTime = common.Timestamp(startTime)
		}

		// duration

		if !flags.Changed("duration") {
			log.Fatal("missing required 'duration' flag")
		}

		var duration time.Duration
		if duration, err = flags.GetDuration("duration"); err != nil {
			log.Fatalf("parsing 'duration' flag: %v", err)
		}
		add.Duration = duration

		// friquency

		if !flags.Changed("friquency") {
			log.Fatal("missing required 'friquency' flag")
		}

		var friquency time.Duration
		if friquency, err = flags.GetDuration("friquency"); err != nil {
			log.Fatalf("parsing 'friquency' flag: %v", err)
		}
		add.Friquency = friquency

		// destinations

		if !flags.Changed("d") {
			log.Fatal("missing required 'd' flag")
		}

		var destinations []string
		destinations, err = flags.GetStringSlice("d")
		if err != nil {
			log.Fatalf("parsing 'destinations' flags: %v", err)
		}
		add.Destinations = toKeys(destinations)

		// amount

		if !flags.Changed("amount") {
			log.Fatal("missing required 'amount' flag")
		}

		var amount float64
		if amount, err = flags.GetFloat64("amount"); err != nil {
			log.Fatalf("parsing 'amount' flag: %v", err)
		}
		add.Amount = common.Balance(zcncore.ConvertToValue(amount))

		// fee
		var fee int64
		if flags.Changed("fee") {
			var feef float64
			if feef, err = flags.GetFloat64("fee"); err != nil {
				log.Fatalf("can't get 'fee' flag: %v", err)
			}
			fee = zcncore.ConvertToValue(feef)
		}

		var lock int64
		if flags.Changed("lock") {
			var lockf float64
			if lockf, err = flags.GetFloat64("lock"); err != nil {
				log.Fatalf("can't get 'lock' flag: %v", err)
			}
			lock = zcncore.ConvertToValue(lockf)
		}

		var (
			statusBar = NewZCNStatus()
			txn       zcncore.TransactionScheme
		)
		if txn, err = zcncore.NewTransaction(statusBar, fee); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		if err = txn.VestingAdd(&add, common.Balance(lock)); err != nil {
			log.Fatal(err)
		}
		statusBar.Wait()

		if statusBar.success {
			statusBar.success = false

			statusBar.Begin()
			if err = txn.Verify(); err != nil {
				log.Fatal(err)
			}
			statusBar.Wait()

			if statusBar.success {
				log.Println("\nVesting pool added successfully.")
				return
			}
		}

		log.Fatalf("\nFailed to add vesting pool: %s\n", statusBar.errMsg)
	},
}

var vestingPoolDeleteCmd = &cobra.Command{
	Use:   "vp-delete",
	Short: "Delete a vesting pool",
	Long:  "Delete a vesting pool.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags  = cmd.Flags()
			poolID string
			err    error
		)

		if !flags.Changed("pool_id") {
			log.Fatal("missing required 'pool_id' flag")
		}

		if poolID, err = flags.GetString("pool_id"); err != nil {
			log.Fatalf("parsing 'pool_id' flag: %v", err)
		}

		var (
			statusBar = NewZCNStatus()
			txn       zcncore.TransactionScheme
		)
		if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		if err = txn.VestingDelete(common.Key(poolID)); err != nil {
			log.Fatal(err)
		}
		statusBar.Wait()

		if statusBar.success {
			statusBar.success = false

			statusBar.Begin()
			if err = txn.Verify(); err != nil {
				log.Fatal(err)
			}
			statusBar.Wait()

			if statusBar.success {
				log.Println("\nVesting pool deleted successfully.")
				return
			}
		}

		log.Fatalf("\nFailed to delete vesting pool: %s\n", statusBar.errMsg)
	},
}

var vestingPoolLockCmd = &cobra.Command{
	Use:   "vp-lock",
	Short: "Lock some tokens for a vesting pool",
	Long:  "Lock some tokens for a vesting pool.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags  = cmd.Flags()
			poolID string
			err    error
		)

		if !flags.Changed("pool_id") {
			log.Fatal("missing required 'pool_id' flag")
		}

		if poolID, err = flags.GetString("pool_id"); err != nil {
			log.Fatalf("parsing 'pool_id' flag: %v", err)
		}

		var lock int64
		if !flags.Changed("lock") {
			log.Fatal("missing required 'lock' flag")
		}
		var lockf float64
		if lockf, err = flags.GetFloat64("lock"); err != nil {
			log.Fatalf("can't get 'lock' flag: %v", err)
		}
		lock = zcncore.ConvertToValue(lockf)

		var fee int64
		if flags.Changed("fee") {
			var feef float64
			if feef, err = flags.GetFloat64("fee"); err != nil {
				log.Fatalf("can't get 'fee' flag: %v", err)
			}
			fee = zcncore.ConvertToValue(feef)
		}

		var (
			statusBar = NewZCNStatus()
			txn       zcncore.TransactionScheme
		)
		if txn, err = zcncore.NewTransaction(statusBar, fee); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		err = txn.VestingLock(common.Key(poolID), common.Balance(lock))
		if err != nil {
			log.Fatal(err)
		}
		statusBar.Wait()

		if statusBar.success {
			statusBar.success = false

			statusBar.Begin()
			if err = txn.Verify(); err != nil {
				log.Fatal(err)
			}
			statusBar.Wait()

			if statusBar.success {
				log.Println("\nTokens locked successfully.")
				return
			}
		}

		log.Fatalf("\nFailed to lock tokens: %s\n", statusBar.errMsg)
	},
}

var vestingPoolUnlockCmd = &cobra.Command{
	Use:   "vp-unlock",
	Short: "Unlock all tokens of a vesting pool",
	Long:  "Unlock all tokens of a vesting pool.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags  = cmd.Flags()
			poolID string
			err    error
		)

		if !flags.Changed("pool_id") {
			log.Fatal("missing required 'pool_id' flag")
		}

		if poolID, err = flags.GetString("pool_id"); err != nil {
			log.Fatalf("parsing 'pool_id' flag: %v", err)
		}

		var (
			statusBar = NewZCNStatus()
			txn       zcncore.TransactionScheme
		)
		if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		err = txn.VestingUnlock(common.Key(poolID))
		if err != nil {
			log.Fatal(err)
		}
		statusBar.Wait()

		if statusBar.success {
			statusBar.success = false

			statusBar.Begin()
			if err = txn.Verify(); err != nil {
				log.Fatal(err)
			}
			statusBar.Wait()

			if statusBar.success {
				log.Println("\nTokens unlocked successfully.")
				return
			}
		}

		log.Fatalf("\nFailed to unlock tokens: %s\n", statusBar.errMsg)
	},
}

var vestingPoolTriggerCmd = &cobra.Command{
	Use:   "vp-trigger",
	Short: "Trigger a vesting pool work.",
	Long: `Developers command that performs vesting of a vesting pool. 
This transaction used by a Vesting server. It can be used by
a configured trigger for development and debugging.`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags  = cmd.Flags()
			poolID string
			err    error
		)

		if !flags.Changed("pool_id") {
			log.Fatal("missing required 'pool_id' flag")
		}

		if poolID, err = flags.GetString("pool_id"); err != nil {
			log.Fatalf("parsing 'pool_id' flag: %v", err)
		}

		var (
			statusBar = NewZCNStatus()
			txn       zcncore.TransactionScheme
		)
		if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		err = txn.VestingTrigger(common.Key(poolID))
		if err != nil {
			log.Fatal(err)
		}
		statusBar.Wait()

		if statusBar.success {
			statusBar.success = false

			statusBar.Begin()
			if err = txn.Verify(); err != nil {
				log.Fatal(err)
			}
			statusBar.Wait()

			if statusBar.success {
				log.Println("\nVesting triggered successfully.")
				return
			}
		}

		log.Fatalf("\nFailed to trigger vesting: %s\n", statusBar.errMsg)
	},
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	rootCmd.AddCommand(getVestingPoolConfigCmd)
	rootCmd.AddCommand(getVestingPoolInfoCmd)
	rootCmd.AddCommand(getVestingClientPoolsCmd)
	rootCmd.AddCommand(vestingPoolUpdateConfigCmd)
	rootCmd.AddCommand(vestingPoolAddCmd)
	rootCmd.AddCommand(vestingPoolDeleteCmd)
	rootCmd.AddCommand(vestingPoolLockCmd)
	rootCmd.AddCommand(vestingPoolUnlockCmd)
	rootCmd.AddCommand(vestingPoolTriggerCmd)

	getVestingPoolInfoCmd.PersistentFlags().String("pool_id", "",
		"pool identifier")
	getVestingClientPoolsCmd.MarkFlagRequired("pool_id")

	getVestingClientPoolsCmd.PersistentFlags().String("client_id", "",
		"client_id, default is current client")

	var addFlags = vestingPoolAddCmd.PersistentFlags()
	addFlags.String("description", "", "pool description, optional")
	addFlags.Int64("start_time", 0, "start_time, Unix seconds, default is now")
	addFlags.Duration("duration", 0, "vesting duration till end, required")
	addFlags.Duration("friquency", 0, "vesting friquency, required")
	addFlags.StringSlice("d", nil, "list of destinations, at list one required")
	addFlags.Float64("amount", 0.0,
		"amount of tokens to vest at once per destination, required")
	addFlags.Float64("fee", 0.0, "transaction fee, optional")
	addFlags.Float64("lock", 0.0, "lock tokens for the pool, optional")

	vestingPoolDeleteCmd.PersistentFlags().String("pool_id", "",
		"pool identifier, required")
	vestingPoolDeleteCmd.MarkFlagRequired("pool_id")

	vestingPoolLockCmd.PersistentFlags().String("pool_id", "",
		"pool identifier, required")
	vestingPoolLockCmd.PersistentFlags().Float64("lock", 0.0,
		"amount of tokens to lock, required")
	vestingPoolLockCmd.PersistentFlags().Float64("fee", 0.0,
		"transaction fee")
	vestingPoolLockCmd.MarkFlagRequired("pool_id")
	vestingPoolLockCmd.MarkFlagRequired("lock")

	vestingPoolUnlockCmd.PersistentFlags().String("pool_id", "",
		"pool identifier, required")

	vestingPoolTriggerCmd.PersistentFlags().String("pool_id", "",
		"pool identifier, required")
	vestingPoolTriggerCmd.MarkFlagRequired("pool_id")
}
