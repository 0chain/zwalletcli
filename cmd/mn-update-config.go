package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var updateMinerScConfigCmd = &cobra.Command{
	Use:    "mn-update-config",
	Short:  "Update the miner smart contract",
	Long:   `Update the miner smart contract.`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			hash string
			err  error
		)

		input := new(zcncore.InputMap)
		input.Fields = setupInputMap(cmd.Flags(), "keys", "values")
		if err != nil {
			log.Fatal(err)
		}

		if hash, _, _, _, err = zcncore.MinerScUpdateConfig(input); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("minersc smart contract settings updated\nHash: %v\n", hash)
	},
}

func init() {
	rootCmd.AddCommand(updateMinerScConfigCmd)
	updateMinerScConfigCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateMinerScConfigCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
