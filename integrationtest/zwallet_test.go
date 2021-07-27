package integrationtest

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	remotePath   = "/"
	fileName     = "1.txt"
	allocationID = "f212f81ec0208c3cbc21ca0524c13e27ff62f37763a10d7a6cc8c60eb1302f1b"
	clientID     = "0ae17e887ea887f7293d59741db68dabcbef28996a0fd7c2c7d49f020a7ac4e0"
)

var (
	_, b, _, _ = runtime.Caller(0)
	dirPath    = strings.TrimSuffix(filepath.Dir(b), "/integrationtest")
)

func Test_Register(t *testing.T) {
	cmd := exec.Command("./zwallet", "register")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Equal(t, "Wallet registered\n", string(out))
}

func Test_Faucet(t *testing.T) {
	cmd := exec.Command("./zwallet", "faucet", "--methodName", "pour", "--input", "'{Pay day}'")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Execute faucet smart contract success with txn")
}

func Test_GetBalance(t *testing.T) {
	cmd := exec.Command("./zwallet", "getbalance")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Balance:")
}

func Test_RecoverWallet(t *testing.T) {
	cmd := exec.Command("./zwallet", "recoverwallet", "--wallet", "recovered_wallet.json", "--mnemonic", "round rather extra common student valve connect review aerobic struggle sniff jacket peace nominee pill liquid coach slow tree ensure hand regret violin arrive")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Wallet recovered!!")
}

func Test_SendToken(t *testing.T) {
	cmd := exec.Command("./zwallet", "send", "--to_client_id", clientID, "--token", ".2", "--desc", "'gift'", "--fee", "0.1")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Send tokens success")
}

func Test_LockConfig(t *testing.T) {
	cmd := exec.Command("./zwallet", "lockconfig")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_Lock(t *testing.T) {
	cmd := exec.Command("./zwallet", "lock", "--durationHr", "0", "--durationMin", "5", "--tokens", "0.2", "--fee", "0.1")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Tokens (0.200000) locked successfully")
}

func Test_GetLockTokens(t *testing.T) {
	cmd := exec.Command("./zwallet", "getlockedtokens")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_UnlockingTokens(t *testing.T) {
	type Message struct {
		Stats []struct {
			PoolID       string  `json:"pool_id"`
			StartTime    int     `json:"start_time"`
			Duration     int64   `json:"duration"`
			TimeLeft     int64   `json:"time_left"`
			Locked       bool    `json:"locked"`
			Apr          float64 `json:"apr"`
			TokensEarned int     `json:"tokens_earned"`
			Balance      int     `json:"balance"`
		} `json:"stats"`
	}
	cmd := exec.Command("./zwallet", "getlockedtokens")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var msg string
	fmt.Sscanf(string(out), "\nLocked tokens:\n %s\n", &msg)
	fmt.Println(msg)
	m := Message{}
	err = json.Unmarshal([]byte(msg), &m)
	if err != nil {
		log.Fatal(err)
	}
	poolID := (m.Stats[0].PoolID)

	cmd = exec.Command("./zwallet", "unlock", "--pool_id", poolID)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Unlock tokens success")
}

func Test_GetBlobbers(t *testing.T) {
	cmd := exec.Command("./zwallet", "getblobbers")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_GetID(t *testing.T) {
	cmd := exec.Command("./zwallet", "getid", "--url", "http://198.18.0.81:7171")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_ListSharders(t *testing.T) {
	cmd := exec.Command("./zwallet", "ls-sharders")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_ListMiners(t *testing.T) {
	cmd := exec.Command("./zwallet", "ls-miners")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	fmt.Println("-------------------------")
	fmt.Println(string(out))
	fmt.Println("-------------------------")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_GetMinerSCInfo(t *testing.T) {
	cmd := exec.Command("./zwallet", "ls-miners")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var minerID string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- ID:        %s", &minerID)

	cmd = exec.Command("./zwallet", "mn-info", "--id", minerID)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_StakeLock(t *testing.T) {
	cmd := exec.Command("./zwallet", "ls-miners")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var minerID string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- ID:        %s", &minerID)

	cmd = exec.Command("./zwallet", "mn-lock", "--id", minerID, "--tokens", "0.2")
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "locked with:")
}

func Test_MinerSCUserInfo(t *testing.T) {
	cmd := exec.Command("./zwallet", "mn-user-info")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_MinerSCPoolInfo(t *testing.T) {
	cmd := exec.Command("./zwallet", "mn-user-info")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var minerID string
	var poolID string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- node: %s", &minerID)
	fmt.Sscanf(s[1], "  - pool_id:       %s", &poolID)

	cmd = exec.Command("./zwallet", "mn-pool-info", "--id", minerID, "--pool_id", poolID)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_StakeUnlock(t *testing.T) {
	cmd := exec.Command("./zwallet", "mn-user-info")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var minerID string
	var poolID string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- node: %s", &minerID)
	fmt.Sscanf(s[1], "  - pool_id:       %s", &poolID)

	cmd = exec.Command("./zwallet", "mn-unlock", "--id", minerID, "--pool_id", poolID)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "tokens will be unlocked next VC:")
}

func Test_UpdateStakeConfig(t *testing.T) {
	cmd := exec.Command("./zwallet", "ls-miners")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var minerID string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- ID:        %s", &minerID)

	cmd = exec.Command("./zwallet", "mn-update-settings", "--id", minerID, "--max_stake", "100000")
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "settings updated")
}

func Test_AddVestingPool(t *testing.T) {
	cmd := exec.Command("./zwallet", "vp-add", "--duration", "5m", "--lock", "5", "--d", "9842dc9200d738504c71bee02570d09233675b381d67df501cba29c1d97e221c:1")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Vesting pool added successfully:")
}

func Test_ListVestingPool(t *testing.T) {
	cmd := exec.Command("./zwallet", "vp-list")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_VestingPoolInfo(t *testing.T) {
	cmd := exec.Command("./zwallet", "vp-list")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var poolId string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- %s", &poolId)

	cmd = exec.Command("./zwallet", "vp-info", "--pool_id", poolId)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
}

func Test_TriggerVestingPool(t *testing.T) {
	cmd := exec.Command("./zwallet", "vp-list")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var poolId string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- %s", &poolId)

	cmd = exec.Command("./zwallet", "vp-trigger", "--pool_id", poolId)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Vesting triggered successfully.")
}

func Test_UnlockVestingPool(t *testing.T) {
	cmd := exec.Command("./zwallet", "vp-list")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var poolId string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- %s", &poolId)

	cmd = exec.Command("./zwallet", "vp-unlock", "--pool_id", poolId)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Tokens unlocked successfully.")
}

func Test_StopVesting(t *testing.T) {
	cmd := exec.Command("./zwallet", "vp-list")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var poolId string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- %s", &poolId)

	cmd = exec.Command("./zwallet", "vp-stop", "--d", "9842dc9200d738504c71bee02570d09233675b381d67df501cba29c1d97e221c", "--pool_id", poolId)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Stop vesting for 9842dc9200d738504c71bee02570d09233675b381d67df501cba29c1d97e221c")
}

func Test_DeleteVestingPool(t *testing.T) {
	cmd := exec.Command("./zwallet", "vp-list")
	cmd.Dir = dirPath
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	var poolId string
	s := strings.Split(string(out), "\n")
	fmt.Sscanf(s[0], "- %s", &poolId)

	cmd = exec.Command("./zwallet", "vp-delete", "--pool_id", poolId)
	cmd.Dir = dirPath
	out, err = cmd.Output()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("%s", out)
	}
	require.Contains(t, string(out), "Vesting pool deleted successfully.")
}
