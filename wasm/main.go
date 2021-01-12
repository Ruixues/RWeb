package main

import (
	clientGo "client.go/client"
	"syscall/js"

)

func p(data ...string) {
	alert := js.Global().Get("alert")
	for _, i := range data {
		alert.Invoke(i)
	}
}
func main() {
	p("Hello")
	client, err := clientGo.NewRWebsocketClient("ws://127.0.0.1:1111/t")
	if err != nil {
		p(err.Error())
		return
	}
	p("Success")
	client.Bind("test", func(Arg string) {
		p(Arg)
		//clientGo.R().Reply("Ok")
	})
	data, _ := client.Call("test")
	//fmt.Println(data)
	p(data.(string))
}
