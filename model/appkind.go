// 创建时间: 2019/10/17
// 作者: zjy
// 功能介绍:
// 1.常量的定义
// 2.
// 3.
package model

type AppKind int32

const (
	APP_NONE        AppKind = 0 // 无类型
	APP_Client              = 1 // 客户端
	APP_MsgServer           = 2 // 聊天服
	APP_LoginServer         = 3 // 登陆服
	APP_GATEWAY             = 4 // 网关
	APP_DataCenter          = 5 // 数据中心一般由这个服务器操作数据库
	APP_GameServer          = 6 // 游戏服 各种场景处理
	APP_MAX                 = 7
)


var kindArr  =[...]AppKind{
	APP_NONE,
	APP_Client,
	APP_MsgServer,
	APP_LoginServer,
	APP_GATEWAY,
	APP_DataCenter,
	APP_GameServer,
	APP_MAX,
}

// 整数变为AppKind
func ItoAppKind(val int32) AppKind {
	if val >= 0 && val < int32(len(kindArr)) {
		return kindArr[val]
	}
	return APP_NONE
}

var appNames = [...]string{
	"none",
	"appclient",
	"msgsv",
	"loginsv",
	"gateway",
	"datacenter",
	"gamesv",
	"none",
}
func ToKindString(ak int32) string {
	if ak >= int32(APP_NONE)  && ak < int32(len(appNames)) {
		return  appNames[ak]
	}
	return  appNames[APP_NONE]
}
