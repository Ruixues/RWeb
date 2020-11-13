package WebsocketDealer

import (
	"github.com/fasthttp/websocket"
	"sync"
)

type StandardCall struct {
	Function string        `json:"function"`
	Argument []interface{} `json:"argument"`
	Id       uint64        `json:"id"`
	IsReply  bool          `json:"reply"` //只是保留，客户端不需要实现 若为Reply时，需要设置为true
}
type StandardReply struct {
	Id   uint64      `json:"id"`
	Data interface{} `json:"data"`
}
type Replier struct { //回复者
	conn      *websocket.Conn
	idCounter *uint64
	fa        *WebsocketDealer
	id        uint64
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
	id   uint64      `json:"id"`
	data interface{} `json:"data"`
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
