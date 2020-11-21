//  创建时间: 2019/10/23
//  作者: zjy
//  功能介绍:
//  app 初始化工作
package app

import (
	"flag"
	"fmt"
	"wengo/app/appclient"
	"wengo/app/appgate"
	"wengo/app/apploginsv"
	"wengo/app/datacenter"
	"wengo/appdata"
	"wengo/conf"
	"wengo/csvdata"
	"wengo/global"
	"wengo/model"
	"wengo/xengine"
	"wengo/xlog"
	_ "net/http/pprof"
	"time"
)

// 这里app 的初始化工作
func init() {

}

// 获取命令行启动
// 1. app相关的配置文件初始化
// 2. 设置app参数
func (this *App)AppStart() {
	fmt.Println("App AppStart")
	//拉起宕机
	defer global.GrecoverToStd()
	appdata.InitAppData()
	this.ParseAppArgs() //获取命令行
	this.InitApp()      //解析完命令再启动对应程序
	this.AppWG.Add(1)
	// timersys.NewWheelTimer(time.Second,AppRun,nil)
	// go this.signalListen() //监听退出事件
	this.AppWG.Wait() // 等待app退出
	//CloseApp()   // 关闭app退出所有程序
}

// 程序启动获取命令行参数
func (this *App) ParseAppArgs() {
	var appid int
	flag.IntVar(&appid, "AppID", 0, "请输入app id")
	flag.Parse()
	appdata.AppID = int32(appid)
	for {
		appdata.NetConf = csvdata.GetNetworkconfPtr(appdata.AppID)
		if appdata.NetConf == nil {
			fmt.Println( "serverID 未找到")
		} else {
			break
		}
		time.Sleep(time.Second * 5)
	}
	appdata.WorldNetConf = csvdata.GetNetworkconfPtr(conf.SvJson.WorldServerId)
	if appdata.WorldNetConf == nil {
		panic("appdata.WorldNetConf  == nil")
	}
	fmt.Println( "ParseAppArgs success")
}

// 根据配置启动对应服务器
func (this *App)InitApp() {
	this.newLog()
	appdata.InitKind()
	appdata.SetAppFactory(NewAppFactory(appdata.GetAppKind()))
	//初始化协程池
	err := appdata.SetWorkPool()
	if err != nil {
		xlog.DebugLog("",err.Error())
		return
	}
	// 初始化app相关
	this.appBehavior = appdata.AppFactory.CreateAppBehavor()
	if this.appBehavior == nil {
		panic("appBehavior == nil ")
	}
	this.appBehavior.OnStart()
	// 执行对应
	this.SetAppOpen()
	// //读取控制台命令 测试的时候才用
	// this.AppWG.Add(1)
	// go this.ReadConsle()
	
	//查看程序状态 一个端口只能监听一次
	// go func() {
	// 	ip := fmt.Sprintf("0.0.0.0:%v", appdata.NetConf.Sys_stauts_port)
	// 	if err := http.ListenAndServe(ip, nil); err != nil {
	// 		fmt.Printf("start pprof failed on %s  erro = %v\n", ip,err)
	// 		os.Exit(1)
	// 	}
	// }()
}

func (this *App)newLog() {
	logInit := &xlog.LogInitModel{
		ServerName: appdata.NetConf.App_name,
		LogsPath:   appdata.PathModelPtr.LogsPath,
		Volatile:   conf.SvJson.LogConf,
	}
	if !xlog.NewXlog(logInit)   {
		panic("New Xlog errors ")
	}
	fmt.Println( "NewXlog success")
}

// app 逻辑参数根据服务器启动的参数创建对应的服务器工厂
func NewAppFactory(svKind model.AppKind) xengine.AppFactory {
	switch svKind {
	case model.APP_NONE:
		return nil
	case model.APP_Client:
		return new(appclient.ClientFactory)
	case model.APP_GATEWAY:
		return nil
	case model.APP_LoginServer: // 工厂
		return new(apploginsv.LoginServerFactory)
	case model.APP_GameServer:
		return new(appgatesv.GateSvFactory)
	case model.APP_MsgServer:
		return nil
	case model.APP_DataCenter:
		return new(datacenter.DataCenterFactory)
	default:
		return nil
	}
}
