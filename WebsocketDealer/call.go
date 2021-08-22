package WebsocketDealer

import (
	"reflect"

	"github.com/modern-go/gls"
)

func (z *WebsocketDealer) callBind(Dealer WebsocketDealFunction, SMessage StandardCall, makeReplier func() *Replier, s *Session) {
	replier := makeReplier()
	defer replierPool.Put(replier)
	replier.id = SMessage.Id
	replier.haveReplied = false
	goId := gls.GoID()
	replierMap.Store(goId, replier)
	defer replierMap.Delete(goId)
	sessionMap.Store(goId, s)
	defer sessionMap.Delete(goId)
	// 开始调用 拦截器
	for _, f := range z.interceptor {
		if !f(SMessage.Function, replier, s, SMessage.Argument) {
			return
		}
	}
	f := reflect.ValueOf(Dealer)
	// Fill the arguments
	if f.Type().NumIn() != len(SMessage.Argument) && f.Type().NumIn() != len(SMessage.Argument)+2 {
		z.log.FrameworkPrintMessage(ModuleName, "unmatched argument num", 5)
		return
	}
	var args []reflect.Value
	if f.Type().In(0) == replierType && f.Type().In(1) == sessionType {
		args = make([]reflect.Value, len(SMessage.Argument)+2)
		args[0] = reflect.ValueOf(replier)
		args[1] = reflect.ValueOf(s)
		for i, v := range SMessage.Argument {
			args[i+2] = reflect.ValueOf(v)
			realV := args[i+2].Float()
			var nextV interface{} = realV
			switch f.Type().In(i + 2) {
			case int64Type:
				nextV = int64(realV)
			case intType:
				nextV = int(realV)
			}
			args[i+2] = reflect.ValueOf(nextV)
		}
	} else {
		args = make([]reflect.Value, len(SMessage.Argument))
		for i, v := range SMessage.Argument {
			args[i] = reflect.ValueOf(v)
			if args[i].Type() == float64Type { //数字类型，判断一下接收函数
				realV := args[i].Float()
				var nextV interface{} = realV
				switch f.Type().In(i) {
				case int64Type:
					nextV = int64(realV)
				case intType:
					nextV = int(realV)
				}
				args[i] = reflect.ValueOf(nextV)
			}
		}
	}
	// 开始调用
	f.Call(args)
	if !replier.haveReplied {
		replier.Return(nil)
	}
}
