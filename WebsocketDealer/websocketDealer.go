// @Title  WebsocketDealer.go
// @Description  实现websocket低成本call 需要客户端配合
// @Author  Dorbmon
package WebsocketDealer

import (
	"errors"
	"github.com/Ruixues/RWeb"
	"github.com/Ruixues/RWeb/event"
	"github.com/fasthttp/websocket"
	"github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"sync"
)


const (
	ModuleName    = "WebsocketDealer"
	ModuleVersion = 0.2
)
var json = jsoniter.ConfigFastest	//最快速度
type WebsocketDealFunction func(replier *Replier, session *Session, arguments []interface{})
type WebsocketDealer struct {
	link            map[string]WebsocketDealFunction
	linkLock        *sync.RWMutex
	upgrade         websocket.FastHTTPUpgrader // use default options
	OriginCheck     func(ctx *RWeb.Context) bool
	log             RWeb.Log
	callReplyBind   map[uint64]chan StandardReply
	Events          event.System
	connections     []*ConnectData
	lockConnections *sync.RWMutex
	connectionNum   uint64
}

func New() (r WebsocketDealer) {
	Once.Do(InitWebsocketDealer)
	r.link = make(map[string]WebsocketDealFunction)
	r.linkLock = &sync.RWMutex{}
	r.upgrade = websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			if r.OriginCheck != nil {
				return r.OriginCheck(RWeb.CtxToContext(ctx))
			}
			//默认全部允许
			return true
		},
		EnableCompression: true,
	}
	r.log = &RWeb.DefaultLog{}
	r.connections = make([]*ConnectData, 0)
	r.lockConnections = &sync.RWMutex{}
	r.Events = event.New(EventNum)
	return
}

type Ranger func(data *ConnectData, replier *Replier) error

func (z *WebsocketDealer) BroadCast(ranger Ranger) error {
	z.lockConnections.RLock()
	defer z.lockConnections.RUnlock()
	for _, v := range z.connections {
		replier := replierPool.Get().(*Replier)
		replier.id = jsoniter.Number(0)
		replier.fa = z
		if err := ranger(v, replier); err != nil {
			return err
		}
	}
	return nil
}
func (z *WebsocketDealer) BroadCastIdRange(ranger Ranger, ids []uint64) error {
	z.lockConnections.RLock()
	defer z.lockConnections.RUnlock()
	for _, id := range ids {
		replier := replierPool.Get().(*Replier)
		replier.id = jsoniter.Number(0)
		replier.fa = z
		if err := ranger(z.connections[id], replier); err != nil {
			return err
		}
	}
	return nil
}

/**
  使用此函数作为引擎的绑定函数
*/
func (z *WebsocketDealer) Handler(context *RWeb.Context) {
	err := z.upgrade.Upgrade(context.RawCtx, func(ws *websocket.Conn) {
		var myId = uint64(0)
		defer func() {
			z.lockConnections.Lock()
			defer z.lockConnections.Unlock()
			_ = ws.Close()
			if myId == 0 {
				return
			}
			if myId > (z.connectionNum >> 1) { //应该把后面的合并到前面去
				z.connections = append(z.connections[:myId-1], z.connections[myId:]...)
			} else {
				z.connections = append(z.connections[myId:], z.connections[:myId-1]...)
			}
			z.connectionNum = z.connectionNum - 1
		}()
		conn := newConn(ws)
		defer removeConn(conn)
		var MessageId = uint64(1)
		s := NewSession()
		defer sessionPool.Put(s)
		ok := func() bool {
			data := NewConnectDataPool.Get().(*ConnectData)
			data.Session = s
			data.Context = context
			if err := z.Events.RunEvent(EventNewConnection, func(message event.OnMessage) error {
				ok := message(data).(bool)
				if !ok {
					return errors.New("unexpected error when run event listener")
				}
				return nil
			}); err != nil {
				z.log.FrameworkPrintMessage(ModuleName, err.Error(), -2)
				return false
			}
			return true
		}()
		if !ok {
			return
		}
		func() {
			z.lockConnections.Lock()
			defer z.lockConnections.Unlock()
			z.connections = append(z.connections, &ConnectData{
				Session: s,
				Context: context,
			})
			myId = uint64(len(z.connections))
		}()
		for {
			var SMessage StandardCall
			_, message, err := ws.ReadMessage()
			if err != nil {
				z.log.FrameworkPrintMessage(ModuleName, err.Error(), -2)
				break
			}
			//开始处理到标准格式
			SMessage.IsReply = false
			err = json.Unmarshal(message,&SMessage)
			if err != nil {
				z.log.FrameworkPrintMessage(ModuleName, err.Error(), -2)
				break
			}
			if SMessage.IsReply { //回复
				func() {
					var real StandardReply
					err := json.Unmarshal(message, &real)
					if err != nil {
						z.log.FrameworkPrintMessage(ModuleName, err.Error(), -2)
						return
					}
					id, err := real.Id.Int64()
					if err != nil {
						z.log.FrameworkPrintMessage(ModuleName, err.Error(), -2)
						return
					}
					c, ok := z.callReplyBind[uint64(id)]
					if !ok {
						return
					}
					c <- real
					delete(z.callReplyBind, uint64(id))
				}()
				continue
			}
			//否则开始处理
			var Dealer WebsocketDealFunction
			func() {
				z.linkLock.RLock()
				defer z.linkLock.RUnlock()
				if _, ok := z.link[SMessage.Function]; !ok {
					z.log.FrameworkPrintMessage(ModuleName, "unexpected function:"+SMessage.Function, -2)
					return
				}
				Dealer = z.link[SMessage.Function]
			}()
			if Dealer != nil {
				go func(Dealer WebsocketDealFunction, SMessage StandardCall) {
					replier := replierPool.Get().(*Replier)
					defer replierPool.Put(replier)
					replier.conn = ws
					replier.idCounter = &MessageId
					replier.id = SMessage.Id
					Dealer(replier, s, SMessage.Argument)
				}(Dealer, SMessage)
			}
		}
	})

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			z.log.FrameworkPrintMessage(ModuleName, err.Error(), -1)
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
func (z *WebsocketDealer) BindReplyId(id uint64, c chan StandardReply) {
	z.callReplyBind[id] = c
}
func (z *WebsocketDealer) RemoveBindReplyId(id uint64) {
	delete(z.callReplyBind, id)
}
