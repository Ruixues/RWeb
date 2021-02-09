package WebsocketDealer

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Data       interface{}
	CreateTime int64
	lock       *sync.RWMutex
	Id         uuid.UUID
}

var sessionPool = sync.Pool{
	New: func() interface{} {
		ret := new(Session)
		ret.lock = &sync.RWMutex{}
		return ret
	}}

//NewSession Build a session
func NewSession() *Session {
	s := sessionPool.Get().(*Session)
	s.CreateTime = time.Now().Unix()
	s.Data = nil
	s.Id = uuid.New()
	return s
}
