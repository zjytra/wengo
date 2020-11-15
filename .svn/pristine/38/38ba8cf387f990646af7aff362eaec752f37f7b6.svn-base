/*
创建时间: 2020/5/21
作者: zjy
功能介绍:
//自增id生成器暂时不要缓存
*/

package model

type IncrementId struct {
	id *AtomicUInt32FlagModel
	// freeIds *queue.SafeQueue  //回收的id列表
}

func NewIncrementId() *IncrementId  {
	return &IncrementId{
		id: NewAtomicUInt32Flag(),
		// freeIds:queue.NewSafeQueue(),
	}
}

func (this *IncrementId)GetId() (id uint32) {
	//先查看容器里面有没有
	// val,erro := this.freeIds.PopFront()
	// if erro == nil {//成功了就返回回收后的id
	// 	var ok bool
	// 	id,ok = val.(uint32)
	// 	if ok {
	// 		return
	// 	}
	// }
	id = this.id.AddUint32()
	if id == 0 {
		id = this.id.AddUint32()
	}
	return
}

// //回收id
// func (this *IncrementId)RecycleId(id uint32){
// 	// this.freeIds.PushBack(id)
// }

//回收id
func (this *IncrementId)Release(){
	// this.freeIds.Clear()
	// this.freeIds = nil
	this.id = nil
}