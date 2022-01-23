package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var bridgeClientInit = &cobra.Command{
	Use:   "bridge-client-init",
	Short: "init bridge client config (bridge.yaml) in HOME (~/.zcn) folder",
	Long:  `init bridge client config (bridge.yaml) in HOME (~/.zcn) folder`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			path         = GetConfigDir()
			bridgeConfig = ConfigBridgeFileName
		)

		fflags := cmd.Flags()

		// Flags

		check(cmd,
			"password",
			"ethereumaddress",
			"bridgeaddress",
			"wzcnaddress",
			"ethereumnodeurl",
			"gaslimit",
			"consensusthreshold")

		// Reading flags

		password := cmd.Flag("password").Value.String()
		ethereumaddress := cmd.Flag("ethereumaddress").Value.String()
		bridgeaddress := cmd.Flag("bridgeaddress").Value.String()
		wzcnaddress := cmd.Flag("wzcnaddress").Value.String()
		ethereumnodeurl := cmd.Flag("ethereumnodeurl").Value.String()
		gaslimit, err := fflags.GetInt64("gaslimit")
		if err != nil {
			ExitWithError(err)
		}
		value, err := fflags.GetInt64("value")
		if err != nil {
			ExitWithError(err)
		}
		consensusthreshold, err := fflags.GetFloat64("consensusthreshold")
		if err != nil {
			ExitWithError(err)
		}

		path, err = fflags.GetString("path")
		if err != nil {
			fmt.Printf("Flag 'path' not found, defaulting to %s\n", GetConfigDir())
		}

		bridgeConfig, err = fflags.GetString("bridge_config")
		if err != nil {
			bridgeConfig = ConfigBridgeFileName
			fmt.Printf("Flag 'bridge_config' not found, defaulting to %s\n", ConfigBridgeFileName)
		}

		// Action

		zcnbridge.CreateInitialClientConfig(
			bridgeConfig,
			path,
			ethereumaddress,
			bridgeaddress,
			wzcnaddress,
			ethereumnodeurl,
			password,
			gaslimit,
			value,
			consensusthreshold,
		)
	},
}

//goland:noinspection GoUnhandledErrorResult
func init() {
	f := bridgeClientInit
	rootCmd.AddCommand(f)

	f.PersistentFlags().String("path", GetConfigDir(), "Configuration dir")
	f.PersistentFlags().String("bridge_config", ConfigBridgeFileName, "Bridge config file name")
	f.PersistentFlags().String("password", "", "Password be used to unlock private key stored in local storage")
	f.PersistentFlags().String("ethereumaddress", "", "Client Ethereum address")
	f.PersistentFlags().String("bridgeaddress", "", "Bridge contract address")
	f.PersistentFlags().String("wzcnaddress", "", "WZCN token address")
	f.PersistentFlags().String("ethereumnodeurl", "", "Ethereum Node URL (Infura/Alchemy)")
	f.PersistentFlags().Int64("gaslimit", 300000, "appr. Gas limit to execute Ethereum transaction")
	f.PersistentFlags().Float64("consensusthreshold", 0.75, "Consensus threshold required to reach consensus for burn tickets")
	f.PersistentFlags().Int64("value", 0, "Value sent along with Ethereum transaction")

	f.MarkFlagRequired("path")
	f.MarkFlagRequired("password")
	f.MarkFlagRequired("ethereumaddress")
	f.MarkFlagRequired("bridgeaddress")
	f.MarkFlagRequired("wzcnaddress")
	f.MarkFlagRequired("ethereumnodeurl")
	f.MarkFlagRequired("gaslimit")
	f.MarkFlagRequired("consensusthreshold")
	// value not required
}
