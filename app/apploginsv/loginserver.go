/*
创建时间: 2019/11/24
作者: zjy
功能介绍:
登录服
*/

package apploginsv

import (
	"github.com/zjytra/wengo/timersys"
	"time"
)



type LogionServer struct {
	oneMinuteTimeID uint32//定时器id
}

// 程序启动
func (this *LogionServer)OnStart() {
	initOK := this.OnInit()
	if !initOK {
		panic("LogionServer 初始化失败")
	}
	this.AddTimer()
}

//初始化
func (this *LogionServer)OnInit() bool{
	// csvdata.LoadLoginCsvData()
	InitData()
	return true
}

// 程序运行
func (this *LogionServer)OnUpdate(){
	return
}
// 关闭
func (this *LogionServer)OnRelease(){
	ReleaseData()
}

func (this *LogionServer)AddTimer(){
	this.oneMinuteTimeID = timersys.NewWheelTimer(time.Minute, this.PerOneMinuteTimer,DispSys) //每分钟调用
}

func (this *LogionServer)ReleaseTimer(){
}





