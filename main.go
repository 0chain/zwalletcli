package main

import (
	"log"
	"os"

	"github.com/0chain/zwalletcli/cmd"
)

var VersionStr string

func main() {
	cmd.VersionStr = VersionStr
	log.SetOutput(os.Stdout)
	log.SetFlags(0)
	cmd.Execute()
}
