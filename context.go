package RWeb

import "sync"

type Context struct {
	RequestUri string
	Method     int
}

var contextPool sync.Pool

func NewContext() (ret *Context) {
	ret = contextPool.Get().(*Context)
	if ret == nil {
		ret = &Context{}
	}
	return ret
}
func RemoveContext (context *Context) {
	contextPool.Put(context)
}