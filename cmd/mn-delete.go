package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk_common/zcncore"
	"github.com/spf13/cobra"
)

var minerDelete = &cobra.Command{
	Use:   "mn-delete",
	Short: "delete miner",
	Long:  "delete miner",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			id    string
			err   error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		hash, _, _, _, err := zcncore.DeleteMiner(id)
		if err != nil {
			log.Fatal("Delete miner failed : ", err.Error())
		}

		fmt.Println("delete miner success with transaction hash : ", hash)
	},
}

func init() {
	rootCmd.AddCommand(minerDelete)
	minerDelete.PersistentFlags().String("id", "", "miner ID to delete")
	_ = minerDelete.MarkFlagRequired("id")

}
