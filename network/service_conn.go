/*
创建时间: 2020/7/5
作者: zjy
功能介绍:
服务器连接对象 主要与客户端端之间的连接进行区分
复用 连接部分方法
重写连接建立 通知不同的接口
重写读取数据 服务器内部通讯不加密
重写连接关闭
*/

package network

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/ants/v2"
	"net"
	"wengo/app/appdata"
	"wengo/csvdata"
	"wengo/xlog"
	"wengo/xutil/timeutil"
)

//复用连接部分方法
type ServiceConn struct {
	Connecter
	svNet ServiceNetEvent //服务器相关的网络事件
}


func newServiceConn(conn net.Conn, connId uint32, svnet ServiceNetEvent, netconf *csvdata.Networkconf, pool *ants.Pool,msgParser  *MsgParser) *ServiceConn {
	if svnet == nil {
		xlog.ErrorLogNoInScene("svnet 网络处理接口为空")
		return nil
	}
	if msgParser == nil {
		xlog.ErrorLogNoInScene("msgParser 数据解析对象为null")
		return nil
	}
	serviceConn := new(ServiceConn)
	// tcp的方法
	erro := serviceConn.initConnData(conn,connId,netconf,pool,msgParser)
	if erro != nil {
		return nil
	}
	serviceConn.svNet = svnet
	serviceConn.onCloseFun = serviceConn.OnClose  //这里绑定一个回调
	erro = serviceConn.svNet.OnServiceLink(serviceConn) // 通知其他模块已经连接 这里要用服务器的通知模块
	if erro != nil {
		xlog.ErrorLogNoInScene( "newServiceConn %v", erro)
		return nil
	}
	
	return serviceConn
}


// 服务器 连接实现的不一样的
func (this *ServiceConn) OnClose() {
	xlog.DebugLogNoInScene( "ServiceConn 关闭服务器连接")
	// 通知其他模块
	erro := this.svNet.OnServiceClose(this)
	if erro != nil {
		xlog.DebugLogNoInScene( "doClose OnServiceClose %v", erro)
	}
}


// 服务器读取 不用加密
func (this *ServiceConn) ReadMsg() error {
	// 查看连接每秒发多少
	data, err := this._msgParser.Read(this)
	if err != nil {
		return err
	}
	maincmd, subcmd, msgdat, erro := this._msgParser.ServiceUnpackOne(data)
	if erro != nil {
		return erro
	}
	//查看命令是否错误
	erro = this.Connecter.CMDIsErro(maincmd, subcmd)
	if erro != nil {
		return erro
	}
	//設置當前收包时间
	//客户端连接不一样
	this.lastRecTime = timeutil.GetCurrentTimeS()
	// TODO 优化对象创建
	msgData := NewMsgData(this, maincmd, subcmd, msgdat)
	if msgData == nil {
		return errors.New("NewMsgData is nil")
	}
	// 这里应该进入队列
	erro = this.svNet.OnServiceMsg(msgData)
	return erro
}

// 写单个消息
func (this *ServiceConn) WritePBMsg(maincmd, subcmd uint16, pb proto.Message) error {
	data, erro := this.GetPBByteArr(maincmd, subcmd, pb)
	if erro != nil {
		xlog.DebugLogNoInScene( "WritePBMsg erro : %v", erro)
		return erro
	}
	// 向写通道投递数据
	this.Write(data)
	return nil
}

// 将消息体构建为[]byte数组，最终要发出去的单包
func  (this *ServiceConn) GetPBByteArr(maincmd, subcmd uint16,pb proto.Message) (sendMsg []byte,erro error) {
	if pb != nil {
		sendMsg, erro = proto.Marshal(pb)
	}
	if erro != nil {
		xlog.ErrorLog(appdata.GetSecenName(), "GetPBByteArr %v", erro)
		return nil,erro
	}
	sendMsg, erro = this._msgParser.ServicePackOne(maincmd, subcmd, sendMsg)
	return
}

// 写单个消息
func (this *ServiceConn) WriteOneMsg(maincmd, subcmd uint16, msg []byte) error {
	data, erro := this.GetOneMsgByteArr(maincmd, subcmd, msg)
	if erro != nil {
		xlog.DebugLogNoInScene( "WriteOneMsgByConnID erro : %v", erro)
		return erro
	}
	// 向写通道投递数据
	this.Write(data)
	return nil
}

// 将消息体构建为[]byte数组，最终要发出去的单包
func  (this *ServiceConn) GetOneMsgByteArr(maincmd, subcmd uint16, msg []byte) ([]byte, error) {
	return this._msgParser.ServicePackOne(maincmd, subcmd, msg)
}
