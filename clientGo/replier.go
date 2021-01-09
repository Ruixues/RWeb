package clientGo

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"strconv"
	"sync/atomic"
	"time"
)

type Replier struct {
	id      jsoniter.Number
	isReply bool
	father  *RWebsocketClient
}

func (z *Replier) Reply (Data interface{}) error {
	if z.id == "" {
		return errors.New("it is not a replier but a pure caller")
	}
	err := z.father.conn.WriteJSON(StandardReply{
		Id:   z.id,
		Data: Data,
	})
	if err != nil {
		return err
	}
	return nil
}
func (z *Replier) Call(FunctionName string,Arguments ...interface{}) (interface{},error) {
	// 准备调用服务器
	callId := atomic.AddInt64(&z.father.requestId,1)
	// 挂载chan
	ch := make(chan interface{})
	z.father.replyConn [callId] = ch
	err := z.father.conn.WriteJSON(StandardCall{
		Function: FunctionName,
		Argument: Arguments,
		Id:       jsoniter.Number(strconv.FormatInt(callId,10)),
		IsReply:  false,
	})
	if err != nil {
		return nil,err
	}
	select {
	case <-time.After(time.Second * 10):
		return nil,errors.New("timeout")
	case data := <- ch:
		return data,nil
	}
}