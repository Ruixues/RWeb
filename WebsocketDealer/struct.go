package WebsocketDealer

import (
	"github.com/fasthttp/websocket"
	"sync"
)

type StandardCall struct {
	Function string `json:"function"`
	Argument []interface{} `json:"argument"`
	Id int64 `json:"id"`
}
type Conn struct {
	conn *websocket.Conn
}
var connPool = &sync.Pool{
	New: func () interface{} {
		return new (Conn)
	},
}
func newConn (conn *websocket.Conn) (ret *Conn){
	ret = connPool.Get().(*Conn)
	ret.conn = conn
	return
}
func removeConn (conn *Conn) {
	connPool.Put(conn)
}
type WebsocketResponse struct {
	id int64 `json:"id"`
	data interface {} `json:"data"`
}
var ResponsePool = &sync.Pool{
	New: func() interface {} {
		return new (WebsocketResponse)
	},
}
var RequestPool = &sync.Pool{
	New: func() interface{} {
		return &StandardCall {}
	},
}
var Once sync.Once
func InitWebsocketDealer () {

}