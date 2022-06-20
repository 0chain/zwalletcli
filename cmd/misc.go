package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

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
		if fflags.Changed("url") == false {
			ExitWithError("Error: url flag is missing")
		}
		url := cmd.Flag("url").Value.String()

		id := zcncore.GetIdForUrl(url)
		if id == "" {
			ExitWithError("Error: ID not found")
		}
		fmt.Printf("\nURL: %v \nID: %v\n", url, id)
		return
	},
}

type Terms struct {
	ReadPrice        int64         `json:"read_price"`
	WritePrice       int64         `json:"write_price"`
	MinLockDemand    float64       `json:"min_lock_demand"`
	MaxOfferDuration time.Duration `json:"max_offer_duration"`
}

type BlobberInfo struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	Terms     Terms  `json:"terms"`
	Capacity  int64  `json:"capacity"`
	Allocated int64  `json:"allocated"`
}

type BlobberNodes struct {
	Nodes []BlobberInfo `json:"Nodes"`
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

func printBlobberList(nodes BlobberNodes) {
	fmt.Println("Blobbers:")
	header := []string{
		"URL", "ID", "CAP", "R / W PRICE", "DEMAND",
	}
	data := make([][]string, len(nodes.Nodes))
	for idx, child := range nodes.Nodes {
		data[idx] = []string{
			child.Url,
			child.Id,
			fmt.Sprintf("%s / %s",
				byteCountIEC(child.Allocated), byteCountIEC(child.Capacity)),
			fmt.Sprintf("%f / %f",
				zcncore.ConvertToToken(child.Terms.ReadPrice),
				zcncore.ConvertToToken(child.Terms.WritePrice)),
			fmt.Sprint(child.Terms.MinLockDemand),
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
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetBlobbers(statusBar)
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}
		if statusBar.success {
			var blobberNodes BlobberNodes
			err = json.Unmarshal([]byte(statusBar.errMsg), &blobberNodes)
			if err == nil {
				printBlobberList(blobberNodes)
			} else {
				fmt.Println(err)
				fmt.Printf("Blobbers: %v", statusBar.errMsg)
			}
		} else {
			ExitWithError("\nERROR: Get blobbers failed. " + statusBar.errMsg + "\n")
		}
		return
	},
}

func readFile(fileName string) (string, error) {
	w, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(w), nil
}

func writeToaFile(fileNameAndPath string, content string) error {

	file, err := os.Create(fileNameAndPath)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer file.Close()
	fmt.Fprintf(file, content)
	return nil
}

func init() {
	rootCmd.AddCommand(getidcmd)
	rootCmd.AddCommand(getblobberscmd)
	getidcmd.PersistentFlags().String("url", "", "URL to get the ID")
	getidcmd.MarkFlagRequired("url")
}
