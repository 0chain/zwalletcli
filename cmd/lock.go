package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var lockconfigcmd = &cobra.Command{
	Use:   "lockconfig",
	Short: "Get lock configuration",
	Long:  `Get lock configuration`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetLockConfig(statusBar)
		if err != nil {
			ExitWithError(err)
		}
		wg.Wait()
		if statusBar.success {
			fmt.Printf("\nConfiguration:\n %v\n", statusBar.errMsg)
			return
		}
		ExitWithError("\nFailed to get lock config." + statusBar.errMsg + "\n")
		return
	},
}

var getlockedtokenscmd = &cobra.Command{
	Use:   "getlockedtokens",
	Short: "Get locked tokens",
	Long:  `Get locked tokens`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetLockedTokens(statusBar)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wg.Wait()
		if statusBar.success {
			fmt.Printf("\nLocked tokens:\n %v\n", statusBar.errMsg)
			return
		}
		ExitWithError("\nFailed to get locked tokens." + statusBar.errMsg + "\n")
	},
}

var lockcmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock tokens",
	Long: `Lock tokens .
	        <tokens> <[durationHr] [durationMin]> [transaction fee]`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("tokens") == false {
			ExitWithError("Error: tokens flag is missing")
		}
		if (fflags.Changed("durationHr") == false) &&
			(fflags.Changed("durationMin") == false) {
			ExitWithError("Error: durationHr and durationMin flag is missing. atleast one is required")
		}
		token, err := cmd.Flags().GetFloat64("tokens")
		if err != nil {
			ExitWithError("Error: invalid number of tokens")
		}
		durationHr := int64(0)
		durationHr, err = cmd.Flags().GetInt64("durationHr")
		durationMin := int(0)
		durationMin, err = cmd.Flags().GetInt("durationMin")
		if (durationHr < 1) && (durationMin < 1) {
			ExitWithError("Error: invalid duration")
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
		err = txn.LockTokens(zcncore.ConvertToValue(token), durationHr, durationMin)
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
				fmt.Printf("\nTokens (%f) locked successfully\n", token)
				return
			}
		}
		ExitWithError("\nFailed to lock tokens. " + statusBar.errMsg + "\n")
	},
}

func init() {
	rootCmd.AddCommand(lockconfigcmd)
	rootCmd.AddCommand(lockcmd)
	rootCmd.AddCommand(getlockedtokenscmd)
	//lockcmd.PersistentFlags().Float64("tokens", 0, "Number to tokens to lock")
	//lockcmd.PersistentFlags().Int64("durationHr", 0, "Duration Hours to lock")
	//lockcmd.PersistentFlags().Int("durationMin", 0, "Duration Mins to lock")
	//lockcmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	//lockcmd.MarkFlagRequired("tokens")
}
