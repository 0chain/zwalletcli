package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			"bridge-get-zcn-burn",
			"get the confirmed burn ticket for zcn burn transaction",
			"get transaction ticket with the given ZCN transaction hash",
			commandGetZCNBurnTicket,
			hashOption,
		))
}

func commandGetZCNBurnTicket(b *zcnbridge.Bridge, args ...*Arg) {
	hash := GetHash(args)

	payload, err := b.QueryZChainMintPayload(hash)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Println("ZCN burn ticket the completed consensus")
	fmt.Printf("Transaction nonce: %d\n", payload.Nonce)
	fmt.Printf("Transaction amount: %d\n", payload.Amount)
	fmt.Printf("Ethereum transaction ID: %s\n", payload.EthereumTxnID)
	fmt.Printf("ZCN receiving client ID: %s\n", payload.ReceivingClientID)
}
