package cmd

import (
	"encoding/json"
	"fmt"
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

		buffer, err := json.MarshalIndent(authorizers, "", "   ")
		if err != nil {
			ExitWithError("\nFailed to unmarshall\n", err)
		}

		fmt.Println(string(buffer))
	},
}

func init() {
	rootCmd.AddCommand(listAuthorizers)
}
