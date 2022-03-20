package RWeb

import (
	"errors"
	"sync"
)

type bindData struct {
	Method  int
	Handler Handler
	Tags    []string
}
type DefaultRouter struct {
	bindLock     *sync.RWMutex
	linker       map[string][]bindData
	subRouter    []Router
	subRouterBuf map[string]Router
	basicRoute   string
	interceptors []Interceptor
}

func bindData2RouterHandler(address string, bindData *bindData) *RouterHandler {
	return &RouterHandler{
		Address: address,
		Handler: bindData.Handler,
		Method:  bindData.Method,
		Tags:    bindData.Tags,
	}
}
func NewDefaultRouter() (r DefaultRouter) {
	r.bindLock = &sync.RWMutex{}
	r.linker = make(map[string][]bindData)
	r.subRouter = make([]Router, 0)
	r.subRouterBuf = make(map[string]Router)
	r.interceptors = make([]Interceptor, 0)
	return
}
func (z *DefaultRouter) SetBasicRoute(basic string) {
	z.basicRoute = basic
}
func (z *DefaultRouter) Bind(address string, method int, handler Handler, tag []string) error {
	address += z.basicRoute
	z.bindLock.Lock()
	defer z.bindLock.Unlock()
	if _, ok := z.linker[address]; !ok {
		z.linker[address] = make([]bindData, 0)
	}
	for _, v := range z.linker[address] {
		if v.Method == method {
			return errors.New("you are trying to bind an existed address")
		}
	}
	z.linker[address] = append(z.linker[address], bindData{
		Method:  method,
		Handler: handler,
		Tags:    tag,
	})
	return nil
}
func (z *DefaultRouter) GetHandler(context *Context) *RouterHandler {
	var handler *RouterHandler
	handlers, ok := z.linker[context.RequestUri]
	if !ok { // 那就在子路由内
		handler = z.GetFromSubRouter(context)
	} else {
		for _, v := range handlers {
			if v.Method != MethodAll && v.Method != v.Method {
				continue
			}
			handler = bindData2RouterHandler(context.RequestUri, &v)
			break
		}
	}
	if handler == nil {
		return nil
	}
	for _, v := range z.interceptors {
		if !v(context, handler) {
			return nil
		}
	}
	return handler
}
func (z *DefaultRouter) GetFromSubRouter(context *Context) *RouterHandler {
	r, ok := z.subRouterBuf[context.RequestUri]
	if ok {
		if handler := r.GetHandler(context); handler != nil {
			return handler
		}
		delete(z.subRouterBuf, context.RequestUri)
	}
	for _, router := range z.subRouter {
		if handler := router.GetHandler(context); handler != nil {
			z.subRouterBuf[context.RequestUri] = router
			return handler
		}
	}
	return nil
}
func (z *DefaultRouter) BindSubRouter(other Router) {
	z.bindLock.Lock()
	defer z.bindLock.Unlock()
	// 开始挂载
	z.subRouter = append(z.subRouter, other)
	return
}
func (z *DefaultRouter) LoadFromRouter(other Router) {
	list := other.OutputRules()
	for _, v := range list {
		z.Bind(v.Address, v.Method, v.Handler, v.Tags)
	}
}
func (z *DefaultRouter) OutputRules() []*RouterHandler {
	ret := make([]*RouterHandler, 0)
	for k, v := range z.linker {
		for _, rule := range v {
			ret = append(ret, bindData2RouterHandler(k, &rule))
		}
	}
	return ret
}
func (z *DefaultRouter) AddInterceptor(interceptor Interceptor) {
	z.interceptors = append(z.interceptors, interceptor)
}
