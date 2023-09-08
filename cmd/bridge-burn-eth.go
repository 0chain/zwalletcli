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
			WithAmount("WZCN token amount to be burned"),
			WithRetries("Num of seconds a transaction status check should run"),
		))
}

func commandBurnEth(b *zcnbridge.BridgeClient, args ...*Arg) {
	retries := GetRetries(args)
	amount := GetAmount(args)

	fmt.Println("Starting IncreaseBurnerAllowance transaction")
	transaction, err := b.IncreaseBurnerAllowance(context.Background(), zcnbridge.Wei(amount))
	if err != nil {
		ExitWithError(err, "failed to execute IncreaseBurnerAllowance")
	}

	hash := transaction.Hash().Hex()
	status, err := zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
	if err != nil {
		ExitWithError(fmt.Sprintf("Failed to confirm IncreaseBurnerAllowance: hash = %s, error = %v", hash, err))
	}

	if status == 1 {
		fmt.Printf("Verification: IncreaseBurnerAllowance [OK]: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("Verification: IncreaseBurnerAllowance [FAILED]: %s\n", hash))
	}

	// Burn Eth

	fmt.Println("Starting WZCN burn transaction")
	transaction, err = b.BurnWZCN(context.Background(), amount)
	if err != nil {
		ExitWithError(err, "failed to burn WZCN tokens")
	}
	hash = transaction.Hash().String()
	fmt.Printf("Confirming WZCN burn transaction %s\n", hash)

	status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
	if err != nil {
		ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("Verification: WZCN burn [OK]: %s\n", hash)
	}

	if status == 0 {
		ExitWithError(fmt.Sprintf("Verification: WZCN burn [PENDING]: %s\n", hash))
	}

	if status == -1 {
		ExitWithError(fmt.Sprintf("Verification: WZCN burn not started, please, check later [FAILED]: %s\n", hash))
	}
}
