package WebsocketDealer

import (
	"github.com/modern-go/gls"
	"reflect"
)

func (z *WebsocketDealer) Call (Dealer WebsocketDealFunction, SMessage StandardCall,makeReplier func()*Replier,s *Session) {
	replier := makeReplier()
	defer replierPool.Put(replier)
	replier.id = SMessage.Id
	goId := gls.GoID()
	replierMap.Store(goId,replier)
	defer replierMap.Delete(goId)
	sessionMap.Store(goId,s)
	defer sessionMap.Delete(goId)
	f := reflect.ValueOf(Dealer)
	// Fill the arguments
	if f.Type().NumIn() != len(SMessage.Argument) {
		z.log.FrameworkPrintMessage(ModuleName,"unmatched argument num",5)
		return
	}
	var args []reflect.Value
	if f.Type().In(0) == replierType && f.Type().In(1) == sessionType {
		args = make([]reflect.Value,len(SMessage.Argument) + 2)
		args [0] = reflect.ValueOf(replier)
		args [1] = reflect.ValueOf(s)
		for i,v := range SMessage.Argument {
			args [i + 2] = reflect.ValueOf(v)
		}
	} else {
		args = make([]reflect.Value,len(SMessage.Argument))
		for i,v := range SMessage.Argument {
			args [i] = reflect.ValueOf(v)
		}
	}
	// 开始调用
	f.Call(args)
	return
}