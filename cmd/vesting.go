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
		var (
			info = new(zcncore.VestingPoolInfo)
			cb   = NewJSONInfoCB(info)
		)
		err = zcncore.GetVestingPoolInfo(poolID, cb)
		if err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		var earned, pending, vested common.Balance
		for _, d := range info.Destinations {
			pending += d.Wanted - d.Vested
			vested += d.Vested
			earned += d.Earned
		}

		fmt.Println("pool_id:     ", info.ID)
		fmt.Println("balance:     ", info.Balance)
		fmt.Println("can unlock:  ", info.Left, "(excess)")
		fmt.Println("sent:        ", vested, "(real value)")
		fmt.Println("pending:     ", pending, "(not sent, real value)")
		fmt.Println("vested:      ", earned+vested, "(virtual, time based value)")
		fmt.Println("description: ", info.Description)
		fmt.Println("start_time:  ", info.StartTime.ToTime())
		fmt.Println("expire_at:   ", info.ExpireAt.ToTime())
		fmt.Println("destinations:")
		for _, d := range info.Destinations {
			fmt.Println("  - id:         ", d.ID)
			fmt.Println("    vesting:    ", d.Wanted)
			fmt.Println("    can unlock: ", d.Earned, "(virtual, time based value)")
			fmt.Println("    sent:       ", d.Vested, "(real value)")
			fmt.Println("    pending:    ", d.Wanted-d.Vested, "(not sent, real value)")
			fmt.Println("    vested:     ", d.Earned+d.Vested, "(virtual, time based value)")
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
			list     []common.Key
			err      error
		)
		if flags.Changed("client_id") {
			if clientID, err = flags.GetString("client_id"); err != nil {
				log.Fatalf("error in 'client_id' flag: %v", err)
			}
		}
		var (
			vcl zcncore.VestingClientList
			cb  = NewJSONInfoCB(&vcl)
		)
		err = zcncore.GetVestingClientList(clientID, cb)
		if err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}
		list = vcl.Pools
		if len(list) == 0 {
			log.Println("no vesting pools")
			return
		}
		for _, pool := range list {
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
			ID:     id,
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

		var lock uint64
		if flags.Changed("lock") {
			var lockf float64
			if lockf, err = flags.GetFloat64("lock"); err != nil {
				log.Fatalf("can't get 'lock' flag: %v", err)
			}
			lock = zcncore.ConvertToValue(lockf)
		} else {
			log.Fatal("missing required 'lock' flag")
		}

		var (
			statusBar = NewZCNStatus()
			txn       zcncore.TransactionScheme
		)
		if txn, err = zcncore.NewTransaction(statusBar, getTxnFee(), nonce); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		if err = txn.VestingAdd(&add, lock); err != nil {
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
				switch txn.GetVerifyConfirmationStatus() {
				case zcncore.ChargeableError:
					ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
				case zcncore.Success:
					log.Printf("\nVesting pool added successfully:%v:vestingpool:%v\nHash: %v",
						zcncore.VestingSmartContractAddress, txn.GetTransactionHash(), txn.GetTransactionHash())
				default:
					ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
						strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
				}
				return
			} else {
				log.Fatalf("\nFailed to add vesting pool: %s\n", statusBar.errMsg)
			}
		}
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
		if txn, err = zcncore.NewTransaction(statusBar, getTxnFee(), nonce); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		if err = txn.VestingDelete(poolID); err != nil {
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
				switch txn.GetVerifyConfirmationStatus() {
				case zcncore.ChargeableError:
					ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
				case zcncore.Success:
					log.Printf("\nVesting pool deleted successfully.\nHash: %v", txn.GetTransactionHash())
				default:
					ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
						strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
				}
				return
			} else {
				log.Fatalf("\nFailed to delete vesting pool: %s\n", statusBar.errMsg)
			}
		}
	},
}

