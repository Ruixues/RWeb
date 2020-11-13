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
func CtxToContext(ctx *fasthttp.RequestCtx) *Context {
	context := NewContext()
	defer RemoveContext(context)
	context.Method = ToRWebMethod(ctx.Method())
	context.RequestUri = string(ctx.Request.URI().Path())
	context.RawCtx = ctx
	return context
}
func (z *Engine) SetRouter(router Router) {
	z.router = router
}
func (z *Engine) handler(ctx *fasthttp.RequestCtx) {
	context := CtxToContext(ctx)
	handler := z.router.GetHandler(context)
	if handler == nil {
		return
	}
	handler(context)
}
func (z *Engine) RunAndServe(address string) error {
	if z.router == nil {
		return errors.New("the router of engine hasn't been set")
	}
	z.server = &fasthttp.Server{
		Handler: z.handler,
	}
	return z.server.ListenAndServe(address)
}
