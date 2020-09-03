package WebsocketDealer

import (
	"github.com/fasthttp/websocket"
	"sync"
)

type StandardCall struct {
	Function string `json:"function"`
	Argument []interface{} `json:"argument"`
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