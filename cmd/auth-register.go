package cmd

import (
	"context"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcnbridge/transaction"
	"github.com/0chain/gosdk/zcncore"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
	"strings"
)

type addAuthorizerPayload struct {
	URL           string
	ClientID      string
	ClientKey     string
	MinStake      int64
	MaxStake      int64
	NumDelegates  int
	ServiceCharge float64
}

var getAuthorizerRegisterCmd = &cobra.Command{
	Use:   "auth-register",
	Short: "Register an authorizer manually",
	Long:  `Register an authorizer manually.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			err   error
		)

		payload := &addAuthorizerPayload{}
		if flags.Changed("url") {
			if payload.URL, err = flags.GetString("url"); err != nil {
				log.Fatalf("error in 'url' flag: %v", err)
			}
		} else {
			ExitWithError("Error: url flag is missing")
		}

		if flags.Changed("client_id") {
			if payload.ClientID, err = flags.GetString("client_id"); err != nil {
				log.Fatalf("error in 'client_id' flag: %v", err)
			}
		} else {
			ExitWithError("Error: client_id flag is missing")
		}

		if flags.Changed("client_key") {
			if payload.ClientKey, err = flags.GetString("client_key"); err != nil {
				log.Fatalf("error in 'client_key' flag: %v", err)
			}
		} else {
			ExitWithError("Error: client_key flag is missing")
		}

		if payload.MaxStake, err = flags.GetInt64("max_stake"); err != nil {
			log.Fatalf("error in 'max_stake' flag: %v", err)
		}

		if payload.MinStake, err = flags.GetInt64("min_stake"); err != nil {
			log.Fatalf("error in 'min_stake' flag: %v", err)
		}

		if payload.NumDelegates, err = flags.GetInt("num_delegates"); err != nil {
			log.Fatalf("error in 'num_delegates' flag: %v", err)
		}

		if payload.ServiceCharge, err = flags.GetFloat64("service_charge"); err != nil {
			log.Fatalf("error in 'service_charge' flag: %v", err)
		}

		registerAuthorizerInChain(payload)
	},
}

//goland:noinspection ALL
func init() {
	cmd := getAuthorizerRegisterCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String("url", "", "authorizer endpoint url")
	cmd.PersistentFlags().String("client_id", "", "the client_id of the wallet")
	cmd.PersistentFlags().String("client_key", "", "the client_key which is the public key of the wallet")
	cmd.PersistentFlags().Int64("min_stake", 1, "the minimum stake value for the stake pool")
	cmd.PersistentFlags().Int64("max_stake", 10, "the maximum stake value for the stake pool")
	cmd.PersistentFlags().Int("num_delegates", 5, "the number of delegates in the authorizer stake pool")
	cmd.PersistentFlags().Float64("service_charge", 0.0, "the service charge for the authorizer stake pool")
}

// registerAuthorizerInChain registers a new authorizer
func registerAuthorizerInChain(addAuthorizerPayload *addAuthorizerPayload) {
	input := &zcncore.AddAuthorizerPayload{
		PublicKey: addAuthorizerPayload.ClientKey,
		URL:       addAuthorizerPayload.URL,
		StakePoolSettings: zcncore.AuthorizerStakePoolSettings{
			DelegateWallet: addAuthorizerPayload.ClientID,
			MinStake:       common.Balance(addAuthorizerPayload.MinStake),
			MaxStake:       common.Balance(addAuthorizerPayload.MaxStake),
			NumDelegates:   addAuthorizerPayload.NumDelegates,
			ServiceCharge:  addAuthorizerPayload.ServiceCharge,
		},
	}

	trx, err := transaction.AddAuthorizer(context.Background(), input)
	if err != nil {
		log.Fatal(err, "failed to add authorizer with transaction: '%s'", trx.Hash)
	}

	log.Printf("Authorizer submitted OK... " + trx.Hash)
	log.Printf("Starting verification: " + trx.Hash)

	err = trx.Verify(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Print("Authorizer has already been added to 0Chain...  Continue")
		} else {
			ExitWithError(errors.Wrapf(err, "failed to verify transaction: '%s'", trx.Hash))
		}
	}
}