var vestingPoolStopCmd = &cobra.Command{
	Use:   "vp-stop",
	Short: "Stop vesting for one of destinations and unlock tokens not vested",
	Long:  "Stop vesting for one of destinations and unlock tokens not vested",
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

		var dest string
		if !flags.Changed("d") {
			log.Fatal("missing required 'd' flag")
		}

		if dest, err = flags.GetString("d"); err != nil {
			log.Fatalf("parsing 'd' flag: %v", err)
		}

		var (
			statusBar = NewZCNStatus()
			txn       zcncore.TransactionScheme
		)
		if txn, err = zcncore.NewTransaction(statusBar, getTxnFee(), nonce); err != nil {
			log.Fatal(err)
		}

		var sr zcncore.VestingStopRequest
		sr.PoolID = poolID
		sr.Destination = dest

		statusBar.Begin()
		if err = txn.VestingStop(&sr); err != nil {
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
				switch txn.GetVerifyConfirmationStatus() {
				case zcncore.ChargeableError:
					ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
				case zcncore.Success:
					log.Printf("\nStop vesting for %s.\nHash: %v", dest, txn.GetTransactionHash())
				default:
					ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
						strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
				}
			} else {
				log.Fatalf("\nFailed to stop vesting: %s\n", statusBar.errMsg)
			}
		}
	},
}

var vestingPoolUnlockCmd = &cobra.Command{
	Use:   "vp-unlock",
	Short: "Unlock tokens of a vesting pool",
	Long:  "Unlock tokens of a vesting pool.",
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
		if txn, err = zcncore.NewTransaction(statusBar, getTxnFee(), nonce); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		err = txn.VestingUnlock(poolID)
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
				switch txn.GetVerifyConfirmationStatus() {
				case zcncore.ChargeableError:
					ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
				case zcncore.Success:
					log.Printf("\nTokens unlocked successfully.\nHash: %v", txn.GetTransactionHash())
				default:
					log.Fatalf("\nFailed to unlock tokens: %s\n", statusBar.errMsg)

				}
				return
			} else {
				fmt.Printf("Pour request failed\n")
			}
		}
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
		if txn, err = zcncore.NewTransaction(statusBar, getTxnFee(), nonce); err != nil {
			log.Fatal(err)
		}

		statusBar.Begin()
		err = txn.VestingTrigger(poolID)
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
				switch txn.GetVerifyConfirmationStatus() {
				case zcncore.ChargeableError:
					ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
				case zcncore.Success:
					log.Printf("\nVesting triggered successfully.\nHash: %v", txn.GetTransactionHash())
				default:
					ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
						strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
				}
			} else {
				log.Fatalf("\nFailed to trigger vesting: %s\n", statusBar.errMsg)
			}
		}
	},
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	rootCmd.AddCommand(getVestingPoolInfoCmd)
	rootCmd.AddCommand(getVestingClientPoolsCmd)
	rootCmd.AddCommand(vestingPoolAddCmd)
	rootCmd.AddCommand(vestingPoolDeleteCmd)
	rootCmd.AddCommand(vestingPoolStopCmd)
	rootCmd.AddCommand(vestingPoolUnlockCmd)
	rootCmd.AddCommand(vestingPoolTriggerCmd)

	getVestingPoolInfoCmd.PersistentFlags().String("pool_id", "",
		"pool identifier")
	getVestingPoolInfoCmd.MarkFlagRequired("pool_id")

	getVestingClientPoolsCmd.PersistentFlags().String("client_id", "",
		"client_id, default is current client")
	getVestingClientPoolsCmd.MarkFlagRequired("client_id")

	var addFlags = vestingPoolAddCmd.PersistentFlags()
	addFlags.String("description", "", "pool description, optional")
	addFlags.Int64("start_time", 0, "start_time, Unix seconds, default is now")
	addFlags.Duration("duration", 0, "vesting duration till end, required")
	addFlags.StringSlice("d", nil, `list of colon separated 'destination:amount' values,
use -d flag many times to provide few destinations, for example 'dst:1.2'`)
	addFlags.Float64("fee", 0.0, "transaction fee, optional")
	addFlags.Float64("lock", 0.0, "lock tokens for the pool")

	vestingPoolDeleteCmd.PersistentFlags().String("pool_id", "",
		"pool identifier, required")
	vestingPoolDeleteCmd.MarkFlagRequired("pool_id")

	vestingPoolStopCmd.PersistentFlags().String("pool_id", "",
		"pool identifier, required")
	vestingPoolStopCmd.PersistentFlags().String("d", "",
		"destination to stop vesting, required")
	vestingPoolStopCmd.MarkFlagRequired("pool_id")
	vestingPoolStopCmd.MarkFlagRequired("d")

	vestingPoolUnlockCmd.PersistentFlags().String("pool_id", "",
		"pool identifier, required")

	vestingPoolTriggerCmd.PersistentFlags().String("pool_id", "",
		"pool identifier, required")
	vestingPoolTriggerCmd.MarkFlagRequired("pool_id")
}
