package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
)

const (
	DefaultRetries = 60
)

const (
	DefaultConfigChainFileName = "config.yaml"
	DefaultWalletFileName      = "wallet.json"
)

const (
	OptionHash            = "hash"          // OptionHash hash passed to cmd
	OptionAmount          = "amount"        // OptionAmount amount passed to cmd
	OptionToken           = "token"         // OptionToken token in SAS passed to cmd
	OptionRetries         = "retries"       // OptionRetries retries
	OptionConfigFolder    = "path"          // OptionConfigFolder config folder
	OptionChainConfigFile = "chain_config"  // OptionChainConfigFile sdk config filename
	OptionMnemonic        = "mnemonic"      // OptionMnemonic bridge config filename
	OptionKeyPassword     = "password"      // OptionKeyPassword bridge config filename
	OptionAccountIndex    = "account_index" // OptionAccountIndex ethereum account index
	OptionAddressIndex    = "address_index" // OptionAddressIndex ethereum address index
	OptionBip32           = "bip32"         // OptionBip32 use bip32 derivation path
	OptionClientKey       = "client_key"
	OptionClientID        = "client_id"
	OptionEthereumAddress = "ethereum_address"
	OptionURL             = "url"
	OptionMinStake        = "min_stake"
	OptionMaxStake        = "max_stake"
	OptionNumDelegates    = "num_delegates"
	OptionServiceCharge   = "service_charge"
	OptionWalletFile      = "wallet"
)

type CommandWithBridge func(*zcnbridge.BridgeClient, ...*Arg)
type Command func(...*Arg)

type Option struct {
	name         string
	value        interface{}
	typename     string
	usage        string
	missingError string
	required     bool
}

type Arg struct {
	typeName  string
	fieldName string
	value     interface{}
}

var (
	walletFileOption = &Option{
		name:         OptionWalletFile,
		value:        "wallet.json",
		typename:     "string",
		usage:        "Wallet file",
		missingError: "Wallet file not specified",
		required:     false,
	}

	configFolderOption = &Option{
		name:         OptionConfigFolder,
		value:        GetConfigDir(),
		typename:     "string",
		usage:        "Config home folder",
		missingError: "Config home folder not specified",
		required:     false,
	}

	configChainFileOption = &Option{
		name:         OptionChainConfigFile,
		value:        DefaultConfigChainFileName,
		typename:     "string",
		usage:        "Chain config file name",
		missingError: "Chain config file name not specified",
		required:     false,
	}
)

func WithRetries(usage string) *Option {
	return &Option{
		name:         OptionRetries,
		value:        DefaultRetries,
		typename:     "int",
		usage:        usage,
		missingError: "Retries count should be provided",
		required:     false,
	}
}

func WithToken(usage string) *Option {
	return &Option{
		name:         OptionToken,
		value:        float64(0),
		usage:        usage,
		typename:     "float64",
		missingError: "Token should be provided",
		required:     true,
	}
}

func WithAmount(usage string) *Option {
	return &Option{
		name:         OptionAmount,
		value:        int64(0),
		usage:        usage,
		typename:     "int64",
		missingError: "Amount should be provided",
		required:     true,
	}
}

func WithHash(usage string) *Option {
	return &Option{
		name:         OptionHash,
		value:        "",
		usage:        usage,
		typename:     "string",
		missingError: "hash of the transaction should be provided",
		required:     true,
	}
}

func GetChainConfigFile(args []*Arg) string {
	return getString(args, OptionChainConfigFile)
}

func GetConfigFolder(args []*Arg) string {
	return getString(args, OptionConfigFolder)
}

func GetHash(args []*Arg) string {
	return getString(args, OptionHash)
}

func GetAmount(args []*Arg) uint64 {
	return uint64(getInt64(args, OptionAmount))
}

func GetToken(args []*Arg) float64 {
	return getFloat64(args, OptionToken)
}

func GetRetries(args []*Arg) int {
	return getInt(args, OptionRetries)
}

func GetClientID(args []*Arg) string {
	return getString(args, OptionClientID)
}

func GetClientKey(args []*Arg) string {
	return getString(args, OptionClientKey)
}

func GetEthereumAddress(args []*Arg) string {
	return getString(args, OptionEthereumAddress)
}

func GetURL(args []*Arg) string {
	return getString(args, OptionURL)
}

func GetMinStake(args []*Arg) int64 {
	return getInt64(args, OptionMinStake)
}

func GetMaxStake(args []*Arg) int64 {
	return getInt64(args, OptionMaxStake)
}

func GetNumDelegates(args []*Arg) int {
	return getInt(args, OptionNumDelegates)
}

func GetServiceCharge(args []*Arg) float64 {
	return getFloat64(args, OptionServiceCharge)
}

func GetWalletFile(args []*Arg) string {
	return getString(args, OptionWalletFile)
}

func getString(args []*Arg, name string) string {
	if len(args) == 0 {
		ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(string)
		}
	}

	ExitWithError("failed to get " + name)

	return ""
}

func getInt(args []*Arg, name string) int {
	if len(args) == 0 {
		ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(int)
		}
	}

	ExitWithError("failed to get " + name)

	return 0
}

func getFloat64(args []*Arg, name string) float64 {
	if len(args) == 0 {
		ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(float64)
		}
	}

	ExitWithError("failed to get " + name)

	return 0
}

