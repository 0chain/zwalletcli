package cmd

import (
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var getInterestPoolConfigCmd = &cobra.Command{
	Use:   "ip-config",
	Short: "Show interest pool configurations.",
	Long:  `Show interest pool configurations.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			fields = new(zcncore.InputMap)
			cb     = NewJSONInfoCB(fields)
			err    error
		)
		if err = zcncore.GetInterestPoolSCConfig(cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		printMap(fields.Fields)
	},
}

func init() {
	rootCmd.AddCommand(getInterestPoolConfigCmd)
}
