package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"time"
)

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			VerifyEthereumTransaction,
			"bridge-verify",
			"verify ethereum transaction ",
			`verify transaction.
	        <hash>`,
		))
}

func VerifyEthereumTransaction(_ *zcnbridge.Bridge, hash string) {
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
