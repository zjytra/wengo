package network

import (
	"encoding/binary"
	"errors"
	"wengo/xlog"
	"wengo/xutil"
	"io"
	"math"
)

// 包格式
// --------------
//  | 包总长 | 加密key| 主命令 | 字命令 | datalen | data |
// --------------
type MsgParser struct {
	msgLenSize uint32 // 包长度字节大小
	msgheadLen uint32 // 包头默认需要最小字节数
	minMsgLen  uint32 // 最小长度
	maxMsgLen  uint32 // 最大长度
	isEncrypt  bool   // 是否加密
}


var (
	ByteOrder binary.ByteOrder // 设置网址字节序
)
func init() {
	// 设置大端序列或小端序列
	if xutil.IsLittleEndian() {
		ByteOrder = binary.LittleEndian
	} else {
		ByteOrder = binary.BigEndian
	}
}

func NewMsgParser(msgLenSize uint8, maxMsgLen uint32, isencrypt bool) *MsgParser {
	p := new(MsgParser)
	p.msgLenSize = uint32(msgLenSize)
	if p.msgLenSize == 0 {
		p.msgLenSize = 2
	}
	var head TcpMsgHead
	headLen :=  uint32(binary.Size(head)) //包头需要长度
	if isencrypt {
		headLen += 1 // 加密才需要加上1byte长度
	}
	p.msgheadLen = headLen
	p.setMsgLen(maxMsgLen)
	p.isEncrypt = isencrypt
	return p
}

