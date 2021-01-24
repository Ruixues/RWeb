package main

import "syscall/js"

func main() {
	js.Global().Set("ttt", js.FuncOf(Test))
	select {}
}
