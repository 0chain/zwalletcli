package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk_common/zboxcore/commonsdk"
	"github.com/0chain/gosdk_common/zcncore"
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
			log.Fatal("missing id flag")
		}

		if id, err = flags.GetString("id"); err != nil {
			log.Fatal(err)
		}

		if !flags.Changed("provider-type") {
			log.Fatal("missing provider-type flag")
		}

		nodeType, err := flags.GetString("provider-type")
		if err != nil {
			log.Fatal(err)
		}

		var pt commonsdk.ProviderType
		switch nodeType {
		case "miner":
			pt = commonsdk.ProviderMiner
		case "sharder":
			pt = commonsdk.ProviderSharder
		default:
			log.Fatalf("unknown provider type: %v", nodeType)
		}

		hash, _, _, _, err := zcncore.VcRegisterNode(id, pt)
		if err != nil {
			log.Fatal("Vc register node : ", err)
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
