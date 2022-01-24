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
			ownerConfig = DefaultConfigOwnerFileName
		)

		fflags := cmd.Flags()

		// Flags

		check(
			cmd,
			OptionKeyPassword,
			"ethereumaddress",
			"bridgeaddress",
			"wzcnaddress",
			"authorizersaddress",
			"ethereumnodeurl",
			"gaslimit",
			"value",
		)

		// Reading flags

		password := cmd.Flag(OptionKeyPassword).Value.String()
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

		path, err = fflags.GetString(OptionConfigFolder)
		if err != nil {
			fmt.Printf("Flag '%s' not found, defaulting to %s\n", OptionConfigFolder, GetConfigDir())
		}

		ownerConfig, err = fflags.GetString(OptionOwnerConfigFile)
		if err != nil {
			ownerConfig = DefaultConfigOwnerFileName
			fmt.Printf("Flag '%s' not found, defaulting to %s\n", OptionOwnerConfigFile, DefaultConfigOwnerFileName)
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

	f.PersistentFlags().String(OptionConfigFolder, GetConfigDir(), "Configuration dir")
	f.PersistentFlags().String(OptionOwnerConfigFile, DefaultConfigOwnerFileName, "Owner config file name")
	f.PersistentFlags().String(OptionKeyPassword, "", "Password to unlock private key stored in local storage")
	f.PersistentFlags().String("ethereumaddress", "", "Client Ethereum address")
	f.PersistentFlags().String("bridgeaddress", "", "Bridge smart contract address")
	f.PersistentFlags().String("wzcnaddress", "", "WZCN token address")
	f.PersistentFlags().String("authorizersaddress", "", "Authorizers smart contract address")
	f.PersistentFlags().String("ethereumnodeurl", "", "Ethereum Node URL (Infura/Alchemy)")
	f.PersistentFlags().Int64("gaslimit", 0, "Appr. gas limit to execute Ethereum transaction")
	f.PersistentFlags().Int64("value", 0, "Value sent along with Ethereum transaction")

	f.MarkFlagRequired(OptionConfigFolder)
	f.MarkFlagRequired(OptionKeyPassword)
	f.MarkFlagRequired("ethereumaddress")
	f.MarkFlagRequired("bridgeaddress")
	f.MarkFlagRequired("wzcnaddress")
	f.MarkFlagRequired("authorizersaddress")
	f.MarkFlagRequired("ethereumnodeurl")
	f.MarkFlagRequired("gaslimit")
}
