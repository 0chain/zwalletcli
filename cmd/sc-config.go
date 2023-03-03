package cmd

import (
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
	"log"
)

// scConfig shows SC configurations
var scConfig = &cobra.Command{
	Use:   "sc-config",
	Short: "Show storage SC configuration.",
	Long:  `Show storage SC configuration.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			fields = new(zcncore.InputMap)
			cb     = NewJSONInfoCB(fields)
			err    error
		)
		if err = zcncore.GetStorageSCConfig(cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		doJSON, _ := cmd.Flags().GetBool("json")
		if doJSON {
			util.PrintJSON(fields.Fields)
			return
		}
		printMap(fields.Fields)
	},
}

func init() {
	rootCmd.AddCommand(scConfig)
	scConfig.Flags().Bool("json", false, "pass this option to print response as json data")
}
