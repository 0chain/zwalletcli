package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/0chain/zwalletcli/cmd"
)

var VersionStr string
var MinTxFee string

func main() {
	if MinTxFee != "" {
		fee, err := strconv.ParseFloat(MinTxFee, 10)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid min transaction fee, expecting float: %v", err)
			os.Exit(1)
		}
		cmd.MinTxFee = fee
	}

	cmd.VersionStr = VersionStr
	cmd.Execute()
}
