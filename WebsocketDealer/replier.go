package WebsocketDealer

import (
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
)

var chanBoolPool = sync.Pool{New: func() interface{} {
	return make(chan StandardReply, 1)
}}

//Call 调用客户端函数
func (z *Replier) Call(functionName string, args ...interface{}) (interface{}, error) {
	call := requestPool.Get().(*StandardCall)
	defer requestPool.Put(call)
	call.Function = functionName
	call.Argument = args
	id := atomic.AddUint64(z.idCounter, 1)
	call.Id = jsoniter.Number(strconv.FormatUint(id, 10))
	call.IsReply = false
	byte, err := jsoniter.Marshal(call)
	if err != nil {
		return nil, err
	}
	byte = z.encode(byte)
	c := chanBoolPool.Get().(chan StandardReply)
	z.bindReplyId(id, c)
	z.conn.WriteMessage(websocket.TextMessage, byte)
	// 开始检测超时
	select {
	case <-time.After(time.Second * 3):
		// 开始移除绑定
		z.removeBindReplyId(id)
		return nil, errors.New("timeout")
	case reply := <-c:
		return reply.Data, nil
	}
}

// CallNotWait 不等待回复
func (z *Replier) CallNotWait(functionName string, args ...interface{}) error {
	call := requestPool.Get().(*StandardCall)
	defer requestPool.Put(call)
	call.Function = functionName
	call.Argument = args
	id := atomic.AddUint64(z.idCounter, 1)
	call.Id = jsoniter.Number(strconv.FormatUint(id, 10))
	call.IsReply = false
	byte, err := jsoniter.Marshal(call)
	if err != nil {
		return err
	}
	byte = z.encode(byte)
	z.conn.WriteMessage(websocket.TextMessage, byte)
	return nil
}

//Return 返回值
func (z *Replier) Return(data interface{}) error {
	defer func() {
		recover()
	}()
	if z.id.String() == "0" {
		panic("Unexpected call to reply .It is not a real replier.Just a caller")
	}
	if data == nil {
		return nil
	}
	res := responsePool.Get().(*WebsocketResponse)
	defer responsePool.Put(res)
	res.Data = data
	res.Id = z.id
	res.IsReply = true
	byte, err := json.Marshal(*res)
	if err != nil {
		return err
	}
	byte = z.encode(byte)
	err = z.conn.WriteMessage(websocket.TextMessage, byte)
	if err != nil {
		return err
	}
	return nil
}

// ToCaller With this function,you can turn a replier to a pure caller
// which means that you can keep this new caller while the replier will be collected and reused by RWeb.
func (z *Replier) ToCaller() *Replier {
	ret := replierPool.Get().(*Replier)
	ret.idCounter = z.idCounter
	ret.bindReplyId = z.bindReplyId
	ret.removeBindReplyId = z.removeBindReplyId
	ret.conn = z.conn
	return ret
}