// 设置消息长度
func (p *MsgParser) setMsgLen(maxMsgLen uint32) {
	// 数据包最小长度 | 包总长大小 | 主命令大小 | 字命令大小 | datalen大小 |
	p.minMsgLen = uint32(p.msgLenSize) + p.msgheadLen
	if maxMsgLen != 0 {
		p.maxMsgLen = maxMsgLen
	}
	var max uint32
	switch p.msgLenSize {
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if p.minMsgLen > max { // 不能超过设置的
		p.minMsgLen = max
	}
	if p.maxMsgLen > max {
		p.maxMsgLen = max
	}
}

// goroutine safe
func (p *MsgParser) Read(conn Conner) ([]byte, error) {
	// 根据长度字节大小解析第一个长度
	msgLenBuf := make([]byte, p.msgLenSize) // TODO 可以优化接收消息的buf
	n, err := conn.Read(msgLenBuf)
	if err != nil {
		return nil, err
	}
	if p.msgLenSize != uint32(n) {
		return nil, errors.New("包头长度出错")
	}
	// parse len
	msgLen := p.byteArrToMsgLen(msgLenBuf)
	// check len
	if msgLen > p.maxMsgLen {
		return nil, errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil, errors.New("message too short")
	}
	// data 读取
	// 已经提取了消息长度的字节 剩余 加密key| 主命令 | 字命令 | datalen | data | 这里扣除长度的字节
	datalen := msgLen - uint32(p.msgLenSize)
	msgdata := make([]byte, datalen) // TODO 可以优化接收消息的buf
	// conn.Read(msgdata)
	readlen, err := io.ReadFull(conn,msgdata)
	if err != nil {
		return nil, err
	}
	unreadlen := xutil.MaxUint32(0, datalen-uint32(readlen))
	if unreadlen == 0 {
		return msgdata, nil
	}
	// // 出现断包再读一次  不处理断包避免不明连接占用连接
	// unreadData := msgdata[readlen:datalen:datalen]  //获取未填满的容器
	// rereadlen, err := tcpconn.Read(unreadData) // 未读数据
	// if err != nil {
	// 	return nil, err
	// }
	// //如果还不能填满长度
	// if uint32(rereadlen) != unreadlen {
	// 	return nil,errors.New("断包处理失败")
	// }
	// return msgdata, nil
	return nil, errors.New("读取数据出现断包")
}

// 根据字节数组解析长度
func (p *MsgParser) byteArrToMsgLen(bArr []byte) uint32 {
	switch p.msgLenSize {
	case 2: // uint16 解析
		return uint32(ByteArrToUint16(bArr))
	case 4: // uint32 解析
		return ByteArrToUint32(bArr)
	}
	xlog.DebugLogNoInScene("not set msg size ")
	return 0 // 无效
}

func ByteArrToUint32(bArr []byte) uint32 {
	return ByteOrder.Uint32(bArr)
}

func ByteArrToUint16(bArr []byte) uint16 {
	return ByteOrder.Uint16(bArr)
}

// param |加密key| 主命令 | 字命令 | datalen | data |
func (p *MsgParser) UnpackOne(readb []byte) (maincmd uint16, subcmd uint16, msg []byte, err error) {
	msglen := len(readb) // 包头的长度已经被解掉
	if readb == nil || msglen == 0 {
		return
	}
	
	var parseBuf []byte
	if p.isEncrypt {  //加密处理
		parseBuf = make([]byte, msglen-1) // 扣除key  取后面的数据
		key := readb[0]                    // 解析加密key
		copy(parseBuf, readb[1:])          // 将后面的命令放到一个需要解析的数组中
		setEncrypt(parseBuf, key)          // 将数据 解密  用key解密
	}else{
		parseBuf = readb         // 将后面的命令放到一个需要解析的数组中
	}
	
	var start, end int
	start, end = GetNextIndex(end, maincmd)
	maincmd = ByteOrder.Uint16(parseBuf[start:end]) // 解析主命令
	start, end = GetNextIndex(end, subcmd)
	subcmd = ByteOrder.Uint16(parseBuf[start:end]) // 解析字命令
	var datalen uint32
	start, end = GetNextIndex(end, datalen)
	datalen = ByteOrder.Uint32(parseBuf[start:end]) // 解析长度
	if  datalen > 0{ //当消息只有命令没有消息的就不考虑数据了
		msg = make([]byte, datalen)
		copy(msg, parseBuf[end:])
	}
	return
}

// 打单包
func (p *MsgParser) PackOne(maincmd, subcmd uint16, msg []byte) ([]byte, error) {

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
	// 加密key
	var bkey byte
	if p.isEncrypt { // 只有加密了才生成
		bkey = byte(xutil.RandInterval(1, 254))
		msgLen += 1  //需要加一个加密字节的长度
	}
	writeBuf := make([]byte, msgLen)
	var start, end, keylastIndex int
	switch p.msgLenSize { // 写长度
	case 2: // uint16
		end = 2
		ByteOrder.PutUint16(writeBuf[start:end], uint16(msgLen))
	case 4: // uint32
		end = 4
		ByteOrder.PutUint32(writeBuf[start:end], msgLen)
	}
	if p.isEncrypt { // 如果加密了就写加密key
		start, end = GetNextIndex(end, bkey)
		writeBuf[start] = bkey // 只有一个字节
		keylastIndex = end     // 记录加密key的下一个位置 方便取key后面的数据
	}
	start, end = GetNextIndex(end, maincmd)
	ByteOrder.PutUint16(writeBuf[start:end], maincmd) // 主命令
	start, end = GetNextIndex(end, subcmd)
	ByteOrder.PutUint16(writeBuf[start:end], subcmd) // 字命令
	
	if datalen > 0 {
		start, end = GetNextIndex(end, datalen)
		ByteOrder.PutUint32(writeBuf[start:end], datalen) // protobuf数据的长度
		copy(writeBuf[end:], msg)  //最后写proto数据
	}

	if p.isEncrypt {
		setEncrypt(writeBuf[keylastIndex:], bkey) // 将key后面的数据加密
	}
	return writeBuf, nil
}

// 将多个包合并成一个
func (p *MsgParser) MorePackageToOne(args ...[]byte) ( []byte,error ){
	if args == nil {
		return  nil,errors.New("MorePackageToOne args is nil")
	}
	var msgLen uint32
	// 计算消息长度
	for i := 0; i < len(args); i++ {
		if args[i] == nil {
			continue
		}
		msgLen += uint32(len(args[i]))
	}
	// check len
	if msgLen > p.maxMsgLen {
		return nil,errors.New("message too long")
	} else if msgLen < p.minMsgLen {
		return nil,errors.New("message too short")
	}
	// 构建所有数据包
	buf := make([]byte,msgLen)
	l := 0
	for i := 0; i < len(args); i++ {
		if args[i] == nil {
			continue
		}
		copy(buf[l:], args[i])
		l += len(args[i])
	}
	return buf,nil
}

//根据数据类型从字节数组中取数据下标
func GetNextIndex(end int, data interface{}) (head, tail int) {
	head = end                            // 尾部变成头部
	tail = head + xutil.IntDataSize(data) // 新的尾部=头加上数据的长度
	return
}


// 加密函数
func setEncrypt(data []byte, bkey byte) {
	datalen := xutil.MinInt(len(data),30) //加密优化，避免加密太长
	for i := 0; i < datalen; i++ {
		data[i] = data[i] ^ bkey
	}
}


// // param |加密key| 主命令 | 字命令 | datalen | data |
// func (p *MsgParser) UnpackOne(readb []byte) (maincmd uint16, subcmd uint16, msg []byte, err error) {
// 	msglen := len(readb)
// 	if readb == nil || msglen == 0 {
// 		return
// 	}
// 	parseBuf := make([]byte, msglen-1) // 存儲key 后面的数据
// 	copy(parseBuf, readb[1:])          // 将后面的命令放到一个桶中
// 	key := readb[0]                    // 加密key
// 	if key > 0 { //由加密key才解密
// 		setEncrypt(parseBuf, key)          // 将数据 解密 第一个是key 用key解密
// 	}
//
// 	// 保存数据到缓冲区
// 	reader := p.tcpconn.reader
// 	reader.Reset()
// 	reader.Write(parseBuf)
// 	err = binary.Read(reader, ByteOrder, &maincmd) // 解析主命令
// 	if err != nil {
// 		return
// 	}
// 	err = binary.Read(reader, ByteOrder, &subcmd) // 子命令
// 	if err != nil {
// 		return
// 	}
// 	var datalen uint32
// 	err = binary.Read(reader, ByteOrder, &datalen)
// 	if err != nil {
// 		return
// 	}
// 	msg = make([]byte, datalen)
// 	err = binary.Read(reader, ByteOrder, &msg)
// 	return
// }
//
// // 打单包
// func (p *MsgParser) PackOne(maincmd, subcmd uint16, msg []byte) ([]byte, error) {
// 	datalen := uint32(len(msg))
// 	// 加密key
// 	bkey := byte(xutil.RandInterval(1, 254))
// 	// head := TcpMsgHead{
// 	// 	MKey: bkey,
// 	// 	MainCmd: maincmd,
// 	// 	SubCmd:  subcmd,
// 	// 	Datalen: datalen,
// 	// }
// 	msgLen := p.minMsgLen + datalen
// 	// check len
// 	if msgLen > p.maxMsgLen {
// 		return nil, errors.New("message too long")
// 	} else if msgLen < p.minMsgLen {
// 		return nil, errors.New("message too short")
// 	}
// 	writer := p.tcpconn.writer
// 	writer.Reset()
// 	// 为了加密一个一个写
// 	erro := binary.Write(writer, ByteOrder, maincmd) // 写头
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	erro = binary.Write(writer, ByteOrder, subcmd) // 写头
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	erro = binary.Write(writer, ByteOrder, datalen) // 写头
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	erro = binary.Write(writer, ByteOrder, msg) // 写数据
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	databuf := make([]byte, writer.Len())
// 	copy(databuf, writer.Bytes())
// 	setEncrypt(databuf, bkey)                    // 将数据 加密
// 	writer.Reset()                               // 重置buf变量
// 	p.WriteLen(writer, msgLen)                   // 写长度  长度与key 不加密
// 	erro = binary.Write(writer, ByteOrder, bkey) // 写加密数据
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	erro = binary.Write(writer, ByteOrder, databuf) // 最后写数据
// 	if erro != nil {
// 		return nil, erro
// 	}
// 	return writer.Bytes(), erro
// }
// func (p *MsgParser) WriteLen(writer *bytes.Buffer, msgLen uint32) {
// 	switch p.msgLenSize {
// 	case 2: // uint16
// 		binary.Write(writer, ByteOrder, uint16(msgLen))
// 	case 4: // uint32
// 		binary.Write(writer, ByteOrder, uint32(msgLen))
// 	}
// }
