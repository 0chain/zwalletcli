package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var verifyEthereumTrxCmd = &cobra.Command{
	Use:   "bridge-verify",
	Short: "verify ethereum transaction ",
	Long: `verify transaction.
	        <hash>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			logPath    = "logs"
			configFile string
			configDir  string
			hash       string
			err        error
		)

		fflags := cmd.Flags()

		if !fflags.Changed("hash") {
			ExitWithError("Error: hash flag is missing")
		}
		if !fflags.Changed("file") {
			ExitWithError("Error: file flag is missing")
		}

		configFlag, err := fflags.GetString("file")
		if err != nil {
			ExitWithError(err)
		}

		hash, err = fflags.GetString("hash")
		if err != nil {
			ExitWithError(err)
		}

		configDir, configFile = filepath.Split(configFlag)
		configFile = strings.TrimSuffix(configFile, path.Ext(configFlag))

		clientConfig, _ := conf.GetClientConfig()
		if clientConfig.EthereumNode == "" {
			ExitWithError("ethereum_node_url must be setup in config")
		}

		bridge := zcnbridge.SetupBridge(configDir, configFile, false, logPath)
		bridge.SetupEthereumWallet()

		status, err := zcnbridge.ConfirmEthereumTransaction(hash, 5, time.Second)
		if err != nil {
			ExitWithError(err)
		}

		if status == 1 {
			fmt.Printf("\nTransaction verification success: %s\n", hash)
		} else {
			ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyEthereumTrxCmd)

	verifyEthereumTrxCmd.PersistentFlags().String("file", "", "bridge config file")
	verifyEthereumTrxCmd.PersistentFlags().String("log", "", "bridge log file")
	verifyEthereumTrxCmd.PersistentFlags().String("hash", "", "hash of the ethereum transaction")

	_ = verifyEthereumTrxCmd.MarkFlagRequired("hash")
	_ = verifyEthereumTrxCmd.MarkFlagRequired("file")
}
