package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var providerRegister = &cobra.Command{
	Use:   "vc-add",
	Short: "add node to view change",
	Long:  "add node to view change, add a miner or sharder to the register list so that they can join the MB. Only chainowner can do this.",
	// Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			id    string
			err   error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if !flags.Changed("provider-type") {
			log.Fatal("missing provider-type flag")
		}

		nodeType, err := flags.GetString("provider-type")
		if err != nil {
			log.Fatal(err)
		}

		var pt sdk.ProviderType
		switch nodeType {
		case "miner":
			pt = sdk.ProviderMiner
		case "sharder":
			pt = sdk.ProviderSharder
		default:
			log.Fatalf("unknown provider type: %v", nodeType)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		txn, err := zcncore.NewTransaction(statusBar, 0, nonce)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		err = txn.MinerSCVCAdd(id, pt)
		if err != nil {
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
				fmt.Println("vc add: ", id)
			default:
				ExitWithError("\nvc add " + id + " failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
			return
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}
	},
}

func init() {
	rootCmd.AddCommand(providerRegister)
	providerRegister.PersistentFlags().String("id", "", "provider ID to add to view change")
	_ = providerRegister.MarkFlagRequired("id")

	providerRegister.PersistentFlags().String("provider-type", "", "provider type")
	_ = providerRegister.MarkFlagRequired("provider-type")

}
