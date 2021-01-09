package WebsocketDealer

import (
	"sync"
	"time"
)

type Session struct {
	Data       interface{}
	CreateTime int64
	lock       *sync.RWMutex
}

var sessionPool = sync.Pool{
	New: func() interface{} {
		ret := new(Session)
		ret.lock = &sync.RWMutex{}
		return ret
	}}

func NewSession() *Session {
	s := sessionPool.Get().(*Session)
	s.CreateTime = time.Now().Unix()
	s.Data = nil
	return s
}
