// 创建时间: 2019-10-2019/10/17
// 作者: zjy
// 功能介绍:
// 1.主要入口
// 2.
// 3.
package main

import (
	"wengo/app"
	"runtime"
)

// main 初始化工作
func init() {
}
// 各服务器主入口
func main() {
	
	// pro := profile.Start(profile.MemProfile,profile.ProfilePath("./profiles"))
	// 设置最大运行核数
	runtime.GOMAXPROCS(runtime.NumCPU())
	app.NewApp().AppStart()
	// 等待退出 在app 退出后整个程序退出
	// pro.Stop()
}



