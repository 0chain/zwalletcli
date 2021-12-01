package cmd

import (
	"context"
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"time"
)

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			"bridge-mint-eth",
			"mint WZCN tokens using the hash of ZCN burn transaction",
			"mint WZCN tokens after burning ZCN tokens in ZCN chain",
			commandMintEth,
			hashOption,
		))
}

func commandMintEth(b *zcnbridge.Bridge, args ...*Arg) {
	hash := GetHash(args)

	fmt.Printf("Query ticket for ZCN transaction hash: %s\n", hash)

	payload, err := b.QueryEthereumMintPayload(hash)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Printf("Sending mint transaction to Ethereum\n")
	fmt.Printf("Payload amount: %d\n", payload.Amount)
	fmt.Printf("Payload nonce: %d\n", payload.Nonce)
	fmt.Printf("ZCN transaction ID: %s\n", payload.ZCNTxnID)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*20)
	defer cancelFunc()

	tx, err := b.MintWZCN(ctx, payload)
	if err != nil {
		ExitWithError(err)
	}

	hash = tx.Hash().String()
	fmt.Printf("Confirming Ethereum mint transaction: %s\n", hash)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 5, time.Second)
	if err != nil {
		ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("\nTransaction verification success: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
	}

	fmt.Println("Done.")
}
