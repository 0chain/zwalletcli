package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/0chain/gosdk/core/zcncrypto"
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

var clientConfig string
var minSubmit int
var minCfm int
var CfmChainLength int

var (
	cfgConfig  *viper.Viper
	cfgNetwork *viper.Viper
	cfgWallet  string
)

var rootCmd = &cobra.Command{
	Use:   "zwallet",
	Short: "Use Zwallet to store, send and execute smart contract on 0Chain platform",
	Long: `Use Zwallet to store, send and execute smart contract on 0Chain platform.
			Complete documentation is available at https://0chain.net`,
	PersistentPreRun: initCmdContext,
}

var clientWallet *zcncrypto.Wallet

func init() {
	cobra.OnInitialize(loadConfigs)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().StringVar(&networkFile, "network", "", "network file to overwrite the network details (if required, default is network.yaml)")
	rootCmd.PersistentFlags().StringVar(&walletFile, "wallet", "", "wallet file (default is wallet.json)")
	rootCmd.PersistentFlags().StringVar(&cDir, "configDir", "", "configuration directory (default is $HOME/.zcn)")
	rootCmd.PersistentFlags().Int64Var(&nonce, "withNonce", 0, "nonce that will be used in transaction (default is 0)")
	rootCmd.PersistentFlags().BoolVar(&bSilent, "silent", false, "Do not print sdk logs in stderr (prints by default)")
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
	configDir = home + "/.zcn"
	return configDir
}

func initZCN() {

	// set the log file
	zcncore.SetLogFile("cmdlog.log", !bSilent)

	miners := cfgNetwork.GetStringSlice("miners")
	sharders := cfgNetwork.GetStringSlice("sharders")
	if len(miners) > 0 && len(sharders) > 0 {
		zcncore.SetNetwork(miners, sharders)
	}

	blockWorker := cfgConfig.GetString("block_worker")
	signScheme := cfgConfig.GetString("signature_scheme")
	chainID := cfgConfig.GetString("chain_id")
	ethereumNodeURL := cfgConfig.GetString("ethereum_node_url")

	err := zcncore.InitZCNSDK(blockWorker, signScheme,
		zcncore.WithChainID(chainID),
		zcncore.WithMinSubmit(minSubmit),
		zcncore.WithMinConfirmation(minCfm),
		zcncore.WithConfirmationChainLength(CfmChainLength),
		zcncore.WithEthereumNode(ethereumNodeURL))
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
	if &cfgFile != nil && len(cfgFile) > 0 {
		cfgConfig.SetConfigFile(configDir + "/" + cfgFile)
	} else {
		cfgConfig.SetConfigFile(configDir + "/" + "config.yaml")
	}

	minSubmit = cfgConfig.GetInt("min_submit")
	minCfm = cfgConfig.GetInt("min_confirmation")
	CfmChainLength = cfgConfig.GetInt("confirmation_chain_length")

	// ~/.zcn/network.yaml
	cfgNetwork.AddConfigPath(configDir)
	if &networkFile != nil && len(networkFile) > 0 {
		cfgNetwork.SetConfigFile(configDir + "/" + networkFile)
	} else {
		cfgNetwork.SetConfigFile(configDir + "/" + "network.yaml")
	}

	if err := cfgConfig.ReadInConfig(); err != nil {
		ExitWithError("Can't read config:", err)
	}

	cfgNetwork.ReadInConfig() //nolint: errcheck

	// TODO: move the private key storage to the keychain or secure storage
	// ~/.zcn/wallet.json
	if &walletFile != nil && len(walletFile) > 0 {
		cfgWallet = configDir + "/" + walletFile
	} else {
		cfgWallet = configDir + "/wallet.json"
	}
}

var zcnIsConnected bool
var hasWallet bool

func initCmdContext(cmd *cobra.Command, args []string) {

	// connect to zcn
	if !zcnIsConnected {
		_, ok := withoutZCNCmds[cmd]
		if !ok {
			initZCN()
			zcnIsConnected = true
		}
	}

	// create wallet
	if !hasWallet {
		_, ok := withoutWalletCmds[cmd]
		if !ok {
			createWallet()
			loadWallet()

			hasWallet = true
		}
	}

}

func createWallet() {
	if _, err := os.Stat(cfgWallet); os.IsNotExist(err) {
		fmt.Println("No wallet in path ", cfgWallet, "found. Creating wallet...")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}

		wg.Add(1)

		err = zcncore.CreateWallet(statusBar)
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}

		if len(statusBar.walletString) == 0 || !statusBar.success {
			ExitWithError("Error creating the wallet." + statusBar.errMsg)
		}

		fmt.Println("ZCN wallet created!!")

		err = os.WriteFile(cfgWallet, []byte(statusBar.walletString), 0644)
		if err != nil {
			ExitWithError("Error creating the wallet." + err.Error())
		}

		log.Print("Creating related read pool for storage smart-contract...")
		if err = createReadPool(); err != nil {
			log.Fatalf("Failed to create read pool: %v", err)
		}
		log.Printf("Read pool created successfully")

	}
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
	err = zcncore.SetWalletInfo(clientConfig, false)
	if err == nil {
		wg.Wait()
	} else {
		ExitWithError(err.Error())
	}
}
