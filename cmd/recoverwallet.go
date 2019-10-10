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
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.RecoverWallet(mnemonic, statusBar)
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}
		if len(statusBar.walletString) == 0 || !statusBar.success {
			ExitWithError("Error recovering the wallet." + statusBar.errMsg)
		}
		var walletFilePath string
		if &walletFile != nil && len(walletFile) > 0 {
			walletFilePath = getConfigDir() + "/" + walletFile
		} else {
			walletFilePath = getConfigDir() + "/wallet.txt"
		}
		clientConfig = string(statusBar.walletString)
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
	rootCmd.AddCommand(recoverwalletcmd)
	recoverwalletcmd.PersistentFlags().String("mnemonic", "", "mnemonic")
	recoverwalletcmd.MarkFlagRequired("mnemonic")
}
