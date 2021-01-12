package clientGo

import (
	"errors"
	"nhooyr.io/websocket"
	"reflect"
	"strconv"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type StandardCall struct {
	Function string          `json:"function"`
	Argument []interface{}   `json:"argument"`
	Id       jsoniter.Number `json:"id"`
	IsReply  bool            `json:"reply"` //只是保留，客户端不需要实现 若为Reply时，需要设置为true
}

// 处理服务器调用客户端
func (z *RWebsocketClient) dealWithCall(data []byte) error {
	var call StandardCall
	err := json.Unmarshal(data, &call)
	if err != nil {
		return err
	}
	f, ok := z.functionBind[call.Function]
	if !ok {
		return errors.New("no such function named:" + call.Function)
	}
	v := reflect.ValueOf(f)
	if v.Type().NumIn() != len(call.Argument) && v.Type().NumIn() != len(call.Argument)+1 {
		return errors.New("unexpected argument num")
	}
	replier := Replier{
		id:         call.Id,
		isReply:    true,
		father:     z,
		hasReplied: false,
	}
	replierMap[goid()] = &replier
	defer delete(replierMap, goid())
	var args []reflect.Value
	if v.Type().NumIn() == len(call.Argument)+1 {
		args = make([]reflect.Value, len(call.Argument)+1)
		args[0] = reflect.ValueOf(&replier)
		for i, v := range call.Argument {
			args[i+1] = reflect.ValueOf(v)
		}
	} else {
		args = make([]reflect.Value, len(call.Argument))
		for i, v := range call.Argument {
			args[i] = reflect.ValueOf(v)
		}
	}
	v.Call(args)
	if replier.hasReplied == false { //自动回复nil
		replier.Reply(nil)
	}
	return nil
}
func (z *RWebsocketClient) Call(FunctionName string, Arguments ...interface{}) (interface{}, error) {
	// 准备调用服务器
	callId := atomic.AddInt64(&z.requestId, 1)
	// 挂载chan
	ch := make(chan interface{})
	z.replyConn[callId] = ch
	b,_ := json.Marshal(StandardCall{
		Function: FunctionName,
		Argument: Arguments,
		Id:       jsoniter.Number(strconv.FormatInt(callId, 10)),
		IsReply:  false,
	})
	err := z.conn.Write(z.ctx,websocket.MessageText,b)
	if err != nil {
		return nil, err
	}
	select {
	case <-time.After(time.Second * 10):
		return nil, errors.New("timeout")
	case data := <-ch:
		return data, nil
	}
}
