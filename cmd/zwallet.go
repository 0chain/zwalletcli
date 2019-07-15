package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"sync"
)

var recoverwalletcmd = &cobra.Command{
	Use:   "recoverwallet",
	Short: "Recover wallet",
	Long:  `Recover wallet from the previously stored mnemonic`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("mnemonic") == false {
			fmt.Println("Error: Mnemonic not provided")
			return
		}
		mnemonic := cmd.Flag("mnemonic").Value.String()
		if zcncore.IsMnemonicValid(mnemonic) == false {
			fmt.Println("Error: Invalid mnemonic")
			return
		}
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.RecoverWallet(mnemonic, numKeys, statusBar)
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if len(statusBar.walletString) == 0 || !statusBar.success {
			fmt.Println("Error recovering the wallet." + statusBar.errMsg)
			os.Exit(1)
		}
		var walletFilePath string
		if &walletFile != nil && len(walletFile) > 0 {
			walletFilePath = getConfigDir() + "/" + walletFile
		} else {
			walletFilePath = getConfigDir() + "/wallet.txt"
		}
		clientConfig = string(statusBar.walletString)
		file, err := os.Create(walletFilePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		defer file.Close()
		fmt.Fprintf(file, clientConfig)
		fmt.Println("Wallet recovered!!")
		return
	},
}

var getbalancecmd = &cobra.Command{
	Use:   "getbalance",
	Short: "Get balance from sharders",
	Long:  `Get balance from sharders`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetBalance(statusBar)
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if statusBar.success {
			fmt.Printf("\nBalance: %v\n", zcncore.ConvertToToken(statusBar.balance))
		} else {
			fmt.Println("\nGet balance failed. " + statusBar.errMsg + "\n")
		}
		return
	},
}

var sendcmd = &cobra.Command{
	Use:   "send",
	Short: "Send ZCN token to another wallet",
	Long: `Send ZCN token to another wallet.
	        <toclientID> <token> <description>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("toclientID") == false {
			fmt.Println("Error: toclientID flag is missing")
			return
		}
		if fflags.Changed("token") == false {
			fmt.Println("Error: token flag is missing")
			return
		}
		if fflags.Changed("desc") == false {
			fmt.Println("Error: Description flag is missing")
			return
		}
		toclientID := cmd.Flag("toclientID").Value.String()
		token, err := cmd.Flags().GetFloat64("token")
		desc := cmd.Flag("desc").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar)
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		err = txn.Send(toclientID, zcncore.ConvertToValue(token), desc)
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if statusBar.success {
			statusBar.success = false
			wg.Add(1)
			err := txn.Verify()
			if err == nil {
				wg.Wait()
			} else {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if statusBar.success {
				fmt.Println("Send token success")
				return
			}
		}
		fmt.Println("Send token failed. " + statusBar.errMsg)
		return
	},
}

var faucetcmd = &cobra.Command{
	Use:   "faucet",
	Short: "Faucet smart contract",
	Long: `Faucet smart contract.
	        <methodName> <input>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("methodName") == false {
			fmt.Println("Error: Methodname flag is missing")
			return
		}
		if fflags.Changed("input") == false {
			fmt.Println("Error: Input flag is missing")
			return
		}

		methodName := cmd.Flag("methodName").Value.String()
		input := cmd.Flag("input").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar)
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		err = txn.ExecuteFaucetSC(methodName, []byte(input))
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if statusBar.success {
			statusBar.success = false
			wg.Add(1)
			err := txn.Verify()
			if err == nil {
				wg.Wait()
			} else {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if statusBar.success {
				fmt.Println("Execute faucet smart contract success")
				return
			}
		}
		fmt.Println("\nExecute faucet smart contract failed. " + statusBar.errMsg + "\n")
		return
	},
}

var lockcmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock tokens",
	Long: `Lock tokens .
	        <tokens> <[durationHr] [durationMin]>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("token") == false {
			fmt.Println("Error: token flag is missing")
			return
		}
		if (fflags.Changed("durationHr") == false) &&
			(fflags.Changed("durationMin") == false) {
			fmt.Println("Error: durationHr and durationMin flag is missing. atleast one is required")
			return
		}
		token, err := cmd.Flags().GetFloat64("token")
		if err != nil {
			fmt.Println("Error: invalid number of tokens")
			return
		}
		durationHr := int64(0)
		durationHr, err = cmd.Flags().GetInt64("durationHr")
		durationMin := int(0)
		durationMin, err = cmd.Flags().GetInt("durationMin")
		if (durationHr < 1) && (durationMin < 1) {
			fmt.Println("Error: invalid duration")
		}
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar)
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		err = txn.LockTokens(zcncore.ConvertToValue(token), durationHr, durationMin)
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if statusBar.success {
			statusBar.success = false
			wg.Add(1)
			err := txn.Verify()
			if err == nil {
				wg.Wait()
			} else {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if statusBar.success {
				fmt.Printf("\nTokens (%f) locked successfully\n", token)
				return
			}
		}
		fmt.Println("\nFailed to lock tokens. " + statusBar.errMsg + "\n")
		return
	},
}

var unlockcmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock tokens",
	Long: `Unlock previously locked tokens .
	        <poolid>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("poolid") == false {
			fmt.Println("Error: poolid flag is missing")
			return
		}
		poolid := cmd.Flag("poolid").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar)
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		err = txn.UnlockTokens(poolid)
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if statusBar.success {
			statusBar.success = false
			wg.Add(1)
			err := txn.Verify()
			if err == nil {
				wg.Wait()
			} else {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			if statusBar.success {
				fmt.Printf("\nUnlock token success\n")
				return
			}
		}
		fmt.Println("\nFailed to unlock tokens. " + statusBar.errMsg + "\n")
		return
	},
}

var lockconfigcmd = &cobra.Command{
	Use:   "lockconfig",
	Short: "Get lock configuration",
	Long:  `Get lock configuration`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetLockConfig(statusBar)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wg.Wait()
		if statusBar.success {
			fmt.Printf("\nConfiguration:\n %v\n", statusBar.errMsg)
			return
		}
		fmt.Println("\nFailed to get lock config." + statusBar.errMsg + "\n")
		return
	},
}

var getlockedtokenscmd = &cobra.Command{
	Use:   "getlockedtokens",
	Short: "Get locked tokens",
	Long:  `Get locked tokens`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetLockedTokens(statusBar)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		wg.Wait()
		if statusBar.success {
			fmt.Printf("\nLocked tokens:\n %v\n", statusBar.errMsg)
			return
		}
		fmt.Println("\nFailed to get locked tokens." + statusBar.errMsg + "\n")
		return
	},
}

var createmswalletcmd = &cobra.Command{
	Use:   "createmswallet",
	Short: "create multisig wallet",
	Long:  `create multisig wallet`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.CreateMSWallet(statusBar)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		wg.Wait()
		if statusBar.success {

			err := registerMSWallets(statusBar.wallets)
			if err != nil {
				fmt.Printf("Error while registering ms sub wallets. The error is:\n %v\n", err)
				return
			}

			err = registerMultiSig(statusBar.wallets[0], statusBar.walletString)
			if err != nil {
				fmt.Printf("Error while registering ms group wallet. The error is:\n %v\n", err)
				return
			}

			var walletFilePath string
			if &walletFile != nil && len(walletFile) > 0 {
				walletFilePath = getConfigDir() + "/" + walletFile
			} else {
				walletFilePath = getConfigDir() + "/msgrpwallet.txt"
			}
			writeToaFile(walletFilePath, statusBar.walletString)
			fmt.Printf("\nMS Wallet created and saved in:\n %s\n", walletFilePath)
			fileName := getConfigDir() + "/msgroupwallet.txt"

			writeToaFile(fileName, statusBar.wallets[0])
			fmt.Printf("Created file:%v\n\n", fileName)
			for i := 1; i < len(statusBar.wallets); i++ {
				fileName := fmt.Sprintf("%s/mssubwallet%d.txt", getConfigDir(), i)
				writeToaFile(fileName, statusBar.wallets[i])
				fmt.Printf("Created file: %v\n\n", fileName)
			}

			return
		}

		fmt.Println("\nFailed to create MS Wallet." + statusBar.errMsg + "\n")
		return
	},
}

