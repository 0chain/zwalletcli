package cmd

import (
	"context"
	"fmt"

	"github.com/0chain/gosdk/zcncore"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-burn-zcn",
			"burn zcn tokens",
			"burn zcn tokens that will be minted for WZCN tokens",
			commandBurnZCN,
			WithToken("ZCN tokens quantity to be burned"),
		))
}

func commandBurnZCN(b *zcnbridge.BridgeClient, args ...*Arg) {
	amount := GetToken(args)

	fmt.Println("Starting burn transaction")
	transaction, err := b.BurnZCN(context.Background(), zcncore.ConvertToValue(amount))
	if err == nil {
		fmt.Printf("Submitted burn transaction %s\n", transaction.Hash)
	} else {
		ExitWithError(err)
	}

	fmt.Printf("Starting transaction verification %s\n", transaction.Hash)
	err = transaction.Verify(context.Background())
	if err != nil {
		ExitWithError(err)
	}

	fmt.Printf("Transaction completed successfully: %s\n", transaction.Hash)
}
