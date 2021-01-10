package clientGo

import jsoniter "github.com/json-iterator/go"

type StandardReply struct {
	Id    jsoniter.Number `json:"id"`
	Data  interface{}     `json:"data"`
	Reply bool            `json:"reply"`
}

func (z *RWebsocketClient) dealWithReply(data []byte) {
	var reply StandardReply
	err := json.Unmarshal(data, &reply)
	if err != nil {
		return
	}
	id, err := reply.Id.Int64()
	if err != nil {
		return
	}
	c, ok := z.replyConn[id]
	if !ok { //已经超时，自行删除了
		return
	}
	c <- reply.Data
	close(c)
	delete(z.replyConn, id)
	return
}
