package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
	"sync"
)

var updateAuthorizerConfigCmd = &cobra.Command{
	Use:   "bridge-auth-config-update",
	Short: "Update ZCNSC authorizer settings by ID",
	Long:  `Update ZCNSC authorizer settings by ID.`,
	Args:  cobra.MinimumNArgs(0),
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

		var wg sync.WaitGroup
		statusBar := &ZCNStatus{wg: &wg}
		txn, err := zcncore.NewTransaction(statusBar, transactionFee(), nonce)
		if err != nil {
			log.Fatal(err)
		}

		wg.Add(1)
		if err = txn.ZCNSCUpdateAuthorizerConfig(node); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		statusBar.success = false
		wg.Add(1)
		if err = txn.Verify(); err != nil {
			log.Fatal(err)
		}
		wg.Wait()

		if statusBar.success {
			//fmt.Printf("Nonce:%v\n", txn.GetTransactionNonce())
			switch txn.GetVerifyConfirmationStatus() {
			case zcncore.ChargeableError:
				ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
			case zcncore.Success:
				fmt.Printf("global settings updated\nHash: %v\n", txn.GetTransactionHash())
			default:
				ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
			return
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}

	},
}

//goland:noinspection GoUnhandledErrorResult
func init() {
	cmd := updateAuthorizerConfigCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().String("fee", "", "fee")
	cmd.MarkFlagRequired("fee")

	cmd.PersistentFlags().String("id", "", "authorizer ID")
	cmd.MarkFlagRequired("id")
}
