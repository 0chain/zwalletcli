package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var sharderDelete = &cobra.Command{
	Use:   "sh-delete",
	Short: "delete sharder",
	Long:  "delete sharder",
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

		hash, _, _, _, err := zcncore.DeleteSharder(id)
		if err != nil {
			ExitWithError("Delete sharder failed : ", err.Error())
		}

		fmt.Println("delete sharder success with transaction hash : ", hash)
	},
}

func init() {
	rootCmd.AddCommand(sharderDelete)
	sharderDelete.PersistentFlags().String("id", "", "sharder ID to delete")
	_ = sharderDelete.MarkFlagRequired("id")

}
