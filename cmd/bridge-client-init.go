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

		if fflags.Changed("password") == false {
			ExitWithError("Error: 'password' flag is missing")
		}
		if fflags.Changed("ethereumaddress") == false {
			ExitWithError("Error: 'ethereumaddress' flag is missing")
		}
		if fflags.Changed("bridgeaddress") == false {
			ExitWithError("Error: 'bridgeaddress' flag is missing")
		}
		if fflags.Changed("wzcnaddress") == false {
			ExitWithError("Error: 'wzcnaddress' flag is missing")
		}
		if fflags.Changed("ethereumnodeurl") == false {
			ExitWithError("Error: 'ethereumnodeurl' flag is missing")
		}
		if fflags.Changed("gaslimit") == false {
			ExitWithError("Error: 'gaslimit' flag is missing")
		}
		if fflags.Changed("consensusthreshold") == false {
			ExitWithError("Error: 'consensusthreshold' flag is missing")
		}

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

	bridgeClientInit.PersistentFlags().String("password", "", "password")
	bridgeClientInit.PersistentFlags().String("ethereumaddress", "", "ethereumaddress")
	bridgeClientInit.PersistentFlags().String("bridgeaddress", "", "bridgeaddress")
	bridgeClientInit.PersistentFlags().String("wzcnaddress", "", "wzcnaddress")
	bridgeClientInit.PersistentFlags().String("ethereumnodeurl", "", "ethereumnodeurl")
	bridgeClientInit.PersistentFlags().Int64("gaslimit", 300000, "gaslimit")
	bridgeClientInit.PersistentFlags().Int64("value", 0, "value")
	bridgeClientInit.PersistentFlags().Float64("consensusthreshold", 0.75, "consensusthreshold")

	_ = bridgeClientInit.MarkFlagRequired("password")
	_ = bridgeClientInit.MarkFlagRequired("ethereumaddress")
	_ = bridgeClientInit.MarkFlagRequired("bridgeaddress")
	_ = bridgeClientInit.MarkFlagRequired("wzcnaddress")
	_ = bridgeClientInit.MarkFlagRequired("ethereumnodeurl")
	_ = bridgeClientInit.MarkFlagRequired("gaslimit")
	_ = bridgeClientInit.MarkFlagRequired("consensusthreshold")
}
