package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var getFaucetConfigCmd = &cobra.Command{
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
		if err = zcncore.GetFaucetSCConfig(cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		for key, value := range fields.Fields {
			fmt.Println(key, value)
		}
	},
}

func init() {
	rootCmd.AddCommand(getFaucetConfigCmd)
}
