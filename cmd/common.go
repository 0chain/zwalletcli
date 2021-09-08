package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/0chain/gosdk/core/common"
	"github.com/0chain/gosdk/zcncore"
	"gopkg.in/cheggaaa/pb.v1"
)

type StatusBar struct {
	b  *pb.ProgressBar
	wg *sync.WaitGroup
}

type ZCNStatus struct {
	walletString string
	wg           *sync.WaitGroup
	success      bool
	errMsg       string
	balance      common.Balance
	wallets      []string
	clientID     string
}

func NewZCNStatus() (zcns *ZCNStatus) {
	return &ZCNStatus{wg: new(sync.WaitGroup)}
}

func (zcns *ZCNStatus) Begin() { zcns.wg.Add(1) }
func (zcns *ZCNStatus) Wait()  { zcns.wg.Wait() }

func (zcn *ZCNStatus) OnBalanceAvailable(status int, value int64, info string) {
	defer zcn.wg.Done()
	if status == zcncore.StatusSuccess {
		zcn.success = true
	} else {
		zcn.success = false
	}
	zcn.balance = common.Balance(value)
}

func (zcn *ZCNStatus) OnTransactionComplete(t *zcncore.Transaction, status int) {
	defer zcn.wg.Done()
	if status == zcncore.StatusSuccess {
		zcn.success = true
	} else {
		zcn.errMsg = t.GetTransactionError()
	}
	// fmt.Println("Txn Hash:", t.GetTransactionHash())
}

func (zcn *ZCNStatus) OnVerifyComplete(t *zcncore.Transaction, status int) {
	defer zcn.wg.Done()
	if status == zcncore.StatusSuccess {
		zcn.success = true
	} else {
		zcn.errMsg = t.GetVerifyError()
	}
	// fmt.Println(t.GetVerifyOutput())
}

func (zcn *ZCNStatus) OnAuthComplete(t *zcncore.Transaction, status int) {
	fmt.Println("Authorization complete on zauth.", status)
}

func (zcn *ZCNStatus) OnWalletCreateComplete(status int, wallet string, err string) {
	defer zcn.wg.Done()
	if status != zcncore.StatusSuccess {
		zcn.success = false
		zcn.errMsg = err
		zcn.walletString = ""
		return
	}
	zcn.success = true
	zcn.errMsg = ""
	zcn.walletString = wallet
	return
}

func (zcn *ZCNStatus) OnInfoAvailable(Op int, status int, config string, err string) {
	defer zcn.wg.Done()
	if status != zcncore.StatusSuccess {
		zcn.success = false
		zcn.errMsg = err
		return
	}
	zcn.success = true
	zcn.errMsg = config
}

func (zcn *ZCNStatus) OnSetupComplete(status int, err string) {
	defer zcn.wg.Done()
}

func (zcn *ZCNStatus) OnAuthorizeSendComplete(status int, toClienID string, val int64, desc string, creationDate int64, signature string) {
	defer zcn.wg.Done()
	fmt.Println("Status:", status)
	fmt.Println("Timestamp:", creationDate)
	fmt.Println("Signature:", signature)
}

//OnVoteComplete callback when a multisig vote is completed
func (zcn *ZCNStatus) OnVoteComplete(status int, proposal string, err string) {
	defer zcn.wg.Done()
	if status != zcncore.StatusSuccess {
		zcn.success = false
		zcn.errMsg = err
		zcn.walletString = ""
		return
	}
	zcn.success = true
	zcn.errMsg = ""
	zcn.walletString = proposal
}

func PrintError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func ExitWithError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}
