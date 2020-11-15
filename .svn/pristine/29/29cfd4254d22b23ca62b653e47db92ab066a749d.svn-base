/*
创建时间: 2020/4/23
作者: zjy
功能介绍:

*/

package timersys

import (
	"time"
)


type EveryScheduler struct {
	Interval time.Duration
}

func (s *EveryScheduler) Next(prev time.Time) time.Time {
	return prev.Add(s.Interval)
}

