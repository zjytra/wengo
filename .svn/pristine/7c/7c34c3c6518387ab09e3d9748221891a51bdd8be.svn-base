/*
创建时间: 2020/2/21
作者: zjy
功能介绍:

*/

package network

// // Tcp自定义包格式
type TcpMsgHead struct {
	// PackageLen uint32 // 数据包总长度
	// MKey    uint8 // 加密key对后面的数据进行异或 这里只有加密了才增加1字节
	MainCmd uint16 // 数据包主命令
	SubCmd  uint16 // 数据包子命令
	Datalen uint32 // 这里是protobuf数据长度
	// Msgdata       []byte // protobuf数据
}

// 投递的解析数据
type MsgData struct {
	Conn    Conner
	MainCmd uint16 // 数据包主命令
	SubCmd  uint16 // 数据包子命令
	Msgdata []byte // protobuf数据
}

func NewMsgData(conn Conner, maincmd, subcmd uint16, msgdata []byte) *MsgData {
	msg := new(MsgData)
	msg.Conn = conn
	msg.MainCmd = maincmd
	msg.SubCmd = subcmd
	msg.Msgdata = msgdata
	return msg
}

// 群发消息
type GroupMessage struct {
	ConnIDs []uint32
	Msgdata []byte
}
