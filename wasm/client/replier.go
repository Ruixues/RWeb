package clientGo

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"nhooyr.io/websocket"
)

type Replier struct {
	id         jsoniter.Number
	isReply    bool
	father     *RWebsocketClient
	hasReplied bool
}

func (z *Replier) Reply(Data interface{}) error {
	if z.id == "" {
		return errors.New("it is not a replier but a pure caller")
	}
	byte, err := json.Marshal(&StandardReply{
		Id:    z.id,
		Data:  Data,
		Reply: true,
	})
	if err != nil {
		return err
	}
	err = z.father.conn.Write(z.father.ctx,websocket.MessageText, byte)
	if err != nil {
		return err
	}
	return nil
}
func (z *Replier) Call(FunctionName string, Arguments ...interface{}) (interface{}, error) {
	return z.father.Call(FunctionName, Arguments)
}
