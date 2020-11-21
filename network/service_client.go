package network

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/panjf2000/ants/v2"
	"wengo/csvdata"
	"wengo/model"
	"wengo/xlog"
	"net"
	"sync"
	"time"
)

var (
	colseErro = errors.New("已经关闭连接")
)

// 每個对象维持一个连接,连接其他服务器使用
type ServeiceClient struct {
	sync.Mutex
	Addr            string        // 服务器连接地址
	ConnectInterval time.Duration // 重连时间
	PendingWriteNum int
	AutoReconnect   bool
	svNet           ServiceNetEvent // 服务器事件观察者
	conn            *ServiceConn    // 连接对象
	wg              sync.WaitGroup
	closeFlag       *model.AtomicBool
	netConf         *csvdata.Networkconf
	workPool        *ants.Pool // 协程池
	msgParser       *MsgParser // 数据包解析对象 共用一个对象
}

// 创建tcp 客戶端
func NewServeiceClient(svNet ServiceNetEvent, netconf *csvdata.Networkconf, pool *ants.Pool) *ServeiceClient {
	if netconf == nil {
		xlog.WarningLogNoInScene("server conf is nil")
		return nil
	}
	if svNet == nil {
		xlog.WarningLogNoInScene("服务器 消息处理 svNet is nil")
		return nil
	}
	client := new(ServeiceClient)
	client.svNet = svNet
	client.netConf = netconf
	client.workPool = pool
	client.closeFlag = model.NewAtomicBool()
	client.msgParser = NewMsgParser(netconf.Msglen_size, netconf.Max_msglen, netconf.Msg_isencrypt)
	return client
}

func (client *ServeiceClient) Start() {
	client.init()
	client.wg.Add(1)
	go client.connect() // 开启一个连接
}

func (client *ServeiceClient) init() {
	if client.ConnectInterval <= 0 {
		client.ConnectInterval = 5 * time.Second
		xlog.DebugLogNoInScene("invalid ConnectInterval, reset to %v", client.ConnectInterval)
	}
	if client.PendingWriteNum <= 0 {
		client.PendingWriteNum = 100
		xlog.DebugLogNoInScene("invalid PendingWriteNum, reset to %v", client.PendingWriteNum)
	}
	client.closeFlag.SetFalse()
	client.AutoReconnect = true
	
	client.Addr = fmt.Sprintf("%s:%s", client.netConf.Out_addr, client.netConf.Out_prot)
	xlog.DebugLogNoInScene("client.Addr connet %v ", client.Addr)
}

func (client *ServeiceClient) dial() net.Conn {
	for {
		if client.closeFlag.IsTrue() {
			return nil
		}
		conn, err := net.Dial("tcp", client.Addr)
		if err == nil {
			return conn
		}
		
		xlog.DebugLogNoInScene("ServeiceClient dial to %v error: %v", client.Addr, err)
		time.Sleep(client.ConnectInterval)
		continue
	}
}

func (client *ServeiceClient) connect() {
	defer client.wg.Done()
redial:
	client.doconnect() // 执行连接及读
	// 没有关闭才进行重连
	if client.closeFlag.IsFalse() && client.AutoReconnect {
		time.Sleep(client.ConnectInterval)
		goto redial
	}
}

func (client *ServeiceClient) doconnect() {
	conn := client.dial()
	if conn == nil {
		return
	}
	if !client.setConn(conn) {
		return
	}
}

// 添加链接信息
func (client *ServeiceClient) setConn(conn net.Conn) bool {
	if client.closeFlag.IsTrue() {
		conn.Close()
		return false
	}
	svconn := newServiceConn(conn, connIDs.GetId(), client.svNet, client.netConf, client.workPool, client.msgParser)
	client.conn = svconn
	xlog.DebugLogNoInScene("连接远程 %v 地址成功", conn.RemoteAddr())
	// 连接成功,将阻塞读取数据
	client.ReceiveData(client.conn)
	xlog.DebugLogNoInScene("TCPClient结束读取")
	return true
}

// 连接中读取数据
func (client *ServeiceClient) ReceiveData(svcon *ServiceConn) {
	for {
		err := svcon.ReadMsg()
		if err != nil { // 这里读到错误消息,关闭
			xlog.WarningLogNoInScene("read message: ", err)
			break // 关闭连接
		}
	}
	// cleanup
	client.closeConn(svcon)
}

func (client *ServeiceClient) closeConn(conn *ServiceConn) {
	conn.Close()
	xlog.DebugLogNoInScene("关闭远程服务器连接")
}

// 写单个消息
func (client *ServeiceClient) WriteOneMsg(maincmd, subcmd uint16, msg []byte) error {
	if client.closeFlag.IsTrue() {
		return colseErro
	}
	if client.conn == nil {
		return errors.New(fmt.Sprintf("ServeiceClient WriteOneMsg未建立连接 %v", client.Addr))
	}
	return client.conn.WriteOneMsg(maincmd, subcmd, msg)
}

// 将消息体构建为[]byte数组，最终要发出去的单包
func (client *ServeiceClient) GetOneMsgByteArr(maincmd, subcmd uint16, msg []byte) ([]byte, error) {
	if client.closeFlag.IsTrue() {
		return nil, colseErro
	}
	if client.conn == nil {
		return nil,errors.New(fmt.Sprintf("GetOneMsgByteArr未建立连接 %v", client.Addr))
	}
	return client.conn.GetOneMsgByteArr(maincmd, subcmd, msg)
}

// 写单个消息pb实现
func (client *ServeiceClient) WritePBMsg(maincmd, subcmd uint16, pb proto.Message) error {
	if client.closeFlag.IsTrue() {
		return colseErro
	}
	if client.conn == nil {
		return errors.New(fmt.Sprintf("WritePBMsg未建立连接 %v", client.Addr))
	}
	return client.conn.WritePBMsg(maincmd, subcmd, pb)
}

// 将消息体构建为[]byte数组，最终要发出去的单包 pb实现
func (client *ServeiceClient) GetPBByteArr(maincmd, subcmd uint16, pb proto.Message) ([]byte, error) {
	if client.closeFlag.IsTrue() {
		return nil, colseErro
	}
	if client.conn == nil {
		return nil,errors.New(fmt.Sprintf("GetOneMsgByteArr未建立连接 %v", client.Addr))
	}
	return client.conn.GetPBByteArr(maincmd, subcmd, pb)
}

// 一起写多个数据包
// 每个包的数据 由GetOneMsgByteArr构建
func (client *ServeiceClient) WriteMsg(args ...[]byte) error {
	if client.closeFlag.IsTrue() {
		return colseErro
	}
	return client.conn.WriteMsg(args...)
}

// 是否存活 没有存活会被提下线
func (client *ServeiceClient) IsAlive() bool {
	return client.conn.IsAlive()
}

// 是否关闭
func (client *ServeiceClient) IsClose() bool {
	return client.conn.IsClose()
}

// 获取连接对象id
func (client *ServeiceClient) GetConnID() uint32 {
	return client.conn.GetConnID()
}

// 关闭服务
func (client *ServeiceClient) Close() {
	client.closeFlag.SetTrue()
	client.conn.Close()
	client.wg.Wait()
}

// 关闭连接
func (client *ServeiceClient) DoCloseConn() {
	client.conn.Close()
}


func (client *ServeiceClient) GetServiceConn() *ServiceConn{
	return client.conn
}