package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var updateAuthorizerConfigCmd = &cobra.Command{
	Use:    "bridge-auth-config-update",
	Short:  "Update ZCNSC authorizer settings by ID",
	Long:   `Update ZCNSC authorizer settings by ID.`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		const (
			IDFlag  = "id"
			FeeFlag = "fee"
			URLFlag = "url"
		)

		var (
			flags      = cmd.Flags()
			ID         string
			Fee        string
			FeeBalance int64
			URL        string
			hash       string
			err        error
		)

		if flags.Changed(IDFlag) {
			if ID, err = flags.GetString(IDFlag); err != nil {
				log.Fatalf("error in '%s' flag: %v", IDFlag, err)
			}
		}

		if flags.Changed(FeeFlag) {
			if Fee, err = flags.GetString(FeeFlag); err != nil {
				log.Fatalf("error in '%s' flag: %v", FeeFlag, err)
			}
		}

		FeeBalance, err = strconv.ParseInt(Fee, 10, 64)
		if err != nil {
			log.Fatalf("error in '%s' flag: %v", FeeFlag, err)
		}

		if flags.Changed(URLFlag) {
			if URL, err = flags.GetString(URLFlag); err != nil {
				log.Fatalf("error in '%s' flag: %v", FeeFlag, err)
			}
		}

		node := &zcncore.AuthorizerNode{
			ID:  ID,
			URL: URL,
			Config: &zcncore.AuthorizerConfig{
				Fee: common.Balance(FeeBalance),
			},
		}

		if hash, _, _, _, err = zcncore.ZCNSCUpdateAuthorizerConfig(node); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("global settings updated\nHash: %v\n", hash)
	},
}

//goland:noinspection GoUnhandledErrorResult
func init() {
	cmd := updateAuthorizerConfigCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String("id", "", "authorizer ID")
	cmd.MarkFlagRequired("id")
}
