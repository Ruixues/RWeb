package WebsocketDealer

import (
	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"sync"
)

type StandardCall struct {
	Function string          `json:"function"`
	Argument []interface{}   `json:"argument"`
	Id       jsoniter.Number `json:"id"`
	IsReply  bool            `json:"reply"` //只是保留，客户端不需要实现 若为Reply时，需要设置为true
}
type StandardReply struct {
	Id   jsoniter.Number `json:"id"`
	Data string          `json:"data"`
}
type Replier struct { //回复者
	conn              *websocket.Conn
	idCounter         *uint64
	id                jsoniter.Number
	bindReplyId       func(id uint64, c chan StandardReply)
	removeBindReplyId func(id uint64)
}

var replierPool = &sync.Pool{
	New: func() interface{} {
		return new(Replier)
	},
}

type Conn struct {
	conn *websocket.Conn
}

var connPool = &sync.Pool{
	New: func() interface{} {
		return new(Conn)
	},
}

func newConn(conn *websocket.Conn) (ret *Conn) {
	ret = connPool.Get().(*Conn)
	ret.conn = conn
	return
}
func removeConn(conn *Conn) {
	connPool.Put(conn)
}

type WebsocketResponse struct {
	Id   jsoniter.Number `json:"id"`
	Data interface{}     `json:"data"`
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
