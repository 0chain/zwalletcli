package cmd

import (
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var mnGlobalsCmd = &cobra.Command{
	Use:   "global-config",
	Short: "Show global configurations.",
	Long:  `Show global configurations.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			fields = new(zcncore.InputMap)
			cb     = NewJSONInfoCB(fields)
			err    error
		)
		if err = zcncore.GetMinerSCGlobals(cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		printMap(fields.Fields)
	},
}

func init() {
	rootCmd.AddCommand(mnGlobalsCmd)
}