func init() {
	rootCmd.AddCommand(recoverwalletcmd)
	rootCmd.AddCommand(getbalancecmd)
	rootCmd.AddCommand(sendcmd)
	rootCmd.AddCommand(faucetcmd)
	rootCmd.AddCommand(lockcmd)
	rootCmd.AddCommand(unlockcmd)
	rootCmd.AddCommand(lockconfigcmd)
	rootCmd.AddCommand(getlockedtokenscmd)
	rootCmd.AddCommand(createmswalletcmd)
	recoverwalletcmd.PersistentFlags().String("mnemonic", "", "mnemonic")
	recoverwalletcmd.MarkFlagRequired("mnemonic")
	sendcmd.PersistentFlags().String("toclientID", "", "toclientID")
	sendcmd.PersistentFlags().Float64("token", 0, "Token to send")
	sendcmd.PersistentFlags().String("desc", "", "Description")
	sendcmd.MarkFlagRequired("toclientID")
	sendcmd.MarkFlagRequired("token")
	sendcmd.MarkFlagRequired("desc")
	faucetcmd.PersistentFlags().String("methodName", "", "methodName")
	faucetcmd.PersistentFlags().String("input", "", "input")
	faucetcmd.MarkFlagRequired("methodName")
	faucetcmd.MarkFlagRequired("input")
	lockcmd.PersistentFlags().Float64("token", 0, "Number to tokens to lock")
	lockcmd.PersistentFlags().Int64("durationHr", 0, "Duration Hours to lock")
	lockcmd.PersistentFlags().Int("durationMin", 0, "Duration Mins to lock")
	lockcmd.MarkFlagRequired("token")
	unlockcmd.PersistentFlags().String("poolid", "", "Poolid - hash of the locked transaction")
	unlockcmd.MarkFlagRequired("poolid")
}

func readFile(fileName string) (string, error) {
	w, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(w), nil
}

func registerMultiSig(grw string, msw string) error {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	txn, err := zcncore.NewMSTransaction(grw, statusBar)
	if err != nil {
		fmt.Println(err)
		return err
	}
	wg.Add(1)
	err = txn.RegisterMultiSig(grw, msw)
	if err == nil {
		wg.Wait()
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if statusBar.success {
		fmt.Printf("\nMultisig wallet SC registration requested. verifying status")
		statusBar.success = false
		wg.Add(1)
		err := txn.Verify()
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println("error in verifying: ", err.Error())
			os.Exit(1)
		}
		if statusBar.success {
			fmt.Printf("\nMultisigSC success\n")
			return nil
		}
	}
	fmt.Println("\nFailed to register multisigsc. " + statusBar.errMsg + "\n")
	return fmt.Errorf(statusBar.errMsg)

}

func registerAWallet(w string) error {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	wg.Add(1)
	zcncore.RegisterWallet(w, statusBar)
	wg.Wait()
	if statusBar.success {
		return nil
	}
	return fmt.Errorf(statusBar.errMsg)

}

func registerMSWallets(wallets []string) error {

	fmt.Printf("\n registering %d wallets \n", len(wallets))
	i := 0
	for _, wallet := range wallets {

		var walletName string
		if i == 0 {
			walletName = "group wallet"
		} else {
			walletName = fmt.Sprintf("sub wallet number %d \n", i)
		}
		err := registerAWallet(wallet)
		if err != nil {
			fmt.Printf("\nFailed ot register wallet number %s\nAborting...", walletName)
			return err
		}
		fmt.Printf("\nSuccessfully registered %s\n", walletName)

		i++
	}
	return nil
}
func writeToaFile(fileNameAndPath string, content string) error {

	file, err := os.Create(fileNameAndPath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer file.Close()
	fmt.Fprintf(file, content)
	return nil
}
