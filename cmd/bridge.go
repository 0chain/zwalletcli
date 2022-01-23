package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

const (
	DefaultRetries = 60
)

const (
	ConfigBridgeFileName = "bridge.yaml"
	ConfigOwnerFileName  = "owner.yaml"
	ConfigChainFileName  = "config.yaml"
)

const (
	OptionHash             = "hash"          // OptionHash hash passed to cmd
	OptionAmount           = "amount"        // OptionAmount amount passed to cmd
	OptionRetries          = "retries"       // OptionRetries retries
	OptionConfigFolder     = "path"          // OptionConfigFolder config folder
	OptionChainConfigFile  = "chain_config"  // OptionChainConfigFile sdk config filename
	OptionBridgeConfigFile = "bridge_config" // OptionBridgeConfigFile bridge config filename
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
		value:        ConfigChainFileName,
		typename:     "string",
		usage:        "Chain config file name",
		missingError: "Chain config file name not specified",
		required:     false,
	}

	configBridgeFileOption = &Option{
		name:         OptionBridgeConfigFile,
		value:        ConfigBridgeFileName,
		typename:     "string",
		usage:        "Bridge config file name",
		missingError: "Bridge config file name not specified",
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

func GetBridgeConfigFile(args []*Arg) string {
	return getString(args, OptionBridgeConfigFile)
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

func GetAmount(args []*Arg) int64 {
	return getInt64(args, OptionAmount)
}

func GetRetries(args []*Arg) int {
	return getInt(args, OptionRetries)
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

// createCommand Function to initialize bridge commands with DRY principle
func createCommand(use, short, long string, functor Command, opts ...*Option) *cobra.Command {
	fn := func(parameters ...*Arg) {
		functor(parameters...)
	}

	opts = append(opts, configFolderOption)
	opts = append(opts, configBridgeFileOption)
	opts = append(opts, configChainFileOption)

	command := createBridgeComm(use, short, long, fn, opts)
	AppendOptions(opts, command)
	return command
}

// createCommandWithBridge Function to initialize bridge commands with DRY principle
func createCommandWithBridge(use, short, long string, functor CommandWithBridge, opts ...*Option) *cobra.Command {
	fn := func(parameters ...*Arg) {
		folder := GetConfigFolder(parameters)
		chainConfigFile := GetChainConfigFile(parameters)
		bridgeConfigFile := GetBridgeConfigFile(parameters)

		bridge := initBridge(folder, chainConfigFile, bridgeConfigFile)
		functor(bridge, parameters...)
	}

	opts = append(opts, configFolderOption)
	opts = append(opts, configBridgeFileOption)
	opts = append(opts, configChainFileOption)

	command := createBridgeComm(use, short, long, fn, opts)
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
) *cobra.Command {
	var cobraCommand = &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fflags := cmd.Flags()

			var parameters []*Arg

			for _, opt := range opts {
				if !fflags.Changed(opt.name) && opt.required {
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

func initBridge(overrideConfigFolder, overrideConfigFile, overrideBridgeFile string) *zcnbridge.BridgeClient {
	var (
		configDir            = GetConfigDir()
		configBridgeFileName = ConfigBridgeFileName
		configChainFileName  = ConfigChainFileName
		logPath              = "logs"
		loglevel             = "info"
		development          = false
	)

	if overrideConfigFolder != "" {
		configDir = overrideConfigFolder
	}

	configBridgeFileName = overrideBridgeFile
	configChainFileName = overrideConfigFile

	configDir, err := filepath.Abs(configDir)
	if err != nil {
		ExitWithError(err)
	}

	cfg := &zcnbridge.BridgeSDKConfig{
		ConfigDir:        &configDir,
		ConfigBridgeFile: &configBridgeFileName,
		ConfigChainFile:  &configChainFileName,
		LogPath:          &logPath,
		LogLevel:         &loglevel,
		Development:      &development,
	}

	bridge := zcnbridge.SetupBridgeClientSDK(cfg)
	bridge.RestoreChain()

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

func verify(hash string) {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	txn, err := zcncore.NewTransaction(statusBar, 0)
	if err != nil {
		ExitWithError(err)
	}
	_ = txn.SetTransactionHash(hash)
	wg.Add(1)
	err = txn.Verify()
	if err == nil {
		wg.Wait()
	} else {
		ExitWithError(err.Error())
	}
	if statusBar.success {
		statusBar.success = false
		fmt.Printf("\nTransaction verification success\n")
		return
	}
	ExitWithError("\nVerification failed." + statusBar.errMsg + "\n")
}
