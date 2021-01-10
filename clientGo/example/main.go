package main

import (
	"fmt"
	"github.com/Ruixues/RWeb/clientGo"
)

func main() {
	client, err := clientGo.NewRWebsocketClient("ws://127.0.0.1:1111/t")
	if err != nil {
		fmt.Println(err)
		return
	}
	client.Bind("test", func(Arg string) {
		fmt.Println("Got:", Arg)
		//clientGo.R().Reply("Ok")
	})
	data, _ := client.Call("test")
	fmt.Println(data)
}
