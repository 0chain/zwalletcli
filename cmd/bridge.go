package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
	"os"
)

const (
	DefaultRetries = 60
)

const (
	ConfigFileName = "bridge"
)

const (
	OptionHash    = "hash"
	OptionAmount  = "amount"
	OptionRetries = "retries"
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
	hashOption = &Option{
		name:         OptionHash,
		value:        "",
		usage:        "hash of the transaction",
		typename:     "string",
		missingError: "hash of the transaction should be provided",
		required:     true,
	}

	amountOption = &Option{
		name:         OptionAmount,
		value:        0,
		usage:        "amount",
		typename:     "int64",
		missingError: "amount should be provided",
		required:     true,
	}

	retriesOption = &Option{
		name:         OptionRetries,
		value:        DefaultRetries,
		typename:     "int",
		usage:        "",
		missingError: "retries count should be provided",
		required:     false,
	}
)

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

	command := createBridgeComm(use, short, long, fn, opts)
	AppendOptions(opts, command)
	return command
}

// createCommandWithBridge Function to initialize bridge commands with DRY principle
func createCommandWithBridge(use, short, long string, functor CommandWithBridge, opts ...*Option) *cobra.Command {
	fn := func(parameters ...*Arg) {
		bridge := initBridge()
		functor(bridge, parameters...)
	}

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

func initBridge() *zcnbridge.BridgeClient {
	var (
		configDir   = GetConfigDir()
		configFile  = ConfigFileName
		logPath     = "logs"
		loglevel    = "info"
		development = false
	)

	cfg := &zcnbridge.BridgeSDKConfig{
		LogLevel:    &loglevel,
		LogPath:     &logPath,
		ConfigFile:  &configFile,
		ConfigDir:   &configDir,
		Development: &development,
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
