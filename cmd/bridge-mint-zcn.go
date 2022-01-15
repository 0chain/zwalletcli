package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-mint-zcn",
			"mint zcn tokens using the hash of Ethereum burn transaction",
			"mint zcn tokens after burning WZCN tokens in Ethereum chain",
			commandMintZCN,
			hashOption,
		))
}

func commandMintZCN(b *zcnbridge.BridgeClient, args ...*Arg) {
	hash := GetHash(args)

	fmt.Printf("Query ticket for Ethereum transaction hash: %s\n", hash)

	payload, err := b.QueryZChainMintPayload(hash)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Printf("Sending mint transaction to ZCN\n")
	fmt.Printf("Ethereum transaction ID: %s\n", payload.EthereumTxnID)
	fmt.Printf("Payload amount: %d\n", payload.Amount)
	fmt.Printf("Payload nonce: %d\n", payload.Nonce)
	fmt.Printf("Receiving ZCN ClientID: %s\n", payload.ReceivingClientID)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*20)
	defer cancelFunc()

	fmt.Println("Starting to mint ZCN")

	tx, err := b.MintZCN(ctx, payload)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Println("Completed ZCN mint transaction")
	fmt.Printf("Transaction hash: %s\n", tx.Hash)

	fmt.Println("Done.")
}
