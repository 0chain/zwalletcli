package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/0chain/gosdk/zboxcore/sdk"
	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcli/util"
	"github.com/spf13/cobra"
)

var getidcmd = &cobra.Command{
	Use:   "getid",
	Short: "Get Miner or Sharder ID from its URL",
	Long:  `Get Miner or Sharder ID from its URL`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if !fflags.Changed("url") {
			ExitWithError("Error: url flag is missing")
		}
		url := cmd.Flag("url").Value.String()

		id := zcncore.GetIdForUrl(url)
		if id == "" {
			ExitWithError("Error: ID not found")
		}
		fmt.Printf("\nURL: %v \nID: %v\n", url, id)
	},
}

func byteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func printBlobberList(nodes []*sdk.Blobber) {
	fmt.Println("Blobbers:")
	header := []string{
		"URL", "ID", "CAP", "R / W PRICE", "DEMAND",
	}
	data := make([][]string, len(nodes))
	for idx, child := range nodes {
		data[idx] = []string{
			child.BaseURL,
			string(child.ID),
			fmt.Sprintf("%s / %s",
				byteCountIEC(int64(child.Allocated)), byteCountIEC(int64(child.Capacity))),
			fmt.Sprintf("%f / %f",
				zcncore.ConvertToToken(int64(child.Terms.ReadPrice)),
				zcncore.ConvertToToken(int64(child.Terms.WritePrice))),
		}
	}
	util.WriteTable(os.Stdout, header, []string{}, data)
	fmt.Println("")
}

var getblobberscmd = &cobra.Command{
	Use:   "getblobbers",
	Short: "Get registered blobbers from sharders",
	Long:  `Get registered blobbers from sharders`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		active, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatal(err)
		}
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}

		var blobberList []*sdk.Blobber
		limit, offset := 20, 0

		for {
			wg.Add(1)
			zcncore.GetBlobbers(statusBar, limit, offset, !active)
			wg.Wait()

			type nodes struct {
				Nodes []*sdk.Blobber
			}

			var wrap nodes

			err := json.Unmarshal([]byte(statusBar.errMsg), &wrap)
			if err != nil {
				log.Fatal("error unmarshalling blobbers")
			}
			if len(wrap.Nodes) == 0 {
				break
			}

			blobberList = append(blobberList, wrap.Nodes...)

			offset += limit
		}

		printBlobberList(blobberList)
	},
}

func init() {
	rootCmd.AddCommand(getidcmd)
	rootCmd.AddCommand(getblobberscmd)
	getidcmd.PersistentFlags().String("url", "", "URL to get the ID")
	getidcmd.MarkFlagRequired("url")
	getblobberscmd.PersistentFlags().Bool("all", false, "Gets all blobbers, including inactive blobbers")
}
