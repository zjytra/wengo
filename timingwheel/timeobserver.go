/*
创建时间: 2020/5/24
作者: zjy
功能介绍:
定时器触发观察者
*/

package timingwheel

type TimerEventObserver interface {
	PostTimerEvent(cb func()) error
}
