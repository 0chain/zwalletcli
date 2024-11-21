package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var minerKill = &cobra.Command{
	Use:   "mn-kill",
	Short: "kill miner",
	Long:  "kill miner",
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

		_, _, _, _, err = zcncore.MinerSCKill(id, zcncore.ProviderMiner)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("killed :", id)
	},
}

func init() {
	rootCmd.AddCommand(minerKill)
	minerKill.PersistentFlags().String("id", "", "miner ID to update")
	_ = minerKill.MarkFlagRequired("id")

}
