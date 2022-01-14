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
			"bridge-burn-eth",
			"burn eth tokens",
			"burn eth tokens that will be minted on ZCN chain",
			commandBurnEth,
			amountOption,
		))
}

func commandBurnEth(b *zcnbridge.BridgeClient, args ...*Arg) {
	amount := GetAmount(args)

	// Increase Allowance

	// Example: https://ropsten.etherscan.io/tx/0xa28266fb44cfc2aa27b26bd94e268e40d065a05b1a8e6339865f826557ff9f0e
	transaction, err := b.IncreaseBurnerAllowance(context.Background(), zcnbridge.Wei(amount))
	if err != nil {
		ExitWithError(err, "failed to execute IncreaseBurnerAllowance")
	}

	hash := transaction.Hash().Hex()
	res, err := zcnbridge.ConfirmEthereumTransaction(hash, 60, time.Second)
	if err != nil {
		ExitWithError(fmt.Sprintf("failed to confirm transaction: hash = %s, error = %v", hash, err))
	}

	if res == 0 {
		ExitWithError(fmt.Sprintf("failed to confirm transaction: %s, status = failed", transaction.Hash().String()))
	}

	// Burn Eth

	fmt.Println("Starting WZCN burn transaction")
	transaction, err = b.BurnWZCN(context.Background(), amount)
	if err != nil {
		ExitWithError(err, "failed to burn WZCN tokens")
	}
	hash = transaction.Hash().String()
	fmt.Printf("Confirming WZCN burn transaction %s\n", hash)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 50, time.Second)
	if err != nil {
		ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("Transaction verification success: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("Verification failed: %s\n", hash))
	}
}
