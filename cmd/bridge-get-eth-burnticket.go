package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			commandGetETHBurnTicket,
			"bridge-get-eth-burn",
			"get confirmed burn ticket for ethereum burn transaction",
			"get confirmed burn ticket for ethereum burn transaction",
		))
}

func commandGetETHBurnTicket(b *zcnbridge.Bridge, hash string) {
	payload, err := b.QueryEthereumMintPayload(hash)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Println(payload.Nonce)
	fmt.Println(payload.Amount)
	fmt.Println(payload.ZCNTxnID)
}
