package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcnbridge/wallet"
	"github.com/0chain/gosdk/zcncore"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-mint-wzcn",
			"mint WZCN tokens using the hash of ZCN burn transaction",
			"mint WZCN tokens after burning ZCN tokens in ZCN chain",
			commandMintEth,
		))
}

func commandMintEth(b *zcnbridge.BridgeClient, args ...*Arg) {
	userNonce, err := b.GetUserNonceMinted(context.Background(), b.EthereumAddress)
	if err != nil {
		ExitWithError(err)
	}

	var burnTickets []zcncore.BurnTicket
	cb := wallet.NewZCNStatus(&burnTickets)

	cb.Begin()

	err = zcncore.GetNotProcessedZCNBurnTickets(b.EthereumAddress, userNonce.String(), cb)
	if err != nil {
		ExitWithError(err)
	}

	if err := cb.Wait(); err != nil {
		ExitWithError(err)
	}

	if !cb.Success {
		ExitWithError(cb.Err)
	}

	fmt.Printf("Found %d not processed ZCN burn transactions\n", len(burnTickets))

	for _, burnTicket := range burnTickets {
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

	if len(burnTickets) > 0 {
		fmt.Println("Done.")
	} else {
		fmt.Println("Failed.")
	}
}
