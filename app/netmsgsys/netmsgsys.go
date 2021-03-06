/*
创建时间: 2020/4/25
作者: zjy
功能介绍:
注册消息处理消息的
*/

package netmsgsys

import (
	"errors"
	"fmt"
	"github.com/zjytra/wengo/app/appdata"
	"github.com/zjytra/wengo/dispatch"
	"github.com/zjytra/wengo/network"
	"github.com/zjytra/wengo/xlog"
	"github.com/zjytra/wengo/xutil"
	"github.com/zjytra/wengo/xutil/timeutil"
)

type  NetMsgSys struct {
	msgHandler map[uint32]dispatch.HandleMsg
}

func NewMsgHandler() *NetMsgSys {
	hd := new(NetMsgSys)
	hd.msgHandler = make(map[uint32]dispatch.HandleMsg)
	return hd
}

func  (this *NetMsgSys)RegisterMsgHandle(MainCmd,SubCmd uint16,handler dispatch.HandleMsg) {
	cmd := xutil.MakeUint32(MainCmd,SubCmd)
	this.msgHandler[cmd] = handler
}
func (this *NetMsgSys)Release() {
	this.msgHandler = nil
}

func (this *NetMsgSys) OnNetWorkMsgHandle(msgdata *network.MsgData) error{
	handle,err := this.GetHandler(msgdata.MainCmd,msgdata.SubCmd)
	if  err != nil {
		//踢掉
		msgdata.Conn.Close()
		return err
	}
	startT := timeutil.GetCurrentTimeMs()		//计算当前时间
	err = handle(msgdata.Conn,msgdata.Msgdata) //查看是否注册对应的处理函数
	since := xutil.MaxInt64(0,timeutil.GetCurrentTimeMs() - startT)
	if since > 20 { //大于20毫秒
		xlog.ErrorLog(appdata.GetSecenName(), "主命令 = %d,子命令 = %d  耗时 %v",msgdata.MainCmd,msgdata.SubCmd, since)
	}
	
	return err
}

func (this *NetMsgSys)GetHandler(MainCmd,SubCmd uint16) (dispatch.HandleMsg,error) {
	cmd := xutil.MakeUint32(MainCmd,SubCmd)
	handle,ok := this.msgHandler[cmd]
	if !ok {
		return nil, errors.New(fmt.Sprintf("主命令 = %d,子命令 = %d  未处理",MainCmd,SubCmd))
	}
	return handle,nil
}

