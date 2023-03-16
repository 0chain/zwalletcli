package cmd

import (
	"context"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcnbridge/transaction"
	"github.com/0chain/gosdk/zcncore"
)

//goland:noinspection ALL
func init() {
	rootCmd.AddCommand(
		createCommandWithBridgeOwner(
			"authorizer-hc",
			"",
			"",
			authorizerHc,
		))
}

// registerAuthorizerInSC registers a new authorizer to token bridge SC
func authorizerHc(bo *zcnbridge.BridgeOwner, args ...*Arg) {
	payload := &zcncore.AuthorizerHealthCheckPayload{
		ID: "",
	}

	trx, err := transaction.AuthorizerHealthCheck(context.Background(), payload)
	if err != nil {
		ExitWithError(err)
	}

	if err := trx.Verify(context.Background()); err != nil {
		ExitWithError(err)
	}
}
