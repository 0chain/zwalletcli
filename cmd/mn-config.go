package cmd

import (
	"log"

	"github.com/0chain/gosdk_common/core/transaction"

	"github.com/spf13/cobra"
)

var mnConfigCmd = &cobra.Command{
	Use:   "mn-config",
	Short: "Show miner SC configuration.",
	Long:  `Show miner SC configuration.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			fields = new(transaction.InputMap)
			err    error
		)
		if fields, err = transaction.GetConfig("miner_sc_configs"); err != nil {
			log.Fatal(err)
		}

		printMap(fields.Fields)
	},
}

func init() {
	rootCmd.AddCommand(mnConfigCmd)
}
