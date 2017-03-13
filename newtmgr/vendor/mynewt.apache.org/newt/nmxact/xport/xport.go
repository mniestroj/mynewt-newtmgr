package xport

import (
	"mynewt.apache.org/newt/nmxact/sesn"
)

type RxFn func(data []byte)

type Xport interface {
	Start() error
	Stop() error
	BuildSesn(cfg sesn.SesnCfg) (sesn.Sesn, error)

	Tx(data []byte) error
}
