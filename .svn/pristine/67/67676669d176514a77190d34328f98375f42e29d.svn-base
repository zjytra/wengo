/*
创建时间: 2019/11/24
作者: zjy
功能介绍:
日期相关辅助函数
*/

package timeutil

import (
	"fmt"
	"strings"
	"time"
)

var (
	TimeAllTemplate    string // 带毫秒的模板
	DateTemplate       string // 日期模板
	HHmmssTemplate     string // 时分秒
	YMDHTemplate       string
	FileTemplateDate   string // 日期模板
	FileTemplatemsYMDH string
	diffUnixNano int64 //与数据中心相差时间
)

func init() {
	TimeAllTemplate = "2006-01-02 15:04:05.000" // 常规类型
	DateTemplate = "2006-01-02"                 // 只有日期
	HHmmssTemplate = "15:04:05"                 // 时分秒
	YMDHTemplate = "2006-01-02 15:04:05"        // 常规类型
	FileTemplateDate = "20060102"
	FileTemplatemsYMDH ="2006010215"
}

func SetDiffUnixNano(diff int64){
	if diff == 0 {
		return
	}
	diffUnixNano = diff
}

func GetTimeNow() time.Time {
	now := time.Now()
	if diffUnixNano == 0 { //不差
		return now
	}
	nowunix := now.UnixNano() + diffUnixNano
	return 	time.Unix(0,nowunix)
}
//当前时间增加时间
func NowAddDate(years int, months int, days int) time.Time {
	return GetTimeNow().AddDate(years,months,days)
}

//获取当前时间戳秒
func GetCurrentTimeS() int64 {
	now := GetTimeNow()
	return now.Unix()
}
//获取当前时间戳毫秒
func GetCurrentTimeMs() int64 {
	return int64(GetCurrentTimeNano() / int64(time.Millisecond))
}
//获取当前时间戳纳秒
func GetCurrentTimeNano() int64 {
	now := GetTimeNow()
	return now.UnixNano()
}

func TimeStrToTime(timestr string) time.Time  {
	formatTime,_:=time.Parse(TimeAllTemplate,timestr)
	return formatTime
}

func GetTimeStrByTime(tm time.Time,template string) string  {
	return tm.Format(template)
}

func getTimeStrByTimeNano(timeNano int64,template string) string  {
	return time.Unix(0,timeNano).Format(template)
}

func GetTimeALLStr(t time.Time) string {
	return GetTimeStrByTime(t, TimeAllTemplate)
}

func GetYearMonthFromatStrByTimeString(t string) string {
	if len(t) < 8 {
		return ""
	}
	return strings.Replace(t[0:8],"-","",-1)
}

func GetYearMonthDayFromatStrByTimeString(t string) string {
	if len(t) < 10 {
		return ""
	}
	return strings.Replace(t[0:10],"-","",-1)
}



func GetYearMonthFromatStr(t time.Time) string {
	datestr := fmt.Sprintf("%d%02d",
		t.Year(),
		t.Month())
	return datestr
}

func GetYearMonthDayFromatStr(nowTime time.Time) string {
	datestr := fmt.Sprintf("%d%02d%02d",
		nowTime.Year(),
		nowTime.Month(),
		nowTime.Day())
	return datestr
}
func GetYearMonthDayHourFromatStr(nowTime time.Time) string {
	datestr := fmt.Sprintf("%d%02d%02d_%2d",
		nowTime.Year(),
		nowTime.Month(),
		nowTime.Day(),
		nowTime.Hour())
	return datestr
}
// 获取时间间隔的字符串
func GetYearMonthDayHourFromatStrBySpan(nowTime time.Time,span int) string {
	hour := nowTime.Hour()
	switch span {
	case 2,3,4,5,6,9,12,24: //日志间隔
		if !IsSpanTime(hour,span) {
			hour -= (hour % span)
		}
	default:
	
	}
	
	datestr := fmt.Sprintf("%d%02d%02d_%d",
		nowTime.Year(),
		nowTime.Month(),
		nowTime.Day(),
		hour)
	return datestr
}

func IsSpanTime(hour,span int) bool  {
	return  hour % span == 0
}

func GetDateFileName(timeNano int64) string {
	return getTimeStrByTimeNano(timeNano, FileTemplateDate)
}

//将数据库字符串解析成秒
func TimeParseInToSecond(timestr string)int64{
	tm,erro := TimeParseIn(timestr)
	if erro != nil {
		fmt.Printf("TimeParseInToSecon erro %v \n",erro)
		return 0
	}
	return tm.Unix()
}

//将数据库字符串解析成go time
func TimeParseIn(timestr string)(tm time.Time,erro error) {
	tm,erro = time.ParseInLocation(YMDHTemplate,timestr,time.Local)
	if erro != nil {
		fmt.Printf("TimeParseIn erro %v\n",erro)
		return
	}
	return
}