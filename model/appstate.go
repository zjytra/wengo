/*
创建时间: 2020/5/20
作者: zjy
功能介绍:

*/

package model

//app的状态常量
const (
	AppState_Run int32 = 1      // app运行中
	AppState_Maintain int32 = 2 // app维护
	AppState_Close int32 = 3    // app 关闭
)


type AppState struct {
	flag *AtomicInt32FlagModel // app 狀態
}


func (this *AppState)InitAppState()  {
	this.flag = NewAtomicInt32Flag()
}

func (this *AppState)ChangeAppState(astate int32)  {
	this.flag.SetInt32(astate)
}

//获取app状态
func (this *AppState)GetAppSate() int32 {
	return this.flag.GetInt32()
}
//app 是否在运行中
func (this *AppState)AppIsRun() bool {
	return this.GetAppSate() == AppState_Run
}
//app 是否在维护
func (this *AppState)AppIsMaintain() bool {
	return this.GetAppSate() == AppState_Maintain
}
//app 是否关闭
func (this *AppState)AppIsClose() bool {
	return this.GetAppSate() == AppState_Close
}

// 设置启动标志位
func (this *AppState)AppOpen() {
	this.ChangeAppState(AppState_Run)
}
// 设置程序维护
func (this *AppState)AppMaintain() {
	this.ChangeAppState(AppState_Maintain)
}
// 设置程序维护
func (this *AppState)AppClose() {
	this.ChangeAppState(AppState_Close)
}