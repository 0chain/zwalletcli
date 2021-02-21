package cmd

import (
	"errors"
	"sync"

	"github.com/0chain/gosdk/zcncore"
)

import "fmt"

func createReadPool() (err error) {
  fmt.Println("d1 a")
  var (
    txn       zcncore.TransactionScheme
    wg        sync.WaitGroup
    statusBar = &ZCNStatus{wg: &wg}
  )

  fmt.Println("d2.1", statusBar.errMsg)

  if txn, err = zcncore.NewTransaction(statusBar, 0); err != nil {
    fmt.Println("d1 b")
    return
  }

  fmt.Println("d2.2", statusBar.errMsg)

  wg.Add(1)
  fmt.Println("d2.2.1", statusBar.errMsg)
  if err = txn.CreateReadPool(0); err != nil {
    fmt.Println("d1 c")
    return
  }
  fmt.Println("d2.2.2", statusBar.errMsg)
  wg.Wait()

  fmt.Println("d2.3", statusBar.errMsg)
  if statusBar.success {
    statusBar.success = false

    wg.Add(1)
    if err = txn.Verify(); err != nil {
      fmt.Println("d1 d")
      return
    }
    wg.Wait()

    if statusBar.success {
      fmt.Println("d1 e")
      return // nil
    }
  }

  fmt.Println("d1 f", statusBar.errMsg)
  return errors.New(statusBar.errMsg)
}
