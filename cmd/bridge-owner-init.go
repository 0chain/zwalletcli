package cmd

import (
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var bridgeOwnerInit = &cobra.Command{
	Use:   "bridge-owner-init",
	Short: "init bridge owner config (owner.yaml) in HOME (~/.zcn) folder",
	Long:  `init bridge owner config (owner.yaml) in HOME (~/.zcn) folder`,
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
		if fflags.Changed("authorizersaddress") == false {
			ExitWithError("Error: 'authorizersaddress' flag is missing")
		}
		if fflags.Changed("ethereumnodeurl") == false {
			ExitWithError("Error: 'ethereumnodeurl' flag is missing")
		}
		if fflags.Changed("gaslimit") == false {
			ExitWithError("Error: 'gaslimit' flag is missing")
		}
		if fflags.Changed("value") == false {
			ExitWithError("Error: 'value' flag is missing")
		}

		// Reading flags

		password := cmd.Flag("password").Value.String()
		ethereumaddress := cmd.Flag("ethereumaddress").Value.String()
		bridgeaddress := cmd.Flag("bridgeaddress").Value.String()
		wzcnaddress := cmd.Flag("wzcnaddress").Value.String()
		authorizersaddress := cmd.Flag("authorizersaddress").Value.String()
		ethereumnodeurl := cmd.Flag("ethereumnodeurl").Value.String()
		gaslimit, err := fflags.GetInt64("gaslimit")
		if err != nil {
			ExitWithError(err)
		}
		value, err := fflags.GetInt64("value")
		if err != nil {
			ExitWithError(err)
		}

		// Action

		zcnbridge.CreateInitialOwnerConfig(
			"owner.yaml",
			ethereumaddress,
			bridgeaddress,
			wzcnaddress,
			authorizersaddress,
			ethereumnodeurl,
			password,
			gaslimit,
			value,
		)
	},
}

func init() {
	rootCmd.AddCommand(bridgeOwnerInit)

	bridgeOwnerInit.PersistentFlags().String("password", "", "password")
	bridgeOwnerInit.PersistentFlags().String("ethereumaddress", "", "ethereumaddress")
	bridgeOwnerInit.PersistentFlags().String("bridgeaddress", "", "bridgeaddress")
	bridgeOwnerInit.PersistentFlags().String("wzcnaddress", "", "wzcnaddress")
	bridgeOwnerInit.PersistentFlags().String("ethereumnodeurl", "", "ethereumnodeurl")
	bridgeOwnerInit.PersistentFlags().Int64("gaslimit", 0, "gaslimit")
	bridgeOwnerInit.PersistentFlags().Int64("value", 0, "value")

	_ = bridgeOwnerInit.MarkFlagRequired("password")
	_ = bridgeOwnerInit.MarkFlagRequired("ethereumaddress")
	_ = bridgeOwnerInit.MarkFlagRequired("bridgeaddress")
	_ = bridgeOwnerInit.MarkFlagRequired("wzcnaddress")
	_ = bridgeOwnerInit.MarkFlagRequired("ethereumnodeurl")
	_ = bridgeOwnerInit.MarkFlagRequired("gaslimit")
}
