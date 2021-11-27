package cmd

import (
	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
	"strings"
)

var getZCNBurnTicket = &cobra.Command{
	Use:   "bridge-get-zcn-burnticket",
	Short: "get zcn burn ticket ",
	Long: `returns burn ticket for zcn chain.
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
		bridge.RestoreChain()
		bridge.SetupEthereumWallet()

		// get burn t
	},
}

func init() {
	rootCmd.AddCommand(getZCNBurnTicket)

	verifyEthereumTrxCmd.PersistentFlags().String("file", "", "bridge config file")
	verifyEthereumTrxCmd.PersistentFlags().String("log", "", "bridge log file")
	verifyEthereumTrxCmd.PersistentFlags().String("hash", "", "hash of the ethereum transaction")

	_ = verifyEthereumTrxCmd.MarkFlagRequired("hash")
	_ = verifyEthereumTrxCmd.MarkFlagRequired("file")
}
