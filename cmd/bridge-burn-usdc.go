package cmd

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-burn-usdc",
			"burn usdc tokens",
			"burn usdc tokens that will be minted on ZCN chain",
			commandBurnUsdc,
			false,
			WithAmount("WZCN token amount to be burned"),
			WithRetries("Num of seconds a transaction status check should run"),
		))
}

func commandBurnUsdc(b *zcnbridge.BridgeClient, args ...*Arg) {
	retries := GetRetries(args)
	amount := GetAmount(args)

	var (
		transaction *types.Transaction
		err         error
		hash        string
		status      int
	)

	var balanceRaw *big.Int

	balanceRaw, err = b.GetTokenBalance()
	if err != nil {
		ExitWithError(err, "failed to GetTokenBalance")
	}

	balance := balanceRaw.Uint64()

	if balance < amount {
		target := amount - balance

		transaction, err = b.ApproveUSDCSwap(context.Background(), 0)
		if err != nil {
			ExitWithError(err, "failed to execute ApproveUSDCSwap")
		}

		hash = transaction.Hash().Hex()
		status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
		if err != nil {
			ExitWithError(fmt.Sprintf("Failed to confirm ApproveUSDCSwap: hash = %s, error = %v", hash, err))
		}

		if status == 1 {
			fmt.Printf("Verification: ApproveUSDCSwap [OK]: %s\n", hash)
		} else {
			ExitWithError(fmt.Sprintf("Verification: ApproveUSDCSwap [FAILED]: %s\n", hash))
		}

		transaction, err = b.ApproveUSDCSwap(context.Background(), target)
		if err != nil {
			ExitWithError(err, "failed to execute ApproveUSDCSwap")
		}

		hash = transaction.Hash().Hex()
		status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
		if err != nil {
			ExitWithError(fmt.Sprintf("Failed to confirm ApproveUSDCSwap: hash = %s, error = %v", hash, err))
		}

		if status == 1 {
			fmt.Printf("Verification: ApproveUSDCSwap [OK]: %s\n", hash)
		} else {
			ExitWithError(fmt.Sprintf("Verification: ApproveUSDCSwap [FAILED]: %s\n", hash))
		}

		transaction, err = b.SwapUSDC(context.Background(), target, target)
		if err != nil {
			ExitWithError(err, "failed to execute SwapUSDC")
		}

		hash = transaction.Hash().Hex()
		status, err = zcnbridge.ConfirmEthereumTransaction(hash, retries, time.Second)
		if err != nil {
			ExitWithError(fmt.Sprintf("Failed to confirm SwapUSDC: hash = %s, error = %v", hash, err))
		}

		if status == 1 {
			fmt.Printf("Verification: SwapUSDC [OK]: %s\n", hash)
		} else {
			ExitWithError(fmt.Sprintf("Verification: SwapUSDC [FAILED]: %s\n", hash))
		}
	}

	fmt.Println("Starting IncreaseBurnerAllowance transaction")
	transaction, err = b.IncreaseBurnerAllowance(context.Background(), amount)
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
