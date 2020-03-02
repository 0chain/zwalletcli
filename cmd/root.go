package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/0chain/gosdk/core/zcncrypto"
	"github.com/0chain/gosdk/zcncore"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var walletFile string
var cDir string
var bVerbose bool

var sharders []string
var miners []string
var clientConfig string
var minSubmit int
var minCfm int
var CfmChainLength int

var rootCmd = &cobra.Command{
	Use:   "zwallet",
	Short: "Use Zwallet to store, send and execute smart contract on 0Chain platform",
	Long: `Use Zwallet to store, send and execute smart contract on 0Chain platform.
			Complete documentation is available at https://0chain.net`,
}

var clientWallet *zcncrypto.Wallet

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")
	rootCmd.PersistentFlags().StringVar(&walletFile, "wallet", "", "wallet file (default is wallet.json)")
	rootCmd.PersistentFlags().StringVar(&cDir, "configDir", "", "configuration directory (default is $HOME/.zcn)")
	rootCmd.PersistentFlags().BoolVar(&bVerbose, "verbose", false, "prints sdk log in stderr (default false)")
	initConfig()
	fmt.Printf("%s", cfgFile)
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
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	configDir = home + "/.zcn"
	return configDir
}

func initConfig() {
	nodeConfig := viper.New()
	var configDir string
	if cDir != "" {
		configDir = cDir
	} else {
		configDir = getConfigDir()
	}
	nodeConfig.AddConfigPath(configDir)
	if &cfgFile != nil && len(cfgFile) > 0 {
		nodeConfig.SetConfigName(cfgFile)
	} else {
		nodeConfig.SetConfigName("config")
	}
	if err := nodeConfig.ReadInConfig(); err != nil {
		ExitWithError("Can't read config:", err)
	}
	sharders = nodeConfig.GetStringSlice("sharders")
	miners = nodeConfig.GetStringSlice("miners")
	signScheme := nodeConfig.GetString("signature_scheme")
	chainID := nodeConfig.GetString("chain_id")
	minSubmit = nodeConfig.GetInt("min_submit")
	minCfm = nodeConfig.GetInt("min_confirmation")
	CfmChainLength = nodeConfig.GetInt("confirmation_chain_length")
	//chainID := nodeConfig.GetString("chain_id")

	//TODO: move the private key storage to the keychain or secure storage
	var walletFilePath string
	if &walletFile != nil && len(walletFile) > 0 {
		walletFilePath = configDir + "/" + walletFile
	} else {
		walletFilePath = configDir + "/wallet.json"
	}
	//set the log file
	zcncore.SetLogFile("cmdlog.log", bVerbose)
	err := zcncore.InitZCNSDK(miners, sharders, signScheme,
		zcncore.WithChainID(chainID),
		zcncore.WithMinSubmit(minSubmit),
		zcncore.WithMinConfirmation(minCfm),
		zcncore.WithConfirmationChainLength(CfmChainLength))
	if err != nil {
		ExitWithError(err.Error())
	}

	if _, err := os.Stat(walletFilePath); os.IsNotExist(err) {
		fmt.Println("No wallet in path ", walletFilePath, "found. Creating wallet...")
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
		clientConfig = string(statusBar.walletString)
		file, err := os.Create(walletFilePath)
		if err != nil {
			ExitWithError(err.Error())
		}
		defer file.Close()
		fmt.Fprintf(file, clientConfig)
	} else {
		f, err := os.Open(walletFilePath)
		if err != nil {
			ExitWithError("Error opening the wallet", err)
		}
		clientBytes, err := ioutil.ReadAll(f)
		if err != nil {
			ExitWithError("Error reading the wallet", err)
		}
		clientConfig = string(clientBytes)
	}

	wallet := &zcncrypto.Wallet{}
	err = json.Unmarshal([]byte(clientConfig), wallet)
	clientWallet = wallet
	if err != nil {
		ExitWithError("Invalid wallet at path:" + walletFilePath)
	}
	wg := &sync.WaitGroup{}
	err = zcncore.SetWalletInfo(clientConfig, false)
	if err == nil {
		wg.Wait()
	} else {
		ExitWithError(err.Error())
	}
}
