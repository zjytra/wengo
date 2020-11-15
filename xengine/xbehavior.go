//  创建时间: 2019/10/23
//  作者: zjy
//  功能介绍:
//  组件接口
package xengine

type ServerBehavior interface {
	Behavior
	// Updater
}

//一般行为接口
type Behavior interface {
	OnStart()     // 启动
	OnInit() bool //初始化
	OnRelease()   // 关闭On
}
