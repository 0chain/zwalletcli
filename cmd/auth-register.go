package cmd

import (
	"context"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcnbridge/transaction"
	"github.com/0chain/gosdk/zcncore"
	"github.com/pkg/errors"
	"log"
	"strings"
)

//goland:noinspection ALL
func init() {
	rootCmd.AddCommand(
		createCommandWithBridgeOwner(
			"auth-register",
			"Register an authorizer manually",
			"Register an authorizer manually",
			registerAuthorizerInChain,
			&Option{
				name:     "url",
				typename: "string",
				value:    "",
				usage:    "authorizer endpoint url",
				required: true,
			},
			&Option{
				name:     "client_id",
				typename: "string",
				value:    "",
				usage:    "the client_id of the wallet",
				required: true,
			},
			&Option{
				name:     "client_key",
				typename: "string",
				value:    "",
				usage:    "the client_key which is the public key of the wallet",
				required: true,
			},
			&Option{
				name:     "min_stake",
				typename: "int64",
				value:    int64(1),
				usage:    "the minimum stake value for the stake pool",
				required: false,
			},
			&Option{
				name:     "max_stake",
				typename: "int64",
				value:    int64(10),
				usage:    "the maximum stake value for the stake pool",
				required: false,
			},
			&Option{
				name:     "num_delegates",
				typename: "int",
				value:    5,
				usage:    "the number of delegates in the authorizer stake pool",
				required: false,
			},
			&Option{
				name:     "service_charge",
				typename: "float64",
				value:    0.0,
				usage:    "the service charge for the authorizer stake pool",
				required: false,
			},
		))
}

// registerAuthorizerInChain registers a new authorizer
// addAuthorizerPayload *addAuthorizerPayload
func registerAuthorizerInChain(bo *zcnbridge.BridgeOwner, args ...*Arg) {
	clientID := GetClientID(args)
	clientKey := GetClientKey(args)
	url := GetURL(args)
	minStake := GetMinStake(args)
	maxStake := GetMaxStake(args)
	numDelegates := GetNumDelegates(args)
	serviceCharge := GetServiceCharge(args)

	input := &zcncore.AddAuthorizerPayload{
		PublicKey: clientKey,
		URL:       url,
		StakePoolSettings: zcncore.AuthorizerStakePoolSettings{
			DelegateWallet: clientID,
			MinStake:       common.Balance(minStake),
			MaxStake:       common.Balance(maxStake),
			NumDelegates:   numDelegates,
			ServiceCharge:  serviceCharge,
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
			ExitWithError("Authorizer has already been added to 0Chain...  Continue")
		} else {
			ExitWithError(errors.Wrapf(err, "failed to verify transaction: '%s'", trx.Hash))
		}
	}
}
