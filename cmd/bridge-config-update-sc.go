package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/0chain/gosdk/zcnbridge"
	"github.com/0chain/zwalletcli/util"
	"github.com/ethereum/go-ethereum/core/types"
)

//goland:noinspection ALL
func init() {
	rootCmd.AddCommand(
		createCommandWithBridge(
			"bridge-config-set-uint256",
			"Set a uint256 key value in NFT Config contract",
			"Set a uint256 key value in NFT Config contract",
			setUint256InNFTConfig,
			&Option{
				name:     "key",
				typename: "string",
				value:    "",
				usage:    "key",
				required: true,
			},
			&Option{
				name:     "value",
				typename: "int64",
				value:    int64(0),
				usage:    "value",
				required: true,
			},
			&Option{
				name:     "plan",
				typename: "int64",
				value:    int64(0),
				usage:    " plan id, used when setting royalty fees",
				required: false,
			},
		),
		createCommandWithBridge(
			"bridge-config-get-uint256",
			"Get a uint256 value in NFT Config contract",
			"Get a uint256 value in NFT Config contract",
			getUint256InNFTConfig,
			&Option{
				name:     "key",
				typename: "string",
				value:    "",
				usage:    "key",
				required: true,
			},
			&Option{
				name:     "plan",
				typename: "int64",
				value:    int64(0),
				usage:    "royalty plan id, used when setting royalty fees",
				required: false,
			},
		),
		createCommandWithBridge(
			"bridge-config-set-address",
			"Set an address key value in NFT Config contract",
			"Set an address key value in NFT Config contract",
			setAddressInNFTConfig,
			&Option{
				name:     "key",
				typename: "string",
				value:    "",
				usage:    "key",
				required: true,
			},
			&Option{
				name:     "address",
				typename: "string",
				value:    "",
				usage:    "address",
				required: true,
			},
		),
		createCommandWithBridge(
			"bridge-config-get-address",
			"Get an address value in NFT Config contract",
			"Get an address value in NFT Config contract",
			getAddressInNFTConfig,
			&Option{
				name:     "key",
				typename: "string",
				value:    "",
				usage:    "key",
				required: true,
			},
		))
}

// setUint256InNFTConfig sets a uint256 key value in NFT Config contract
func setUint256InNFTConfig(bc *zcnbridge.BridgeClient, args ...*Arg) {
	key := GetNFTConfigKey(args)
	value := GetNFTConfigValue(args)
	plan := GetNFTConfigRoyaltyPlanID(args)

	var (
		tx  *types.Transaction
		err error
	)

	if plan > 0 {
		k := zcnbridge.EncodePackInt64(key, plan)
		tx, err = bc.NFTConfigSetUint256Raw(context.Background(), k, value)
		if err != nil {
			ExitWithError(err)
		}
	} else {
		tx, err = bc.NFTConfigSetUint256(context.Background(), key, value)
		if err != nil {
			ExitWithError(err)
		}
	}

	hash := tx.Hash().String()
	fmt.Printf("Confirming Ethereum transaction: %s\n", hash)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 100, time.Second*5)
	if err != nil {
		ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("\nTransaction verification success: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
	}
}

func getUint256InNFTConfig(bc *zcnbridge.BridgeClient, args ...*Arg) {
	key := GetNFTConfigKey(args)
	plan := GetNFTConfigRoyaltyPlanID(args)

	var (
		k   string
		v   int64
		err error
	)
	if plan > 0 {
		k, v, err = bc.NFTConfigGetUint256(context.Background(), key, plan)
		if err != nil {
			ExitWithError(err)
		}
	} else {
		k, v, err = bc.NFTConfigGetUint256(context.Background(), key)
		if err != nil {
			ExitWithError(err)
		}
	}

	var response = struct {
		Key   string `json:"key"`
		Value int64  `json:"value"`
	}{
		Key:   k,
		Value: v,
	}

	util.PrettyPrintJSON(response)
}

func setAddressInNFTConfig(bc *zcnbridge.BridgeClient, args ...*Arg) {
	key := GetNFTConfigKey(args)
	address := GetNFTConfigAddress(args)
	tx, err := bc.NFTConfigSetAddress(context.Background(), key, address)
	if err != nil {
		ExitWithError(err)
	}

	hash := tx.Hash().String()
	fmt.Printf("Confirming Ethereum transaction: %s\n", hash)

	status, err := zcnbridge.ConfirmEthereumTransaction(hash, 100, time.Second*5)
	if err != nil {
		ExitWithError(err)
	}

	if status == 1 {
		fmt.Printf("\nTransaction verification success: %s\n", hash)
	} else {
		ExitWithError(fmt.Sprintf("\nVerification failed: %s\n", hash))
	}
}

func getAddressInNFTConfig(bc *zcnbridge.BridgeClient, args ...*Arg) {
	key := GetNFTConfigKey(args)
	k, address, err := bc.NFTConfigGetAddress(context.Background(), key)
	if err != nil {
		ExitWithError(err)
	}

	var response = struct {
		Key     string `json:"key"`
		Address string `json:"address"`
	}{
		Key:     k,
		Address: address,
	}

	util.PrettyPrintJSON(response)
}
