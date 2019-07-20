package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
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
		err := zcncore.RecoverWallet(mnemonic, statusBar)
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
	        <to_client_id> <token> <description> [transaction fee]`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("to_client_id") == false {
			fmt.Println("Error: to_client_id flag is missing")
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
		to_client_id := cmd.Flag("to_client_id").Value.String()
		token, err := cmd.Flags().GetFloat64("token")
		if err != nil {
			fmt.Println("Error: invalid token.", err)
			return
		}
		desc := cmd.Flag("desc").Value.String()
		fee := float64(0)
		fee, err = cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee))
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		err = txn.Send(to_client_id, zcncore.ConvertToValue(token), desc)
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
		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		token := float64(0)
		token, err = cmd.Flags().GetFloat64("token")
		wg.Add(1)
		err = txn.ExecuteSmartContract(zcncore.FaucetSmartContractAddress, methodName, input, zcncore.ConvertToValue(token))
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
	        <tokens> <[durationHr] [durationMin]> [transaction fee]`,
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
		fee := float64(0)
		fee, err = cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee))
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
	        <pool_id> [transaction fee]`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("pool_id") == false {
			fmt.Println("Error: pool_id flag is missing")
			return
		}
		pool_id := cmd.Flag("pool_id").Value.String()
		fee := float64(0)
		fee, err := cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee))
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		err = txn.UnlockTokens(pool_id)
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
var verifycmd = &cobra.Command{
	Use:   "verify",
	Short: "verify transaction",
	Long: `verify transaction.
	        <hash>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("hash") == false {
			fmt.Println("Error: hash flag is missing")
			return
		}
		hash := cmd.Flag("hash").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		txn.SetTransactionHash(hash)
		wg.Add(1)
		err = txn.Verify()
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if statusBar.success {
			statusBar.success = false
			fmt.Printf("\nTransaction verification success\n")
			return
		}
		fmt.Println("\nVerification failed." + statusBar.errMsg + "\n")
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

