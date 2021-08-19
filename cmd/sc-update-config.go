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

var updateStoragScConfigCmd = &cobra.Command{
	Use:   "sc-update-config",
	Short: "Update the Faucet smart contract",
	Long:  `Update the Faucet smart contract.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			err   error

			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)

		var keys []string
		if flags.Changed("keys") {
			keys, err = flags.GetStringSlice("keys")
			if err != nil {
				log.Fatal(err)
			}
		}

		var values []string
		if flags.Changed("values") {
			values, err = flags.GetStringSlice("values")
			if err != nil {
				log.Fatal(err)
			}
		}

		var input = new(zcncore.InputMap)
		input.Fields = make(map[string]interface{})
		if len(keys) != len(values) {
			log.Fatal("number keys must equal the number values")
		}
		for i := 0; i < len(keys); i++ {
			v := strings.TrimSpace(values[i])
			k := strings.TrimSpace(keys[i])
			switch v {
			case "true":
				input.Fields[k], err = strconv.ParseBool(v)
			case "false":
				input.Fields[k], err = strconv.ParseBool(v)
			default:
				input.Fields[k], err = strconv.ParseFloat(v, 64)
			}
			if err != nil {
				log.Fatal(values[i] + "cannot be converted to boolean or numeric value")
			}
		}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if err = txn.StorageScUpdateConfig(input); err != nil {
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

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		fmt.Printf("storagesc smart contract settings updated\nHash: %v\n", txn.GetTransactionHash())
	},
}

func init() {
	rootCmd.AddCommand(updateStoragScConfigCmd)
	updateStoragScConfigCmd.PersistentFlags().StringSlice("keys", nil, "list of keys")
	updateStoragScConfigCmd.PersistentFlags().StringSlice("values", nil, "list of new values")
}
