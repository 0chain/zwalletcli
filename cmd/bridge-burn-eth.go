package cmd

import (
	"context"
	"fmt"
	"github.com/0chain/common/core/currency"
	"github.com/ethereum/go-ethereum/core/types"
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

	tokenBalance, err := b.GetTokenBalance()
	if err != nil {
		ExitWithError(err, "failed to retrieve current token balance")
	}

	tokenBalanceZCN, err := currency.Coin(tokenBalance.Int64()).ToZCN()
	if err != nil {
		ExitWithError(err, "failed to convert current token balance to ZCN")
	}

	var (
		transaction *types.Transaction
		hash        string
		status      int
	)

	if tokenBalanceZCN < float64(amount) {
		transaction, err = b.Swap(context.Background(), amount, time.Now().Add(time.Minute*3))
		if err != nil {
			ExitWithError(err, "failed to execute Swap")
		}

		hash = transaction.Hash().Hex()
		status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
		if err != nil {
			ExitWithError(fmt.Sprintf("Failed to confirm Swap: hash = %s, error = %v", hash, err))
		}

		if status == 1 {
			fmt.Printf("Verification: Swap [OK]: %s\n", hash)
		} else {
			ExitWithError(fmt.Sprintf("Verification: Swap [FAILED]: %s\n", hash))
		}
	}

	fmt.Println("Starting IncreaseBurnerAllowance transaction")
	transaction, err = b.IncreaseBurnerAllowance(context.Background(), zcnbridge.Wei(amount))
	if err != nil {
		ExitWithError(err, "failed to execute IncreaseBurnerAllowance")
	}

	hash = transaction.Hash().Hex()
	status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
	if err != nil {
		ExitWithError(fmt.Sprintf("Failed to confirm IncreaseBurnerAllowance: hash = %s, error = %v", hash, err))
	}

	if status == 1 {
		fmt.Printf("Verification: IncreaseBurnerAllowance [OK]: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("Verification: IncreaseBurnerAllowance [FAILED]: %s\n", hash))
	}

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
