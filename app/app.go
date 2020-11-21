// 创建时间: 2019/10/17
// 作者: zjy
// 功能介绍:
// 程序最外层 ,这里给main的入口,以及整个进程退出的控制
// 包含程序的启动,停止
package app

import (
	"fmt"
	"sync"
	"time"
	"wengo/app/appdata"
	"wengo/model"
	"wengo/timersys"
	"wengo/xengine"
	"wengo/xlog"
)



type App struct {
	model.AppState // app 狀態
	appBehavior xengine.ServerBehavior
	AppWG        sync.WaitGroup      // app进程结束标志
	appCloseTime time.Duration               // app关闭倒计时
	tc *timersys.TimeTicker
}

func NewApp() *App {
	return new(App)
}

//关闭程序倒计时
func (this *App)SetAppCloseTime(time time.Duration)  {
	this.appCloseTime = time
}

func (this *App)SetAppOpen()  {
	this.InitAppState()
	this.AppOpen()
}

func (this *App)AppRun() error   {
	// 运行app 逻辑
	// appBehavior.SendHeartToWS()
	for   {
		xlog.DebugLogNoInScene( "Main update")
		
		
		// app 关闭需要通知所有连接倒计时关闭时间
		if this.AppIsClose() {
			
			
			//进入关闭状态倒计时
			this.appCloseTime -= (time.Millisecond * 100)
			if this.appCloseTime <= 0 {
				//倒计时没有真的关闭程序了
				this.onCloseApp()
			}
		}
		
		
		time.Sleep(time.Millisecond * 100)
	}

	return  nil
}


// 关闭app 程序
func (this *App)CloseApp() {
	this.AppClose()
}

// 执行关闭
func (this *App)onCloseApp() {
	this.appBehavior.OnRelease() // 进程结束
	appdata.RealseAppData()
	xlog.CloseLog() // 退出日志
	timersys.Release()
	this.AppWG.Done()
	fmt.Println("App Close")
}


