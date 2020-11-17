package event

import "fmt"

type OnMessage func(interface{}) interface{}
type System struct {
	link     [][]OnMessage
	eventNum uint64
}

func New(eventNum uint64) (ret System) {
	ret.link = make([][]OnMessage, eventNum)
	ret.eventNum = eventNum
	return
}
func (z *System) On(Event uint64, Handler OnMessage) {
	if Event >= z.eventNum {
		panic(fmt.Sprintf("unexpected event:%d for the max of event is %d", Event, z.eventNum))
	}
	if z.link[Event] == nil {
		z.link[Event] = make([]OnMessage, 0)
	}
	z.link[Event] = append(z.link[Event], Handler)
}

// If runner returns false,the RunEvent will quit with an error from that function
func (z *System) RunEvent(Event uint64, Runner func(message OnMessage) error) error {
	for _, k := range z.link[Event] {
		if err := Runner(k); err != nil {
			return err
		}
	}
	return nil
}
