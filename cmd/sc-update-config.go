package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var updateStoragScConfigCmd = &cobra.Command{
	Use:    "sc-update-config",
	Short:  "Update the storage smart contract",
	Long:   `Update the storage smart contract.`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		input := new(zcncore.InputMap)
		input.Fields = setupInputMap(cmd.Flags(), "keys", "values")
		if err != nil {
			log.Fatal(err)
		}

		hash, _, _, _, err := zcncore.StorageScUpdateConfig(input)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("storagesc smart contract settings updated\nHash: %v\n", hash)
	},
}

func init() {
	rootCmd.AddCommand(updateStoragScConfigCmd)
	updateStoragScConfigCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateStoragScConfigCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
