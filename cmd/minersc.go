package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

var minerscUpdateNodeSettings = &cobra.Command{
	Use:   "mn-update-node-settings",
	Short: "Change miner/sharder settings in Miner SC.",
	Long:  "Change miner/sharder settings in Miner SC by delegate wallet.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags   = cmd.Flags()
			id      string
			sharder bool
			err     error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if sharder, err = flags.GetBool("sharder"); err != nil {
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

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if sharder {
			err = txn.MinerSCSharderSettings(miner)
		} else {
			err = txn.MinerSCMinerSettings(miner)
		}
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

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Println("settings updated")
	},
}

var minerscDeleteNode = &cobra.Command{
	Use:   "minersc-delete-node",
	Short: "Delete a miner or sharder node from Miner SC.",
	Long:  "Delete a miner or sharder node from Miner SC.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags   = cmd.Flags()
			id      string
			sharder bool
			err     error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if sharder, err = flags.GetBool("sharder"); err != nil {
			log.Fatal(err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)

		// remove not settings fields
		miner := &zcncore.MinerSCMinerInfo{SimpleMinerSCMinerInfo: &zcncore.SimpleMinerSCMinerInfo{
			ID: id,
		},
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if sharder {
			err = txn.MinerSCDeleteSharder(miner)
		} else {
			err = txn.MinerSCDeleteMiner(miner)
		}
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

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Println("settings updated")
	},
}

var minerscUpdateSettings = &cobra.Command{
	Use:   "minersc-update-settings",
	Short: "Change settings in Miner SC.",
	Long:  "Change settings in Miner SC by owner wallet.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags      = cmd.Flags()
			err        error
			globalNode = new(zcncore.MinerSCConfig)
			wg         sync.WaitGroup
			statusBar  = &ZCNStatus{wg: &wg}
		)

		if flags.Changed("max_n") {
			globalNode.MaxN, err = flags.GetInt("max_n")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("min_n") {
			globalNode.MinN, err = flags.GetInt("min_n")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_s") {
			globalNode.MaxS, err = flags.GetInt("max_s")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("min_s") {
			globalNode.MinS, err = flags.GetInt("min_s")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_delegates") {
			globalNode.MaxDelegates, err = flags.GetInt("max_delegates")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("t_percent") {
			if globalNode.TPercent, err = flags.GetFloat64("t_percent"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("k_percent") {
			if globalNode.KPercent, err = flags.GetFloat64("k_percent"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_stake") {
			var max float64
			if max, err = flags.GetFloat64("max_stake"); err != nil {
				log.Fatal(err)
			}
			globalNode.MaxStake = common.Balance(zcncore.ConvertToValue(max))
		}
		if flags.Changed("min_stake") {
			var min float64
			if min, err = flags.GetFloat64("min_stake"); err != nil {
				log.Fatal(err)
			}
			globalNode.MinStake = common.Balance(zcncore.ConvertToValue(min))
		}
		if flags.Changed("interest_rate") {
			if globalNode.InterestRate, err = flags.GetFloat64("interest_rate"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("reward_rate") {
			if globalNode.RewardRate, err = flags.GetFloat64("reward_rate"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("share_ratio") {
			if globalNode.ShareRatio, err = flags.GetFloat64("share_ratio"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("block_reward") {
			var blockReward float64
			if blockReward, err = flags.GetFloat64("block_reward"); err != nil {
				log.Fatal(err)
			}
			globalNode.BlockReward = common.Balance(zcncore.ConvertToValue(blockReward))
		}
		if flags.Changed("max_charge") {
			if globalNode.MaxCharge, err = flags.GetFloat64("max_charge"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("epoch") {
			if globalNode.Epoch, err = flags.GetInt64("epoch"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("reward_decline_rate") {
			if globalNode.RewardDeclineRate, err = flags.GetFloat64("reward_decline_rate"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("interest_decline_rate") {
			if globalNode.InterestDeclineRate, err = flags.GetFloat64("interest_decline_rate"); err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_mint") {
			var maxMint float64
			if maxMint, err = flags.GetFloat64("max_mint"); err != nil {
				log.Fatal(err)
			}
			globalNode.MaxMint = common.Balance(zcncore.ConvertToValue(maxMint))
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.MinerSCSettings(globalNode); err != nil {
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

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Println("miner smart contract settings updated")
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

		if len(info.Nodes) == 0 {
			fmt.Println("no miners in Miner SC")
			return
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

		var (
			flags = cmd.Flags()
			err   error
			info  = new(zcncore.MinerSCNodes)
			cb    = NewJSONInfoCB(info)
		)
		if err = zcncore.GetSharders(cb); err != nil {
			log.Fatal(err)
		}

		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		if len(info.Nodes) == 0 {
			fmt.Println("no sharders in Miner SC")
			return
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

		for _, node := range info.Nodes {
			fmt.Println("- ID:        ", node.Miner.ID)
			fmt.Println("- Host:      ", node.Miner.Host)
			fmt.Println("- Port:      ", node.Miner.Port)
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
		txn, err := zcncore.NewTransaction(statusBar, 0)
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

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Println("locked with:", txn.GetTransactionHash())
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
		txn, err := zcncore.NewTransaction(statusBar, 0)
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

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Println("tokens will be unlocked next VC")
	},
}

var minerConfig = &cobra.Command{
	Use:   "mn-config",
	Short: "Get miner SC global info.",
	Long:  "Get miner SC global info.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			conf = new(zcncore.MinerSCConfig)
			cb   = NewJSONInfoCB(conf)
			err  error
		)

		if err = zcncore.GetMinerSCConfig(cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("view_change:          ", conf.ViewChange)
		fmt.Println("max_n:                ", conf.MaxN)
		fmt.Println("min_n:                ", conf.MinN)
		fmt.Println("max_s:                ", conf.MaxS)
		fmt.Println("min_s:                ", conf.MinS)
		fmt.Println("t_percent:            ", conf.TPercent)
		fmt.Println("k_percent:            ", conf.KPercent)
		fmt.Println("last_round:           ", conf.LastRound)
		fmt.Println("max_stake:            ", conf.MaxStake)
		fmt.Println("min_stake:            ", conf.MinStake)
		fmt.Println("interest_rate:        ", conf.InterestRate)
		fmt.Println("reward_rate:          ", conf.RewardRate)
		fmt.Println("share_ratio:          ", conf.ShareRatio)
		fmt.Println("block_reward:         ", conf.BlockReward)
		fmt.Println("max_charge:           ", conf.MaxCharge)
		fmt.Println("epoch:                ", conf.Epoch)
		fmt.Println("reward_decline_rate:  ", conf.RewardDeclineRate)
		fmt.Println("interest_decline_rate:", conf.InterestDeclineRate)
		fmt.Println("max_mint:             ", conf.MaxMint)
		fmt.Println("minted:               ", conf.Minted)
		fmt.Println("max_delegates:        ", conf.MaxDelegates)
	},
}

func init() {
	rootCmd.AddCommand(minerscUpdateNodeSettings)
	rootCmd.AddCommand(minerscUpdateSettings)
	rootCmd.AddCommand(minerscInfo)
	rootCmd.AddCommand(minerscUserInfo)
	rootCmd.AddCommand(minerscPoolInfo)
	rootCmd.AddCommand(minerscLock)
	rootCmd.AddCommand(minerscUnlock)
	rootCmd.AddCommand(minerConfig)
	rootCmd.AddCommand(minerscMiners)
	rootCmd.AddCommand(minerscSharders)
	rootCmd.AddCommand(minerscDeleteNode)

	minerscMiners.PersistentFlags().Bool("json", false, "as JSON")
	minerscSharders.PersistentFlags().Bool("json", false, "as JSON")

	minerscUpdateNodeSettings.PersistentFlags().String("id", "", "miner/sharder ID to update")
	minerscUpdateNodeSettings.PersistentFlags().Bool("sharder", false, "set true for sharder node")
	minerscUpdateNodeSettings.PersistentFlags().Int("num_delegates", 0, "max number of delegate pools")
	minerscUpdateNodeSettings.PersistentFlags().Float64("min_stake", 0.0, "min stake allowed")
	minerscUpdateNodeSettings.PersistentFlags().Float64("max_stake", 0.0, "max stake allowed")
	minerscUpdateNodeSettings.MarkFlagRequired("id")

	minerscUpdateSettings.PersistentFlags().Int("max_n", 0, "max number of miner nodes")
	minerscUpdateSettings.PersistentFlags().Int("min_n", 0, "minimum number of miner nodes")
	minerscUpdateSettings.PersistentFlags().Int("max_s", 0, "max number of sharder nodes")
	minerscUpdateSettings.PersistentFlags().Int("min_s", 0, "minimum number of sharder nodes")
	minerscUpdateSettings.PersistentFlags().Int("max_delegates", 0, "max number of delegate pools")
	minerscUpdateSettings.PersistentFlags().Float64("t_percent", 0.0, "threshold percent for miners dkg threshold")
	minerscUpdateSettings.PersistentFlags().Float64("k_percent", 0.0, "percent of miners needed to finish dkg")
	minerscUpdateSettings.PersistentFlags().Float64("max_stake", 0.0, "max stake for nodes")
	minerscUpdateSettings.PersistentFlags().Float64("min_stake", 0.0, "minimum stake for nodes")
	minerscUpdateSettings.PersistentFlags().Float64("interest_rate", 0.0, "interest rate earned")
	minerscUpdateSettings.PersistentFlags().Float64("reward_rate", 0.0, "reward rate")
	minerscUpdateSettings.PersistentFlags().Float64("share_ratio", 0.0, "share ratio")
	minerscUpdateSettings.PersistentFlags().Float64("block_reward", 0.0, "block reward")
	minerscUpdateSettings.PersistentFlags().Float64("max_charge", 0.0, "max charge")
	minerscUpdateSettings.PersistentFlags().Int64("epoch", 0.0, "epoch")
	minerscUpdateSettings.PersistentFlags().Float64("reward_decline_rate", 0.0, "reward decline rate")
	minerscUpdateSettings.PersistentFlags().Float64("interest_decline_rate", 0.0, "interest decline rate")
	minerscUpdateSettings.PersistentFlags().Float64("max_mint", 0.0, "max mint allowed")

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

	minerscDeleteNode.PersistentFlags().String("id", "", "miner/sharder ID of node to delete")
	minerscDeleteNode.PersistentFlags().Bool("sharder", false, "set for true if you delete sharder")
	minerscDeleteNode.MarkFlagRequired("id")
}
