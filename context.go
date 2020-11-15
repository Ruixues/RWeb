package RWeb

import (
	"github.com/valyala/fasthttp"
	"sync"
)

type Context struct {
	RequestUri string
	Method     int
	RawCtx     *fasthttp.RequestCtx
}
func (z *Context) FormValue (key string) []byte {
	return z.RawCtx.FormValue(key)
}
var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

func NewContext() (ret *Context) {
	ret = contextPool.Get().(*Context)
	return ret
}
func RemoveContext(context *Context) {
	contextPool.Put(context)
}
