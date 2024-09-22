package main

import (
	"github.com/0chain/zwalletcli/cmd"
	"log"
)

var VersionStr string

func main() {
	cmd.VersionStr = VersionStr
	log.SetFlags(0)
	cmd.Execute()
}
