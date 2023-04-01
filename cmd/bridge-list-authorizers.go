package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

// listAuthorizers prints all authorizers
var listAuthorizers = &cobra.Command{
	Use:   "bridge-list-auth",
	Short: "List authorizers",
	Long:  `List available authorizers registered in 0Chain defined in config`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(*cobra.Command, []string) {
		var (
			response = new(zcnbridge.AuthorizerNodesResponse)
			cb       = NewJSONInfoCB(response)
			err      error
		)
		if err = zcnbridge.GetAuthorizers(true, cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}
		if len(response.Nodes) == 0 {
			fmt.Println("no response found")
			return
		}

		util.PrettyPrintJSON(response.Nodes)
	},
}

func init() {
	rootCmd.AddCommand(listAuthorizers)
}
