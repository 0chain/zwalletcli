package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/core/client"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

var getbalancecmd = &cobra.Command{
	Use:   "getbalance",
	Short: "Get balance from sharders",
	Long:  `Get balance from sharders`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		doJSON, _ := cmd.Flags().GetBool("json")

		bal, err := client.GetBalance()
		if err != nil {
			ExitWithError(err)
			return
		}
		token, err := bal.ToToken()
		if err != nil {
			ExitWithError(err)
			return
		}
		usd, err := zcncore.ConvertTokenToUSD(token)

		if doJSON {
			j := map[string]string{
				"usd": fmt.Sprintf("%f", usd),
				"zcn": fmt.Sprintf("%f", token),
				"fmt": fmt.Sprintf("%d", bal.Balance)}
			util.PrintJSON(j)
			return
		}
		if err != nil {
			ExitWithError(fmt.Sprintf("\nBalance: %v (Failed to get USD: %v)\n", bal.Balance, err))
			return
		}
		fmt.Printf("\nBalance: %v ZCN (%.2f USD)\n", token, usd)
	},
}

func init() {
	rootCmd.AddCommand(getbalancecmd)
	getbalancecmd.Flags().Bool("json", false, "pass this option to print response as json data")
}
