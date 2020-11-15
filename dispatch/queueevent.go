/*
创建时间: 2020/4/26
作者: zjy
功能介绍:
队列事件，接收各种事件进队，可以多个协程处理，也可以单个协程处理
*/

package dispatch

import (
	"errors"
	"github.com/wengo/xcontainer/queue"
	"github.com/wengo/xlog"
	"sync"
)

type OnQueueEventFunc func(eventdata *EventData)

type QueueEvent struct {
	eventQueue *queue.SyncQueue // 存放事件的的队列
	wg         sync.WaitGroup
}

func NewQueueEvent() *QueueEvent {
	qe := new(QueueEvent)
	qe.eventQueue = queue.NewSyncQueue()
	return qe
}

// 添加事件处理者，可以添加多个，代表多线程处理事件队列
func (this *QueueEvent) AddEventDealer(eventfunc OnQueueEventFunc) error {
	if eventfunc == nil {
		return errors.New("AddEventDealer deal is nil")
	}
	this.wg.Add(1)
	go this.run(eventfunc)
	return nil
}

func (this *QueueEvent) AddEvent(event *EventData) {
	this.eventQueue.PushBack(event)
}

func (this *QueueEvent) run(onQueueEventfunc OnQueueEventFunc) {
	
	// 拉起错误避免宕机
	defer xlog.RecoverToLog(func() {
		this.wg.Add(1)
		go	this.run(onQueueEventfunc)
	})
	defer this.wg.Done()
	for {
		// 当事件处理完并且没有数据
		qlen := this.eventQueue.Len()
		if qlen > 100 {
			xlog.WarningLogNoInScene( "QueueEvent.eventQueue剩余需要处理的数据 = %d", qlen)
		}
		if qlen == 0 && this.eventQueue.IsClose() {
			break
		}
		// 取队列数据
		eventdata := this.eventQueue.WaitPop() // 等待取出
		if eventdata == nil {
			xlog.DebugLogNoInScene( " eventQueue取数据为nil ")
			continue
		}
		// 解析数据
		data, ok := eventdata.(*EventData)
		if !ok {
			xlog.DebugLogNoInScene(  " OnQueueEvent 解析数据失败")
			continue
		}
		onQueueEventfunc(data)
	}
	
}

// 先关闭队列等待任务处理完
func (this *QueueEvent) Release() {
	this.eventQueue.Close()
	xlog.DebugLogNoInScene( "还有未处理的事件%d", this.eventQueue.Len())
	this.wg.Wait()
	this.eventQueue.Clear()
	this.eventQueue = nil
}
