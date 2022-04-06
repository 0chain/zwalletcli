package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/0chain/gosdk/core/zcncrypto"
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
		if !fflags.Changed("mnemonic") {
			ExitWithError("Error: Mnemonic not provided")
		}

		mnemonic := cmd.Flag("mnemonic").Value.String()
		if !zcncore.IsMnemonicValid(mnemonic) {
			ExitWithError("Error: Invalid mnemonic")
		}

		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)

		skipWalletRegistration, err := cmd.Flags().GetBool("offline")
		if err != nil {
			fmt.Println(err.Error())
		}

		if (!skipWalletRegistration) {
			err := zcncore.RecoverWallet(mnemonic, statusBar)

			if err == nil {
				wg.Wait()
			} else {
				ExitWithError(err.Error())
			}
		} else {
			sigScheme := zcncrypto.NewSignatureScheme(SignScheme)

			wallet, err := sigScheme.RecoverKeys(mnemonic)
			if err != nil {
				ExitWithError(err.Error())
			}

			w, err := wallet.Marshal()
			if err != nil {
				ExitWithError(err.Error())
			}

			statusSuccess := 0
			statusBar.OnWalletCreateComplete(statusSuccess, w, "")
		}
		
		if len(statusBar.walletString) == 0 || !statusBar.success {
			ExitWithError("Error recovering the wallet." + statusBar.errMsg)
		}

		var walletFilePath string
		if &walletFile != nil && len(walletFile) > 0 {
			walletFilePath = getConfigDir() + "/" + walletFile
		} else {
			walletFilePath = getConfigDir() + "/wallet.json"
		}

		clientConfig = string(statusBar.walletString)
		file, err := os.Create(walletFilePath)
		if err != nil {
			ExitWithError(err.Error())
		}

		defer file.Close()
		fmt.Fprint(file, clientConfig)
		fmt.Println("Wallet recovered!!")
	},
}

func init() {
	rootCmd.AddCommand(recoverwalletcmd)
	recoverwalletcmd.PersistentFlags().String("mnemonic", "", "mnemonic")
	recoverwalletcmd.PersistentFlags().Bool("offline", false, "if true, wallet won't be registered to miners again")
	recoverwalletcmd.MarkFlagRequired("mnemonic")
}
