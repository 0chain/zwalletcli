package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-mint-wzcn",
			"mint WZCN tokens using the hash of ZCN burn transaction",
			"mint WZCN tokens after burning ZCN tokens in ZCN chain",
			commandMintEth,
			WithHash("Ethereum transaction hash"),
		))
}

func commandMintEth(b *zcnbridge.BridgeClient, args ...*Arg) {
	ethereumAddress := GetEthereumAddress(args)
	userNonce, err := b.GetUserNonceMinted(context.Background(), ethereumAddress)
	if err != nil {
		ExitWithError(err)
	}

	var cb zcncore.GetZCNNotProcessedBurnTicketsCallbackStub
	err = zcncore.GetZCNNotProcessedBurnTickets(ethereumAddress, userNonce.Int64(), &cb)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Printf("Found %d not processed ZCN burn transactions", len(cb.Value))
	return

	for _, burnTicket := range cb.Value {
		fmt.Printf("Query ticket for ZCN transaction hash: %s\n", burnTicket.Hash)

		payload, err := b.QueryEthereumMintPayload(burnTicket.Hash)
		if err != nil {
			ExitWithError(err)
		}

		fmt.Printf("Sending mint transaction to Ethereum\n")
		fmt.Printf("Payload amount: %d\n", payload.Amount)
		fmt.Printf("Payload nonce: %d\n", payload.Nonce)
		fmt.Printf("ZCN transaction ID: %s\n", payload.ZCNTxnID)

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*20)
		defer cancelFunc()

		fmt.Println("Starting to mint WZCN")

		tx, err := b.MintWZCN(ctx, payload)
		if err != nil {
			ExitWithError(err)
		}

		hash := tx.Hash().String()
		fmt.Printf("Confirming Ethereum mint transaction: %s\n", hash)

		status, err := zcnbridge.ConfirmEthereumTransaction(hash, 20, time.Second*5)
		if err != nil {
			ExitWithError(err)
		}

		if status == 1 {
			fmt.Printf("\nTransaction verification success: %s\n", hash)
		} else {
			ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
		}
	}

	fmt.Println("Done.")
}
