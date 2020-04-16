package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

		fmt.Println("min_lock:", conf.MinLock)
		fmt.Println("min_duration:", conf.MinDuration)
		fmt.Println("max_duration:", conf.MaxDuration)
		fmt.Println("max_destinations:", conf.MaxDestinations)
		fmt.Println("max_description_length:", conf.MaxDescriptionLength)
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

		fmt.Println("pool_id:     ", info.ID)
		fmt.Println("balance:     ", info.Balance)
		fmt.Println("can unlock:  ", info.Left)
		fmt.Println("description: ", info.Description)
		fmt.Println("start_time:  ", info.StartTime.ToTime())
		fmt.Println("expire_at:   ", info.ExpireAt.ToTime())
		fmt.Println("destinations:")
		for _, d := range info.Destinations {
			fmt.Println("  - id:         ", d.ID)
			fmt.Println("    vesting:    ", d.Wanted)
			fmt.Println("    can unlock: ", d.Earned)
			fmt.Println("    last unlock:", d.Last.ToTime())
		}
		fmt.Println("client_id:   ", info.ClientID)
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
			list     *zcncore.VestingClientList
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
		if len(list.Pools) == 0 {
			log.Println("no vesting pools")
			return
		}
		for _, pool := range list.Pools {
			log.Println("- ", pool)
		}
	},
}

func toKeys(ss []string) (keys []common.Key) {
	keys = make([]common.Key, 0, len(ss))
	for _, s := range ss {
		keys = append(keys, common.Key(s))
	}
	return
}

func vestingDests(dd []string) (vds []*zcncore.VestingDest, err error) {
	vds = make([]*zcncore.VestingDest, 0, len(dd))
	for _, d := range dd {
		var ss = strings.Split(d, ":")
		if len(ss) != 2 {
			return nil, fmt.Errorf("invalid destination: %q", d)
		}
		var id, amounts = ss[0], ss[1]
		if len(id) != 64 {
			return nil, fmt.Errorf("invalid destination id: %q", id)
		}
		var amount float64
		if amount, err = strconv.ParseFloat(amounts, 64); err != nil {
			return nil, fmt.Errorf("invalid destination amount %q: %v",
				amounts, err)
		}
		if amount < 0 {
			return nil, fmt.Errorf("negative amount: %f", amount)
		}
		vds = append(vds, &zcncore.VestingDest{
			ID:     common.Key(id),
			Amount: common.Balance(zcncore.ConvertToValue(amount)),
		})
	}
	return
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

		// destinations

		if !flags.Changed("d") {
			log.Fatal("missing required 'd' flag")
		}

		var destinations []string
		destinations, err = flags.GetStringSlice("d")
		if err != nil {
			log.Fatalf("parsing 'destinations' flags: %v", err)
		}
		if add.Destinations, err = vestingDests(destinations); err != nil {
			log.Fatalf("parsing destinations: %v", err)
		}

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
	Long:  `Move current vested tokens to destinations`,
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
	addFlags.StringSlice("d", nil, `list of colon separated 'destination:amount' values,
use -d flag many times to provide few destinations, for example 'dst:1.2'`)
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
