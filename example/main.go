package main

import (
	"fmt"

	"github.com/Ruixues/RWeb"
	"github.com/Ruixues/RWeb/WebsocketDealer"
)

func t(replier *WebsocketDealer.Replier, session *WebsocketDealer.Session) {
	res, _ := replier.Call("test", "Hi")
	fmt.Println(res)
	err := replier.Return("Hello World!!")
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	fmt.Println("Server is running")
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
