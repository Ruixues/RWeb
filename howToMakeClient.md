在建立websocket连接后，所有数据请使用JSON格式传递。

之后，请以以下格式传递参数:
{
"function" : "函数名称",
"argument" : [参数列表],
"id": 请求ID } 请求ID用来标记本次请求，方便客户端实现对应回调。 调用服务端函数后，服务器将返回此数据

{
"id" : 请求ID,
"data": 数据,
"reply": true }

服务器可能也会主动调用客户端函数.调用形式为:
{
"function" : "函数名称",
"argument" : [参数列表],
"id": 请求ID }

请你返回:
{
"id" : 请求ID,
"data": 数据
"reply": true }

默认超时为3s

客户端实现语言如果有async await,推荐此类实现.