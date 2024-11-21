package cmd

import (
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
			false,
			WithToken("ZCN tokens quantity to be burned"),
		))
}

func commandBurnZCN(b *zcnbridge.BridgeClient, args ...*Arg) {
	amount := GetToken(args)

	fmt.Println("Starting burn transaction")
	hash, _, err := b.BurnZCN(zcncore.ConvertToValue(amount))
	if err == nil {
		fmt.Printf("Submitted burn transaction %s\n", hash)
	} else {
		ExitWithError(err)
	}

	fmt.Printf("Transaction completed successfully: %s\n", hash)
}
