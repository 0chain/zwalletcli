package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
)

var deleteAuthorizerConfigCmd = &cobra.Command{
	Use:    "bridge-auth-delete",
	Short:  "Delete ZCNSC authorizer by ID",
	Long:   `Delete ZCNSC authorizer by ID`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		const (
			IDFlag = "id"
		)
		var (
			flags = cmd.Flags()
			err   error
			ID    string
			hash  string
		)

		if flags.Changed(IDFlag) {
			if ID, err = flags.GetString(IDFlag); err != nil {
				log.Fatalf("error in '%s' flag: %v", IDFlag, err)
			}
		} else {
			ExitWithError("Error: id flag is missing")
		}

		payload := &zcncore.DeleteAuthorizerPayload{
			ID: ID,
		}
		if hash, _, _, _, err = zcncore.ZCNSCDeleteAuthorizer(payload); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("global settings updated\nHash: %v\n", hash)
	},
}

//goland:noinspection GoUnhandledErrorResult
func init() {
	cmd := deleteAuthorizerConfigCmd
	rootCmd.AddCommand(cmd)
	cmd.PersistentFlags().String("id", "", "authorizer ID")
	cmd.MarkFlagRequired("id")
}
