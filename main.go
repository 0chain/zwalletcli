package main

import (
	"fmt"
	"github.com/0chain/zwalletcli/cmd"
	"os"
	"strconv"
)

var VersionStr string
var MinTxFee string

func main() {
	// TODO: stop throwing error, capture it from blockchain.
	if MinTxFee == "" {
		fmt.Fprintf(os.Stderr, "need a min tx fee")
		os.Exit(1)
	}

	fee, err := strconv.ParseUint(MinTxFee, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid min transaction fee, expecting a non negative integer: %v", err)
		os.Exit(1)
	}

	cmd.MinTxFee = fee
	cmd.VersionStr = VersionStr
	cmd.Execute()
}
