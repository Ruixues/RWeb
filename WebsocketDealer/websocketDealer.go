// @Title  WebsocketDealer.go
// @Description  实现websocket低成本call 需要客户端配合
// @Author  Dorbmon
package WebsocketDealer

import (
	"errors"
	"strconv"
	"sync"

	"github.com/Ruixues/RWeb"
	"github.com/Ruixues/RWeb/event"
	"github.com/fasthttp/websocket"
	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/gls"
	"github.com/valyala/fasthttp"
)

const (
	ModuleName    = "WebsocketDealer"
	ModuleVersion = 0.4
)

var sessionMap sync.Map
var replierMap sync.Map
var json = jsoniter.ConfigFastest //最快速度
type WebsocketDealFunction interface{}
type WebsocketCallInterceptor func(functionName string, replier *Replier, session *Session, args []interface{}) bool
type WebsocketDealer struct {
	link            map[string]WebsocketDealFunction
	linkLock        *sync.RWMutex
	upgrade         websocket.FastHTTPUpgrader // use default options
	OriginCheck     func(ctx *RWeb.Context) bool
	log             RWeb.Log
	Events          event.System
	connections     []*ConnectData
	lockConnections *sync.RWMutex
	connectionNum   uint64
	idPool          NumberPool
	interceptor     []WebsocketCallInterceptor
}

// S get the session
func S() *Session {
	id := gls.GoID()
	s, ok := sessionMap.Load(id)
	if !ok {
		return nil
	}
	return s.(*Session)
}

// R get the replier
func R() *Replier {
	id := gls.GoID()
	s, ok := replierMap.Load(id)
	if !ok {
		return nil
	}
	return s.(*Replier)
}

//New new websocket dealer
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
	r.idPool = NewNumberPool(10)
	r.interceptor = make([]WebsocketCallInterceptor, 0)
	return
}

//Ranger broadCast ranger
type Ranger func(data *ConnectData) error

//BroadCast 广播
func (z *WebsocketDealer) BroadCast(ranger Ranger) error {
	z.lockConnections.RLock()
	defer z.lockConnections.RUnlock()
	for _, v := range z.connections {
		if err := ranger(v); err != nil {
			return err
		}
	}
	return nil
}
func (z *WebsocketDealer) BroadCastIdRange(ranger Ranger, ids []uint64) error {
	z.lockConnections.RLock()
	defer z.lockConnections.RUnlock()
	for _, id := range ids {
		if uint64(len(z.connections)) <= id { //不存在
			return errors.New("unused id:" + strconv.FormatUint(id, 10))
		}
		connection := z.connections[id]
		if connection == nil {
			return errors.New("unused id:" + strconv.FormatUint(id, 10))
		}
		if err := ranger(connection); err != nil {
			return err
		}
	}
	return nil
}

/**
  使用此函数作为引擎的绑定函数
*/
func (z *WebsocketDealer) Handler(context *RWeb.Context) {
	s := NewSession()
	defer sessionPool.Put(s)
	err := z.upgrade.Upgrade(context.RawCtx, func(ws *websocket.Conn) {
		var myId = uint64(0)
		defer func() {
			z.lockConnections.Lock()
			defer z.lockConnections.Unlock()
			_ = ws.Close()
			if myId == 0 {
				return
			}
			replierPool.Put(z.connections[myId].Caller)
			z.connections[myId] = nil
			z.connectionNum = z.connectionNum - 1
		}()
		var MessageId = uint64(1)
		callReplyBind := make(map[uint64]chan StandardReply)
		bindReplyId := func(id uint64, c chan StandardReply) {
			callReplyBind[id] = c
		}
		removeBindReplyId := func(id uint64) {
			delete(callReplyBind, id)
		}
		makeReplier := func() *Replier {
			replier := replierPool.Get().(*Replier)
			replier.conn = ws
			replier.bindReplyId = bindReplyId
			replier.removeBindReplyId = removeBindReplyId
			replier.idCounter = &MessageId
			return replier
		}
		ok := func() bool {
			data := NewConnectDataPool.Get().(*ConnectData)
			data.Session = s
			data.Context = context
			if err := z.Events.RunEvent(EventNewConnection, func(message event.OnMessage) error {
				data.Caller = makeReplier()
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
			myId = z.idPool.Get()
			data := &ConnectData{
				Session: s,
				Context: context,
				Caller:  makeReplier(),
			}
			if uint64(len(z.connections)) <= myId {
				z.connections = append(z.connections, data)
			} else {
				z.connections[myId] = data
			}
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
			err = json.Unmarshal(message, &SMessage)
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
					c, ok := callReplyBind[uint64(id)]
					if !ok {
						return
					}
					c <- real
					delete(callReplyBind, uint64(id))
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
				go z.callBind(Dealer, SMessage, makeReplier, s)
			}
		}
	})

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			z.log.FrameworkPrintMessage(ModuleName, err.Error(), -1)
		}
		// 连接断开，开始调用 拦截器
		data := NewConnectDataPool.Get().(*ConnectData)
		defer NewConnectDataPool.Put(data)
		data.Session = s
		data.Context = context
		z.Events.RunEvent(EventConnectionClose, func(message event.OnMessage) error {
			ok := message(data).(bool)
			if !ok {
				return errors.New("unexpected error when run event listener")
			}
			return nil
		})
		return
	}
}

//BindFunction 绑定函数
func (z *WebsocketDealer) BindFunction(FunctionName string, Function WebsocketDealFunction) error {
	z.linkLock.Lock()
	defer z.linkLock.Unlock()
	if _, ok := z.link[FunctionName]; ok {
		return errors.New("there is already a function called " + FunctionName)
	}
	z.link[FunctionName] = Function
	return nil
}

//AddInterceptor 添加拦截器0
func (z *WebsocketDealer) AddInterceptor(interceptor WebsocketCallInterceptor) {
	z.interceptor = append(z.interceptor, interceptor)
}
