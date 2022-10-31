package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"strconv"
	"strings"
	"sync"
)

var deleteAuthorizerConfigCmd = &cobra.Command{
	Use:   "bridge-auth-delete",
	Short: "Delete ZCNSC authorizer by ID",
	Long:  `Delete ZCNSC authorizer by ID`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		const (
			IDFlag = "id"
		)
		var (
			flags = cmd.Flags()
			err   error
			ID    string
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
		var wg sync.WaitGroup
		statusBar := &ZCNStatus{wg: &wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(txFee), nonce)
		if err != nil {
			log.Fatal(err)
		}

		if err := txn.AdjustTransactionFee(txVelocity.toZCNFeeType()); err != nil {
			log.Fatal("failed to adjust transaction fee: ", err)
		}

		wg.Add(1)
		if err = txn.ZCNSCDeleteAuthorizer(payload); err != nil {
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
	cmd := deleteAuthorizerConfigCmd
	rootCmd.AddCommand(cmd)
	cmd.PersistentFlags().String("id", "", "authorizer ID")
	cmd.MarkFlagRequired("id")
}
