package clientGo

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

type StandardCall struct {
	Function string          `json:"function"`
	Argument []interface{}   `json:"argument"`
	Id       jsoniter.Number `json:"id"`
	IsReply  bool            `json:"reply"` //只是保留，客户端不需要实现 若为Reply时，需要设置为true
}
func (z *RWebsocketClient) dealWithCall(data []byte)  error {
	var call StandardCall
	err := json.Unmarshal(data,&call)
	if err != nil {
		return err
	}
	f,ok := z.functionBind [call.Function]
	if !ok {
		return errors.New("no such function named:" + call.Function)
	}
	replier := Replier{
		id:      call.Id,
		isReply: true,
	}
	v := reflect.ValueOf(f)
	args := make([]reflect.Value,len(call.Argument) + 1)
	args [0] = reflect.ValueOf(&replier)
	for i,v := range call.Argument {
		args [i] = reflect.ValueOf(v)
	}
	v.Call(args)
	return nil
}