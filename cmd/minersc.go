package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

var minerscUpdateSettings = &cobra.Command{
	Use:   "mn-update-settings",
	Short: "Change miner/sharder settings in Miner SC.",
	Long:  "Change miner/sharder settings in Miner SC by delegate wallet.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags = cmd.Flags()
			id    string
			err   error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		var (
			miner     *zcncore.MinerSCMinerInfo
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		if err = zcncore.GetMinerSCNodeInfo(id, statusBar); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		miner = new(zcncore.MinerSCMinerInfo)
		err = json.Unmarshal([]byte(statusBar.errMsg), miner)
		if err != nil {
			log.Fatal(err)
		}

		// remove not settings fields
		miner = &zcncore.MinerSCMinerInfo{SimpleMinerSCMinerInfo: &zcncore.SimpleMinerSCMinerInfo{
			NumberOfDelegates: miner.NumberOfDelegates,
			MinStake:          miner.MinStake,
			MaxStake:          miner.MaxStake,
			ID:                id,
		},
		}

		if flags.Changed("num_delegates") {
			miner.NumberOfDelegates, err = flags.GetInt("num_delegates")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("min_stake") {
			var min float64
			if min, err = flags.GetFloat64("min_stake"); err != nil {
				log.Fatal(err)
			}
			miner.MinStake = common.Balance(zcncore.ConvertToValue(min))
		}

		if flags.Changed("max_stake") {
			var max float64
			if max, err = flags.GetFloat64("max_stake"); err != nil {
				log.Fatal(err)
			}
			miner.MaxStake = common.Balance(zcncore.ConvertToValue(max))
		}

		txn, err := zcncore.NewTransaction(statusBar, 0, nonce)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.MinerSCMinerSettings(miner); err != nil {
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
				fmt.Printf("settings updated\nHash: %v", txn.GetTransactionHash())
			default:
				ExitWithError("\nExecute settings update update smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
			return
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}
	},
}

var minerscInfo = &cobra.Command{
	Use:   "mn-info",
	Short: "Get miner/sharder info from Miner SC.",
	Long:  "Get miner/sharder info from Miner SC.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags = cmd.Flags()
			id    string
			err   error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		if err = zcncore.GetMinerSCNodeInfo(id, statusBar); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Println(statusBar.errMsg)
	},
}

var minerscMiners = &cobra.Command{
	Use:   "ls-miners",
	Short: "Get list of all active miners fro Miner SC",
	Long:  "Get list of all active miners from Miner SC",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			err   error
			info  = new(zcncore.MinerSCNodes)
			cb    = NewJSONInfoCB(info)
		)

		if err = zcncore.GetMiners(cb); err != nil {
			log.Fatal(err)
		}

		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		if flags.Changed("json") {
			var j bool
			if j, err = flags.GetBool("json"); err != nil {
				log.Fatal(err)
			}
			if j {
				util.PrintJSON(info)
				return
			}
		}

		if len(info.Nodes) == 0 {
			fmt.Println("no miners in Miner SC")
			return
		}

		for _, node := range info.Nodes {
			fmt.Println("- ID:        ", node.Miner.ID)
			fmt.Println("- Host:      ", node.Miner.Host)
			fmt.Println("- Port:      ", node.Miner.Port)
		}
	},
}

var minerscSharders = &cobra.Command{
	Use:   "ls-sharders",
	Short: "Get list of all active sharders fro Miner SC",
	Long:  "Get list of all active sharders from Miner SC",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		flags := cmd.Flags()

		var err error
		var jsonFlag, allFlag bool

		if flags.Changed("json") {
			jsonFlag, err = flags.GetBool("json")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("all") {
			allFlag, err = flags.GetBool("all")
			if err != nil {
				log.Fatal(err)
			}
		}

		mb, err := zcncore.GetLatestFinalizedMagicBlock(context.Background(), 1)
		if err != nil {
			log.Fatalf("Failed to get MagicBlock: %v", err)
		}

		if mb != nil && mb.Sharders != nil {
			fmt.Println("MagicBlock Sharders")
			if jsonFlag {
				util.PrettyPrintJSON(mb.Sharders.Nodes)
			} else {
				for _, node := range mb.Sharders.Nodes {
					fmt.Println("ID:", node.ID)
					fmt.Println("  - N2NHost:", node.N2NHost)
					fmt.Println("  - Host:", node.Host)
					fmt.Println("  - Port:", node.Port)
				}
			}
			fmt.Println()
		}

		if allFlag {
			sharders := new(zcncore.MinerSCNodes)
			callback := NewJSONInfoCB(sharders)
			if err = zcncore.GetSharders(callback); err != nil {
				log.Fatalf("Failed to get registered sharders: %v", err)
			}
			if err = callback.Waiting(); err != nil {
				log.Fatalf("Failed to get registered sharders: %v", err)
			}
			fmt.Println("Registered Sharders")
			if jsonFlag {
				util.PrettyPrintJSON(sharders.Nodes)
			} else {
				for _, node := range sharders.Nodes {
					fmt.Println("ID:", node.Miner.ID)
					fmt.Println("  - N2NHost:", node.Miner.N2NHost)
					fmt.Println("  - Host:", node.Miner.Host)
					fmt.Println("  - Port:", node.Miner.Port)
				}
			}
		}
	},
}

