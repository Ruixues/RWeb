// Package RWeb @Title  Engine.go
// @Description  引擎的主要实现文件
// @Author  Dorbmon
package RWeb

import (
	"errors"
	"fmt"
	"os"
	"os/signal"

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
	defer RemoveContext(context)
	handler := z.router.GetHandler(context)
	if handler == nil {
		return
	}
	handler.Handler(context)
}
func (z *Engine) RunAndServe(address string) error {
	if z.router == nil {
		return errors.New("the router of engine hasn't been set")
	}
	z.server = &fasthttp.Server{
		Handler: z.handler,
	}
	go z.waitShutdown()
	return z.server.ListenAndServe(address)
}
func (z *Engine) waitShutdown() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := z.server.Shutdown(); err != nil {
		fmt.Println(err)
	}
}
func (z *Engine) RunAndServeTLS(address string, certData, keyData []byte) error {
	if z.router == nil {
		return errors.New("the router of engine hasn't been set")
	}
	z.server = &fasthttp.Server{
		Handler: z.handler,
	}
	go z.waitShutdown()
	return z.server.ListenAndServeTLSEmbed(address, certData, keyData)
}
