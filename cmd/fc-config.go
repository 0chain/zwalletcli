package cmd

//TODO: JAYASH UNCOMMENT
//import (
//	"log"
//
//	"github.com/0chain/gosdk/zcncore"
//	"github.com/spf13/cobra"
//)
//
//var getFaucetConfigCmd = &cobra.Command{
//	Use:   "fc-config",
//	Short: "Show Faucet configurations.",
//	Long:  `Show Faucet configurations.`,
//	Args:  cobra.MinimumNArgs(0),
//	Hidden: true,
//	Run: func(cmd *cobra.Command, args []string) {
//
//		var (
//			fields = new(zcncore.InputMap)
//			cb     = NewJSONInfoCB(fields)
//			err    error
//		)
//		if err = zcncore.GetFaucetSCConfig(cb); err != nil {
//			log.Fatal(err)
//		}
//		if err = cb.Waiting(); err != nil {
//			log.Fatal(err)
//		}
//
//		printMap(fields.Fields)
//	},
//}
//
//func init() {
//	rootCmd.AddCommand(getFaucetConfigCmd)
//}
