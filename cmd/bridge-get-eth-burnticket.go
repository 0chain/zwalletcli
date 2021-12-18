package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			"bridge-get-eth-burn",
			"get confirmed burn ticket for ethereum burn transaction",
			"get transaction ticket with the given Ethereum transaction hash",
			commandGetETHBurnTicket,
			hashOption,
		))
}

func commandGetETHBurnTicket(b *zcnbridge.BridgeClient, args ...*Arg) {
	hash := GetHash(args)

	payload, err := b.QueryEthereumMintPayload(hash)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Println("Ethereum burn ticket the completed consensus")
	fmt.Printf("Transaction nonce: %d\n", payload.Nonce)
	fmt.Printf("Transaction amount: %d\n", payload.Amount)
	fmt.Printf("ZCN transaction ID: %s\n", payload.ZCNTxnID)
}
