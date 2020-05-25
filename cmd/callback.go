package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/0chain/gosdk/zcncore"
)

//
// on JSON info available
//

type OnJSONInfoCb struct {
	value interface{}
	err   error
	got   chan struct{}
}

func (ojsonic *OnJSONInfoCb) OnInfoAvailable(op int, status int,
	info string, errMsg string) {

	defer close(ojsonic.got)

	if status != zcncore.StatusSuccess {
		ojsonic.err = errors.New(errMsg)
		return
	}
	if info == "" || info == "{}" {
		ojsonic.err = errors.New("empty response from sharders")
		return
	}
	var err error
	if err = json.Unmarshal([]byte(info), ojsonic.value); err != nil {
		ojsonic.err = fmt.Errorf("decoding response: %v", err)
	}
}

// Wait for info.
func (ojsonic *OnJSONInfoCb) Waiting() (err error) {
	<-ojsonic.got
	return ojsonic.err
}

func NewJSONInfoCB(val interface{}) (cb *OnJSONInfoCb) {
	cb = new(OnJSONInfoCb)
	cb.value = val
	cb.got = make(chan struct{})
	return
}
