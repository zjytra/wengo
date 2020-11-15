/*
创建时间: 2019/12/22
作者: zjy
功能介绍:
原子开关用于检测
*/

package model

import "sync/atomic"

type AtomicInt32FlagModel struct {
	 checkFlag int32
}
func NewAtomicInt32Flag() *AtomicInt32FlagModel {
	return new(AtomicInt32FlagModel)
}
// 检测是否关闭
func (af *AtomicInt32FlagModel) GetInt32() int32 {
	return atomic.LoadInt32(&af.checkFlag)
}
// 检测是否开启
func (af *AtomicInt32FlagModel) SetInt32(flagval int32) {
	atomic.StoreInt32(&af.checkFlag,flagval)
}
//自增
func (af *AtomicInt32FlagModel) AddInt32()(new int32) {
	return atomic.AddInt32(&af.checkFlag,1)
}

//自减
func (af *AtomicInt32FlagModel) SubInt32()(new int32) {
	return atomic.AddInt32(&af.checkFlag,-1)
}



type AtomicUInt32FlagModel struct {
	checkFlag uint32
}

func NewAtomicUInt32Flag() *AtomicUInt32FlagModel {
	return new(AtomicUInt32FlagModel)
}

func (af *AtomicUInt32FlagModel) GetUInt32() uint32 {
	return atomic.LoadUint32(&af.checkFlag)
}

func (af *AtomicUInt32FlagModel) SetUInt32(flagval uint32) {
	atomic.StoreUint32(&af.checkFlag,flagval)
}
//自增
func (af *AtomicUInt32FlagModel) AddUint32()(new uint32) {
	return atomic.AddUint32(&af.checkFlag,1)
}

//自减
func (af *AtomicUInt32FlagModel) SubUint32()(new uint32) {
	d := int32(-1)
	return atomic.AddUint32(&af.checkFlag,uint32(d))
}