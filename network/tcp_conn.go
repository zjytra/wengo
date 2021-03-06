/*
创建时间: 2020/7/14
作者: zjy
功能介绍:
客户端连接对象
//这里需要验证
1.自己主动关闭连接流程 向写的通道里写空,在写的线程关闭连接
2.远端关闭，被动关闭
3.远端直接关机
*/

package network

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/ants/v2"
	"net"
	"github.com/zjytra/wengo/app/appdata"
	"github.com/zjytra/wengo/csvdata"
	"github.com/zjytra/wengo/model"
	"github.com/zjytra/wengo/xlog"
	"github.com/zjytra/wengo/xutil/timeutil"
)


type ConnSet map[*TcpConn]struct{}
//对端连接服务器对象
type TcpConn struct {
	Connecter //继承 连接对象
	netObserver NetWorkObserver
	recMsgPs    int          // 每秒收包个数
}

func newTcpConn(conn net.Conn, connId uint32,  netOb NetWorkObserver, netconf *csvdata.Networkconf, pool *ants.Pool,msgParser  *MsgParser) *TcpConn {
	if netOb == nil {
		xlog.ErrorLogNoInScene("netOb 网络处理接口为空")
		return nil
	}
	if msgParser == nil {
		xlog.ErrorLogNoInScene("msgParser 数据解析对象为null")
		return nil
	}
	clntConn := new(TcpConn)
	// tcp的方法
	erro := clntConn.initConnData(conn,connId,netconf,pool,msgParser)
	if erro != nil {
		return nil
	}
	clntConn.netObserver = netOb
	clntConn.onCloseFun = clntConn.OnClose  //这里绑定一个回调
	erro = clntConn.netObserver.OnNetWorkConnect(clntConn) // 通知其他模块已经连接 这里是客户端接口
	if erro != nil {
		xlog.ErrorLogNoInScene("newClientConn %v", erro)
		return nil
	}
	return clntConn
}

// 用解析对象读,单协程调用
func (this *TcpConn) ReadMsg() error {

	data, err := this._msgParser.Read(this)
	if err != nil {
		return err
	}
	// 查看连接每秒发多少
	if  !this.checkCanRead() {
		return errors.New("每秒超过最大包数")
	}
	
	maincmd, subcmd, msgdat, erro := this._msgParser.UnpackOne(data)
	if erro != nil {
		return erro
	}
	//查看命令是否错误
	erro = this.Connecter.CMDIsErro(maincmd, subcmd)
	if erro != nil {
		return erro
	}
	// TODO 优化对象创建
	msgData := NewMsgData(this, maincmd, subcmd, msgdat)
	if msgData == nil {
		return errors.New("NewMsgData is nil")
	}
	// 这里应该进入队列
	erro = this.netObserver.OnNetWorkRead(msgData)
	return erro
}


// 查看是否可以读取
func (this *TcpConn) checkCanRead() bool {
	
	if !this.isCheck() { // 查看是否检查每秒包量
		this.lastRecTime = timeutil.GetCurrentTimeS() //这里要设置下时间
		return  true
	}
	currentTime := timeutil.GetCurrentTimeS()
	// 在同一秒内
	if this.lastRecTime == currentTime {
		// 收包的数量超过每秒最大限制数量
		if this.recMsgPs >= this.netConf.Max_rec_msg_ps {
			xlog.ErrorLogNoInScene("每秒收包 %v 个超过 %v个",this.recMsgPs,this.netConf.Max_rec_msg_ps)
			return false
		}
		this.recMsgPs ++
		return true
	}
	// 过了一秒重置变量
	this.recMsgPs = 0
	this.lastRecTime = currentTime
	return true
}

//查看是否检查每秒包量
//服务器内部通信不需要检测
func (this *TcpConn) isCheck() bool {
	return this.netConf.App_kind != model.APP_DataCenter && this.netConf.App_kind != model.APP_GameServer
}


// 写单个消息
func (this *TcpConn) WritePBMsg(maincmd, subcmd uint16, pb proto.Message) error {
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
func  (this *TcpConn) GetPBByteArr(maincmd, subcmd uint16,pb proto.Message) (sendMsg []byte, erro error) {
	if pb != nil {
		sendMsg, erro = proto.Marshal(pb)
	}
	if erro != nil {
		xlog.ErrorLog(appdata.GetSecenName(), "GetPBByteArr %v", erro)
		return nil,erro
	}
	sendMsg, erro = this._msgParser.PackOne(maincmd, subcmd, sendMsg)
	return
}

// 写单个消息
func (this *TcpConn) WriteOneMsg(maincmd, subcmd uint16, msg []byte) error {
	data, erro := this.GetOneMsgByteArr(maincmd, subcmd, msg)
	if erro != nil {
		xlog.DebugLogNoInScene( "WriteOneMsg erro : %v", erro)
		return erro
	}
	// 向写通道投递数据
	this.Write(data)
	return nil
}


// 将消息体构建为[]byte数组，最终要发出去的单包
func  (this *TcpConn) GetOneMsgByteArr(maincmd, subcmd uint16, msg []byte) ([]byte, error) {
	return this._msgParser.PackOne(maincmd, subcmd, msg)
}


// 服务器 连接实现的不一样的
func (this *TcpConn) OnClose() {
	xlog.DebugLogNoInScene( "TcpConn OnClose")
	// 通知其他模块
	erro := this.netObserver.OnNetWorkClose(this)
	if erro != nil {
		xlog.DebugLogNoInScene( " TcpConn Close %v", erro)
	}
}