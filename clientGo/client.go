package clientGo

import (
	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
)
type BindFunction interface {}
// The first
type RWebsocketClient struct {
	server string
	conn *websocket.Conn
	functionBind map[string]BindFunction	//The first argument of functionBind must be a *replier
}
var json = jsoniter.ConfigCompatibleWithStandardLibrary
func NewRWebsocketClient (Address string) (ret RWebsocketClient,_ error){
	ret.server = Address
	ret.functionBind = make(map[string]BindFunction)
	c, _, err := websocket.DefaultDialer.Dial(Address, nil)
	if err != nil {
		return ret,err
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
			isReply bool	`json:"isReply"`
		}
		if err := json.Unmarshal(message,&tmpStruct);err != nil {
			log.Println(err)
			continue
		}
		if tmpStruct.isReply {	//是回复消息
			go z.dealWithReply(message)
		} else {
			go z.dealWithCall(message)
		}
	}
}