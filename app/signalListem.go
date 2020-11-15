/*
创建时间: 2020/7/19
作者: zjy
功能介绍:

*/

package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func (this *App)signalListen() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case s := <-c:
		// 收到信号后的处理，这里只是输出信号内容，可以做一些更有意思的事
		fmt.Println("get signal:", s)
		this.CloseApp()
		break
	}
}
