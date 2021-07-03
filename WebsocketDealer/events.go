package WebsocketDealer

import (
	"sync"

	"github.com/Ruixues/RWeb"
)

const EventNum = 3
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
	EventFunctionCall
	/*
		EventConnectionClose call with *FunctionCall
		It will be called when client calls the server
	*/
)

type ConnectData struct {
	Session *Session
	Context *RWeb.Context
	Caller  *Replier
}
type FunctionCall struct {
	FunctionName string
	Argument     []interface{}
	Session      *Session
}

var NewConnectDataPool = sync.Pool{New: func() interface{} {
	return new(ConnectData)
}}
var FunctionCallPool = sync.Pool{New: func() interface{} {
	return new(FunctionCall)
}}
