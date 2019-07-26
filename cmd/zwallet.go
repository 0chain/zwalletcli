package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	MaxExpiryTime = 60 * 60 * 24 * 7 * 30 // Max expiry time is 30 days.
	MinExpiryTime = 1                     // Min expiry time is 1 second.

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
	Long: `create multisig wallet
			<numsigners> <threshold> <testN>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		MaxSigners := 20 //This is the limitation from MultiSigSC
		MinSigners := 2  //This is the limitation from MultiSigSC
		fflags := cmd.Flags()
		if fflags.Changed("numsigners") == false {
			fmt.Println("Error: numsigners flag is missing")
			return
		}
		numsigners, err := cmd.Flags().GetInt("numsigners")
		if err != nil {
			fmt.Println("Error: numsubkeys is not an integer")
			return
		}
		if numsigners > MaxSigners {
			fmt.Printf("Error: too many signers. Maximum numsigners allowed is %d\n", MaxSigners)
			return
		}

		if numsigners < MinSigners {
			fmt.Printf("Error: too few signers. Minimum numsigners required is %d\n", MinSigners)
			return
		}

		if fflags.Changed("threshold") == false {
			fmt.Println("Error: threshold flag is missing")
			return
		}
		threshold, err := cmd.Flags().GetInt("threshold")
		if err != nil {
			fmt.Println("Error: threshold is not an integer")
			return
		}
		if threshold > numsigners {
			fmt.Printf("Error: given threshold (%d) is too high. Threshold has to be less than or equal to numsigners (%d)\n", threshold, numsigners)
			return
		}

		delay, err := cmd.Flags().GetInt("delay")
		if err != nil {
			fmt.Println("Error: delay is not an integer")
			return
		}
		if delay < 0 {
			delay = 0
		}

		if fflags.Changed("expiry") == false {
			fmt.Println("Error: expiry flag is missing")
			return
		}
		expiry, err := cmd.Flags().GetInt64("expiry")
		if err != nil {
			fmt.Printf("Error: expiry is not an integer.\n%v\n", err)
			return
		}

		if expiry < 1 {
			fmt.Printf("Error: expiry %d is invalid. It needs to be at least 1 second.\n", expiry)
			return
		}

		if expiry > MaxExpiryTime {
			fmt.Printf("Error: expiry %d is invalid. Max expiry time allowed is 18144000.\n", expiry)
			return
		}

		testN, err := cmd.Flags().GetBool("testn")
		if err != nil {
			fmt.Println("testn is not used or not set to true. Setting it to false")
		}

		stress, err := cmd.Flags().GetInt("stress")
		if err != nil {
			fmt.Println("stress is not used or not set to true. Setting it to 1")
		}

		wg := &sync.WaitGroup{}

		runner := func() {
			defer wg.Done()
			smsw, groupClientID, wallets, err := zcncore.CreateMSWallet(threshold, numsigners, expiry)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			//register all wallets
			err = registerMSWallets(wallets)
			if err != nil {
				fmt.Printf("Error while registering ms sub wallets. The error is:\n %v\n", err)
				return
			}

			groupWallet := wallets[0]
			signerWallets := wallets[1:]

			err = registerMultiSig(groupWallet, smsw)
			if err != nil {
				fmt.Printf("Error while registering ms group wallet. The error is:\n %v\n", err)
				return
			}

			//if !testMSVoting(msw, grpWallet, grpClientID, signerWallets, threshold, testN) {
			if !testMSVoting(smsw, groupWallet, groupClientID, signerWallets, threshold, delay, testN) {
				fmt.Printf("Failed to test voting\n")
				return
			}
			fmt.Printf("Creating and testing a multisig wallet is successful!\n\n")
			return
		}

		fmt.Printf("running tests for %d\n", stress)

		for i := 1; i < stress+1; i++ {
			fmt.Printf("Launching test#%d\n", i)
			go runner()
			wg.Add(1)
			time.Sleep(1 * time.Second)
		}
		wg.Wait()

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
	createmswalletcmd.PersistentFlags().Int("numsigners", 0, "Number of signers")
	createmswalletcmd.PersistentFlags().Int("threshold", 0, "Threshold number of signers required to sign the proposal")
	createmswalletcmd.PersistentFlags().Int64("expiry", 0, "Expiration time for proposals in seconds. Min 1 and Max 1814400 (30 days)")
	createmswalletcmd.PersistentFlags().Int("delay", 0, "Delay time in seconds before proposals")
	createmswalletcmd.PersistentFlags().Int("stress", 1, "Stress test Multisig with wallets")
	createmswalletcmd.PersistentFlags().Bool("testn", false, "test Multiwallet with all signers. Default is false")
	createmswalletcmd.MarkFlagRequired("expiry")
	createmswalletcmd.MarkFlagRequired("threshold")
	createmswalletcmd.MarkFlagRequired("numsigners")
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
			if statusBar.success {
				fmt.Printf("\nMultisigSC  wallet SC registration request success\n")
				return nil
			}
			fmt.Printf("\nMultisigSC wallet SC registration request failed\n%s", statusBar.errMsg)
			return nil

		}
		fmt.Println("error in verifying multisig wallet registration: ", err.Error())
		os.Exit(1)

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
			walletName = fmt.Sprintf("signer wallet number %d \n", i)
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

func registerMSVote(signerWalletStr string, voteStr string) error {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	txn, err := zcncore.NewMSTransaction(signerWalletStr, statusBar)
	if err != nil {
		fmt.Println(err)
		return err
	}
	wg.Add(1)
	err = txn.RegisterVote(signerWalletStr, voteStr)
	if err == nil {
		wg.Wait()
	} else {
		fmt.Println(err.Error())
		return err
	}
	if statusBar.success {
		fmt.Printf("\nMultisig Vote registration requested. verifying status")
		statusBar.success = false
		wg.Add(1)
		err := txn.Verify()

		if err == nil {
			wg.Wait()
		} else {
			fmt.Println("error in verifying: ", err.Error())
			return err
		}
		if statusBar.success {
			fmt.Printf("\nMultisig Voting success\n")
			return nil
		}
	}
	fmt.Println("\nFailed to register multisig vote. " + statusBar.errMsg + "\n")
	return fmt.Errorf(statusBar.errMsg)

}

func createAWallet() string {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	wg.Add(1)
	err := zcncore.CreateWallet(1, statusBar)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wg.Wait()
	if statusBar.success {
		return statusBar.walletString
	}
	return ""

}

func checkBalance(wallet string) bool {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	wg.Add(1)
	err := zcncore.GetBalanceWallet(wallet, statusBar)
	if err == nil {
		wg.Wait()
	} else {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if statusBar.success {
		fmt.Printf("\nBalance: %v\n", zcncore.ConvertToToken(statusBar.balance))
		if zcncore.ConvertToToken(statusBar.balance) > 0 {
			return true
		}
		return false

	}
	fmt.Println("\nGet balance failed. " + statusBar.errMsg + "\n")
	return false

}

func pourToWallet(wallet string) bool {
	methodName := "pour"
	input := "{fillwallet}"
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	txn, err := zcncore.NewMSTransaction(wallet, statusBar)
	if err != nil {
		fmt.Println(err)
		return false
	}
	wg.Add(1)
	err = txn.ExecuteFaucetSCWallet(wallet, methodName, []byte(input))
	if err == nil {
		wg.Wait()
	} else {
		fmt.Println(err.Error())
		return false
	}
	if statusBar.success {
		statusBar.success = false
		wg.Add(1)
		err := txn.Verify()
		if err == nil {
			wg.Wait()
		} else {
			fmt.Printf("error in faucet transaction:\n%v\n", err.Error())
			return false
		}
		if statusBar.success {
			fmt.Printf("Pour request success\n")
			b := checkBalance(wallet)
			return b
		}
		fmt.Printf("Pour request failed\n")

	}
	return false
}
func testMSVoting(msw string, groupWallet string, groupClientID string, signerWallets []string, t int, delay int, testN bool) bool {
	fmt.Printf("\n\ntesting vote")
	anoWallet := createAWallet()

	fmt.Printf("Recipient test wallet:\n%s\n", anoWallet)
	fmt.Printf("\nActivating group wallet by pouring test tokens\n")
	if !pourToWallet(groupWallet) {
		fmt.Printf("pour failed, for group wallet...")
		return false

	}

	for i, wallet := range signerWallets {
		fmt.Printf("\nActivating signer wallet %d by pouring test tokens\n", i+1)
		if !pourToWallet(wallet) {
			fmt.Printf("pour failed for a signer wallet")
			return false
		}
	}

	fmt.Printf("Checking balance on group wallet with clientID %s before the vote", groupClientID)
	checkBalance(groupWallet)

	toClientID, err := zcncore.GetWalletClientID(anoWallet)
	if err != nil {
		fmt.Printf("Failed to get clientID from the wallet\n%v\nError is:%v\n", anoWallet, err)
		return false
	}

	if !testN {
		if !testMSVotingThreshold(msw, toClientID, groupClientID, signerWallets, t, delay) {
			fmt.Printf("\n Failed in MSVoting test for threshold\n")
			return false
		}
	} else {
		if !testMSVotingForAllN(msw, toClientID, groupClientID, signerWallets) {
			fmt.Printf("\n Failed in MSVoting test for threshold\n")
			return false
		}
	}

	fmt.Printf("\n\nChecking balance on group wallet %s after the vote", groupClientID)
	checkBalance(groupWallet)

	fmt.Printf("\nChecking balance on recipient wallet after the vote")
	checkBalance(anoWallet)
	return true

}

func testMSVotingThreshold(msw string, toClientID string, grpClientID string, signerWallets []string, t int, delay int) bool {

	proposal := "testing MSVoting"
	tokenVal := zcncore.ConvertToValue(0.1)

	cnt := 0

	for _, idx := range rand.Perm(t) {
		signer := signerWallets[idx]

		//for _, signer := range signerWallets {
		if cnt >= t {
			break
		}

		vote, err := zcncore.CreateMSVote(proposal, grpClientID, signer, toClientID, tokenVal)
		if err != nil {
			fmt.Printf("Failed to create a vote. Error is:%v\n", err)
			return false
		}
		fmt.Printf("\nCreated Vote#%d from signer #%d:\n%s\n", cnt+1, idx, vote)
		err = registerMSVote(signer, vote)
		if err != nil {
			fmt.Printf("Failed to create a vote. Error is:%v\n", err)
			return false
		}
		cnt++
		fmt.Printf("sleeping for %d seconds...", delay)
		time.Sleep(time.Duration(delay) * time.Second)
	}

	return true
}

func testMSVotingForAllN(msw string, toClientID string, grpClientID string, signerWallets []string) bool {

	proposal := "testing MSVoting"
	tokenVal := zcncore.ConvertToValue(0.1)

	cnt := 0
	for _, signer := range signerWallets {
		vote, err := zcncore.CreateMSVote(proposal, grpClientID, signer, toClientID, tokenVal)
		if err != nil {
			fmt.Printf("Failed to create a vote. Error is:%v\n", err)
			return false
		}
		fmt.Printf("\nCreated Vote#%d:\n%s\n", cnt+1, vote)
		err = registerMSVote(signer, vote)
		if err != nil {
			fmt.Printf("Failed to create a vote. Error is:%v\n", err)
			return false
		}
		cnt++

	}
	return true
}
