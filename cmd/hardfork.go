package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk_common/zcncore"
	"github.com/spf13/cobra"
)

var addHardForkCmd = &cobra.Command{
	Use:    "add-hardfork",
	Short:  "Add hardfork",
	Long:   `Add hardfork`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			hash string
			err  error
		)

		input := new(zcncore.InputMap)
		input.Fields = setupInputMap(cmd.Flags(), "names", "rounds")
		if err != nil {
			log.Fatal(err)
		}

		if hash, _, _, _, err = zcncore.AddHardfork(input); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("storagesc smart contract settings updated\nHash: %v\n", hash)
	},
}

func init() {
	addHardForkCmd.PersistentFlags().StringSliceP("names", "n", nil, "list of names")
	addHardForkCmd.PersistentFlags().StringSliceP("rounds", "r", nil, "list of rounds")

	rootCmd.AddCommand(addHardForkCmd)

}
