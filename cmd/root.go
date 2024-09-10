package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/0chain/gosdk/core/client"
	"github.com/0chain/gosdk/core/conf"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/0chain/gosdk/core/zcncrypto"
	"github.com/0chain/gosdk/zboxcore/sdk"
	bridge "github.com/0chain/gosdk/zcnbridge/http"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var networkFile string
var walletFile string
var cDir string
var bSilent bool
var nonce int64

// gTxnFee is the user specified fee passed from client/user.
// If the fee is absent/low it is adjusted to the min fee required
// (acquired from miner) for the transaction to write into blockchain.
var gTxnFee float64

var clientConfig string
var minSubmit int
var minCfm int
var CfmChainLength int
var signatureScheme string

var (
	cfgConfig  *viper.Viper
	cfgNetwork *viper.Viper
	cfgWallet  string
)

var rootCmd = &cobra.Command{
	Use:   "zwallet",
	Short: "Use Zwallet to store, send and execute smart contract on 0Chain platform",
	Long: `Use Zwallet to store, send and execute smart contract on 0Chain platform.
			Complete documentation is available at https://docs.zus.network/guides/zwallet-cli`,
	PersistentPreRun: initCmdContext,
}

var clientWallet *zcncrypto.Wallet

func init() {
	InstallDLLs()
	cobra.OnInitialize(loadConfigs)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().StringVar(&networkFile, "network", "", "network file to overwrite the network details (if required, default is network.yaml)")
	rootCmd.PersistentFlags().StringVar(&walletFile, "wallet", "", "wallet file (default is wallet.json)")
	rootCmd.PersistentFlags().StringVar(&cDir, "configDir", "", "configuration directory (default is $HOME/.zcn)")
	rootCmd.PersistentFlags().Int64Var(&nonce, "withNonce", 0, "nonce that will be used in transaction (default is 0)")
	rootCmd.PersistentFlags().BoolVar(&bSilent, "silent", false, "Do not print sdk logs in stderr (prints by default)")

	rootCmd.PersistentFlags().Float64Var(&gTxnFee, "fee", 0, "transaction fee for the given transaction (if unset, it will be set to blockchain min fee)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getConfigDir() string {
	if cDir != "" {
		return cDir
	}
	var configDir string
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configDir = filepath.Join(home, "/.zcn")
	return configDir
}

func initZCNCore() {

	// set the log file
	zcncore.SetLogFile("cmdlog.log", !bSilent)
	bridge.SetLogFile("bridge.log", !bSilent)
	sdk.SetLogFile("cmdlog.log", !bSilent)

	blockWorker := cfgConfig.GetString("block_worker")
	chainID := cfgConfig.GetString("chain_id")
	ethereumNodeURL := cfgConfig.GetString("ethereum_node_url")

	cfg := conf.Config{
		BlockWorker:             blockWorker,
		SignatureScheme:         signatureScheme,
		ChainID:                 chainID,
		MinSubmit:               minSubmit,
		MinConfirmation:         minCfm,
		ConfirmationChainLength: CfmChainLength,
		EthereumNode:            ethereumNodeURL,
	}

	err := client.Init(context.Background(), cfg)
	if err != nil {
		ExitWithError(err.Error())
	}
}

func loadConfigs() {
	cfgConfig = viper.New()
	cfgNetwork = viper.New()
	var configDir string

	if cDir != "" {
		configDir = cDir
	} else {
		configDir = getConfigDir()
	}

	// ~/.zcn/config.yaml
	cfgConfig.AddConfigPath(configDir)
	if cfgFile != "" {
		cfgConfig.SetConfigFile(filepath.Join(configDir, cfgFile))
	} else {
		cfgConfig.SetConfigFile(filepath.Join(configDir, "config.yaml"))
	}

	if err := cfgConfig.ReadInConfig(); err != nil {
		ExitWithError("Can't read config:", err, cDir, configDir, cfgFile)
	}

	minSubmit = cfgConfig.GetInt("min_submit")
	minCfm = cfgConfig.GetInt("min_confirmation")
	CfmChainLength = cfgConfig.GetInt("confirmation_chain_length")
	signatureScheme = cfgConfig.GetString("signature_scheme")

	// ~/.zcn/network.yaml
	cfgNetwork.AddConfigPath(configDir)
	if len(networkFile) > 0 {
		cfgNetwork.SetConfigFile(filepath.Join(configDir, networkFile))
	} else {
		cfgNetwork.SetConfigFile(filepath.Join(configDir, "network.yaml"))
	}

	cfgNetwork.ReadInConfig() //nolint: errcheck

	// TODO: move the private key storage to the keychain or secure storage
	// ~/.zcn/wallet.json
	if len(walletFile) > 0 {
		cfgWallet = filepath.Join(configDir, walletFile)
	} else {
		cfgWallet = filepath.Join(configDir, "/wallet.json")
	}
}

var zcncoreIsInitialized bool
var walletIsLoaded bool

func initCmdContext(cmd *cobra.Command, args []string) {

	_, ok := withoutZCNCoreCmds[cmd]
	if !ok {
		initZCNCoreContext()
	}

	_, ok = withoutWalletCmds[cmd]
	if !ok {
		initZwalletContext()
	}

}

func initZCNCoreContext() {
	// zcncore is initialized , skip any zcncore checking
	if !zcncoreIsInitialized {
		initZCNCore()
		zcncoreIsInitialized = true
	}
}

func initZwalletContext() {
	// create wallet
	if !walletIsLoaded {
		createAndLoadWallet()
		walletIsLoaded = true
	}
}

func createAndLoadWallet() {

	_, err := os.Stat(cfgWallet)

	isNewWallet := os.IsNotExist(err)

	if isNewWallet {
		fmt.Println("No wallet in path ", cfgWallet, "found. Creating wallet...")
		walletString, err := createWallet()
		if err != nil {
			ExitWithError(err)
		}

		if err = os.WriteFile(cfgWallet, []byte(walletString), 0644); err != nil {
			ExitWithError(err.Error())
		}
	}

	loadWallet()

	_, err = sdk.GetReadPoolInfo(clientWallet.ClientID)
	if err != nil {
		if strings.Contains(err.Error(), "resource_not_found") {
			fmt.Println("Creating related read pool for storage smart-contract...")
			if _, _, err = sdk.CreateReadPool(); err != nil {
				fmt.Printf("Failed to create read pool: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Read pool created successfully")
		}
	}
}

func createWallet() (string, error) {
	walletStr, err := zcncore.CreateWalletOffline()
	if err != nil {
		return "", err
	}

	fmt.Println("ZCN wallet created!!")
	return walletStr, nil
}

func loadWallet() {

	clientBytes, err := ioutil.ReadFile(cfgWallet)
	if err != nil {
		ExitWithError("Error reading the wallet", err)
	}
	clientConfig = string(clientBytes)

	wallet := zcncrypto.Wallet{}
	err = json.Unmarshal([]byte(clientConfig), &wallet)

	if err != nil {
		ExitWithError("Invalid wallet at path:" + cfgWallet)
	}

	clientWallet = &wallet

	wg := &sync.WaitGroup{}
	err = zcncore.SetWalletInfo(clientConfig, signatureScheme, false)
	if err == nil {
		wg.Wait()
	} else {
		ExitWithError(err.Error())
	}
}

func getTxnFee() uint64 {
	return zcncore.ConvertToValue(gTxnFee)
}
