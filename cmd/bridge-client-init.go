package cmd

import (
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var bridgeClientInit = &cobra.Command{
	Use:   "bridge-client-init",
	Short: "init bridge client config (bridge.yaml) in HOME (~/.zcn) folder",
	Long:  `init bridge client config (bridge.yaml) in HOME (~/.zcn) folder`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
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

		// Action

		zcnbridge.CreateInitialClientConfig(
			"bridge.yaml",
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

func init() {
	rootCmd.AddCommand(bridgeClientInit)

	bridgeClientInit.PersistentFlags().String("password", "", "password be used to unlock private key stored in local storage")
	bridgeClientInit.PersistentFlags().String("ethereumaddress", "", "client Ethereum address")
	bridgeClientInit.PersistentFlags().String("bridgeaddress", "", "bridge contract address")
	bridgeClientInit.PersistentFlags().String("wzcnaddress", "", "WZCN token address")
	bridgeClientInit.PersistentFlags().String("ethereumnodeurl", "", "Ethereum Node URL (Infura/Alchemy)")
	bridgeClientInit.PersistentFlags().Int64("gaslimit", 300000, "appr. gas limit to execute Ethereum transaction")
	bridgeClientInit.PersistentFlags().Int64("value", 0, "value sent along with Ethereum transaction")
	bridgeClientInit.PersistentFlags().Float64("consensusthreshold", 0.75, "Consensus threshold required to reach consensus for burn tickets")

	_ = bridgeClientInit.MarkFlagRequired("password")
	_ = bridgeClientInit.MarkFlagRequired("ethereumaddress")
	_ = bridgeClientInit.MarkFlagRequired("bridgeaddress")
	_ = bridgeClientInit.MarkFlagRequired("wzcnaddress")
	_ = bridgeClientInit.MarkFlagRequired("ethereumnodeurl")
	_ = bridgeClientInit.MarkFlagRequired("gaslimit")
	_ = bridgeClientInit.MarkFlagRequired("consensusthreshold")
}
