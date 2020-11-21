/*
创建时间: 2020/3/29
作者: zjy
功能介绍:

*/

package main

import (
	"container/list"
	"database/sql"
	"encoding/binary"
	"fmt"
	"wengo/timingwheel"
	"wengo/csvdata"
	"wengo/dbsys"
	"wengo/model/dbmodels"
	"wengo/xcontainer/queue"
	"wengo/xutil"
	"wengo/xutil/timeutil"
	"reflect"
	"strings"
	"sync"
	// "sync/atomic"
	"time"
)

var locker = new(sync.Mutex)
var cond = sync.NewCond(locker)
var queue1 *queue.Queue
var wg sync.WaitGroup

var pool sync.Pool

type Person struct {
	name string
	age  int
}

var i int

func Test(t *timingwheel.Timer) {
	fmt.Println("%p  ", t, time.Now())
	if i == 5 {
		t.Stop()
		fmt.Println("Test stop")
	}
	i ++
}

// param |加密key| 主命令 | 字命令 | datalen | data |
func UnpackOne(readb []byte) (maincmd uint16, subcmd uint16, msg []byte, err error) {
	msglen := len(readb)
	if readb == nil || msglen == 0 {
		return
	}
    var start,end int
	start,end = GetNextIndex(end,maincmd)
	maincmd = binary.LittleEndian.Uint16(readb[start:end])
	start,end = GetNextIndex(end,subcmd)
	subcmd = binary.LittleEndian.Uint16(readb[start:end])
	var datalen uint32
	start,end = GetNextIndex(end,datalen)
	datalen = binary.LittleEndian.Uint32(readb[start:end])
	msg = make([]byte, datalen)
	copy(msg,readb[end:])
	return
}

func GetNextIndex(end int,data interface{})(head,tail int)   {
	head = end //尾部变成头部
	tail = head + xutil.IntDataSize(data) //新的尾部=头加上数据的长度
	return
}

// 打单包
func  PackOne(maincmd, subcmd uint16, msg []byte) ([]byte, error) {
	msglen := uint32(len(msg))
	// 加密key
	var alllen int
	alllen += xutil.IntDataSize(maincmd)
	alllen += xutil.IntDataSize(subcmd)
	alllen += xutil.IntDataSize(msglen)
	alllen += int(msglen)
	writeBuf := make([]byte,alllen)
	var start,end int
	start,end = GetNextIndex(end,maincmd)
	binary.LittleEndian.PutUint16(writeBuf[start:end],maincmd)
	start,end = GetNextIndex(end,subcmd)
	binary.LittleEndian.PutUint16(writeBuf[start:end],subcmd)
	start,end = GetNextIndex(end,msglen)
	binary.LittleEndian.PutUint32(writeBuf[start:end],msglen)
	copy(writeBuf[end:],msg)
	return writeBuf, nil
}

type RowStringMap map[string]interface{}

func testQuery(rows *sql.Rows,to interface{})  {
	acc:= dbsys.RowsToStructSlice(rows,to)
	for _, data := range acc {
		ptr := data.(*dbmodels.Accounts)
		fmt.Println(ptr)
	}
}

func DBTest()  {
	csvdata.SetDbconfMapData("./csv")
	conf := csvdata.GetDbconfPtr("gamedb")
	datasource := dbsys.GetMysqlDataSourceName(conf)
	db, Erro := sql.Open("mysql", datasource)
	if Erro != nil {
		fmt.Printf("%v \n",Erro)
		return
	}
	now := time.Now()
	rows,erro := db.Query("SELECT * FROM Accounts where AccountID = 1")
	if erro != nil {
		fmt.Printf("%v \n",erro)
		return
	}
	end := time.Now()
	fmt.Printf("Query time %v \n",end.Sub(now).Microseconds())
	// pacc := dbsys.RowsToStructSlice(rows,reflect.TypeOf(&dbmodels.Accounts{}))
	acc := new(dbmodels.Accounts)
	dbsys.RowToStruct(rows,acc)
	fmt.Println(acc)
	// testQuery(rows,dbmodels.Accounts{})
	// var t time.Time
	// types , erro:= rows.ColumnTypes()

	// var  nt mysql.NullTime
	// for rows.Next() {
	// 	rows.Scan(&nt)
	// 	fmt.Println("mysql.NullTime = ",nt )
	// }
	// rowStrarr := dbsys.RowToStringSlice(rows)
	// fmt.Println(len(rowStrarr),rowStrarr)
	// end2 := time.Now()
	// fmt.Printf("Next1  time %v \n",end2.Sub(end).Microseconds())
	//
	
	// //下一个结果集
	// if !rows.NextResultSet() {
	// 	return
	// }
	// a := dbsys.RowsToStringsSlices(rows)
	// // rowStrarr = append(rowStrarr,a...)
	// // tp := reflect.TypeOf(rowStrarr)
	// end = time.Now()
	// fmt.Printf("Next2 time %v \n",end.Sub(end).Microseconds())
	// fmt.Println(a)
}



func main() {
	//DBTest()
	fmt.Println(timeutil.GetYearMonthFromatStrByTimeString("2006-01-02 15:04:05.000"))
}

func TestList()  {
	//初始化一个list
	l := list.New()
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)
	
	fmt.Println("Before Removing...")
	//遍历list，删除元素
	var n *list.Element
	for e := l.Front(); e != nil; e = n {
		fmt.Println("removing", e.Value)
		n = e.Next()
		l.Remove(e)
	}
	fmt.Println("After Removing...")
	//遍历删除完元素后的list
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}
func Test2(p interface{}) {
	v := reflect.ValueOf(p)
	typ := reflect.TypeOf(p)
	elmv := reflect.Indirect(v)
	fmt.Println("name", elmv.Type().Name())
	
	allname := typ.String()
	splitname := strings.Split(allname, ".")
	if len(splitname) > 1 {
		fmt.Println("name2", splitname[1])
	}
	
	tagv := v.Elem()
	fmt.Println(v.Kind(), v.Elem().Kind())
	for i := 0; i < tagv.NumField(); i++ {
		feildval := tagv.Field(i)
		fmt.Println(feildval.Type().Name())
	}
}

func TestTaskPool() {
	for {
		fmt.Println("TestTaskPool")
		time.Sleep(time.Second)
	}
	
}
func TestTaskPool2() {
	fmt.Println("TestTaskPool2")
}

func TestPool() {
	pool.New = func() interface{} {
		return new(Person)
	}
	data := pool.Get()
	fmt.Println(data)
}

func TestCond() {
	queue1 = queue.NewQueue()
	_, err := queue1.PopFront()
	if err != nil {
		fmt.Println(err)
	}
	wg.Add(20)
	// 10个消费
	for i := 0; i < 10; i++ {
		go func(x int) {
			defer wg.Done()
			cond.L.Lock() // 获取锁
			for queue1.Len() == 0 {
				cond.Wait() // 等待通知，阻塞当前 goroutine
			}
			cond.L.Unlock() // 释放锁
			val, erro := queue1.PopFront()
			if erro != nil {
				fmt.Println(erro)
				return
			}
			// do something. 这里仅打印
			fmt.Println("队列的值 ", val, "type=", reflect.TypeOf(val))
		}(i)
	}
	for i := 0; i < 10; i++ {
		go func(x int) {
			defer wg.Done()
			queue1.PushBack(x)
			cond.Signal() // 通知其他线程
		}(i)
	}
	
	wg.Wait()
	fmt.Printf("end")
}
