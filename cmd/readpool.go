package cmd

import (
	"errors"
	"sync"

	"github.com/0chain/gosdk/zcncore"
)

func createReadPool() (err error) {
	var (
		txn       zcncore.TransactionScheme
		wg        sync.WaitGroup
		statusBar = &ZCNStatus{wg: &wg}
	)

	if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
		return
	}

	wg.Add(1)
	if err = txn.CreateReadPool(); err != nil {
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
			return // nil
		}
	}

	return errors.New(statusBar.errMsg)
}
