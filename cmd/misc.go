package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/0chain/gosdk/zcncore"
	"github.com/0chain/zwalletcmd/util"
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

var getuserpoolscmd = &cobra.Command{
	Use:   "getuserpools",
	Short: "Get user pools from sharders",
	Long:  `Get user pools from sharders`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetUserPools(statusBar)
		if err == nil {
			wg.Wait()
		} else {
			ExitWithError(err.Error())
		}
		if statusBar.success {
			fmt.Printf("\nUser pools: %v\n", statusBar.errMsg)
		} else {
			ExitWithError("\nERROR: Get user pool failed. " + statusBar.errMsg + "\n")
		}
		return
	},
}

var getuserpooldetailscmd = &cobra.Command{
	Use:   "getuserpooldetails",
	Short: "Get user pool details",
	Long: `Get user pool details for client_id and pool_id.
			<client_id> <pool_id>`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fflags := cmd.Flags()
		if fflags.Changed("client_id") == false {
			ExitWithError("Error: client_id flag is missing")
		}
		if fflags.Changed("pool_id") == false {
			ExitWithError("Error: pool_id flag is missing")
		}
		clientID := cmd.Flag("client_id").Value.String()
		poolID := cmd.Flag("pool_id").Value.String()
		wg := &sync.WaitGroup{}
		statusBar := &ZCNStatus{wg: wg}
		wg.Add(1)
		err := zcncore.GetUserPoolDetails(clientID, poolID, statusBar)
		if err != nil {
			ExitWithError(err)
		}
		wg.Wait()
		if statusBar.success {
			fmt.Printf("\nUser pool details: %v\n", statusBar.errMsg)
		} else {
			ExitWithError("\nERROR: Get user pool details failed. " + statusBar.errMsg + "\n")
		}

	},
}

type BlobberInfo struct {
	Id  string `json:"id"`
	Url string `json:"url"`
}

type BlobberNodes struct {
	Nodes []BlobberInfo `json:"Nodes"`
}

func printBlobberList(nodes BlobberNodes) {
	fmt.Println("Blobbers:")
	header := []string{"URL", "ID"}
	data := make([][]string, len(nodes.Nodes))
	for idx, child := range nodes.Nodes {
		data[idx] = []string{child.Url, child.Id}
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
	rootCmd.AddCommand(getuserpoolscmd)
	rootCmd.AddCommand(getuserpooldetailscmd)
	rootCmd.AddCommand(getblobberscmd)
	getidcmd.PersistentFlags().String("url", "", "URL to get the ID")
	getidcmd.MarkFlagRequired("url")
	getuserpooldetailscmd.PersistentFlags().String("client_id", "", "Miner or Sharder client id")
	getuserpooldetailscmd.PersistentFlags().String("pool_id", "", "Pool ID from user pool matching miner or sharder id")
	getuserpooldetailscmd.MarkFlagRequired("client_id")
	getuserpooldetailscmd.MarkFlagRequired("pool_id")
}
