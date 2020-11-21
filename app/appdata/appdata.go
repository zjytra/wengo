/*
创建时间: 2019/11/23
作者: zjy
功能介绍:

*/

package appdata

import (
	"github.com/panjf2000/ants/v2"
	"os"
	"time"
	"wengo/conf"
	"wengo/csvdata"
	"wengo/model"
	"wengo/xengine"
)

var (
	WorkPool     *ants.Pool
	PathModelPtr *model.PathModel     //最先有路径对象
	AppID        int32                //serverId
	NetConf      *csvdata.Networkconf // 本服务器网络配置在服务器解析参数的时候就获得
	WorldNetConf *csvdata.Networkconf // 世界服务器配置
	appKind      model.AppKind        // app类型 通过外部传递参数确定
	AppFactory   xengine.AppFactory
)

// 创建代理对象
func InitAppData() {
	SetAppPath() // 获取当前路径程序执行路径
	//获取配置
	erro := conf.ReadJson(PathModelPtr.ConfjsonPath)
	if erro != nil {
		panic(erro)
	}
	csvdata.SetCsvPath(PathModelPtr.CsvPath)
	csvdata.LoadCommonCsvData() // 读取公共的csv
}

//设置app程序路径
func SetAppPath() {
	//创建对象在前
	PathModelPtr = model.NewPathModel()
	if PathModelPtr == nil {
		panic("创建 PathModelPtr 失敗")
	}
	pwd, _ := os.Getwd()
	PathModelPtr.SetRootPath(pwd)
	PathModelPtr.InitPathModel()
}

func InitKind() {
	appKind = model.ItoAppKind(NetConf.App_kind)
}

func GetAppKind() model.AppKind {
	return appKind
}

// App 相关数据存放
func SetAppFactory(appFactory xengine.AppFactory) {
	AppFactory = appFactory
}

func SetWorkPool() (erro error) {
	// 协程池,这里要为每个连接开读写线程
	WorkPool, erro = ants.NewPool(NetConf.Goroutines_size,
		ants.WithExpiryDuration(time.Minute*10),
		ants.WithNonblocking(true)) //非阻塞 空闲协程可以存活10分钟
	if erro != nil {
		return erro
	}
	return nil
}

func RealseAppData() {
	//创建对象在前
	PathModelPtr = nil
	WorkPool.Release()
}

func GetSecenName() string {
	switch appKind {
	//gameserver需要区分场景
	case model.APP_GameServer:
		return NetConf.App_name
	//这些服务器器都没有场景名称
	// case model.APP_NONE,model.APP_Client,model.APP_MsgServer,model.APP_DataCenter:
	// 	return ""
	default:
		return ""
	}
	return ""
}
