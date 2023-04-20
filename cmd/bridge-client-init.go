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
			path                 = GetConfigDir()
			bridgeConfigFileName = DefaultConfigBridgeFileName
		)

		fflags := cmd.Flags()

		// Flags

		check(cmd,
			OptionKeyPassword,
			"ethereumaddress",
			"bridgeaddress",
			"wzcnaddress",
			"ethereumnodeurl",
			"gaslimit",
			"consensusthreshold")

		// Reading flags

		password := cmd.Flag(OptionKeyPassword).Value.String()
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

		path, err = fflags.GetString(OptionConfigFolder)
		if err != nil {
			fmt.Printf("Flag '%s' not found, defaulting to %s\n", OptionConfigFolder, GetConfigDir())
		}

		bridgeConfigFileName, err = fflags.GetString(OptionBridgeConfigFile)
		if err != nil {
			bridgeConfigFileName = DefaultConfigBridgeFileName
			fmt.Printf("Flag '%s' not found, defaulting to %s\n", OptionBridgeConfigFile, DefaultConfigBridgeFileName)
		}

		// Action

		zcnbridge.CreateInitialClientConfig(
			bridgeConfigFileName,
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

	f.PersistentFlags().String(OptionConfigFolder, GetConfigDir(), "Configuration dir")
	f.PersistentFlags().String(OptionBridgeConfigFile, DefaultConfigBridgeFileName, "Bridge config file name")
	f.PersistentFlags().String(OptionKeyPassword, "", "Password be used to unlock private key stored in local storage")
	f.PersistentFlags().String("ethereumaddress string", "", "Client Ethereum address")
	f.PersistentFlags().String("bridgeaddress string", "", "Bridge contract address")
	f.PersistentFlags().String("wzcnaddress string", "", "WZCN token address")
	f.PersistentFlags().String("ethereumnodeurl string", "", "Ethereum Node URL (Infura/Alchemy)")
	f.PersistentFlags().Int64("gaslimit int", 300000, "appr. Gas limit to execute Ethereum transaction")
	f.PersistentFlags().Float64("consensusthreshold float", 0.75, "Consensus threshold required to reach consensus for burn tickets")
	f.PersistentFlags().Int64("value int", 0, "Value sent along with Ethereum transaction")

	f.MarkFlagRequired(OptionConfigFolder)
	f.MarkFlagRequired(OptionKeyPassword)
	f.MarkFlagRequired("ethereumaddress")
	f.MarkFlagRequired("bridgeaddress")
	f.MarkFlagRequired("wzcnaddress")
	f.MarkFlagRequired("ethereumnodeurl")
	f.MarkFlagRequired("gaslimit")
	f.MarkFlagRequired("consensusthreshold")
	// value not required
}
