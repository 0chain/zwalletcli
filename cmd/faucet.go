package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var faucetcmd = &cobra.Command{
	Use:   "faucet",
	Short: "Faucet smart contract",
	Long: `Faucet smart contract.
	        <methodName> <input>`,
	Args:   cobra.MinimumNArgs(0),
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("methodName") == false {
			ExitWithError("Error: Methodname flag is missing")
		}
		if fflags.Changed("input") == false {
			ExitWithError("Error: Input flag is missing")
		}

		input := cmd.Flag("input").Value.String()

		token, err := cmd.Flags().GetFloat64("tokens")
		if err != nil {
			ExitWithError(err)
			return
		}

		hash, _, _, _, err := zcncore.Faucet(uint64(token*1e10), input)
		if err != nil {
			ExitWithError(err)
			return
		}

		fmt.Println("Execute faucet smart contract success with txn : ", hash)
	},
}

func init() {
	rootCmd.AddCommand(faucetcmd)
	faucetcmd.PersistentFlags().String("methodName", "", "methodName")
	faucetcmd.PersistentFlags().String("input", "", "input")
	faucetcmd.PersistentFlags().Float64("tokens", 0, "Token request")
	faucetcmd.MarkFlagRequired("input")
}
