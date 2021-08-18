package WebsocketDealer

import (
	"reflect"
	"sync"

	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
)

var sessionType = reflect.TypeOf((*Session)(nil))
var replierType = reflect.TypeOf((*Replier)(nil))

type StandardCall struct {
	Function string          `json:"function"`
	Argument []interface{}   `json:"argument"`
	Id       jsoniter.Number `json:"id"`
	IsReply  bool            `json:"reply"` //只是保留，客户端不需要实现 若为Reply时，需要设置为true
}
type StandardReply struct {
	Id   jsoniter.Number `json:"id"`
	Data interface{}     `json:"data"`
}
type Replier struct { //回复者
	conn              *websocket.Conn
	idCounter         *uint64
	id                jsoniter.Number
	bindReplyId       func(id uint64, c chan StandardReply)
	removeBindReplyId func(id uint64)
	haveReplied       bool //是否已经进行了回复。如果没有进行回复，RWeb将自动回复nil
	coder             Coder
}

var replierPool = &sync.Pool{
	New: func() interface{} {
		return new(Replier)
	},
}

type WebsocketResponse struct {
	Id      jsoniter.Number `json:"id"`
	Data    interface{}     `json:"data"`
	IsReply bool            `json:"reply"`
}

var responsePool = &sync.Pool{
	New: func() interface{} {
		return new(WebsocketResponse)
	},
}
var requestPool = &sync.Pool{
	New: func() interface{} {
		return new(StandardCall)
	},
}
var Once sync.Once

func InitWebsocketDealer() {

}
