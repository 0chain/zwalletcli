package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"time"
)

var createWalletCmd = &cobra.Command{
	Use:   "create-wallet",
	Short: "Create wallet and logs it into stdout",
	Long:  `Create wallet and logs it into standard output`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		walletStr, err := zcncore.CreateWalletOffline()
		if err != nil {
			ExitWithError("failed to generate offline wallet", err)
		}
		walletName := cmd.Flags().Lookup("wallet").Value.String()
		if len(walletName) == 0 {
			walletName = fmt.Sprintf("%d_wallet.json", time.Now().Unix())
		}

		// write wallet into wallet dir
		filename := walletFilename(walletName)
		//if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
		//	// same wallet exists
		//	ExitWithError(fmt.Sprintf("unable to write wallet, file with %q name already exists", filename))
		//}

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
	createWalletCmd.PersistentFlags().Bool("silent", false, "do not print wallet details in the standard output (default false)")
	createWalletCmd.PersistentFlags().String("wallet", "", "give custom name to the wallet")
}

func walletFilename(walletName string) string {
	cfgDir := getConfigDir()
	if len(walletFile) > 0 {
		return filepath.Join(cfgDir, walletFile)
	}

	return filepath.Join(getConfigDir(),
		fmt.Sprintf(walletName))
}
