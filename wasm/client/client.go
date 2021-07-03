/// +build wasm

package clientGo

import (
	"context"
	"errors"
	"log"
	"time"

	jsoniter "github.com/json-iterator/go"
	"nhooyr.io/websocket"
)

type BindFunction interface{}

// The first
type RWebsocketClient struct {
	server       string
	conn         *websocket.Conn
	functionBind map[string]BindFunction //The first argument of functionBind must be a *replier
	requestId    int64
	replyConn    map[int64]chan interface{}
	ctx          context.Context
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var replierMap = make(map[int64]*Replier)

func R() *Replier {
	r, _ := replierMap[goid()]
	return r
}
func NewRWebsocketClient(Address string) (ret RWebsocketClient, _ error) {
	var err error
	ret.server = Address
	ret.functionBind = make(map[string]BindFunction)
	//var cancel context.CancelFunc
	ret.ctx, _ = context.WithTimeout(context.Background(), time.Minute)
	//defer cancel()
	ret.conn, _, err = websocket.Dial(ret.ctx, Address, nil)
	ret.replyConn = make(map[int64]chan interface{})
	if err != nil {
		return ret, err
	}
	go ret.listener()
	return
}

func (z *RWebsocketClient) listener() {
	defer z.conn.Close(0, "client is fucked down")
	for {
		_, message, err := z.conn.Read(z.ctx)
		if err != nil {
			log.Println("read:", err)
			return
		}
		var tmpStruct struct {
			Reply bool `json:"reply"`
		}
		if err := json.Unmarshal(message, &tmpStruct); err != nil {
			log.Println(err)
			continue
		}
		if tmpStruct.Reply { //是回复消息
			go z.dealWithReply(message)
		} else {
			go z.dealWithCall(message)
		}
	}
}
func (z *RWebsocketClient) Bind(FunctionName string, Function BindFunction) error {
	_, ok := z.functionBind[FunctionName]
	if ok {
		return errors.New("the function:" + FunctionName + " has existed")
	}
	z.functionBind[FunctionName] = Function
	return nil
}
