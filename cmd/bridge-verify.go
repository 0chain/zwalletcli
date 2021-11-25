package cmd

import (
	"fmt"
	"sync"

	"github.com/0chain/gosdk/"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var verifyEthereumTrxCmd = &cobra.Command{
	Use:   "bridge-verify",
	Short: "verify ethereum transaction ",
	Long: `verify transaction.
	        <hash>`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("hash") == false {
			ExitWithError("Error: hash flag is missing")
		}

		//cfg := &config.BridgeConfig{
		//	Mnemonic:        "",
		//	BridgeAddress:   "",
		//	WzcnAddress:     "",
		//	EthereumNodeURL: "",
		//	ChainID:         "",
		//	GasLimit:        0,
		//	Value:           0,
		//}
		//
		//zcnbridge.InitBridge(cfg)

		//config.Bridge.BridgeAddress = viper.GetString("bridge.BridgeAddress")
		//config.Bridge.Mnemonic = viper.GetString("bridge.Mnemonic")
		//config.Bridge.EthereumNodeURL = viper.GetString("bridge.EthereumNodeURL")
		//config.Bridge.Value = viper.GetInt64("bridge.Value")
		//config.Bridge.GasLimit = viper.GetUint64("bridge.GasLimit")
		//config.Bridge.WzcnAddress = viper.GetString("bridge.WzcnAddress")
		//config.Bridge.ChainID = viper.GetString("bridge.ChainID")

		hash := cmd.Flag("hash").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			ExitWithError(err)
		}
		txn.SetTransactionHash(hash)
		wg.Add(1)
		err = txn.Verify()
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}

		if statusBar.success {
			statusBar.success = false
			fmt.Printf("\nTransaction verification success\n")
			return
		}
		ExitWithError("\nVerification failed." + statusBar.errMsg + "\n")
	},
}

func init2() {
	rootCmd.AddCommand(verifyEthereumTrxCmd)
	verifyEthereumTrxCmd.PersistentFlags().String("hash", "", "hash of the ethereum transaction")
	err := verifyEthereumTrxCmd.MarkFlagRequired("hash")
	if err != nil {
		return
	}
}

func init() {
	rootCmd.AddCommand(faucetcmd)
	faucetcmd.PersistentFlags().String("methodName", "", "methodName")
	faucetcmd.PersistentFlags().String("input", "", "input")
	faucetcmd.PersistentFlags().Float64("tokens", 0, "Token request")
	faucetcmd.MarkFlagRequired("methodName")
	faucetcmd.MarkFlagRequired("input")
}
