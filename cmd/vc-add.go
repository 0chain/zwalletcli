package cmd

import (
	"fmt"
	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

var providerRegister = &cobra.Command{
	Use:   "vc-add",
	Short: "add node to view change",
	Long:  "add node to view change, add a miner or sharder to the register list so that they can join the MB. Only chainowner can do this.",
	// Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			flags = cmd.Flags()
			id    string
			err   error
		)

		if !flags.Changed("id") {
			ExitWithError("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			ExitWithError(err)
		}

		if !flags.Changed("provider-type") {
			ExitWithError("missing provider-type flag")
		}

		nodeType, err := flags.GetString("provider-type")
		if err != nil {
			ExitWithError(err)
		}

		var pt sdk.ProviderType
		switch nodeType {
		case "miner":
			pt = sdk.ProviderMiner
		case "sharder":
			pt = sdk.ProviderSharder
		default:
			ExitWithErrorf("unknown provider type: %v", nodeType)
		}

		hash, _, _, _, err := zcncore.VcRegisterNode(id, pt)
		if err != nil {
			ExitWithError("Vc register node : ", err)
		}

		fmt.Println("vc add success with transaction hash : ", hash)
	},
}

func init() {
	rootCmd.AddCommand(providerRegister)
	providerRegister.PersistentFlags().String("id", "", "provider ID to add to view change")
	_ = providerRegister.MarkFlagRequired("id")

	providerRegister.PersistentFlags().String("provider-type", "", "provider type")
	_ = providerRegister.MarkFlagRequired("provider-type")

}
