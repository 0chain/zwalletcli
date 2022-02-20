package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

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
		if !fflags.Changed("numsigners") {
			ExitWithError("Error: numsigners flag is missing")
		}
		numsigners, err := cmd.Flags().GetInt("numsigners")
		if err != nil {
			ExitWithError("Error: numsubkeys is not an integer")
		}
		if numsigners > MaxSigners {
			ExitWithError(fmt.Sprintf("Error: too many signers. Maximum numsigners allowed is %d\n", MaxSigners))
		}

		if numsigners < MinSigners {
			ExitWithError(fmt.Sprintf("Error: too few signers. Minimum numsigners required is %d\n", MinSigners))
		}

		if !fflags.Changed("threshold") {
			ExitWithError("Error: threshold flag is missing")
		}

		threshold, err := cmd.Flags().GetInt("threshold")
		if err != nil {
			ExitWithError("Error: threshold is not an integer")
		}
		if threshold > numsigners {
			ExitWithError(fmt.Sprintf("Error: given threshold (%d) is too high. Threshold has to be less than or equal to numsigners (%d)\n", threshold, numsigners))
		}
		if threshold <= 0 {
			ExitWithError("Error: threshold should be bigger than 0")
		}

		testN, err := cmd.Flags().GetBool("testn")
		if err != nil {
			fmt.Println("testn is not used or not set to true. Setting it to false")
		}

		smsw, groupClientID, wallets, err := zcncore.CreateMSWallet(threshold, numsigners)
		if err != nil {
			ExitWithError(err)
		}

		//register all wallets
		err = registerMSWallets(wallets)
		if err != nil {
			ExitWithError(fmt.Sprintf("Error while registering ms sub wallets. The error is:\n %v\n", err))
		}

		groupWallet := wallets[0]
		signerWallets := wallets[1:]

		err = registerMultiSig(groupWallet, smsw)
		if err != nil {
			ExitWithError(fmt.Sprintf("Error while registering ms group wallet. The error is:\n %v\n", err))
		}

		//if !testMSVoting(msw, grpWallet, grpClientID, signerWallets, threshold, testN) {
		if !testMSVoting(smsw, groupWallet, groupClientID, signerWallets, threshold, testN) {
			ExitWithError("Failed to test voting\n")
		}
		fmt.Printf("\nCreating and testing a multisig wallet is successful!\n\n")
		return
	},
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
			switch txn.GetVerifyConfirmationStatus() {
			case zcncore.ChargeableError:
				ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
			case zcncore.Success:
				fmt.Printf("Pour request success\n")
				b := checkBalance(wallet)
				return b
			default:
				ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
		} else {
			fmt.Printf("Pour request failed\n")
		}

	}
	return false
}
func testMSVoting(msw string, groupWallet string, groupClientID string, signerWallets []string, t int, testN bool) bool {
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
		if !testMSVotingThreshold(msw, toClientID, groupClientID, signerWallets, t) {
			fmt.Printf("/nFailed in MSVoting test for threshold\n")
			return false
		}
	} else {
		if !testMSVotingForAllN(msw, toClientID, groupClientID, signerWallets) {
			fmt.Printf("/nFailed in MSVoting test for threshold\n")
			return false
		}
	}

	fmt.Printf("\n\nChecking balance on group wallet %s after the vote", groupClientID)
	checkBalance(groupWallet)

	fmt.Printf("\nChecking balance on recipient wallet after the vote")
	checkBalance(anoWallet)
	return true

}

func testMSVotingThreshold(msw string, toClientID string, grpClientID string, signerWallets []string, t int) bool {

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

			if statusBar.success {
				switch txn.GetVerifyConfirmationStatus() {
				case zcncore.ChargeableError:
					ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
				case zcncore.Success:
					fmt.Printf("\nMultisigSC  wallet SC registration request success\n")
				default:
					ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
						strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
				}
				return nil
			} else {
				fmt.Printf("\nMultisigSC wallet SC registration request failed\n")
				return nil
			}

		}
		fmt.Println("error in verifying multisig wallet registration: ", err.Error())
		os.Exit(1)

	}
	fmt.Println("\nFailed to register multisigsc. " + statusBar.errMsg + "\n")
	return fmt.Errorf(statusBar.errMsg)
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
			switch txn.GetVerifyConfirmationStatus() {
			case zcncore.ChargeableError:
				ExitWithError("\n", strings.Trim(txn.GetVerifyOutput(), "\""))
			case zcncore.Success:
				fmt.Printf("\nMultisig Voting success\n")
			default:
				ExitWithError("\nExecute global settings update smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}

			return nil
		} else {
			fmt.Printf("Pour request failed\n")
		}

	}
	fmt.Println("\nFailed to register multisig vote. " + statusBar.errMsg + "\n")
	return fmt.Errorf(statusBar.errMsg)

}

func checkBalance(wallet string) bool {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	wg.Add(1)
	err := zcncore.GetBalanceWallet(wallet, statusBar)
	if err != nil {
		ExitWithError(err)
		return false
	}
	wg.Wait()
	if !statusBar.success {
		ExitWithError(fmt.Sprintf("\nFailed to get balance: %s\n", statusBar.errMsg))
		return false
	}
	fmt.Printf("\nBalance: %v\n", statusBar.balance)
	return statusBar.balance.ToToken() > 0
}

func createAWallet() string {
	wg := &sync.WaitGroup{}
	statusBar := &ZCNStatus{wg: wg}
	wg.Add(1)
	err := zcncore.CreateWallet(statusBar)
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

func init() {
	rootCmd.AddCommand(createmswalletcmd)
	createmswalletcmd.PersistentFlags().Int("numsigners", 0, "Number of signers")
	createmswalletcmd.PersistentFlags().Int("threshold", 0, "Threshold number of signers required to sign the proposal")
	createmswalletcmd.PersistentFlags().Bool("testn", false, "test Multiwallet with all signers. Default is false")
	createmswalletcmd.MarkFlagRequired("threshold")
	createmswalletcmd.MarkFlagRequired("numsigners")
}
