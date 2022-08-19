package cmd

import (
	"log"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

var getAuthorizerConfigCmd = &cobra.Command{
	Use:   "bridge-auth-config",
	Short: "Show authorizer configurations.",
	Long:  `Show authorizer configurations.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			ID    string
			err   error
		)

		if flags.Changed("id") {
			if ID, err = flags.GetString("id"); err != nil {
				log.Fatalf("error in 'id' flag: %v", err)
			}
		} else {
			ExitWithError("Error: id flag is missing")
		}

		var (
			response = new(zcnbridge.AuthorizerResponse)
			cb       = NewJSONInfoCB(response)
		)
		if err = zcnbridge.GetAuthorizer(ID, cb); err != nil {
			log.Fatal(err)
		}
		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		util.PrettyPrintJSON(response)
	},
}

//goland:noinspection ALL
func init() {
	cmd := getAuthorizerConfigCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String("id", "", "authorizer id")
	cmd.MarkFlagRequired("id")
}
