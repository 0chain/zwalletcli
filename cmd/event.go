package cmd

import (
	"fmt"
	"log"

	"github.com/0chain/gosdk/zcncore"
	"github.com/spf13/cobra"
)

type Event struct {
	BlockNumber int64  `json:"block_number"`
	TxHash      string `json:"tx_hash"`
	Type        string `json:"type"`
	Tag         string `json:"tag"`
	Data        string `json:"data"`
}

type Events struct {
	//	Events []map[string]string `json:"events"`
	Events []Event `json:"events"`
}

var events = &cobra.Command{
	Use:   "events",
	Short: "List 0chain events.",
	Long:  `List 0chain events that match input filter settings.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		input := new(zcncore.InputMap)
		input.Fields = setupInputMap(cmd.Flags(), "filters", "values")
		if err != nil {
			log.Fatal(err)
		}

		//fields := new(zcncore.InputMap)

		var events Events
		cb := NewJSONInfoCB(&events)
		//var (
		//	wg        sync.WaitGroup
		//	statusBar = &ZCNStatus{wg: &wg}
		//)
		//wg.Add(1)
		//err = zcncore.GetEvents(statusBar, input.Fields)
		err = zcncore.GetEvents(cb, input.Fields)
		if err != nil {
			log.Fatal(err)
		}
		//wg.Wait()
		//if !statusBar.success {
		//	log.Fatal("fatal:", statusBar.errMsg)
		//}

		//fmt.Println(statusBar.errMsg)

		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}
		fmt.Println("events", events)
		//printMap(fields.Fields)
	},
}

func init() {
	rootCmd.AddCommand(events)
	events.PersistentFlags().StringSlice("filters", nil, "list of filters")
	events.PersistentFlags().StringSlice("values", nil, "filter values")
}
