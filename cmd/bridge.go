package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/conf"
	"github.com/0chain/gosdk/zcnbridge"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
	"strings"
)

// type HashCommand func(*zcnbridge.Bridge, string)

type Command func(*zcnbridge.BridgeClient, ...*Arg)

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
		name:         "hash",
		value:        "",
		usage:        "hash of the transaction",
		typename:     "string",
		missingError: "hash of the transaction should be provided",
		required:     true,
	}

	amountOption = &Option{
		name:         "amount",
		value:        0,
		usage:        "amount",
		typename:     "int64",
		missingError: "amount should be provided",
		required:     true,
	}
)

func GetHash(args []*Arg) string {
	if len(args) == 0 {
		ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == "hash" {
			return (arg.value).(string)
		}
	}

	ExitWithError("failed to get hash")

	return ""
}

func GetAmount(args []*Arg) int64 {
	if len(args) == 0 {
		ExitWithError("wrong number of arguments")
	}

	for _, arg := range args {
		if arg.fieldName == "amount" {
			return (arg.value).(int64)
		}
	}

	ExitWithError("failed to get hash")

	return 0
}

// createBridgeCommand Function to initialize bridge commands with DRY principle
func createBridgeCommand(use, short, long string, functor Command, opts ...*Option) *cobra.Command {
	var cobraCommand = &cobra.Command{
		Use:   use,
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			var (
				logPath     = "logs"
				configFile  string
				configDir   string
				development = false
				loglevel    = "info"
				err         error
			)

			fflags := cmd.Flags()

			var parameters []*Arg

			for _, opt := range opts {
				if !fflags.Changed(opt.name) {
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
				default:
					ExitWithError(fmt.Printf("unknown argument: %s, value: %v\n", opt.name, opt.value))
				}

				parameters = append(parameters, arg)
			}

			if !fflags.Changed("file") {
				ExitWithError("Error: file flag is missing")
			}

			configFlag, err := fflags.GetString("file")
			if err != nil {
				ExitWithError(err)
			}

			configDir, configFile = filepath.Split(configFlag)
			configFile = strings.TrimSuffix(configFile, path.Ext(configFlag))

			// check SDK EthereumNode
			clientConfig, _ := conf.GetClientConfig()
			if clientConfig.EthereumNode == "" {
				ExitWithError("ethereum_node_url must be setup in config")
			}

			cfg := &zcnbridge.BridgeSDKConfig{
				LogLevel:    &loglevel,
				LogPath:     &logPath,
				ConfigFile:  &configFile,
				ConfigDir:   &configDir,
				Development: &development,
			}

			bridge := zcnbridge.SetupBridgeClientSDK(cfg)
			bridge.RestoreChain()

			functor(bridge, parameters...)
		},
	}

	cobraCommand.PersistentFlags().String("file", "", "bridge config file")
	_ = cobraCommand.MarkFlagRequired("file")

	for _, opt := range opts {
		switch opt.typename {
		case "string":
			cobraCommand.PersistentFlags().String(opt.name, "", opt.usage)
		case "int64":
			cobraCommand.PersistentFlags().Int64(opt.name, 0, opt.usage)
		}

		if opt.required {
			_ = cobraCommand.MarkFlagRequired(opt.name)
		}
	}

	return cobraCommand
}
