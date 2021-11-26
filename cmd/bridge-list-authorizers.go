package cmd

import (
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var listAuthorizers = &cobra.Command{
	Use:   "bridge-list-auth",
	Short: "list authorizers",
	Long:  `list available authorizers`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		authorizers, err := zcnbridge.GetAuthorizers()
		if err != nil || authorizers == nil || len(authorizers.NodeMap) == 0 {
			ExitWithError("\nAuthorizers not found\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(listAuthorizers)
}
