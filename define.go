package RWeb

const (
	MethodPost = iota
	MethodGet
	MethodHead
	MethodPut
	MethodPatch
	MethodDelete
	MethodConnect
	MethodOptions
	MethodTrace
	MethodAll
	MethodUnknown
)

type RouterHandler struct {
	Address string
	Handler Handler
	Method  int
	Tags    []string
}
type Router interface {
	// tag用来标识API类型
	Bind(address string, method int, handler Handler, tag []string) error
	GetHandler(context *Context) *RouterHandler
	BindSubRouter(Router)  // 挂载子路由
	LoadFromRouter(Router) // 复制子路由的数据
	OutputRules() []*RouterHandler
	AddInterceptor(interceptor Interceptor)
}
type Interceptor func(context *Context, handler *RouterHandler) bool //为false的时候，将阻断调用
type Handler func(context *Context)
