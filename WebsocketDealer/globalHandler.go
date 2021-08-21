package WebsocketDealer

import (
	"sync"
	"time"
)

type GlobalHandler struct {
	access        *sync.Map
	maxHandleTime time.Duration
	downFunction  func(key interface{}, userData interface{}, lastBeatTime time.Time) bool // if the function returns false,the user won't be kicked off.
}
type userPackage struct {
	data         interface{}
	lastBeatTime time.Time
}

var userPackagePool = &sync.Pool{New: func() interface{} { return new(userPackage) }}

func defaultDownFunction(key interface{}, userData interface{}, lastBeatTime time.Time) bool {
	return true
}
func NewGlobalHandler() GlobalHandler {
	return GlobalHandler{
		access:        &sync.Map{},
		maxHandleTime: time.Second * 10,
		downFunction:  defaultDownFunction,
	}
}
func (z *GlobalHandler) MountUser(key interface{}, val interface{}) bool {
	rPackage := userPackagePool.Get().(*userPackage)
	rPackage.data = val
	rPackage.lastBeatTime = time.Now()
	_, ok := z.access.LoadOrStore(key, rPackage)
	return !ok
}
func (z *GlobalHandler) UserBeat(key interface{}) bool {
	user, ok := z.access.Load(key)
	if !ok {
		return false
	}
	user.(*userPackage).lastBeatTime = time.Now()
	return true
}

// SetMaxHandleTime when it takes too long for a user to send a heart beat message,the user will be kicked off from the server.
// The MaxHandleTime is the limit.
func (z *GlobalHandler) SetMaxHandleTime(time time.Duration) {
	z.maxHandleTime = time
}
func (z *GlobalHandler) SetDownFunction(f func(key interface{}, userData interface{}, lastBeatTime time.Time) bool) {
	if f == nil {
		f = defaultDownFunction
	}
	z.downFunction = f
}
func (z *GlobalHandler) GC() {
	t := time.Now()
	z.access.Range(func(key, value interface{}) bool {
		p := value.(*userPackage)
		if t.Sub(p.lastBeatTime) > z.maxHandleTime {
			if z.downFunction(key, p.data, p.lastBeatTime) {
				// kick off the user
				z.access.Delete(key)
			}
		}
		return true
	})
}
