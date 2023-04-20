package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var createWalletCmd = &cobra.Command{
	Use:   "create-wallet",
	Short: "Create wallet and logs it into stdout (pass --register to register wallet to blockchain)",
	Long:  `Create wallet and logs it into standard output (pass --register to register wallet to blockchain)`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		online, err := cmd.Flags().GetBool("register")
		if err != nil {
			ExitWithError("invalid register flag")
		}

		var walletStr string
		if online {
			statusBar, err := createWallet()
			if err != nil {
				ExitWithError(fmt.Printf("Failed to create wallet: %v", err))
			}
			walletStr = statusBar.walletString

			// update gosdk wallet to use the newly generate wallet for following operation
			if err = zcncore.SetWalletInfo(walletStr, false); err != nil {
				ExitWithError("failed to use new wallet", err)
			}

			log.Print("Creating related read pool for storage smart-contract...")
			if err := createReadPool(); err != nil {
				log.Fatalf("Failed to create read pool: %v", err)
			}
			log.Printf("Read pool created successfully")

		} else {
			walletStr, err = zcncore.CreateWalletOffline()
			if err != nil {
				ExitWithError("failed to generate offline wallet", err)
			}
		}

		// write wallet into wallet dir
		filename := walletFilename()
		if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
			// same wallet exists
			ExitWithError(fmt.Sprintf("unable to write wallet, file with %q name already exists", filename))
		}

		if err := os.WriteFile(filename, []byte(walletStr), 0644); err != nil {
			// no return just print it
			fmt.Fprintf(os.Stderr, "failed to dump wallet into zcn home directory %v", err)
		} else {
			log.Printf("wallet saved in %s\n", filename)
		}

		if !bSilent {
			fmt.Fprintf(os.Stdout, walletStr)
		}
	},
}

func init() {
	rootCmd.AddCommand(WithoutWallet(createWalletCmd))
	createWalletCmd.PersistentFlags().Bool("register boolean", false, "create wallet with registration on blockchain (default false)")
	createWalletCmd.PersistentFlags().Bool("silent boolean", false, "do not print wallet details in the standard output (default false)")
	createWalletCmd.PersistentFlags().String("wallet string", "", "give custom name to the wallet")
}

func walletFilename() string {
	cfgDir := getConfigDir()
	if len(walletFile) > 0 {
		return filepath.Join(cfgDir, walletFile)
	}
	now := time.Now().UTC()

	return filepath.Join(getConfigDir(),
		fmt.Sprintf("%s_wallet_%x.json", now.Format("2006_01_02"), now.UnixNano()))
}
