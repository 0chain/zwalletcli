package cmd

import (
	"fmt"
	"log"
	"os"

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

			// Lock read pool ? TODO: clarification from @dabasov
		} else {
			walletStr, err = zcncore.CreateWalletOffline()
			if err != nil {
				ExitWithError("failed to generate offline wallet", err)
			}
		}

		fmt.Fprintf(os.Stdout, walletStr)
	},
}

func init() {
	rootCmd.AddCommand(WithoutWallet(createWalletCmd))
	createWalletCmd.PersistentFlags().Bool("register", false, "create wallet with registration on blockchain (default false)")
}
