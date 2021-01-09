package clientGo

import jsoniter "github.com/json-iterator/go"

type Replier struct {
	id jsoniter.Number
	isReply bool
	father *RWebsocketClient
}

