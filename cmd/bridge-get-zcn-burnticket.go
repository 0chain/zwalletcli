package cmd

import (
	"fmt"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-get-zcn-burn",
			"get the confirmed burn ticket for zcn burn transaction",
			"get transaction ticket with the given ZCN transaction hash",
			commandGetZCNBurnTicket,
			WithHash("ZCN transaction hash"),
		))
}

func commandGetZCNBurnTicket(b *zcnbridge.BridgeClient, args ...*Arg) {
	hash := GetHash(args)

	payload, err := b.QueryEthereumMintPayload(hash)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Println("ZCN burn ticket the completed consensus")
	fmt.Printf("Transaction nonce: %d\n", payload.Nonce)
	fmt.Printf("Transaction amount: %d\n", payload.Amount)
	fmt.Printf("ZCN transaction ID: %s\n", payload.ZCNTxnID)
}
