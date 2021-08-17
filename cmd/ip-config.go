package cmd

import (
	"fmt"
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

		fmt.Println("min_lock:", fields.Fields["min_lock"])
		fmt.Println("max_mint:", fields.Fields["max_mint"])
		fmt.Println("min_lock_period:", fields.Fields["min_lock_period"])
		fmt.Println("apr:", fields.Fields["apr"])
	},
}

func init() {
	rootCmd.AddCommand(getInterestPoolConfigCmd)
}
