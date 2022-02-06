package cmd

import (
	"context"
	"fmt"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-burn-zcn",
			"burn zcn tokens",
			"burn zcn tokens that will be minted for WZCN tokens",
			commandBurnZCN,
			WithAmount("ZCN token amount to be burned"),
		))
}

func commandBurnZCN(b *zcnbridge.BridgeClient, args ...*Arg) {
	amount := GetAmount(args)

	fmt.Println("Starting burn transaction")
	transaction, err := b.BurnZCN(context.Background(), amount)
	if err == nil {
		fmt.Printf("Submitted burn transaction %s\n", transaction.Hash)
	} else {
		ExitWithError(err)
	}

	fmt.Printf("Starting transaction verification %s\n", transaction.Hash)
	verify(transaction.Hash)
}
