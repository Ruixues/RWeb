package WebsocketDealer

import (
	"github.com/Ruixues/RWeb"
	"sync"
)

const EventNum = 2
const (
	EventNewConnection = iota
	/*
		 	EventNewConnection call with *NewConnectData
			The caller of NewConnectData won't be collected.
	*/
	EventConnectionClose
	/*
		EventConnectionClose call with *NewConnectData.
		But you shouldn't write any data cause the connection had been close.
	*/
)

type ConnectData struct {
	Session *Session
	Context *RWeb.Context
	Caller  *Replier
}

var NewConnectDataPool = sync.Pool{New: func() interface{} {
	return new(ConnectData)
}}