var getidcmd = &cobra.Command{
	Use:   "getid",
	Short: "Get Miner or Sharder ID from its URL",
	Long:  `Get Miner or Sharder ID from its URL`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("url") == false {
			fmt.Println("Error: url flag is missing")
			return
		}
		url := cmd.Flag("url").Value.String()
		id := zcncore.GetIdForUrl(url)
		if id == "" {
			fmt.Println("Error: ID not found")
			os.Exit(1)
		}
		fmt.Printf("\nURL: %v \nID: %v\n", url, id)
		return
	},
}
var stakecmd = &cobra.Command{
	Use:   "stake",
	Short: "Stake Miners or Sharders",
	Long: `Stake Miners or Sharders using their client ID.
			<client_id> <tokens>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("client_id") == false {
			fmt.Println("Error: client_id flag is missing")
			return
		}
		if fflags.Changed("token") == false {
			fmt.Println("Error: token flag is missing")
			return
		}
		clientID := cmd.Flag("client_id").Value.String()
		token, err := cmd.Flags().GetFloat64("token")
		if err != nil {
			fmt.Println("Error: invalid number of tokens")
			return
		}
		fee := float64(0)
		fee, err = cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee))
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		err = txn.Stake(clientID, zcncore.ConvertToValue(token))
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
				fmt.Println("\nStake success\n")
				return
			}
		}
		fmt.Println("Stake failed. " + statusBar.errMsg)
	},
}
var deletestakecmd = &cobra.Command{
	Use:   "deletestake",
	Short: "Delete Stake from user pool",
	Long: `Delete Stake from user pool client_id and pool_id.
			<client_id> <pool_id>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("client_id") == false {
			fmt.Println("Error: client_id flag is missing")
			return
		}
		if fflags.Changed("pool_id") == false {
			fmt.Println("Error: pool_id flag is missing")
			return
		}
		clientID := cmd.Flag("client_id").Value.String()
		poolID := cmd.Flag("pool_id").Value.String()
		fee := float64(0)
		var err error
		fee, err = cmd.Flags().GetFloat64("fee")
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		txn, err := zcncore.NewTransaction(statusBar, zcncore.ConvertToValue(fee))
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Add(1)
		err = txn.DeleteStake(clientID, poolID)
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
				fmt.Println("\nDelete stake success\n")
				return
			}
		}
		fmt.Println("Delete stake failed. " + statusBar.errMsg)
	},
}
var getuserpoolscmd = &cobra.Command{
	Use:   "getuserpools",
	Short: "Get user pools from sharders",
	Long:  `Get user pools from sharders`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetUserPools(statusBar)
		if err == nil {
			wg.Wait()
		} else {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if statusBar.success {
			fmt.Printf("\nUser pools: %v\n", statusBar.errMsg)
		} else {
			fmt.Println("\nERROR: Get user pool failed. " + statusBar.errMsg + "\n")
		}
		return
	},
}
var getuserpooldetailscmd = &cobra.Command{
	Use:   "getuserpooldetails",
	Short: "Get user pool details",
	Long: `Get user pool details for client_id and pool_id.
			<client_id> <pool_id>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("client_id") == false {
			fmt.Println("Error: client_id flag is missing")
			return
		}
		if fflags.Changed("pool_id") == false {
			fmt.Println("Error: pool_id flag is missing")
			return
		}
		clientID := cmd.Flag("client_id").Value.String()
		poolID := cmd.Flag("pool_id").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetUserPoolDetails(clientID, poolID, statusBar)
		if err != nil {
			fmt.Println(err)
			return
		}
		wg.Wait()
		if statusBar.success {
			fmt.Printf("\nUser pool details: %v\n", statusBar.errMsg)
		} else {
			fmt.Println("\nERROR: Get user pool details failed. " + statusBar.errMsg + "\n")
		}
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
	rootCmd.AddCommand(verifycmd)
	rootCmd.AddCommand(getidcmd)
	rootCmd.AddCommand(stakecmd)
	rootCmd.AddCommand(getuserpoolscmd)
	rootCmd.AddCommand(deletestakecmd)
	rootCmd.AddCommand(getuserpooldetailscmd)
	recoverwalletcmd.PersistentFlags().String("mnemonic", "", "mnemonic")
	recoverwalletcmd.MarkFlagRequired("mnemonic")
	sendcmd.PersistentFlags().String("to_client_id", "", "to_client_id")
	sendcmd.PersistentFlags().Float64("token", 0, "Token to send")
	sendcmd.PersistentFlags().String("desc", "", "Description")
	sendcmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	sendcmd.MarkFlagRequired("to_client_id")
	sendcmd.MarkFlagRequired("token")
	sendcmd.MarkFlagRequired("desc")
	faucetcmd.PersistentFlags().String("methodName", "", "methodName")
	faucetcmd.PersistentFlags().String("input", "", "input")
	faucetcmd.PersistentFlags().Float64("token", 0, "Token request")
	faucetcmd.MarkFlagRequired("methodName")
	faucetcmd.MarkFlagRequired("input")
	lockcmd.PersistentFlags().Float64("token", 0, "Number to tokens to lock")
	lockcmd.PersistentFlags().Int64("durationHr", 0, "Duration Hours to lock")
	lockcmd.PersistentFlags().Int("durationMin", 0, "Duration Mins to lock")
	lockcmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	lockcmd.MarkFlagRequired("token")
	unlockcmd.PersistentFlags().String("pool_id", "", "Poolid - hash of the locked transaction")
	unlockcmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	unlockcmd.MarkFlagRequired("pool_id")
	verifycmd.PersistentFlags().String("hash", "", "hash of the transaction")
	verifycmd.MarkFlagRequired("hash")
	getidcmd.PersistentFlags().String("url", "", "URL to get the ID")
	getidcmd.MarkFlagRequired("url")
	stakecmd.PersistentFlags().String("client_id", "", "Miner or Sharder client id")
	stakecmd.PersistentFlags().Float64("token", 0, "Token to send")
	stakecmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	stakecmd.MarkFlagRequired("client_id")
	stakecmd.MarkFlagRequired("token")
	deletestakecmd.PersistentFlags().String("client_id", "", "Miner or Sharder client id")
	deletestakecmd.PersistentFlags().String("pool_id", "", "Pool ID from user pool matching miner or sharder id")
	deletestakecmd.PersistentFlags().Float64("fee", 0, "Transaction Fee")
	deletestakecmd.MarkFlagRequired("client_id")
	deletestakecmd.MarkFlagRequired("pool_id")
	getuserpooldetailscmd.PersistentFlags().String("client_id", "", "Miner or Sharder client id")
	getuserpooldetailscmd.PersistentFlags().String("pool_id", "", "Pool ID from user pool matching miner or sharder id")
	getuserpooldetailscmd.MarkFlagRequired("client_id")
	getuserpooldetailscmd.MarkFlagRequired("pool_id")
}
