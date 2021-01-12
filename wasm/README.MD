### RWeb Go语言实现的wasm模块

注意，受限于当前无法获取id,在大量请求下，可能会出现Replier重复的bug.

请尽量使用func(*Replier,arguments) 的方式绑定函数。可以避免bug的发生。

