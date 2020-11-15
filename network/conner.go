package network

import (
	"github.com/golang/protobuf/proto"
	"net"
)

//连接接口
type Conner interface {
	Read(b []byte) (int, error)
	ReadMsg() (error)
	//一次性发送多个消息
	WriteMsg(args ...[]byte) error
	WriteOneMsg(maincmd, subcmd uint16, msg []byte) error
	//将命令和内容打成一个包
	GetOneMsgByteArr(maincmd, subcmd uint16, msg []byte) ([]byte, error)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	//获取连接id
	GetConnID() uint32
	//多个消息转换为一个
	ConnMorePackageToOne(args ...[]byte) ([]byte, error)
	//发送protobuf消息
	WritePBMsg(maincmd, subcmd uint16, pb proto.Message) error
	//发送protobuf消息
	GetPBByteArr(maincmd, subcmd uint16,pb proto.Message)(sendMsg []byte, erro error)
}

//一般是客户端的网络事件向外传递
type NetWorkObserver interface {
	OnNetWorkConnect(conn Conner) error
	OnNetWorkClose(conn Conner) error
	OnNetWorkRead(msgdata *MsgData) error
}

//其他服务器触发的相关事件
//提供此接口的主要目的是区分服务器与服务器之间交互 与 服务器与客户端的交互
type ServiceNetEvent interface {
	OnServiceLink(conn Conner) error
	OnServiceClose(conn Conner) error
	OnServiceMsg(msgdata *MsgData) error
}

// 消息接口
type HandlerNetWorkMsg func(conn Conner,maincmd,subcmd uint16,msgdata []byte) error