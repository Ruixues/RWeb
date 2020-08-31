// @Title  Engine.go
// @Description  引擎的主要实现文件
// @Author  Dorbmon
package RWeb

import (
	"errors"
	"github.com/valyala/fasthttp"
)

type Engine struct {
	router Router
	server *fasthttp.Server
}

func NewEngine(router Router) (e Engine) {
	e.SetRouter(router)
	return
}
func (z *Engine) SetRouter(router Router) {
	z.router = router
}
func (z *Engine) handler(ctx *fasthttp.RequestCtx) {
	context := NewContext()
	context.Method = ToRWebMethod(ctx.Method())
	context.RequestUri = string(ctx.RequestURI())
	z.router.GetHandler(context) (context)
	RemoveContext(context)
}
func (z *Engine) RunAndServe(address string) error {
	if z.router == nil {
		return errors.New("the router of engine hasn't been set")
	}
	z.server = &fasthttp.Server{
		Handler: z.handler,
	}
	return z.RunAndServe(address)
}
