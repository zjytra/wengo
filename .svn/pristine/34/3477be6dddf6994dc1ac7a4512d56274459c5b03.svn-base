/*
创建时间: 2020/4/1
作者: zjy
功能介绍:

*/

package model

import "sync/atomic"

type AtomicBool struct {
	boolFlag int32
}

const
(
	falseFlag  = 0
	trueFlag = 1
)

// 默认是false
func NewAtomicBool() *AtomicBool {
	return new(AtomicBool)
}

// 检测是否false
func (af *AtomicBool) IsFalse() bool {
	return atomic.LoadInt32(&af.boolFlag) == falseFlag
}

// 检测是否正确
func (af *AtomicBool) IsTrue() bool {
	return atomic.LoadInt32(&af.boolFlag) == trueFlag
}

func (af *AtomicBool) SetFalse() {
	af.setBool(falseFlag)
}

func (af *AtomicBool) SetTrue() {
	af.setBool(trueFlag)
}

func (af *AtomicBool) setBool(flagval int32) {
	atomic.StoreInt32(&af.boolFlag,flagval)
}

