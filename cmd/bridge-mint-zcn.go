package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-mint-zcn",
			"mint zcn tokens using the hash of Ethereum burn transaction",
			"mint zcn tokens after burning WZCN tokens in Ethereum chain",
			commandMintZCN,
			false,
			&Option{
				name:     "burn-txn-hash",
				typename: "string",
				value:    "",
				usage:    "mint the ZCN tokens for the given Ethereum burn transaction hash",
			},
		))
}

func commandMintZCN(b *zcnbridge.BridgeClient, args ...*Arg) {
	burnHash := getString(args, "burn-txn-hash")

	var mintNonce int64
	res, err := zcncore.GetMintNonce()
	if err != nil {
		ExitWithError(err)
	}

	err = json.Unmarshal(res, &mintNonce)
	if err != nil {
		ExitWithError(err)
	}

	burnTickets, err := b.QueryEthereumBurnEvents(strconv.Itoa(int(mintNonce)))
	if err != nil {
		ExitWithError(err)
	}

	fmt.Printf("Found %d not processed WZCN burn transactions\n", len(burnTickets))

	for _, burnTicket := range burnTickets {
		if len(burnHash) > 0 {
			if burnHash != burnTicket.TransactionHash {
				continue
			}
		}

		fmt.Printf("Query ticket for Ethereum transaction hash: %s\n", burnTicket.TransactionHash)

		payload, err := b.QueryZChainMintPayload(burnTicket.TransactionHash)
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

		txHash, err := b.MintZCN(ctx, payload)
		if err != nil {
			ExitWithError(err)
		}

		fmt.Println("Completed ZCN mint transaction")
		fmt.Printf("Transaction hash: %s\n", txHash)

	}

	if len(burnTickets) > 0 {
		fmt.Println("Done.")
	} else {
		fmt.Println("Failed.")
	}
}
