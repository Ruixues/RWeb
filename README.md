# RWeb

Ruixue web framework

仅为后端制作，不包含前端渲染等。

目标是作为瑞雪Go后端的统一web框架。

服务端为GO语言实现。目前服务端支持两种函数绑定方式：

函数第一个与第二个参数分别为*Replier，*Session

或者为只具有协定的参数，并通过RWeb.R(),RWeb.S()获取这两个值。