package main

import (
	"fmt"
	"github.com/Ruixues/RWeb"
	"github.com/Ruixues/RWeb/WebsocketDealer"
)

func t(replier *WebsocketDealer.Replier, session *WebsocketDealer.Session, arguments []interface{}) {
	replier.Call("test", "Hi")
	err := replier.Return("Hello World!!")
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	fmt.Println("Go")
	router := RWeb.NewDefaultRouter()
	router.Bind("/test", RWeb.MethodAll, func(ctx *RWeb.Context) {
		ctx.Write([]byte("<html><h1>Hello RWeb!</h1></html>"))
	})
	ws := WebsocketDealer.New()
	router.Bind("/t", RWeb.MethodAll, ws.Handler)
	ws.BindFunction("test", t)
	Engine := RWeb.NewEngine(&router)
	err := Engine.RunAndServe("0.0.0.0:1111")
	if err != nil {
		fmt.Println(err)
	}
}
