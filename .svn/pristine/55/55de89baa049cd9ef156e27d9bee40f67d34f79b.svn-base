/*
创建时间: 2019/11/23
作者: zjy
功能介绍:
队列结构,对标准库列表封装
*/

package queue

import (
	"container/list"
	"errors"
	"fmt"
)

type Queue struct {
	qlist  *list.List
}

// 队列构造函数
// return 返回队列
func NewQueue() (queue *Queue) {
	queue = new(Queue)
	queue.qlist = list.New()
	return queue
}

//向队列中添加数据
func (q *Queue)PushBack(v interface{})  {
	q.qlist.PushBack(v)
}

func (q *Queue)Len()  int {
	return q.qlist.Len()
}

func (q *Queue)Front()  *list.Element{
	if q.Len()  == 0 {
		fmt.Println("Front q is nil")
		return nil
	}
	return q.qlist.Front()
}

func (q *Queue)Remove(e *list.Element) (interface{}, error)  {
	if q.Len()  == 0 {
		return  nil,errors.New("Queue Empty")
	}
	return 	q.qlist.Remove(e),nil
}

func (q *Queue)PopFront() (interface{}, error)  {
	if q.Len()  == 0 {
		return  nil,errors.New("Queue Empty")
	}
	return 	q.qlist.Remove(q.qlist.Front()),nil
}

func (q *Queue)Clear() {
	if q.qlist.Len() == 0 {
		return
	}
	var n *list.Element  //下一个数据的变量临时存放
	for e :=  q.qlist.Front(); e != nil; e = n{
		n = e.Next() //先保存下一个数据
		q.qlist.Remove(e) //删除当前的数据
	}
	q.qlist = nil
}