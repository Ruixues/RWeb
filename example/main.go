package main

import (
	"RWeb"
	"fmt"
)

func main() {
	fmt.Println ("Go")
	router := RWeb.NewDefaultRouter()
	router.Bind("/test",RWeb.MethodAll,func (ctx *RWeb.Context) {
		fmt.Println ("Get")
		ctx.RawCtx.WriteString("test")
	})
	Engine := RWeb.NewEngine(&router)
	err := Engine.RunAndServe(":1111")
	if err != nil {
		fmt.Println(err)
	}
}