func getInt64(args []*Arg, name string) int64 {
	if len(args) == 0 {
		ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(int64)
		}
	}

	ExitWithError("failed to get " + name)

	return 0
}

func getUint64(args []*Arg, name string) uint64 {
	if len(args) == 0 {
		ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == name {
			return (arg.value).(uint64)
		}
	}

	ExitWithError("failed to get " + name)

	return 0
}

// createCommand Function to initialize bridge commands with DRY principle
func createCommand(use, short, long string, functor Command, hidden bool, opts ...*Option,) *cobra.Command {
	fn := func(parameters ...*Arg) {
		functor(parameters...)
	}

	opts = append(opts, configFolderOption)
	opts = append(opts, configChainFileOption)

	command := createBridgeComm(use, short, long, fn, opts, hidden)
	AppendOptions(opts, command)
	return command
}

// createCommandWithBridge Function to initialize bridge commands with DRY principle
func createCommandWithBridge(use, short, long string, functor CommandWithBridge, hidden bool, opts ...*Option) *cobra.Command {
	fn := func(parameters ...*Arg) {
		folder := GetConfigFolder(parameters)
		chainConfigFile := GetChainConfigFile(parameters)

		bridge := initBridge(folder, chainConfigFile)
		functor(bridge, parameters...)
	}

	opts = append(opts, configFolderOption)
	opts = append(opts, configChainFileOption)
	command := createBridgeComm(use, short, long, fn, opts, hidden)
	AppendOptions(opts, command)
	return command
}

func AppendOptions(opts []*Option, command *cobra.Command) {
	for _, opt := range opts {
		switch opt.typename {
		case "string":
			command.PersistentFlags().String(opt.name, opt.value.(string), opt.usage)
		case "int64":
			command.PersistentFlags().Int64(opt.name, opt.value.(int64), opt.usage)
		case "float64":
			command.PersistentFlags().Float64(opt.name, opt.value.(float64), opt.usage)
		case "int":
			command.PersistentFlags().Int(opt.name, opt.value.(int), opt.usage)
		}

		if opt.required {
			_ = command.MarkFlagRequired(opt.name)
		}
	}
}

func createBridgeComm(
	use string,
	short string,
	long string,
	functor Command,
	opts []*Option,
	hidden bool,
) *cobra.Command {
	var cobraCommand = &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(0),
		Hidden: hidden,
		Run: func(cmd *cobra.Command, args []string) {
			fflags := cmd.Flags()

			var parameters []*Arg

			for _, opt := range opts {
				if !fflags.Changed(opt.name) && opt.required {
					//TODO: add default missing error
					ExitWithError(opt.missingError)
				}

				var arg *Arg
				switch opt.typename {
				case "string":
					optValue, err := fflags.GetString(opt.name)
					if err != nil {
						ExitWithError(err)
					}
					arg = &Arg{
						typeName:  opt.typename,
						fieldName: opt.name,
						value:     optValue,
					}
				case "int64":
					optValue, err := fflags.GetInt64(opt.name)
					if err != nil {
						ExitWithError(err)
					}
					arg = &Arg{
						typeName:  opt.typename,
						fieldName: opt.name,
						value:     optValue,
					}
				case "float64":
					optValue, err := fflags.GetFloat64(opt.name)
					if err != nil {
						ExitWithError(err)
					}
					arg = &Arg{
						typeName:  opt.typename,
						fieldName: opt.name,
						value:     optValue,
					}
				case "int":
					optValue, err := fflags.GetInt(opt.name)
					if err != nil {
						ExitWithError(err)
					}
					arg = &Arg{
						typeName:  opt.typename,
						fieldName: opt.name,
						value:     optValue,
					}
				default:
					ExitWithError(fmt.Printf("unknown argument: %s, value: %v\n", opt.name, opt.value))
				}

				parameters = append(parameters, arg)
			}

			// check SDK EthereumNode
			clientConfig, _ := conf.GetClientConfig()
			if clientConfig.EthereumNode == "" {
				ExitWithError("ethereum_node_url must be setup in config")
			}

			functor(parameters...)
		},
	}
	return cobraCommand
}

func GetConfigDir() string {
	var configDir string
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configDir = home + "/.zcn"
	return configDir
}

func initBridge(overrideConfigFolder, overrideConfigFile string) *zcnbridge.BridgeClient {
	var (
		configDir           = GetConfigDir()
		configChainFileName = DefaultConfigChainFileName
		logPath             = "logs"
		loglevel            = "info"
		development         = false
	)

	if overrideConfigFolder != "" {
		configDir = overrideConfigFolder
	}

	configChainFileName = overrideConfigFile

	configDir, err := filepath.Abs(configDir)
	if err != nil {
		ExitWithError(err)
	}

	cfg := &zcnbridge.BridgeSDKConfig{
		ConfigDir:       &configDir,
		ConfigChainFile: &configChainFileName,
		LogPath:         &logPath,
		LogLevel:        &loglevel,
		Development:     &development,
	}

	bridge := zcnbridge.SetupBridgeClientSDK(cfg)

	return bridge
}

func check(cmd *cobra.Command, flags ...string) {
	fflags := cmd.Flags()
	for _, flag := range flags {
		if !fflags.Changed(flag) {
			ExitWithError(fmt.Sprintf("Error: '%s' flag is missing", flag))
		}
	}
}
