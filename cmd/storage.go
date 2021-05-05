package cmd

import (
	"fmt"
	"log"
	"sync"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var storagescUpdateConfig = &cobra.Command{
	Use:   "storagesc-update-config",
	Short: "Change settings in Storage SC.",
	Long:  "Change settings in Storage SC by owner wallet.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags      = cmd.Flags()
			err        error
			globalNode = &sdk.StorageSCConfig{ReadPool: &sdk.StorageReadPoolConfig{},
				WritePool: &sdk.StorageWritePoolConfig{},
				StakePool: &sdk.StorageStakePoolConfig{},
			}
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)

		if !flags.Changed("challenge_enabled") {
			log.Fatal("missing challenge_enabled flag")
		}
		globalNode.ChallengeEnabled, err = flags.GetBool("challenge_enabled")
		if err != nil {
			log.Fatal(err)
		}

		if flags.Changed("min_alloc_size") {
			var minAllocSize int64
			minAllocSize, err = flags.GetInt64("min_alloc_size")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MinAllocSize = common.Size(minAllocSize)
		}
		if flags.Changed("min_alloc_duration") {
			globalNode.MinAllocDuration, err = flags.GetDuration("min_alloc_duration")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_challenge_completion_time") {
			globalNode.MaxChallengeCompletionTime, err = flags.GetDuration("max_challenge_completion_time")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("min_offer_duration") {
			globalNode.MinOfferDuration, err = flags.GetDuration("min_offer_duration")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("min_blobber_capacity") {
			var minBlobberCapacity int64
			minBlobberCapacity, err = flags.GetInt64("min_blobber_capacity")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MinBlobberCapacity = common.Size(minBlobberCapacity)
		}
		if flags.Changed("validator_reward") {
			globalNode.ValidatorReward, err = flags.GetFloat64("validator_reward")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("blobber_slash") {
			globalNode.BlobberSlash, err = flags.GetFloat64("blobber_slash")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_read_price") {
			var maxReadPrice float64
			maxReadPrice, err = flags.GetFloat64("max_read_price")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MaxReadPrice = common.Balance(zcncore.ConvertToValue(maxReadPrice))
		}
		if flags.Changed("max_write_price") {
			var maxWritePrice float64
			maxWritePrice, err = flags.GetFloat64("max_write_price")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MaxWritePrice = common.Balance(zcncore.ConvertToValue(maxWritePrice))
		}
		if flags.Changed("max_challenges_per_generation") {
			globalNode.MaxChallengesPerGeneration, err = flags.GetInt("max_challenges_per_generation")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("challenge_rate_per_mb_min") {
			globalNode.ChallengeGenerationRate, err = flags.GetFloat64("challenge_rate_per_mb_min")
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
		if flags.Changed("max_charge") {
			globalNode.MaxCharge, err = flags.GetFloat64("max_charge")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("time_unit") {
			globalNode.TimeUnit, err = flags.GetDuration("time_unit")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("max_mint") {
			var maxMint float64
			maxMint, err = flags.GetFloat64("max_mint")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MaxMint = common.Balance(zcncore.ConvertToValue(maxMint))
		}
		if flags.Changed("min_stake") {
			var minMint float64
			minMint, err = flags.GetFloat64("min_stake")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MinStake = common.Balance(zcncore.ConvertToValue(minMint))
		}
		if flags.Changed("max_stake") {
			var maxStake float64
			maxStake, err = flags.GetFloat64("max_stake")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.MaxStake = common.Balance(zcncore.ConvertToValue(maxStake))
		}
		if flags.Changed("stake_min_lock") {
			var minLock float64
			minLock, err = flags.GetFloat64("stake_min_lock")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.StakePool.MinLock = common.Balance(zcncore.ConvertToValue(minLock))
		}
		if flags.Changed("stake_interest_rate") {
			globalNode.StakePool.InterestRate, err = flags.GetFloat64("stake_interest_rate")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("stake_interest_interval") {
			globalNode.StakePool.InterestInterval, err = flags.GetDuration("stake_interest_interval")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("read_min_lock") {
			var minLock float64
			minLock, err = flags.GetFloat64("read_min_lock")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.ReadPool.MinLock = common.Balance(zcncore.ConvertToValue(minLock))
		}
		if flags.Changed("read_min_lock_period") {
			globalNode.ReadPool.MinLockPeriod, err = flags.GetDuration("read_min_lock_period")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("read_max_lock_period") {
			globalNode.ReadPool.MaxLockPeriod, err = flags.GetDuration("read_max_lock_period")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("write_min_lock") {
			var minLock float64
			minLock, err = flags.GetFloat64("write_min_lock")
			if err != nil {
				log.Fatal(err)
			}
			globalNode.WritePool.MinLock = common.Balance(zcncore.ConvertToValue(minLock))
		}
		if flags.Changed("write_min_lock_period") {
			globalNode.WritePool.MinLockPeriod, err = flags.GetDuration("write_min_lock_period")
			if err != nil {
				log.Fatal(err)
			}
		}
		if flags.Changed("write_max_lock_period") {
			globalNode.WritePool.MaxLockPeriod, err = flags.GetDuration("write_max_lock_period")
			if err != nil {
				log.Fatal(err)
			}
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.StorageSCConfig(globalNode); err != nil {
			log.Fatal(err)
		}
		wg.Wait()
		fmt.Printf("Hash: %v\n", txn.GetTransactionHash())
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

		fmt.Printf("storage smart contract settings updated\nHash: %v\n", txn.GetTransactionHash())
	},
}

func init() {
	rootCmd.AddCommand(storagescUpdateConfig)

	storagescUpdateConfig.PersistentFlags().Int64("min_alloc_size", 0, "minimum allocation size")
	storagescUpdateConfig.PersistentFlags().Duration("min_alloc_duration", 0, "minimum time for an allocation")
	storagescUpdateConfig.PersistentFlags().Duration("max_challenge_completion_time", 0, "max amount of time to complete a challenge")
	storagescUpdateConfig.PersistentFlags().Duration("min_offer_duration", 0, "minimum time a blobber can offer")
	storagescUpdateConfig.PersistentFlags().Int64("min_blobber_capacity", 0, "minimum capacity for a blobber")
	storagescUpdateConfig.PersistentFlags().Float64("validator_reward", 0.0, "percent of reward for validator")
	storagescUpdateConfig.PersistentFlags().Float64("blobber_slash", 0.0, "slash penealty for failed challenge")
	storagescUpdateConfig.PersistentFlags().Float64("max_read_price", 0.0, "max price for read")
	storagescUpdateConfig.PersistentFlags().Float64("max_write_price", 0, "max price for write")
	storagescUpdateConfig.PersistentFlags().Bool("challenge_enabled", false, "enable challenge")
	storagescUpdateConfig.PersistentFlags().Int("max_challenges_per_generation", 0, "max challenges per generation")
	storagescUpdateConfig.PersistentFlags().Float64("challenge_rate_per_mb_min", 0, "challenge rate per mb/minute")
	storagescUpdateConfig.PersistentFlags().Int("max_delegates", 0, "max delegates")
	storagescUpdateConfig.PersistentFlags().Float64("max_charge", 0, "max charge")
	storagescUpdateConfig.PersistentFlags().Duration("time_unit", 0, "time unit")
	storagescUpdateConfig.PersistentFlags().Float64("max_mint", 0, "max mint amount")
	storagescUpdateConfig.PersistentFlags().Float64("min_stake", 0, "minimum stake amount")
	storagescUpdateConfig.PersistentFlags().Float64("max_stake", 0, "max stake amount")
	storagescUpdateConfig.PersistentFlags().Float64("stake_min_lock", 0, "minimum to lock for stake pool")
	storagescUpdateConfig.PersistentFlags().Float64("stake_interest_rate", 0, "stake interest rate")
	storagescUpdateConfig.PersistentFlags().Duration("stake_interest_interval", 0, "stake interest interval")
	storagescUpdateConfig.PersistentFlags().Float64("read_min_lock", 0, "minimum to lock for read pool")
	storagescUpdateConfig.PersistentFlags().Duration("read_min_lock_period", 0, "minimum amount of time to lock for read pool")
	storagescUpdateConfig.PersistentFlags().Duration("read_max_lock_period", 0, "max amount of time to lock for read pool")
	storagescUpdateConfig.PersistentFlags().Float64("write_min_lock", 0, "minimum amount to lock for write pool")
	storagescUpdateConfig.PersistentFlags().Duration("write_min_lock_period", 0, "minimum amount of time to lock for write pool")
	storagescUpdateConfig.PersistentFlags().Duration("write_max_lock_period", 0, "max amount of time to lock for write pool")
	storagescUpdateConfig.MarkFlagRequired("challenge_enabled")
}
