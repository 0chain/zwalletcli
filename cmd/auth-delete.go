package cmd

import (
	"context"
	"log"
	"strings"

	// "0chain.net/chaincore/smartcontractinterface"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcnbridge/transaction"
	"github.com/0chain/gosdk/zcncore"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type deleteAuthorizerPayload struct {
	ID string
}

var getAuthorizerDeleteCmd = &cobra.Command{
	Use:   "auth-delete",
	Short: "Register an authorizer manually",
	Long:  `Register an authorizer manually.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			err   error
		)

		payload := &deleteAuthorizerPayload{}
		if flags.Changed("id") {
			if payload.ID, err = flags.GetString("id"); err != nil {
				log.Fatalf("error in 'id' flag: %v", err)
			}
		} else {
			ExitWithError("Error: id flag is missing")
		}

		deleteAuthorizerInChain(payload)
	},
}

//goland:noinspection ALL
func init() {
	cmd := getAuthorizerDeleteCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String("id", "", "authorizer id")
}

// registerAuthorizerInChain registers a new authorizer
func deleteAuthorizerInChain(payload *deleteAuthorizerPayload) {
	input := &zcncore.DeleteAuthorizerPayload{
		ID: common.Key(payload.ID),
	}

	trx, err := transaction.DeleteAuthorizer(context.Background(), input)
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
