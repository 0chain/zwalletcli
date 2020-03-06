package cmd

import (
	"log"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var storageSCConfigCmd = &cobra.Command{
	Use:   "storageconfig",
	Short: "Get storage SC configurations",
	Long:  `Get storage SC configurations, including read/write pool configs`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		if err := zcncore.GetStorageSCConfig(statusBar); err != nil {
			log.Fatal(err)
		}
		wg.Wait()
		if statusBar.success {
			log.Printf("\nWrite pool configurations:\n %s\n", statusBar.errMsg)
			return
		}
		log.Fatalf("\nFailed to get configurations. %s\n", statusBar.errMsg)
	},
}

func init() {
	rootCmd.AddCommand(storageSCConfigCmd)
}
