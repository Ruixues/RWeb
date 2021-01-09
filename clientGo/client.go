package clientGo

import (
	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/gls"
	"log"
)

type BindFunction interface{}

// The first
type RWebsocketClient struct {
	server       string
	conn         *websocket.Conn
	functionBind map[string]BindFunction //The first argument of functionBind must be a *replier
	requestId int64
	replyConn map [int64]chan interface{}
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var replierMap map [int64]*Replier
func R () *Replier {
	r,_ := replierMap[gls.GoID()]
	return r
}
func NewRWebsocketClient(Address string) (ret RWebsocketClient, _ error) {
	var err error
	ret.server = Address
	ret.functionBind = make(map[string]BindFunction)
	ret.conn, _, err = websocket.DefaultDialer.Dial(Address, nil)
	if err != nil {
		return ret, err
	}
	go ret.listener()
	return
}

func (z *RWebsocketClient) listener() {
	defer z.conn.Close()
	for {
		_, message, err := z.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		var tmpStruct struct {
			isReply bool `json:"isReply"`
		}
		if err := json.Unmarshal(message, &tmpStruct); err != nil {
			log.Println(err)
			continue
		}
		if tmpStruct.isReply { //是回复消息
			go z.dealWithReply(message)
		} else {
			go z.dealWithCall(message)
		}
	}
}
