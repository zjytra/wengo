/*
创建时间: 2020/4/26
作者: zjy
功能介绍:

*/

package timersys

import (
	"time"
)

type TimeTicker struct {
	t       *time.Ticker
	cb      func()
	closech chan int
}

func NewTimeTicker(interval time.Duration, cb func()) (tt *TimeTicker) {
	this := new(TimeTicker)
	this.t = time.NewTicker(interval)
	this.cb = cb
	this.closech = make(chan int, 1)
	go this.tickRun()
	return this
}

func (this *TimeTicker) tickRun() {
	for {
		select {
		case <-this.t.C:
			this.cb()
		case <-this.closech:
			close(this.closech)
			return
		}
	}
}

//不能调用两次哦
func (this *TimeTicker) StopTicker() {
	this.t.Stop()
	this.closech <- 1
}
