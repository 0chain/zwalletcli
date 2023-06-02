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
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

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
	Short: "Get miners from Miner SC",
	Long:  "Get miners from Miner SC",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			err   error
			info  = new(zcncore.MinerSCNodes)
		)

		limit, offset := 20, 0
		active := true

		var allFlag, jsonFlag bool

		if flags.Changed("all") {
			allFlag, err = flags.GetBool("all")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("limit") {
			limit, err = flags.GetInt("limit")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("offset") {
			offset, err = flags.GetInt("offset")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("active") {
			active, err = flags.GetBool("active")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("json") {
			jsonFlag, err = flags.GetBool("json")
			if err != nil {
				log.Fatal(err)
			}
		}

		if !allFlag {
			cb := NewJSONInfoCB(info)
			zcncore.GetMiners(cb, limit, offset, active)

			if err = cb.Waiting(); err != nil {
				log.Fatal(err)
			}

			if jsonFlag {
				util.PrintJSON(info)
				return
			}

			if len(info.Nodes) == 0 {
				fmt.Println("no miners in Miner SC")
				return
			}

			printMinerNodes(info.Nodes)
			return
		} else {
			limit = 20
			offset = 0

			var nodes []zcncore.Node
			for curOff := offset; ; curOff += limit {
				cb := NewJSONInfoCB(info)
				zcncore.GetMiners(cb, limit, curOff, active)

				if err = cb.Waiting(); err != nil {
					log.Fatal(err)
				}

				if len(info.Nodes) == 0 {
					break
				}

				nodes = append(nodes, info.Nodes...)
			}

			if jsonFlag {
				util.PrintJSON(nodes)
			} else {
				printMinerNodes(nodes)
			}
		}
	},
}

func printMinerNodes(nodes []zcncore.Node) {
	for _, node := range nodes {
		fmt.Println("- ID:        ", node.Miner.ID)
		fmt.Println("- Host:      ", node.Miner.Host)
		fmt.Println("- Port:      ", node.Miner.Port)
	}
}

var minerscSharders = &cobra.Command{
	Use:   "ls-sharders",
	Short: "Get sharders from Miner SC",
	Long:  "Get sharders from Miner SC",
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

		limit, offset := 20, 0
		active := true
		if flags.Changed("limit") {
			limit, err = flags.GetInt("limit")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("offset") {
			offset, err = flags.GetInt("offset")
			if err != nil {
				log.Fatal(err)
			}
		}

		if flags.Changed("active") {
			active, err = flags.GetBool("active")
			if err != nil {
				log.Fatal(err)
			}
		}

		if !allFlag {
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
		} else {
			sharders := new(zcncore.MinerSCNodes)

			limit = 20
			offset = 0
			var nodes []zcncore.Node
			for curOff := offset; ; curOff += limit {
				callback := NewJSONInfoCB(sharders)
				zcncore.GetSharders(callback, limit, curOff, active)

				if err = callback.Waiting(); err != nil {
					log.Fatal(err)
				}

				if len(sharders.Nodes) == 0 {
					break
				}

				nodes = append(nodes, sharders.Nodes...)
			}

			if jsonFlag {
				util.PrettyPrintJSON(nodes)
			} else {
				printSharderNodes(nodes)
			}
		}
	},
}

func printSharderNodes(nodes []zcncore.Node) {
	for _, node := range nodes {
		fmt.Println("ID:", node.Miner.ID)
		fmt.Println("  - N2NHost:", node.Miner.N2NHost)
		fmt.Println("  - Host:", node.Miner.Host)
		fmt.Println("  - Port:", node.Miner.Port)
	}
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
		for _, delegates := range info.Pools {
			for _, pool := range delegates {
				total += pool.Balance
			}
		}

		for key, delegates := range info.Pools {
			for _, pool := range delegates {
				fmt.Println("- delegates:", "("+key+")")
				fmt.Println("  - pool_id:            ", pool.ID)
				fmt.Println("    balance:            ", pool.Balance)
				fmt.Println("    rewards uncollected:", pool.Reward)
				fmt.Println("    rewards paid:       ", pool.RewardPaid)
				fmt.Println("    status:             ", strings.ToLower(pool.Status))
				fmt.Println("    stake %:            ",
					float64(pool.Balance)/float64(total)*100.0, "%")
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
			flags = cmd.Flags()
			id    string

			err error
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
		err = zcncore.GetMinerSCNodePool(id, statusBar)
		if err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			fields := map[string]string{}
			err := json.Unmarshal([]byte(statusBar.errMsg), &fields)
			if err != nil {
				log.Fatal("fatal:", statusBar.errMsg)
			}
			ExitWithError(fields["error"])
			return
		}

		fmt.Println(statusBar.errMsg)
	},
}

// spLock locks tokens a stake pool lack
var spLock = &cobra.Command{
	Use:   "sp-lock",
	Short: "Lock tokens lacking in stake pool.",
	Long:  `Lock tokens lacking in stake pool.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags        = cmd.Flags()
			providerID   string
			providerType sdk.ProviderType
			tokens       float64
			fee          float64
			err          error
		)

		if flags.Changed("blobber_id") {
			if providerID, err = flags.GetString("blobber_id"); err != nil {
				log.Fatalf("invalid 'blobber_id' flag: %v", err)
			} else {
				providerType = sdk.ProviderBlobber
			}
		} else if flags.Changed("validator_id") {
			if providerID, err = flags.GetString("validator_id"); err != nil {
				log.Fatalf("invalid 'validator_id' flag: %v", err)
			} else {
				providerType = sdk.ProviderValidator
			}
		}

		if providerType == 0 || providerID == "" {
			log.Fatal("missing flag: one of 'blobber_id' or 'validator_id' is required")
		}

		if !flags.Changed("tokens") {
			log.Fatal("missing required 'tokens' flag")
		}

		if tokens, err = flags.GetFloat64("tokens"); err != nil {
			log.Fatal("invalid 'tokens' flag: ", err)
		}

		if tokens < 0 {
			log.Fatal("invalid token amount: negative")
		}

		if flags.Changed("fee") {
			if fee, err = flags.GetFloat64("fee"); err != nil {
				log.Fatal("invalid 'fee' flag: ", err)
			}
		}

		var hash string
		hash, _, err = sdk.StakePoolLock(providerType, providerID,
			zcncore.ConvertToValue(tokens), zcncore.ConvertToValue(fee))
		if err != nil {
			log.Fatalf("Failed to lock tokens in stake pool: %v", err)
		}
		fmt.Println("tokens locked, txn hash:", hash)
	},
}

var minerscLock = &cobra.Command{
	Use:   "mn-lock",
	Short: "Add miner/sharder stake.",
	Long:  "Add miner/sharder stake.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags        = cmd.Flags()
			providerID   string
			providerType zcncore.Provider
			tokens       float64
			err          error
		)

		if flags.Changed("miner_id") {
			if providerID, err = flags.GetString("miner_id"); err != nil {
				log.Fatalf("invalid 'miner_id' flag: %v", err)
			} else {
				providerType = zcncore.ProviderMiner
			}
		} else if flags.Changed("sharder_id") {
			if providerID, err = flags.GetString("sharder_id"); err != nil {
				log.Fatalf("invalid 'sharder_id' flag: %v", err)
			} else {
				providerType = zcncore.ProviderSharder
			}
		}

		if providerType == 0 || providerID == "" {
			log.Fatal("missing flag: one of 'miner_id' or 'sharder_id' is required")
		}

		if !flags.Changed("tokens") {
			log.Fatal("missing tokens flag")
		}

		if tokens, err = flags.GetFloat64("tokens"); err != nil {
			log.Fatal(err)
		}
		if tokens < 0 {
			log.Fatal("invalid token amount: negative")
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		txn, err := zcncore.NewTransaction(statusBar, getTxnFee(), nonce)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		err = txn.MinerSCLock(providerID, providerType, zcncore.ConvertToValue(tokens))
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
			flags        = cmd.Flags()
			providerID   string
			providerType zcncore.Provider
			err          error
		)

		if flags.Changed("miner_id") {
			if providerID, err = flags.GetString("miner_id"); err != nil {
				log.Fatalf("invalid 'miner_id' flag: %v", err)
			} else {
				providerType = zcncore.ProviderMiner
			}
		} else if flags.Changed("sharder_id") {
			if providerID, err = flags.GetString("sharder_id"); err != nil {
				log.Fatalf("invalid 'sharder_id' flag: %v", err)
			} else {
				providerType = zcncore.ProviderSharder
			}
		}

		if providerType == 0 || providerID == "" {
			log.Fatal("missing flag: one of 'miner_id' or 'sharder_id' is required")
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		txn, err := zcncore.NewTransaction(statusBar, getTxnFee(), nonce)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		err = txn.MinerSCUnlock(providerID, providerType)
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
				fmt.Println("tokens unlocked")
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
	rootCmd.AddCommand(minerscInfo)
	rootCmd.AddCommand(minerscUserInfo)
	rootCmd.AddCommand(minerscPoolInfo)
	rootCmd.AddCommand(minerscLock)
	rootCmd.AddCommand(minerscUnlock)
	rootCmd.AddCommand(minerscMiners)
	rootCmd.AddCommand(minerscSharders)

	minerscMiners.PersistentFlags().Bool("json", false, "as JSON")
	minerscMiners.PersistentFlags().Int("limit", 20, "Limits the amount of miners returned")
	minerscMiners.PersistentFlags().Int("offset", 0, "Skips the number of miners mentioned")
	minerscMiners.PersistentFlags().Bool("active", true, "Gets active miners only, set it false to get all miners")
	minerscMiners.PersistentFlags().Bool("all", false, "include all registered miners, default returns the first page of miners")
	minerscSharders.PersistentFlags().Bool("json", false, "as JSON")
	minerscSharders.PersistentFlags().Int("limit", 20, "Limits the amount of sharders returned")
	minerscSharders.PersistentFlags().Int("offset", 0, "Skips the number of sharders mentioned")
	minerscSharders.PersistentFlags().Bool("all", false, "include all registered sharders, default returns the first page of sharders")
	minerscSharders.PersistentFlags().Bool("active", true, "Gets active sharders only, set it false to get all sharders")

	minerscInfo.PersistentFlags().String("id", "", "miner/sharder ID to get info for")
	minerscInfo.MarkFlagRequired("id")

	minerscUserInfo.PersistentFlags().String("client_id", "", "get info for user, if empty, current user used, optional")
	minerscUserInfo.PersistentFlags().Bool("json", false, "as JSON")

	minerscPoolInfo.PersistentFlags().String("id", "", "miner/sharder ID to get info for")
	minerscPoolInfo.MarkFlagRequired("id")

	minerscLock.PersistentFlags().String("miner_id", "", "miner ID to lock stake for")
	minerscLock.PersistentFlags().String("sharder_id", "", "sharder ID to lock stake for")
	minerscLock.PersistentFlags().Float64("tokens", 0, "tokens to lock")
	minerscLock.MarkFlagRequired("tokens")

	minerscUnlock.PersistentFlags().String("miner_id", "", "miner ID to lock stake for")
	minerscUnlock.PersistentFlags().String("sharder_id", "", "sharder ID to lock stake for")
}
