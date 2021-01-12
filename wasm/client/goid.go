// +build !wasm

package clientGo

import "github.com/modern-go/gls"

func goid() int64 {
	return gls.GoID()
}
