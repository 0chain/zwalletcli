package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var listAuthorizers = &cobra.Command{
	Use:   "bridge-list-auth",
	Short: "List authorizers",
	Long:  `List available authorizers registered in 0Chain defined in config`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(*cobra.Command, []string) {
		authorizers, err := zcnbridge.GetAuthorizers()
		if err != nil || authorizers == nil || len(authorizers) == 0 {
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
