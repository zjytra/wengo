/*
创建时间: 2020/4/6
作者: zjy
功能介绍:

*/

package queue

import "sync"


// Synchronous FIFO queue
type SyncQueue struct {
	lock    sync.Mutex  //锁
	popable *sync.Cond  //条件变量
	buffer  *Queue      // 数据buff
	closed  bool
}

// Create a new SyncQueue
func NewSyncQueue() *SyncQueue {
	ch := &SyncQueue{
		buffer: NewQueue(),
	}
	ch.popable = sync.NewCond(&ch.lock)
	return ch
}

// 提取队列头变量
// 注意:当队列中没有数据的时候会阻塞直到队列中有数据或者等待关闭
func (q *SyncQueue) WaitPop() (v interface{}) {
	q.lock.Lock()
	for q.buffer.Len() == 0 && !q.closed {
		q.popable.Wait()
	}
	if q.buffer.Len() > 0 {
		v,_ = q.buffer.PopFront()
	}
	q.lock.Unlock()
	return
}

// 立即返回,不阻塞,如果队列中数据为空,条件变量为false
// Try to pop an item from SyncQueue, will return immediately with bool=false if SyncQueue is empty
func (q *SyncQueue) TryPop() (v interface{}, ok bool) {
	
	q.lock.Lock()
	
	if q.buffer.Len() > 0 {
		v,_ = q.buffer.PopFront()
		ok = true
	} else if q.closed {
		ok = true
	}
	
	q.lock.Unlock()
	return
}

// Push an item to SyncQueue. Always returns immediately without blocking
func (q *SyncQueue) PushBack(v interface{}) {
	q.lock.Lock()
	if !q.closed {
		q.buffer.PushBack(v)
		q.popable.Signal()
	}
	q.lock.Unlock()
}

// Get the length of SyncQueue
func (q *SyncQueue) Len() (l int) {
	q.lock.Lock()
	l = q.buffer.Len()
	q.lock.Unlock()
	return
}

func (q *SyncQueue) Close() {
	q.lock.Lock()
	if !q.closed {
		q.closed = true
		q.popable.Signal()
	}
	q.lock.Unlock()
}

//查看队列是否关闭
func (q *SyncQueue) IsClose()(b bool) {
	q.lock.Lock()
	b = q.closed
	q.lock.Unlock()
	return
}

//清除队列的数据
func (safeq *SyncQueue)Clear(){
	safeq.lock.Lock()
	safeq.buffer.Clear()
	safeq.lock.Unlock()
}