package cmd

import (
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var getVestingPoolConfigCmd = &cobra.Command{
	Use:   "vp-config",
	Short: "Check out vesting pool configurations.",
	Long:  `Check out vesting pool configurations.`,
	Args:  cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			fields = new(zcncore.InputMap)
			cb     = NewJSONInfoCB(fields)
			err    error
		)
		if err = zcncore.GetVestingSCConfig(cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		printMap(fields.Fields)
	},
}

func init() {
	rootCmd.AddCommand(getVestingPoolConfigCmd)
}
