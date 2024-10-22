package cmd

import (
	"github.com/0chain/gosdk/core/transaction"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
	"github.com/pkg/errors"
	"log"
	"strings"
)

//goland:noinspection ALL
func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"auth-register",
			"Register an authorizer manually",
			"Register an authorizer manually",
			registerAuthorizerInChain,
			true,
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
func registerAuthorizerInChain(bc *zcnbridge.BridgeClient, args ...*Arg) {
	clientID := GetClientID(args)
	clientKey := GetClientKey(args)
	url := GetURL(args)
	numDelegates := GetNumDelegates(args)
	serviceCharge := GetServiceCharge(args)

	input := &zcncore.AddAuthorizerPayload{
		PublicKey: clientKey,
		URL:       url,
		StakePoolSettings: zcncore.AuthorizerStakePoolSettings{
			DelegateWallet: clientID,
			NumDelegates:   numDelegates,
			ServiceCharge:  serviceCharge,
		},
	}

	hash, _, _, txn, err := zcncore.ZCNSCAddAuthorizer(input)
	if err != nil {
		log.Fatal(err, "failed to add authorizer with transaction: '%s'", hash)
	}

	log.Printf("Authorizer submitted OK... " + hash)
	log.Printf("Starting verification: " + hash)

	txn, err = transaction.VerifyTransaction(hash)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			ExitWithError("Authorizer has already been added to 0Chain...  Continue")
		} else {
			ExitWithError(errors.Wrapf(err, "failed to verify transaction: '%s'", txn.Hash))
		}
	}
}
