package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
)

func init() {
	rootCmd.AddCommand(
		createBridgeCommand(
			commandGetZCNBurnTicket,
			"bridge-get-zcn-burn",
			"get confirmed burn ticket for zcn burn transaction",
			"get confirmed burn ticket for zcn burn transaction",
		))
}

func commandGetZCNBurnTicket(b *zcnbridge.Bridge, hash string) {
	payload, err := b.QueryZChainMintPayload(hash)
	if err != nil {
		ExitWithError(err)
	}

	fmt.Println(payload.Nonce)
	fmt.Println(payload.Amount)
	fmt.Println(payload.EthereumTxnID)
	fmt.Println(payload.ReceivingClientID)
}
