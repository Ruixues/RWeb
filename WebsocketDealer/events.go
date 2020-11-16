package WebsocketDealer

import (
	"github.com/Ruixues/RWeb"
	"sync"
)

const EventNum = 1
const (
	EventNewConnection = iota
	/*
	 EventNewConnection call with *NewConnectData
	 */
)
type NewConnectData struct {
	Session *Session
	Context *RWeb.Context
}
var NewConnectDataPool = sync.Pool{New: func() interface{} {
	return new (NewConnectData)
}}