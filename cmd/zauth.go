package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/0chain/gosdk/zcncore"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type splitWallet struct {
	ClientID      string `json:"client_id"`
	ClientKey     string `json:"client_key"`
	PublicKey     string `json:"public_key"`
	PrivateKey    string `json:"private_key"`
	PeerPublicKey string `json:"peer_public_key"`
}

func callZauthSetup(serverAddr string, splitWallet splitWallet) error {
	// Add your code here
	endpoint := serverAddr + "/setup"
	wData, err := json.Marshal(splitWallet)
	if err != nil {
		return errors.Wrap(err, "failed to marshal split wallet")
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(wData))
	if err != nil {
		return errors.Wrap(err, "failed to create HTTP request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send HTTP request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

var zauthCmd = &cobra.Command{
	Use:   "zauth",
	Short: "Enable zauth",
	Long:  `Enable zauth to sign transactions and messages, setup split keys and configure the zauth service.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Add your code here
		serverAddr, err := cmd.Flags().GetString("server")
		if err != nil {
			log.Fatalf("Could not find zauth server address")
		}

		// update or setup the zauth server address
		cfgConfig.Set("zauth.server", serverAddr)
		if err := cfgConfig.WriteConfig(); err != nil {
			log.Fatalf("Could not save config: %v", err)
		}

		// figure out what is the private key to use
		// from the web-apps, we can see the private key is the first key's private key
		if clientWallet == nil {
			log.Fatalf("Wallet is initialized yet")
		}

		sw, err := zcncore.SplitKeysWallet(clientWallet.Keys[0].PrivateKey, 2)
		if err != nil {
			log.Fatalf("Failed to split keys: %v", err)
		}

		// call gosdk to connect and setup the split keys
		// zcncore.SplitKeys()
		if err := callZauthSetup(serverAddr, splitWallet{
			ClientID:      sw.ClientID,
			ClientKey:     sw.ClientKey,
			PublicKey:     sw.Keys[1].PublicKey,
			PrivateKey:    sw.Keys[1].PrivateKey,
			PeerPublicKey: sw.Keys[0].PublicKey,
		}); err != nil {
			log.Fatalf("Failed to setup zauth server: %v", err)
		}

		// remove the keys[1]
		sw.Keys = sw.Keys[:1]
		clientWallet.SetSplitKeys(sw)
		if err := clientWallet.SaveTo(cfgWallet); err != nil {
			log.Fatalf("Failed to save wallet: %v", err)
		}

		log.Printf("Setup zauth server successfully")
	},
}

func init() {
	rootCmd.AddCommand(zauthCmd)
	zauthCmd.PersistentFlags().String("server", "s", "The zauth server address")
	if err := zauthCmd.MarkFlagRequired("server"); err != nil {
		log.Fatalf("Could not mark 'server' flag required: %v", err)
	}
}
