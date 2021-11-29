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
			"bridge-burn-eth",
			"burn eth tokens",
			"burn eth tokens that will be minted on ZCN chain",
			commandBurnEth,
			amountOption,
		))
}

func commandBurnEth(b *zcnbridge.Bridge, args ...*Arg) {
	amount := GetAmount(args)

	// Increase Allowance

	// Example: https://ropsten.etherscan.io/tx/0xa28266fb44cfc2aa27b26bd94e268e40d065a05b1a8e6339865f826557ff9f0e
	transaction, err := b.IncreaseBurnerAllowance(context.Background(), zcnbridge.Wei(amount))
	if err != nil {
		ExitWithError("failed to execute IncreaseBurnerAllowance")
	}

	hash := transaction.Hash().Hex()
	res, err := zcnbridge.ConfirmEthereumTransaction(hash, 60, time.Second)
	if err != nil {
		ExitWithError(fmt.Sprintf("failed to confirm transaction ConfirmEthereumTransaction hash = %s, error = %v", hash, err))
	}

	if res == 0 {
		ExitWithError(fmt.Sprintf("failed to confirm transaction: %s, status = failed", transaction.Hash().String()))
	}

	// Burn Eth

	fmt.Printf("Starting burn transaction in Ethereum")
	transaction, err = b.BurnWZCN(context.Background(), amount)
	if err != nil {
		ExitWithError(err)
	}
	hash = transaction.Hash().String()
	fmt.Printf("Submitted burn transaction %s\n", hash)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 5, time.Second)
	if err != nil {
		ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("\nTransaction verification success: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
	}
}
