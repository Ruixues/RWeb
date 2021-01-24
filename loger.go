package RWeb

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Log interface {
	//Level 为信息等级，绝对值越小，信息越不重要。框架输出使用负数
	//Module 表示模块
	FrameworkPrintMessage(Module, Message string, Level int)
}

type DefaultLog struct {
}

var moduleName = color.New(color.FgRed, color.BgHiGreen)
var sendTime = color.New(color.FgGreen, color.BgHiYellow)
var message = color.New(color.FgCyan)

func (z *DefaultLog) FrameworkPrintMessage(Module, Message string, Level int) {
	moduleName.Print(Module)
	sendTime.Print(time.Now().Format("2006-01-02 15:04:05"))
	message.Print(Message)
	fmt.Println()
}
