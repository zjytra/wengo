// 1.封装tcp连接,改结构体只负责写的工作，
// 2.读的协程 交给另一对象处理,方便其他对象知道

package network

import (
	"errors"
	"fmt"
	"github.com/panjf2000/ants"
	"github.com/wengo/cmdconst"
	"github.com/wengo/csvdata"
	"github.com/wengo/xlog"
	"github.com/wengo/xutil/timeutil"
	"net"
	"sync"
)



type Connecter struct {
	conn        net.Conn
	connID      uint32 // 服务器创建连接时生成id
	sync.RWMutex             // 主要作用,防止向关闭后的通道中写入数据
	closeFlag   bool         // 检测关闭标志 这里不用锁使用原子数据更效率
	writeChan   chan []byte  // 写的通道，我服务器写的消息先写入通道再用连接传出去
	// recMutex    sync.RWMutex // 接收锁
	lastRecTime int64        // 最後一次收包时间
	netConf     *csvdata.Networkconf
	workPool    *ants.Pool // 协程池
	_msgParser  *MsgParser // 数据包解析对象
	onCloseFun  func() //关闭连接通知方法
}
//
// func newConnecter(conn net.Conn, connId uint32, netconf *csvdata.Networkconf, pool *ants.Pool,msgParser  *MsgParser) *Connecter {
// 	tcpConn := new(Connecter)
// 	erro := tcpConn.initConnData(conn,connId,netconf,pool)
// 	if erro != nil {
// 	    return nil
// 	}
// 	return tcpConn
// }

//公共初始化方法
func (this *Connecter) initConnData(conn net.Conn, connId uint32, netconf *csvdata.Networkconf, pool *ants.Pool,msgParser  *MsgParser) error {
	if  msgParser == nil {
		panic("数据解析对象为null")
		return errors.New("数据解析对象为null")
	}
	this.conn = conn
	this.connID = connId
	this.netConf = netconf
	this.writeChan = make(chan []byte, this.netConf.Write_cap_num)
	this.workPool = pool
	this._msgParser = msgParser  //使用外面传入的对象，所有连接共用一个对象，需要分开是才分开
	this.lastRecTime = timeutil.GetCurrentTimeS()
	erro := this.workPool.Submit(func() {
		this.writeChanData()
	}) // 写协程
	if erro != nil {
		xlog.ErrorLogNoInScene( "newConnecter Submit %v", erro)
		return erro
	}
	return nil
}


//获取连接id
func (this *Connecter) GetConnID() uint32 {
	return this.connID
}


func (this *Connecter) Destroy() {
	this.Lock()
	this.doDestroy()
	this.Unlock()
}


func (this *Connecter) doDestroy() {
	erro := this.conn.(*net.TCPConn).SetLinger(0)
	if erro != nil {
		xlog.ErrorLogNoInScene( "doDestroy 错误 %v", erro)
	}
	erro = this.conn.Close()
	if erro != nil {
		xlog.ErrorLogNoInScene( "doDestroy 关闭连接错误 %v", erro)
	}
	
	if !this.closeFlag {
		this.closeFlag = true
		close(this.writeChan)
	}
}

func (this *Connecter) Close() {
	this.Lock()
	// 已经关闭
	if this.closeFlag {
		this.Unlock()
		xlog.DebugLogNoInScene( "Connecter 当前连接已经关闭")
		return
	}

	this.closeFlag = true
	this.doWrite(nil)
	this.Unlock()
}

// b must not be modified by the others goroutines
func (this *Connecter) Write(b []byte) {
	this.Lock()
	defer this.Unlock()
	// 已经关闭
	if this.closeFlag || b == nil {
		xlog.DebugLogNoInScene( "当前连接状态 %v,", this.closeFlag)
		return
	}
	this.doWrite(b)
}

func (this *Connecter) doWrite(b []byte) {
	// 写的队列被撑满时
	if len(this.writeChan) == cap(this.writeChan) {
		xlog.DebugLogNoInScene( "close tcpconn: channel full")
		this.doDestroy() // 这里要主动断开避免阻塞 当前调用协程
		return
	}
	this.writeChan <- b
}

// 取通道的数据给连接
func (this *Connecter) writeChanData() {
	// 这里接收写的通道，没有数据会一直阻塞，直到通道关闭
	for b := range this.writeChan {
		if b == nil {
			break
		}
		_, err := this.conn.Write(b)
		if err != nil {
			break
		}
	}
	
	// 主动关闭,被动关闭关闭最后都要走这里
	erro := this.conn.Close()
	xlog.DebugLogNoInScene( "writeChanData this.conn 连接关闭")
	if erro != nil {
		xlog.ErrorLogNoInScene(  "关闭连接错误 %v", erro)
	}
	// 已经关闭
	this.Lock()
	this.closeFlag = true
	this.Unlock()
	this.onCloseFun() //主要通知其他模块子类实现确保网络关闭只调一次
}

func (this *Connecter) Read(b []byte) (int, error) {
	return this.conn.Read(b)
}

func (this *Connecter) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *Connecter) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}


// 一起写多个数据包
// 每个包的数据 由GetOneMsgByteArr构建
func (this *Connecter) WriteMsg(args ...[]byte) error {
	buf,erro := this._msgParser.MorePackageToOne(args...)
	if erro != nil {
		return erro
	}
	this.Write(buf)
	return nil
}

// 多个包构成成一个包
func (this *Connecter) ConnMorePackageToOne(args ...[]byte) ([]byte, error) {
	buf,erro := this._msgParser.MorePackageToOne(args...)
	if erro != nil {
		return nil,erro
	}
	return buf,nil
}

// 是否存活 没有存活会被提下线
func (this *Connecter) IsAlive() bool {
	currentTime := timeutil.GetCurrentTimeS()
	return (this.lastRecTime + int64(this.netConf.Checklink_s)) > currentTime
}

//判断连接是否关闭
func (this *Connecter)IsClose()bool  {
	this.RLock()
	isClose := this.closeFlag
	this.RUnlock()
	return isClose
}

//查看命令是否错误
func (this *Connecter)CMDIsErro(maincmd,subcmd uint16) error{
	if maincmd == cmdconst.Main_MIN_CMD || maincmd >= cmdconst.Main_MAX_CMD || subcmd == 0 {
		return errors.New(fmt.Sprintf("主命令%d错误 子命令%d",maincmd,subcmd))
	}
	return nil
}
