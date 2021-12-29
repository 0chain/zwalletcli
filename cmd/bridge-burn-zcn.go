package cmd

import (
	"context"
	"fmt"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			"bridge-burn-zcn",
			"burn zcn tokens",
			"burn zcn tokens that will be minted on Ethereum chain",
			commandBurnZCN,
			amountOption,
		))
}

func commandBurnZCN(b *zcnbridge.BridgeClient, args ...*Arg) {
	amount := GetAmount(args)

	fmt.Printf("Starting burn transaction")
	transaction, err := b.BurnZCN(context.Background(), amount)
	fmt.Printf("Submitted burn transaction %s\n", transaction.Hash)

	if err == nil {
		fmt.Printf("Transaction confirmed")
	} else {
		ExitWithError(err)
	}
}
