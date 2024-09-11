package cmd

import (
	"github.com/0chain/gosdk/zboxcore/sdk"
)

func createReadPool() (err error) {
	if _, _, err = sdk.CreateReadPool(); err != nil {
		return
	}
	return nil
}
