package cmd

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/0chain/gosdk/zcncore"
)

func createReadPool() (err error) {
	var (
		txn       zcncore.TransactionScheme
		wg        sync.WaitGroup
		statusBar = &ZCNStatus{wg: &wg}
	)

	if txn, err = zcncore.NewTransaction(statusBar, MinTxFee, nonce); err != nil {
		return
	}

	if err = txn.AdjustTransactionFee(txVelocity.toZCNFeeType()); err != nil {
		log.Printf("failed to adjust transaction fee: %v", err)
		return
	}

	wg.Add(1)
	if err = txn.CreateReadPool(0); err != nil {
		return
	}
	wg.Wait()

	if statusBar.success {
		statusBar.success = false

		wg.Add(1)
		if err = txn.Verify(); err != nil {
			return
		}
		wg.Wait()

		if statusBar.success {
			switch txn.GetVerifyConfirmationStatus() {
			case zcncore.ChargeableError:
				return errors.New(strings.Trim(txn.GetVerifyOutput(), "\""))
			case zcncore.Success:
				return
			default:
				return errors.New("\nExecute global settings update smart contract failed. Unknown status code: " +
					strconv.Itoa(int(txn.GetVerifyConfirmationStatus())))
			}
		}
	}

	return errors.New(statusBar.errMsg)
}
