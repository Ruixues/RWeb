package RWeb

import (
	"errors"
	"sync"
)

type bindData struct {
	Method  int
	Handler Handler
}
type DefaultRouter struct {
	bindLock *sync.RWMutex
	linker   map[string]*bindData
}

func NewDefaultRouter() (r DefaultRouter) {
	r.bindLock = &sync.RWMutex{}
	r.linker = make(map[string]*bindData)
	return
}
func (z *DefaultRouter) Bind(address string, method int, handler Handler) error {
	z.bindLock.Lock()
	defer z.bindLock.Unlock()
	if _, ok := z.linker[address]; ok {
		return errors.New("you are trying to bind a existed address")
	}
	z.linker[address] = &bindData{
		Method:  method,
		Handler: handler,
	}
	return nil
}
