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

type Router interface {
	Bind(address string, method int, handler Handler) error
	GetHandler(context *Context) Handler
}
type Handler func(context *Context)
