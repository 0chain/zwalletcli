package cmd

import (
	"encoding/json"
	"log"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var getBridgeConfigCmd = &cobra.Command{
	Use:    "bridge-config",
	Short:  "Show ZCNBridge configurations.",
	Long:   `Show ZCNBridge configurations.`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			response = new(zcncore.InputMap)
			err      error
			res      []byte
		)
		if res, err = zcnbridge.GetGlobalConfig(); err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(res, response)
		if err != nil {
			log.Fatal(err)
		}

		printMap(response.Fields)
	},
}

func init() {
	rootCmd.AddCommand(getBridgeConfigCmd)
}
