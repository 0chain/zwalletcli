package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"sync"
)

var minerscDeleteNode = &cobra.Command{
	Use:   "mn-delete-node",
	Short: "Delete a miner or sharder node from Miner SC.",
	Long:  "Delete a miner or sharder node from Miner SC.",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var (
			flags   = cmd.Flags()
			id      string
			sharder bool
			err     error
		)

		if !flags.Changed("id") {
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if sharder, err = flags.GetBool("sharder"); err != nil {
			log.Fatal(err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)

		// remove not settings fields
		miner := &zcncore.MinerSCMinerInfo{SimpleMinerSCMinerInfo: &zcncore.SimpleMinerSCMinerInfo{
			ID: id,
		}}

		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			log.Fatal(err)
		}
		wg.Add(1)
		if sharder {
			err = txn.MinerSCDeleteSharder(miner)
		} else {
			err = txn.MinerSCDeleteMiner(miner)
		}
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

		if !statusBar.success {
			log.Fatal("fatal:", statusBar.errMsg)
		}

		if sharder {
			fmt.Printf("sharder %s sussesfully removed from the network", id)
		} else {
			fmt.Printf("minner %s sussesfully removed from the network", id)
		}

	},
}

func init() {
	rootCmd.AddCommand(minerscDeleteNode)
	minerscDeleteNode.PersistentFlags().String("id", "", "miner/sharder ID of node to delete")
	minerscDeleteNode.PersistentFlags().Bool("sharder", false, "set for true if you delete sharder")
	minerscDeleteNode.MarkFlagRequired("id")
}
