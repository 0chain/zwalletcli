package cmd

import (
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var getBridgeConfigCmd = &cobra.Command{
	Use:   "bridge-config",
	Short: "Show ZCNBridge configurations.",
	Long:  `Show ZCNBridge configurations.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			response = new(zcncore.InputMap)
			cb       = NewJSONInfoCB(response)
			err      error
		)
		if err = zcnbridge.GetGlobalConfig(cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		printMap(response.Fields)
	},
}

func init() {
	rootCmd.AddCommand(getBridgeConfigCmd)
}
