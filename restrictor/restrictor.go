package restrictor

// 节流器会自动的清理内存
// 当轮创建时间大于1000秒时，就会被自动释放
import (
	"sync"
	"time"

	"github.com/Ruixues/RWeb/WebsocketDealer"
)

type visit struct {
	turnNum   int64
	turnStart int64
	lock      sync.Mutex
}

//Restrictor 节流器
type Restrictor struct {
	ipRecord    map[interface{}]*visit
	accessCheck func(startTime *int64, nowNum *int64) bool
	lastGCTime  int64
}

func defaultChecker(startTime, nowNum *int64) bool {
	if time.Now().Unix()-*startTime > 5 {
		*startTime = time.Now().Unix()
		*nowNum = 0
		return true
	}
	if *nowNum < 10 {
		return true
	}
	return false
}

//New 创建节流器
func New(accessCheck func(startTime *int64, nowNum *int64) bool) *Restrictor {
	if accessCheck == nil {
		accessCheck = defaultChecker
	}
	ret := new(Restrictor)
	ret.ipRecord = make(map[interface{}]*visit)
	ret.accessCheck = accessCheck
	return ret
}
func (z *Restrictor) avaible(id interface{}) bool {
	data, ok := z.ipRecord[id]
	if !ok {
		return false
	}
	go z.GC()
	return z.accessCheck(&data.turnStart, &data.turnNum)
}

//Record 记录访问
func (z *Restrictor) Record(id interface{}) bool {
	if _, ok := z.ipRecord[id]; !ok {
		// 不存在。那就开始初始化计时器
		z.ipRecord[id] = &visit{
			turnNum:   1,
			turnStart: time.Now().Unix(),
		}
		return z.avaible(id)
	}
	// 那就是已经访问过了。直接叠加
	z.ipRecord[id].lock.Lock()
	defer z.ipRecord[id].lock.Unlock()
	z.ipRecord[id].turnNum++
	return z.avaible(id)
}

//GC 清理无用的数据
func (z *Restrictor) GC() {
	nowTime := time.Now().Unix()
	if nowTime-z.lastGCTime < 100 {
		return
	}
	z.lastGCTime = nowTime
	for k, v := range z.ipRecord {
		if nowTime-v.turnStart > 1000 {
			lock := &(v.lock)
			lock.Lock()
			delete(z.ipRecord, k)
			lock.Unlock()
		}
	}
}

//WebsocketCount For rrpc
func (z *Restrictor) WebsocketCount(functionName string, replier *WebsocketDealer.Replier, session *WebsocketDealer.Session, args []interface{}) bool {
	return z.Record(session.Id)
}
