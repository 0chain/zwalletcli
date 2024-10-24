package cmd

import (
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	command := createCommand(
		"bridge-verify",
		"verify ethereum transaction ",
		`verify transaction.
					<hash>`,
		VerifyEthereumTransaction,
		false,
		WithHash("Ethereum transaction hash"),
	)

	rootCmd.AddCommand(command)
}

func VerifyEthereumTransaction(args ...*Arg) {
	hash := GetHash(args)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 60, time.Second)
	if err != nil {
		ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("\nTransaction verification success: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
	}
}
