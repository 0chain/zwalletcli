package cmd

import (
	"github.com/0chain/gosdk/core/transaction"
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
			res = &transaction.InputMap{}
			err error
		)
		if res, err = transaction.GetConfig("storage_sc_config"); err != nil {
			log.Fatal(err)
		}

		printMap(res.Fields)
	},
}

func init() {
	rootCmd.AddCommand(scConfig)
	scConfig.Flags().Bool("json", false, "pass this option to print response as json data")
}
