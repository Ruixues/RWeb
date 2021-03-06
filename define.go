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
}
type Router interface {
	Bind(address string, method int, handler Handler) error
	GetHandler(context *Context) Handler
	BindSubRouter(Router)  // 挂载子路由
	LoadFromRouter(Router) // 复制子路由的数据
	OutputRules() []*RouterHandler
	AddInterceptor(interceptor Interceptor)
}
type Interceptor func(context *Context) bool //为false的时候，将阻断调用
type Handler func(context *Context)
