package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var sharderKill = &cobra.Command{
	Use:   "sh-kill",
	Short: "kill sharder",
	Long:  "kill sharder",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			id    string
			err   error
		)

		if !flags.Changed("id") {
			ExitWithError("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			ExitWithError(err)
		}

		_, _, _, _, err = zcncore.MinerSCKill(id, zcncore.ProviderSharder)
		if err != nil {
			ExitWithError(err)
		}

		fmt.Println("killed :", id)
	},
}

func init() {
	rootCmd.AddCommand(sharderKill)
	sharderKill.PersistentFlags().String("id", "", "sharder ID to update")
	_ = minerKill.MarkFlagRequired("id")

}
