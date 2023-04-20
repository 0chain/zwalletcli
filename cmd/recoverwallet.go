package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var recoverwalletcmd = &cobra.Command{
	Use:   "recoverwallet",
	Short: "Recover wallet",
	Long:  `Recover wallet from the previously stored mnemonic`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("mnemonic") == false {
			ExitWithError("Error: Mnemonic not provided")
		}
		mnemonic := cmd.Flag("mnemonic").Value.String()
		if zcncore.IsMnemonicValid(mnemonic) == false {
			ExitWithError("Error: Invalid mnemonic")
		}

		offline, err := cmd.Flags().GetBool("offline")
		if err != nil {
			fmt.Println("offline is not used or not set to true. Setting it to false")
		}

		var walletString string
		if offline {
			walletString, err = zcncore.RecoverOfflineWallet(mnemonic)
			if err != nil {
				ExitWithError(err.Error())
			}
		} else {
			initZCNCoreContext()
			wg := &sync.WaitGroup{}
			statusBar := &ZCNStatus{wg: wg}
			wg.Add(1)
			err = zcncore.RecoverWallet(mnemonic, statusBar)
			if err == nil {
				wg.Wait()
			} else {
				ExitWithError(err.Error())
			}
			if len(statusBar.walletString) == 0 || !statusBar.success {
				ExitWithError("Error recovering the wallet." + statusBar.errMsg)
			}

			walletString = statusBar.walletString
		}

		var walletFilePath string
		if len(walletFile) > 0 {
			walletFilePath = getConfigDir() + "/" + walletFile
		} else {
			walletFilePath = getConfigDir() + "/wallet.json"
		}
		clientConfig = string(walletString)
		file, err := os.Create(walletFilePath)
		if err != nil {
			ExitWithError(err.Error())
		}
		defer file.Close()
		fmt.Fprintf(file, clientConfig)
		fmt.Println("Wallet recovered!!")
		return
	},
}

func init() {
	rootCmd.AddCommand(WithoutZCNCore(WithoutWallet(recoverwalletcmd)))
	recoverwalletcmd.PersistentFlags().String("mnemonic string", "", "mnemonic")
	recoverwalletcmd.PersistentFlags().Bool("offline boolean", false, "recover wallet without registration on blockchain")
	recoverwalletcmd.MarkFlagRequired("mnemonic")
}
