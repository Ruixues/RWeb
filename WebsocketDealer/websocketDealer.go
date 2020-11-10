// @Title  WebsocketDealer.go
// @Description  实现websocket低成本call 需要客户端配合
// @Author  Dorbmon
package WebsocketDealer

import (
	"github.com/Ruixues/RWeb"
	"errors"
	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"sync"
)
var json = jsoniter.ConfigCompatibleWithStandardLibrary
const (
	ModuleName = "WebsocketDealer"
	ModuleVersion = 0.1
)

type WebsocketDealFunction func(conn *Conn,arguments []interface{}) interface{}
type WebsocketDealer struct {
	link     map[string]WebsocketDealFunction
	linkLock *sync.RWMutex
	upgrade  websocket.FastHTTPUpgrader // use default options
	OriginCheck func (ctx *RWeb.Context) bool
	log RWeb.Log
}

func New() (r WebsocketDealer) {
	Once.Do(InitWebsocketDealer)
	r.link = make(map[string]WebsocketDealFunction)
	r.linkLock = &sync.RWMutex{}
	r.upgrade = websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			if r.OriginCheck != nil {
				return r.OriginCheck (RWeb.CtxToContext(ctx))
			}
			//默认全部允许
			return true
		},
		EnableCompression: true,
	}
	return
}

/**
	使用此函数作为引擎的绑定函数
*/
func (z*WebsocketDealer) Handler(context *RWeb.Context) {
	err := z.upgrade.Upgrade(context.RawCtx, func(ws *websocket.Conn) {
		conn := newConn (ws)
		defer removeConn(conn)
		defer ws.Close()
		var SMessage StandardCall
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				z.log.FrameworkPrintMessage(ModuleName,err.Error(),-2)
				break
			}
			//开始处理到标准格式
			err = json.Unmarshal(message,&SMessage)
			if err != nil {
				z.log.FrameworkPrintMessage(ModuleName,err.Error(),-2)
				break
			}
			//否则开始处理
			var Dealer WebsocketDealFunction
			func () {
				z.linkLock.RLock()
				defer z.linkLock.RUnlock()
				if _,ok := z.link [SMessage.Function];!ok {
					z.log.FrameworkPrintMessage(ModuleName,"unexpected function:" + SMessage.Function,-2)
					return
				}
				Dealer = z.link [SMessage.Function]
			} ()
			if Dealer != nil {
				data := Dealer (conn,SMessage.Argument)
				if data == nil {
					continue
				}
				func () {
					res := ResponsePool.Get().(*WebsocketResponse)
					defer ResponsePool.Put(res)
					res.data = data
					res.id = SMessage.Id
					byte,err := json.Marshal(res)
					if err != nil {
						message := "Error when call function:" + SMessage.Function + " with data:" + string(message) + ".\n error:" + err.Error();
						z.log.FrameworkPrintMessage(ModuleName,message,-1)
						return
					}
					err = ws.WriteMessage(websocket.TextMessage,byte)
					if err != nil {
						message := "Error when call function:" + SMessage.Function + " with data:" + string(message) + ".\n error:" + err.Error();
						z.log.FrameworkPrintMessage(ModuleName,message,-1)
						return
					}
				} ()

			}
		}
	})

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			z.log.FrameworkPrintMessage(ModuleName,err.Error(),-1)
		}
		return
	}
}
func (z *WebsocketDealer) BindFunction(FunctionName string, Function WebsocketDealFunction) error {
	z.linkLock.Lock()
	defer z.linkLock.Unlock()
	if _, ok := z.link[FunctionName]; ok {
		return errors.New("there is already a function called " + FunctionName)
	}
	z.link[FunctionName] = Function
	return nil
}
