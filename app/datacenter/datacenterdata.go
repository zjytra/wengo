/*
创建时间: 2020/5/17
作者: zjy
功能介绍:
数据中心服务器相关数据
各个服务器信息管理
*/

package datacenter

import (
	"wengo/app/datacenter/dcmodel"
)

var (
	PdataCenter    *DataCenter            //数据中心
	PServerInfoMgr *dcmodel.ServerInfoMgr //服务器管理
	PaccountMgr    *dcmodel.AccountMgr    //账号管理
	PAccountPool   *DBAccountPool         //账号注册池
)

func NewData(dc *DataCenter) {
	PdataCenter = dc
	PServerInfoMgr = dcmodel.NewServerInfoMgr()
	PaccountMgr = dcmodel.NewAccountsMgr(PServerInfoMgr)
	PAccountPool = NewDBAccountPool(100)
}

func ClearAllServerData()  {
	PaccountMgr.ReleaseData()
	PServerInfoMgr.ClearAllServerData()
}
