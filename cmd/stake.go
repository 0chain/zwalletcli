package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var stakecmd = &cobra.Command{
	Use:   "stake",
	Short: "Stake Miners or Sharders",
	Long: `Stake Miners or Sharders using their client ID.
			<client_id> <tokens>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("client_id") == false {
			ExitWithError("Error: client_id flag is missing")
		}
		if fflags.Changed("token") == false {
			ExitWithError("Error: token flag is missing")
		}
		clientID := cmd.Flag("client_id").Value.String()
		token, err := cmd.Flags().GetFloat64("token")
		if err != nil {
			ExitWithError("Error: invalid number of tokens")
		}
		fee := float64(0)
		fee, err = cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee))
		if err != nil {
			ExitWithError(err)
		}
		wg.Add(1)
		err = txn.Stake(clientID, zcncore.ConvertToValue(token))
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}
		if statusBar.success {
			statusBar.success = false
			wg.Add(1)
			err := txn.Verify()
			if err == nil {
				wg.Wait()
			} else {
				ExitWithError(err.Error())
			}
			if statusBar.success {
				fmt.Println("Stake success")
				return
			}
		}
		ExitWithError("Stake failed. " + statusBar.errMsg)
		os.Exit(1)
	},
}

var deletestakecmd = &cobra.Command{
	Use:   "deletestake",
	Short: "Delete Stake from user pool",
	Long: `Delete Stake from user pool client_id and pool_id.
			<client_id> <pool_id>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("client_id") == false {
			ExitWithError("Error: client_id flag is missing")
		}
		if fflags.Changed("pool_id") == false {
			ExitWithError("Error: pool_id flag is missing")
		}
		clientID := cmd.Flag("client_id").Value.String()
		poolID := cmd.Flag("pool_id").Value.String()
		fee := float64(0)
		var err error
		fee, err = cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee))
		if err != nil {
			ExitWithError(err)
		}
		wg.Add(1)
		err = txn.DeleteStake(clientID, poolID)
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}
		if statusBar.success {
			statusBar.success = false
			wg.Add(1)
			err := txn.Verify()
			if err == nil {
				wg.Wait()
			} else {
				ExitWithError(err.Error())

			}
			if statusBar.success {
				fmt.Println("Delete stake success")
				return
			}
		}
		ExitWithError("Delete stake failed. " + statusBar.errMsg)

	},
}

func init() {
	rootCmd.AddCommand(stakecmd)
	rootCmd.AddCommand(deletestakecmd)
	stakecmd.PersistentFlags().String("client_id", "", "Miner or Sharder client id")
	stakecmd.PersistentFlags().Float64("token", 0, "Token to send")
	stakecmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	stakecmd.MarkFlagRequired("client_id")
	stakecmd.MarkFlagRequired("token")
	deletestakecmd.PersistentFlags().String("client_id", "", "Miner or Sharder client id")
	deletestakecmd.PersistentFlags().String("pool_id", "", "Pool ID from user pool matching miner or sharder id")
	deletestakecmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	deletestakecmd.MarkFlagRequired("client_id")
	deletestakecmd.MarkFlagRequired("pool_id")
}
