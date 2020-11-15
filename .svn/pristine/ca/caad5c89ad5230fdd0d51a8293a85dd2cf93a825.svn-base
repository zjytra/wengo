/*
创建时间: 2020/7/5
作者: zjy
功能介绍:
服务器拆包封包不加密
*/

package network

import (
	"errors"
)

// param | 主命令 | 字命令 | datalen | data |
func (p *MsgParser) ServiceUnpackOne(readb []byte) (maincmd uint16, subcmd uint16, msg []byte, err error) {
	msglen := len(readb) // 包头的长度已经被解掉
	if readb == nil || msglen == 0 {
		return
	}
	// 服务器内部之间通讯不加密
	var start, end int
	start, end = GetNextIndex(end, maincmd)
	maincmd = ByteOrder.Uint16(readb[start:end]) // 解析主命令
	start, end = GetNextIndex(end, subcmd)
	subcmd = ByteOrder.Uint16(readb[start:end]) // 解析字命令
	var datalen uint32
	start, end = GetNextIndex(end, datalen)
	datalen = ByteOrder.Uint32(readb[start:end]) // 解析长度
	if datalen > 0 {   //当消息只有命令没有消息的就不考虑数据了
		msg = make([]byte, datalen)
		copy(msg, readb[end:])
	}
	return
}

// 打单包
func (p *MsgParser) ServicePackOne(maincmd, subcmd uint16, msg []byte) ([]byte, error) {
	
	var msgLen,datalen uint32
	if msg != nil { //空消息判断
		datalen = uint32(len(msg))
	}
	msgLen = p.minMsgLen + datalen
	
	// check len
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}
	//服务器之间通讯不加密
	writeBuf := make([]byte, msgLen)
	var start, end int
	switch p.msgLenSize { // 写长度
	case 2: // uint16
		end = 2
		ByteOrder.PutUint16(writeBuf[start:end], uint16(msgLen))
	case 4: // uint32
		end = 4
		ByteOrder.PutUint32(writeBuf[start:end], msgLen)
	}
	start, end = GetNextIndex(end, maincmd)
	ByteOrder.PutUint16(writeBuf[start:end], maincmd) // 主命令
	start, end = GetNextIndex(end, subcmd)
	ByteOrder.PutUint16(writeBuf[start:end], subcmd) // 字命令
	
	if datalen > 0{ //有数据才往里面压
		start, end = GetNextIndex(end, datalen)
		ByteOrder.PutUint32(writeBuf[start:end], datalen) // protobuf数据的长度
		copy(writeBuf[end:], msg)
	}
	return writeBuf, nil
}
