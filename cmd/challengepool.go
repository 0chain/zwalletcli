package cmd

import (
	"log"
	"os"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var getChallengePoolStatCmd = &cobra.Command{
	Use:   "getchallengelockedtokens",
	Short: "Get info about a challenge pool",
	Long:  `Get info about a challenge pool`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var flags = cmd.Flags()
		if flags.Changed("allocation_id") == false {
			log.Fatal("error: allocation_id flag is missing")
		}

		var allocID, err = flags.GetString("allocation_id")
		if err != nil {
			log.Fatal("error: invalid allocation id:", err)
		}

		var (
			wg        sync.WaitGroup
			statusBar = &ZCNStatus{wg: &wg}
		)
		wg.Add(1)
		if err := zcncore.GetChallengePoolStat(statusBar, allocID); err != nil {
			log.Fatal(err)
		}
		wg.Wait()
		if statusBar.success {
			log.Printf("\nChallenge pool info:\n %s\n", statusBar.errMsg)
			return
		}
		log.Fatalf("\nFailed to get info. %s\n", statusBar.errMsg)
	},
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(0)

	rootCmd.AddCommand(getChallengePoolStatCmd)

	getChallengePoolStatCmd.PersistentFlags().String("allocation_id", "",
		"allocation identifier")
	getChallengePoolStatCmd.MarkFlagRequired("blobber_id")
}
