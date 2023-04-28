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

var sharderKill = &cobra.Command{
	Use:   "sh-kill",
	Short: "kill sharder",
	Long:  "kill sharder",
	Args:  cobra.MinimumNArgs(0),
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

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		txn, err := zcncore.NewTransaction(statusBar, 0, nonce)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		err = txn.MinerSCKill(id, zcncore.ProviderSharder)
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
				fmt.Println("killed :", id)
			default:
				ExitWithError("\nkill " + id + " failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
			return
		} else {
			log.Fatal("fatal:", statusBar.errMsg)
		}
	},
}

func init() {
	rootCmd.AddCommand(sharderKill)
	sharderKill.PersistentFlags().String("id", "", "sharder ID to update")
	_ = minerKill.MarkFlagRequired("id")

}
