package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"time"
)

func init() {
	command := createBridgeCommand(
		"bridge-verify",
		"verify ethereum transaction ",
		`verify transaction.
					<hash>`,
		VerifyEthereumTransaction,
		hashOption,
	)

	rootCmd.AddCommand(command)
}

func VerifyEthereumTransaction(_ *zcnbridge.BridgeClient, args ...*Arg) {
	hash := GetHash(args)

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
