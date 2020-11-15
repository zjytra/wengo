/*
创建时间: 2020/09/2020/9/11
作者: Administrator
功能介绍:

*/
package dcmodel

import (
	"github.com/wengo/model"
	"github.com/wengo/xlog"
)

type ServerInfoMgr struct {
	serverInfoSMap map[int32]*SeverInfoModel //使用appid作为key 查找得
	kindConnIDSMap map[int32][]uint32        //<kind,<conids>> 这里是为了方便给同一类服务器发送信息时使用
	connIDMapAppID map[uint32]int32          //<connId,appid>断线时方便查找哪个appid断线
}

func NewServerInfoMgr() *ServerInfoMgr {
	pMgr := new(ServerInfoMgr)
	pMgr.serverInfoSMap = make(map[int32]*SeverInfoModel)
	pMgr.kindConnIDSMap = make(map[int32][]uint32)
	pMgr.connIDMapAppID = make(map[uint32]int32)
	return pMgr
}


//添加服务器信息
func (this *ServerInfoMgr)AddServerInfo(info *SeverInfoModel) bool {
	if info == nil {
		xlog.WarningLogNoInScene("AddServerInfo erro")
		return false
	}
	
	//服务器已经注册
	if _,ok := this.serverInfoSMap[info.AppId];ok{
		xlog.WarningLogNoInScene("serverId%v 已经在中心服注册",info.AppId)
		return false
	}
	
	_,hasLink := this.connIDMapAppID[info.ConnID]
	if hasLink {
		xlog.WarningLogNoInScene("ConnID %v 已经在中心服注册",info.ConnID)
		return false
	}
	
	// 将同一类的服务器放在一起
	if connIDs, ok := this.kindConnIDSMap[info.AppKind]; ok { //已经存在同类
		var isFind bool
		for _, conid := range connIDs {
			if conid == info.ConnID {
				isFind = true
				//已经放在列表里面
				xlog.DebugLogNoInScene("serverId %v 已经在类型列表中",info.AppId)
				break
			}
		}
		//没找到才添加
		if !isFind {
			connIDs = append(connIDs, info.ConnID )
			this.kindConnIDSMap[info.AppKind] = connIDs
			xlog.DebugLogNoInScene("info.AppKind %v 连接列表%v", model.ToKindString(info.AppKind),connIDs)
		}
		
	} else {
		var connIDs []uint32
		connIDs = append(connIDs, info.ConnID )
		this.kindConnIDSMap[info.AppKind] = connIDs
		xlog.DebugLogNoInScene("info.AppKind %v 连接列表%v", model.ToKindString(info.AppKind),connIDs)
	}
	this.connIDMapAppID[info.ConnID] =  info.AppId
	this.serverInfoSMap[info.AppId] = info
	xlog.DebugLogNoInScene("AppId %v  ConnID = %v注册成功",info.AppId,info.ConnID)
	return  true
}


// 移除某个服务器连接
func (this *ServerInfoMgr)RemoveServerInfo(connID uint32) bool {
	
	appID,hasLink := this.connIDMapAppID[connID]
	if hasLink {
		delete(this.connIDMapAppID,connID)
		xlog.DebugLogNoInScene("移除 appid = %v 的服务器",appID)
	}else {
		xlog.ErrorLogNoInScene("未找到 appid = %v 的服务器",appID)
	}
	//服务器未注册
	pServerInfo,ok := this.serverInfoSMap[appID]
	if !ok {
		xlog.ErrorLogNoInScene("未找到 appID = %v 的服务器",appID)
		return false
	}
	if pServerInfo.ConnID != connID {
		xlog.ErrorLogNoInScene(" appID = %v 与连接 connID = %v 未匹配",appID,connID)
	}
	var appid  = pServerInfo.AppId  //服务器id
	var appkind = pServerInfo.AppKind //服务器类型
	pServerInfo = nil //移除变量
	
	var isFind bool
	if 	connIDs,ok := this.kindConnIDSMap[appkind]; ok {
		var svlen = len(connIDs)
		for i := 0 ; i < svlen ; i++ {
			if  connIDs[i] == connID {
				isFind = true
				xlog.DebugLogNoInScene("在同类中移除连接 connID = %v 的服务器",connID)
				connIDs = append(connIDs[:i],connIDs[i+1:]...) //移除找到连接
				this.kindConnIDSMap[appkind] = connIDs
				break
			}
		}
	}
	if !isFind {
		xlog.WarningLogNoInScene("在同类中服务器中未找到 connID = %v 的服务器",connID)
	}
	delete(this.serverInfoSMap,appid)
	xlog.DebugLogNoInScene("移除连接 appid =%v connID = %v 的服务器",appid ,connID)
	return true
}


//根据连接ID获取服务器类型
func (this *ServerInfoMgr)GetServerKindByConnID(connID uint32) int32 {
	serverInfo := this.GetServerInfoByConnID(connID)
	if serverInfo == nil {
		return 0
	}
	return serverInfo.AppKind
}

//根据连接ID获取服务器信息
func (this *ServerInfoMgr)GetServerInfoByConnID(connID uint32) *SeverInfoModel {
	appID := this.GetAppIDByConnID(connID)
	if appID == 0  {
		return nil
	}
	serverInfo := this.GetServerInfoByAppID(appID)
	if serverInfo == nil {
		return nil
	}
	if serverInfo.ConnID != connID  {
		xlog.ErrorLogNoInScene("appid =%v 与连接 connID = %v 未匹配上",appID ,connID)
		return nil
	}
	return serverInfo
}

//根据appId获取连接信息
func (this *ServerInfoMgr)GetAppIDByConnID(conID uint32) int32 {
	appID,ok := this.connIDMapAppID[conID]
	if !ok {
		return 0
	}
	
	return appID
}

//根据服务器id获取服务器类型
func (this *ServerInfoMgr)GetAppKindByAppID(appID int32) int32 {
	serverInfo:= this.GetServerInfoByAppID(appID)
	if serverInfo == nil {
		return 0
	}
	return serverInfo.AppKind
}

//根据连接ID获取服务器信息
func (this *ServerInfoMgr)GetServerInfoByAppID(appID int32) *SeverInfoModel {
	serverInfo,ok := this.serverInfoSMap[appID]
	if !ok {
		xlog.ErrorLogNoInScene("appid=%v服务器信息未找到",appID)
		return nil
	}
	return serverInfo
}



//根据服务器类型id获取同一类型的服务器连接
//@return 返回连接切片
func (this *ServerInfoMgr)GetServerConnIDsByAppKind(appkind int32) []uint32 {
	serverIDs,ok := this.kindConnIDSMap[appkind]
	if !ok {
		return nil
	}
	return serverIDs
}

//清除所有的数据
func (this *ServerInfoMgr)ClearAllServerData()  {
	this.serverInfoSMap = nil
	for _,v := range this.kindConnIDSMap {
		v = v[:0:0] //清空切片
	}
	this.kindConnIDSMap = nil
	this.connIDMapAppID = nil
}

//获取负载最小的网关
func (this *ServerInfoMgr)GetGateWayServerInfo()*SeverInfoModel {
	var tem *SeverInfoModel = nil
	//获取负载最小的
	for _,ServerInfo := range this.serverInfoSMap {
		if ServerInfo.AppKind != model.APP_GATEWAY {
			continue
		}
		tem = ServerInfo
	}
	return tem
}