package cmd

import (
	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
	"strings"
)

type HashCommand func(*zcnbridge.Bridge, string)

// createBridgeCommand Function to initialize bridge commands with DRY principle
func createBridgeCommand(comm HashCommand, use, short, long string) *cobra.Command {
	var cobraCommand = &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(0),
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

			comm(bridge, hash)
		},
	}

	cobraCommand.PersistentFlags().String("file", "", "bridge config file")
	cobraCommand.PersistentFlags().String("log", "", "bridge log file")
	cobraCommand.PersistentFlags().String("hash", "", "hash of the ethereum transaction")

	_ = cobraCommand.MarkFlagRequired("hash")
	_ = cobraCommand.MarkFlagRequired("file")

	return cobraCommand
}
