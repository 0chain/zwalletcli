package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

var bridgeOwnerInit = &cobra.Command{
	Use:   "bridge-owner-init",
	Short: "init bridge owner config (owner.yaml) in HOME (~/.zcn) folder",
	Long:  `init bridge owner config (owner.yaml) in HOME (~/.zcn) folder`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			path        = GetConfigDir()
			ownerConfig = ConfigOwnerFileName
		)

		fflags := cmd.Flags()

		// Flags

		check(
			cmd,
			"password",
			"ethereumaddress",
			"bridgeaddress",
			"wzcnaddress",
			"authorizersaddress",
			"ethereumnodeurl",
			"gaslimit",
			"value",
		)

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

		path, err = fflags.GetString("path")
		if err != nil {
			fmt.Printf("Flag 'path' not found, defaulting to %s\n", GetConfigDir())
		}

		ownerConfig, err = fflags.GetString("owner_config")
		if err != nil {
			ownerConfig = ConfigOwnerFileName
			fmt.Printf("Flag 'owner_config' not found, defaulting to %s\n", ConfigOwnerFileName)
		}

		// Action

		zcnbridge.CreateInitialOwnerConfig(
			ownerConfig,
			path,
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

//goland:noinspection GoUnhandledErrorResult
func init() {
	f := bridgeOwnerInit
	rootCmd.AddCommand(f)

	f.PersistentFlags().String("path", GetConfigDir(), "Configuration dir")
	f.PersistentFlags().String("owner_config", ConfigOwnerFileName, "Owner config file name")
	f.PersistentFlags().String("password", "", "Password to unlock private key stored in local storage")
	f.PersistentFlags().String("ethereumaddress", "", "Client Ethereum address")
	f.PersistentFlags().String("bridgeaddress", "", "Bridge smart contract address")
	f.PersistentFlags().String("wzcnaddress", "", "WZCN token address")
	f.PersistentFlags().String("authorizersaddress", "", "Authorizers smart contract address")
	f.PersistentFlags().String("ethereumnodeurl", "", "Ethereum Node URL (Infura/Alchemy)")
	f.PersistentFlags().Int64("gaslimit", 0, "Appr. gas limit to execute Ethereum transaction")
	f.PersistentFlags().Int64("value", 0, "Value sent along with Ethereum transaction")

	f.MarkFlagRequired("path")
	f.MarkFlagRequired("password")
	f.MarkFlagRequired("ethereumaddress")
	f.MarkFlagRequired("bridgeaddress")
	f.MarkFlagRequired("wzcnaddress")
	f.MarkFlagRequired("authorizersaddress")
	f.MarkFlagRequired("ethereumnodeurl")
	f.MarkFlagRequired("gaslimit")
}
