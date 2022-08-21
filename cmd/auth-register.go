package cmd

import (
	"log"

	// "github.com/0chain/gosdk/zcnbridge"
	// "github.com/0chain/zwalletcli/util"
	"context"

	"github.com/0chain/gosdk/zcnbridge/transaction"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

type AddAuthorizerPayload struct {
	// PublicKey         string                      `json:"public_key"`
	URL string `json:"url"`
	// StakePoolSettings AuthorizerStakePoolSettings `json:"stake_pool_settings"` // Used to initially create stake pool
}

var getAuthorizerRegisterCmd = &cobra.Command{
	Use:   "auth-register",
	Short: "Register an authorizer manually",
	Long:  `Register an authorizer manually.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			ID    string
			URL   string
			err   error
		)

		if flags.Changed("id") {
			if ID, err = flags.GetString("id"); err != nil {
				log.Fatalf("error in 'id' flag: %v", err)
			}
		} else {
			ExitWithError("Error: id flag is missing")
		}

		if flags.Changed("url") {
			if URL, err = flags.GetString("url"); err != nil {
				log.Fatalf("error in 'url' flag: %v", err)
			}
		} else {
			ExitWithError("Error: url flag is missing")
		}

		registerAuthorizerInChain(ID, URL)
		// var (
		// 	response = new(zcnbridge.AuthorizerResponse)
		// 	cb       = NewJSONInfoCB(response)
		// )

		// if err = zcnbridge.GetAuthorizer(ID, cb); err != nil {
		// 	log.Fatal(err)
		// }
		// if err = cb.Waiting(); err != nil {
		// 	log.Fatal(err)
		// }

		// util.PrettyPrintJSON(response)
	},
}

//goland:noinspection ALL
func init() {
	cmd := getAuthorizerRegisterCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String("id", "", "authorizer id")
	cmd.PersistentFlags().String("url", "", "authorizer endpoint url")
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("url")
}

func registerAuthorizerInChain(ID, URL string) {

	input := &zcncore.AddAuthorizerPayload{
		// PublicKey: authorizer.App.PublicKey,
		URL: URL,
		// StakePoolSettings: zcncore.AuthorizerStakePoolSettings{
		// 	DelegateWallet: authorizer.App.ID,
		// 	MinStake:       authorizer.App.MinStake,
		// 	MaxStake:       authorizer.App.MaxStake,
		// 	NumDelegates:   authorizer.App.NumDelegates,
		// 	ServiceCharge:  authorizer.App.ServiceCharge,
		// },
	}

	trx, err := transaction.AddAuthorizer(context.Background(), input)
	if err != nil {
		log.Fatal(err, "failed to add authorizer with transaction: '%s'", trx.Hash)
	}
}
