package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"time"
)

const (
	USE   = "bridge-verify"
	SHORT = "verify ethereum transaction "
	LONG  = `verify transaction.
	        <hash>`
)

func ConfirmEthereumTransaction(hash string) {
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

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			USE,
			SHORT,
			LONG,
			ConfirmEthereumTransaction,
		))
}
