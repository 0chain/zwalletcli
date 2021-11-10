package cmd

import (
	"encoding/json"
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

		var events Events
		cb := NewJSONInfoCB(&events)

		err = zcncore.GetEvents(cb, input.Fields)
		if err != nil {
			log.Fatal(err)
		}

		if err = cb.Waiting(); err != nil {
			log.Fatal(err)
		}

		pretty, err := json.MarshalIndent(events, "", "    ")
		if err != nil {
			log.Fatal("cannot marshal indent result: " + err.Error())
		}
		fmt.Println("events:", string(pretty))

	},
}

func init() {
	rootCmd.AddCommand(events)
	events.PersistentFlags().StringSlice("filters", nil, "list of filters")
	events.PersistentFlags().StringSlice("values", nil, "filter values")
}
