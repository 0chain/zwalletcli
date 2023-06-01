package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	comm "github.com/ethereum/go-ethereum/common"
)

//goland:noinspection ALL
func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"auth-sc-register",
			"Register an authorizer to token bridge SC manually",
			"Register an authorizer to token bridge SC manually",
			registerAuthorizerInSC,
			&Option{
				name:     "ethereum_address",
				typename: "string",
				value:    "",
				usage:    "ethereum address which is authorizer linked to",
				required: true,
			},
		))
}

// registerAuthorizerInSC registers a new authorizer to token bridge SC
func registerAuthorizerInSC(bc *zcnbridge.BridgeClient, args ...*Arg) {
	ethereumAddress := GetEthereumAddress(args)

	tx, err := bc.AddEthereumAuthorizer(context.Background(), comm.HexToAddress(ethereumAddress))
	if err != nil {
		ExitWithError(err)
	}

	hash := tx.Hash().String()
	fmt.Printf("Confirming Ethereum mint transaction: %s\n", hash)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 100, time.Second*5)
	if err != nil {
		ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("\nTransaction verification success: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
	}
}
