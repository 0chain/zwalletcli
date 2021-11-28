package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			commandGetETHBurnTicket,
			"bridge-get-eth-burn",
			"get confirmed burn ticket for ethereum burn transaction",
			"get confirmed burn ticket for ethereum burn transaction",
		))
}

func commandGetETHBurnTicket(b *zcnbridge.Bridge, hash string) {
	payload, err := b.QueryEthereumMintPayload(hash)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Println("Ethereum burn ticket the completed consensus")
	fmt.Printf("Transaction nonce: %d\n", payload.Nonce)
	fmt.Printf("Transaction amount: %d\n", payload.Amount)
	fmt.Printf("ZCN transaction ID: %s\n", payload.ZCNTxnID)
}
