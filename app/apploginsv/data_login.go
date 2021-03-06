/*
创建时间: 2020/5/17
作者: zjy
功能介绍:
登陆服数据处理
*/

package apploginsv

import (
	"github.com/zjytra/wengo/app/appdata"
	"github.com/zjytra/wengo/app/netmsgsys"
	"github.com/zjytra/wengo/dispatch"
)

var (
	Client     *ClientHandler        //处理客户端端相关事件
	DataCenter *DataCenterMsgHandler //数据中心相关连接处理
	DispSys    *dispatch.DispatchSys //一个调度器就是单线程处理业务逻辑
	pLoginAccountMgr *LoginAccountMgr//
)

func InitData() {
	//检查初始化数据
	if appdata.WorkPool == nil {
		panic(" InitData appdata.WorkPool is nil")
	}
	if appdata.NetConf == nil || appdata.WorldNetConf == nil {
		panic("InitData 配置有误")
	}
	DispSys = dispatch.NewDispatchSys() // 系统调度对象
	if DispSys == nil {
		panic("InitData DispSys is nil")
	}
	//客户端连接处理逻辑
	Client = NewClientHandle(appdata.NetConf, appdata.WorkPool, DispSys)
	// 世界服消息处理对象
	DataCenter = NewDataCenterHandle(appdata.WorldNetConf, appdata.WorkPool, DispSys)
	// 初始化消息验证
	netmsgsys.InitMsgVerify()
	//创建账号管理类
	pLoginAccountMgr = NewLoginAccountMgr()
}

//获取客户端连接
func GetClientConnSize() int32 {
	return Client.tcpServer.GetConnectSize()
}

func GetDispatchSys() *dispatch.DispatchSys {
	return DispSys
}

// 释放
func ReleaseData() {
	Client.OnRelease()
	DataCenter.OnRelease()
	DispSys.Release()
	netmsgsys.ReleaseData()
	pLoginAccountMgr.Release()
}