var minerscUserInfo = &cobra.Command{
	Use:   "mn-user-info",
	Short: "Get miner/sharder user pools info from Miner SC.",
	Long:  "Get miner/sharder user pools info from Miner SC.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags    = cmd.Flags()
			clientID string

			err error
		)

		if flags.Changed("client_id") {
			if clientID, err = flags.GetString("client_id"); err != nil {
				log.Fatal(err)
			}
		}

		var (
			info = new(zcncore.MinerSCUserPoolsInfo)
			cb   = NewJSONInfoCB(info)
		)
		if err = zcncore.GetMinerSCUserInfo(clientID, cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		if flags.Changed("json") {
			var j bool
			if j, err = flags.GetBool("json"); err != nil {
				log.Fatal(err)
			}
			if j {
				util.PrintJSON(info)
				return
			}
		}

		if len(info.Pools) == 0 {
			fmt.Println("no user pools in Miner SC")
			return
		}

		var total common.Balance
		for _, nodes := range info.Pools {
			for _, pools := range nodes {
				for _, pool := range pools {
					total += pool.Balance
				}
			}
		}

		for key, nodes := range info.Pools {
			for nit, pools := range nodes {
				fmt.Println("- node:", nit+" ("+key+")")
				for _, pool := range pools {
					fmt.Println("  - pool_id:       ", pool.ID)
					fmt.Println("    balance:       ", pool.Balance)
					fmt.Println("    interests paid:", pool.InterestPaid)
					fmt.Println("    rewards paid:  ", pool.RewardPaid)
					fmt.Println("    status:        ", strings.ToLower(pool.Status))
					fmt.Println("    stake %:       ",
						float64(pool.Balance)/float64(total)*100.0, "%")
				}
			}
		}
	},
}

var minerscPoolInfo = &cobra.Command{
	Use:   "mn-pool-info",
	Short: "Get miner/sharder pool info from Miner SC.",
	Long:  "Get miner/sharder pool info from Miner SC.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags  = cmd.Flags()
			id     string
			poolID string

			err error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if !flags.Changed("pool_id") {
			log.Fatal("missing pool_id flag")
		}

		if poolID, err = flags.GetString("pool_id"); err != nil {
			log.Fatal(err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		err = zcncore.GetMinerSCNodePool(id, poolID, statusBar)
		if err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Println(statusBar.errMsg)
	},
}

var minerscLock = &cobra.Command{
	Use:   "mn-lock",
	Short: "Add miner/sharder stake.",
	Long:  "Add miner/sharder stake.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags  = cmd.Flags()
			id     string
			tokens float64
			err    error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if !flags.Changed("tokens") {
			log.Fatal("missing tokens flag")
		}

		if tokens, err = flags.GetFloat64("tokens"); err != nil {
			log.Fatal(err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		txn, err := zcncore.NewTransaction(statusBar, 0, nonce)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		err = txn.MinerSCLock(id, zcncore.ConvertToValue(tokens))
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

var minerscUnlock = &cobra.Command{
	Use:   "mn-unlock",
	Short: "Unlock miner/sharder stake.",
	Long:  "Unlock miner/sharder stake.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags  = cmd.Flags()
			id     string
			poolID string
			err    error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if !flags.Changed("pool_id") {
			log.Fatal("missing pool_id flag")
		}

		if poolID, err = flags.GetString("pool_id"); err != nil {
			log.Fatal(err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		txn, err := zcncore.NewTransaction(statusBar, 0, nonce)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		err = txn.MienrSCUnlock(id, poolID)
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
				fmt.Println("tokens will be unlocked next VC")
			default:
				ExitWithError("\nExecute miner unlock update smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
			return
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}
	},
}

func init() {
	rootCmd.AddCommand(minerscUpdateSettings)
	rootCmd.AddCommand(minerscInfo)
	rootCmd.AddCommand(minerscUserInfo)
	rootCmd.AddCommand(minerscPoolInfo)
	rootCmd.AddCommand(minerscLock)
	rootCmd.AddCommand(minerscUnlock)
	rootCmd.AddCommand(minerscMiners)
	rootCmd.AddCommand(minerscSharders)

	minerscMiners.PersistentFlags().Bool("json", false, "as JSON")
	minerscSharders.PersistentFlags().Bool("json", false, "as JSON")
	minerscSharders.PersistentFlags().Bool("all", false, "include all registered sharders")

	minerscUpdateSettings.PersistentFlags().String("id", "", "miner/sharder ID to update")
	minerscUpdateSettings.PersistentFlags().Int("num_delegates", 0, "max number of delegate pools")
	minerscUpdateSettings.PersistentFlags().Float64("min_stake", 0.0, "min stake allowed")
	minerscUpdateSettings.PersistentFlags().Float64("max_stake", 0.0, "max stake allowed")
	minerscUpdateSettings.MarkFlagRequired("id")

	minerscInfo.PersistentFlags().String("id", "", "miner/sharder ID to get info for")
	minerscInfo.MarkFlagRequired("id")

	minerscUserInfo.PersistentFlags().String("client_id", "", "get info for user, if empty, current user used, optional")
	minerscUserInfo.PersistentFlags().Bool("json", false, "as JSON")

	minerscPoolInfo.PersistentFlags().String("id", "", "miner/sharder ID to get info for")
	minerscPoolInfo.MarkFlagRequired("id")
	minerscPoolInfo.PersistentFlags().String("pool_id", "", "pool ID to get info for")
	minerscPoolInfo.MarkFlagRequired("pool_id")

	minerscLock.PersistentFlags().String("id", "", "miner/sharder ID to lock stake for")
	minerscLock.MarkFlagRequired("id")
	minerscLock.PersistentFlags().Float64("tokens", 0, "tokens to lock")
	minerscLock.MarkFlagRequired("tokens")

	minerscUnlock.PersistentFlags().String("id", "", "miner/sharder ID to unlock pool of")
	minerscUnlock.MarkFlagRequired("id")
	minerscUnlock.PersistentFlags().String("pool_id", "", "pool ID to unlock")
	minerscUnlock.MarkFlagRequired("pool_id")
}
