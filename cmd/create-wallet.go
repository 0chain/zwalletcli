package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var createWalletCmd = &cobra.Command{
	Use:   "create-wallet",
	Short: "Create wallet and logs it into stdout (pass --register to register wallet to blockchain)",
	Long:  `Create wallet and logs it into standard output (pass --register to register wallet to blockchain)`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		walletStr, err := zcncore.CreateWalletOffline()
		if err != nil {
			ExitWithError("failed to generate offline wallet", err)
		}
		walletName := cmd.Flags().Lookup("wallet").Value.String()

		// write wallet into wallet dir
		filename := walletFilename(walletName)
		fmt.Print(filename)
		if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
			// same wallet exists
			ExitWithError(fmt.Sprintf("unable to write wallet, file with %q name already exists. Please try a different wallet name or backup the current wallet file and delete it.", filename))
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
	createWalletCmd.PersistentFlags().Bool("silent", false, "do not print wallet details in the standard output (default false)")
	createWalletCmd.PersistentFlags().String("wallet", "", "give custom name to the wallet")
}

func walletFilename(walletName string) string {
	cfgDir := getConfigDir()
	if len(walletName) > 0 {
		return filepath.Join(cfgDir, walletName)
	}
	return filepath.Join(cfgDir, "wallet.json")
}
