/*
创建时间: 2020/4/24
作者: zjy
功能介绍:

*/

package timersys

import (
	"github.com/RussellLuo/timingwheel"
	"github.com/wengo/model"
	"sync"
	"time"
)

var (
	Twheel *timingwheel.TimingWheel
	timers  sync.Map
	timeids  *model.IncrementId           // 自增id生成器,用id，有利于连接在其他地方使用，降低包的依赖
)

func init()  {
	Twheel = timingwheel.NewTimingWheel(time.Millisecond, 20)
	Twheel.Start()
	timeids = model.NewIncrementId()
}



//定时器
func NewWheelTimer(interval time.Duration,f func(),ob timingwheel.TimerEventObserver)  uint32{
	if f == nil {
		return 0
	}
	t := Twheel.ScheduleFunc(&EveryScheduler{interval},f,ob)
	timeId := timeids.GetId()
	timers.Store(timeId,t)
	return timeId
}
//多少时间后执行
func AfterFunc(interval time.Duration,f func(),ob timingwheel.TimerEventObserver) uint32 {
	if f == nil {
		return 0
	}
	t := Twheel.AfterFunc(interval,f,ob)
	timeId := timeids.GetId()
	timers.Store(timeId,t)
	return timeId
}

func Close()  {
	//清理
	timers.Range(func(key, value interface{}) bool {
		value.(*timingwheel.Timer).Stop()
		timers.Delete(key)
		return true
	})
	Twheel.Stop()
}

func StopTimer(timeId uint32)  {
	timer,ok := timers.Load(timeId)
	if !ok {
		return
	}
	timer.(*timingwheel.Timer).Stop()
	timers.Delete(timeId)
}

func Release()  {
	Close()
	timeids.Release()
}