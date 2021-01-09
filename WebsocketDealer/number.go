package WebsocketDealer

import "sync"

type NumberPool struct {
	nowMax     uint64
	bufPool    []uint64 //未使用
	lock       *sync.Mutex
	initNumber uint64 //在允许的条件下，总数字量会向initNumber靠拢
}

// newNumberPool 生成新的Pool
// initNumber    初始化的数字数量
func NewNumberPool(initNumber uint64) (ret NumberPool) {
	ret.bufPool = make([]uint64, initNumber)
	ret.lock = &sync.Mutex{}
	ret.initNumber = initNumber
	for i := uint64(0); i < initNumber; i++ {
		ret.bufPool[i] = i
	}
	return
}
func (z *NumberPool) Get() uint64 {
	z.lock.Lock()
	defer z.lock.Unlock()
	if len(z.bufPool) > 0 {
		ret := z.bufPool[0]
		z.bufPool = z.bufPool[1:]
		return ret
	}
	// 看来要生成一个数字了
	z.nowMax++
	return z.nowMax
}
func (z *NumberPool) Put(num uint64) {
	z.lock.Lock()
	defer z.lock.Unlock()
	z.bufPool = append(z.bufPool, num)
}
