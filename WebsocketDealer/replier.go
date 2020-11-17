package WebsocketDealer

import (
	"errors"
	"github.com/fasthttp/websocket"
	"sync"
	"sync/atomic"
	"time"
)

var chanBoolPool = sync.Pool{New: func() interface{} {
	return make(chan StandardReply, 1)
}}

func (z *Replier) Call(functionName string, args ...interface{}) (interface{}, error) {
	if z.id == 0 {
		panic("Unexpected call to reply.It is not a real replier.Just a caller")
	}
	call := requestPool.Get().(*StandardCall)
	defer requestPool.Put(call)
	call.Function = functionName
	call.Argument = args
	call.Id = atomic.AddUint64(z.idCounter, 1)
	call.IsReply = true
	byte, err := json.Marshal(call)
	if err != nil {
		return nil, err
	}
	c := chanBoolPool.Get().(chan StandardReply)
	z.fa.BindReplyId(call.Id, c)
	z.conn.WriteMessage(websocket.TextMessage, byte)
	// 开始检测超时
	select {
	case <-time.After(time.Second * 3):
		// 开始移除绑定
		z.fa.RemoveBindReplyId(call.Id)
		return nil, errors.New("call time out")
	case reply := <-c:
		return reply.Data, nil
	}
}
func (z *Replier) Return(data interface{}) error {
	if data == nil {
		return nil
	}
	res := responsePool.Get().(*WebsocketResponse)
	defer responsePool.Put(res)
	res.data = data
	res.id = z.id
	byte, err := json.Marshal(res)
	if err != nil {
		return errors.New(err.Error())
	}
	err = z.conn.WriteMessage(websocket.TextMessage, byte)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}
