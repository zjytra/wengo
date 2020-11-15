/*
创建时间: 2020/2/29
作者: zjy
功能介绍:

*/

package dispatch

import (
	"errors"
)


// 定时器事件
func (this *DispatchSys) onEventTimer(val interface{}) error {
	Cb,ok:= val.(func())
	if !ok {
		return errors.New("onEventTimer cb erro")
	}
	if Cb == nil {
		return errors.New("onEventTimer Cb is nil")
	}
	Cb() //回调
	return nil
}

