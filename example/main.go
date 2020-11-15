package main

import (
	"fmt"
	"github.com/Ruixues/RWeb"
	"github.com/Ruixues/RWeb/WebsocketDealer"
)
func t (replier *WebsocketDealer.Replier, arguments []interface{}){
	replier.Return("Hello World!!")
}
func main() {
	fmt.Println("Go")
	router := RWeb.NewDefaultRouter()
	router.Bind("/test", RWeb.MethodAll, func(ctx *RWeb.Context) {
		ctx.RawCtx.WriteString("test")
	})
	ws := WebsocketDealer.New()
	router.Bind("/t",RWeb.MethodAll,ws.Handler)
	ws.BindFunction("test",t)
	Engine := RWeb.NewEngine(&router)
	err := Engine.RunAndServe("0.0.0.0:1111")
	if err != nil {
		fmt.Println(err)
	}
}
