/*
创建时间: 2019/11/24
作者: zjy
功能介绍:
安全的队列操作
*/

package queue

import (
	"container/list"
	"sync"
)


// 线程安全的队列
type SafeQueue struct {
	queue  *Queue
	lock sync.RWMutex
}

// 队列构造函数
// return 返回队列
func NewSafeQueue() (sfqueue *SafeQueue) {
	sfqueue = new(SafeQueue)
	sfqueue.queue = NewQueue()
	return sfqueue
}





// 获得队列的长度
func (safeq *SafeQueue)Len()  (len int) {
	safeq.lock.RLock()
	len = safeq.queue.Len()
	safeq.lock.RUnlock()
	return
}

func (safeq *SafeQueue)Front() (element *list.Element)  {
	safeq.lock.RLock()
	element = safeq.queue.Front()
	safeq.lock.RUnlock()
	return
}

//向队列中添加数据
func (safeq *SafeQueue)PushBack(v interface{})  {
	safeq.lock.Lock()
	safeq.queue.PushBack(v)
	safeq.lock.Unlock()
}

func (safeq *SafeQueue)PopFront() (v interface{},erro error)  {
	safeq.lock.Lock()
	v,erro = safeq.queue.PopFront()
	safeq.lock.Unlock()
	return
}

func (safeq *SafeQueue)Clear(){
	safeq.lock.Lock()
	safeq.queue.Clear()
	safeq.lock.Unlock()
}