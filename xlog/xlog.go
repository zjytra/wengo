//  创建时间: 2019/10/22
//  作者: zjy
//  功能介绍:
//  $ log日志,输出封装,这里要调整最好不予App有依赖
// 为了提高逻辑线程处理效率
//  多 gocorutine 传入日志信息 ,单 gocorutine 向文件写入
//  每日一个文件夹,每天两小时对应日志进程文件

package xlog

import (
	"fmt"
	"github.com/zjytra/wengo/model"
	"github.com/zjytra/wengo/xutil/osutil"
	"github.com/zjytra/wengo/xutil/timeutil"
	"io"
	"log"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

// 日志等级定义
const (
	Normal     = 0
	DebugLvl   = 1 << 0
	WarningLvl = 1 << 1
	ErrorLvl   = 1 << 2
)

var (
	_xlog        *Xlog // 日志执行者对象
	loglvlStrMap map[uint16]string
	_wg          sync.WaitGroup // 为保证程序统一退出这里加个等待
)

type Xlog struct {
	initInfo     *LogInitModel  // 初始化
	baseLog      *log.Logger    // 内置log库的处理
	logBufchan   chan *LogModel // 日志信息
	closeFlag    *model.AtomicBool
	logmodelPool sync.Pool
}

// 创建日志对象
func NewXlog(info *LogInitModel) bool {
	_xlog = new(Xlog)
	if _xlog == nil {
		fmt.Println("_NewXlog xlog is nil")
		return false
	}
	_xlog.baseLog = log.New(os.Stdout, "", 0)
	_xlog.logBufchan = make(chan *LogModel, info.Volatile.LogQueueCap)
	_xlog.initInfo = info
	
	loglvlStrMap = make(map[uint16]string)
	initXlog()
	_wg.Add(1)
	go _xlog.run()

	return true
}

func initXlog() {
	
	loglvlStrMap[Normal] = "无"
	loglvlStrMap[DebugLvl] = "调试|"
	loglvlStrMap[WarningLvl] = "警告|"
	loglvlStrMap[ErrorLvl] = "错误|"
	_xlog.closeFlag = model.NewAtomicBool()
	_xlog.closeFlag.SetTrue()
	_xlog.logmodelPool.New = func() interface{} {
		return new(LogModel)
	}
}

// 设置日志等级并设置是否在控制台显示 目前这两个经常改变
func SetShowLogAndStartLog(restmodel VolatileLogModel) bool {
	if _xlog == nil {
		fmt.Println("SetShowLogAndStartLog xlog is nil")
		return false
	}
	_xlog.initInfo.Volatile = restmodel
	return true
}

func DebugLogNoInScene(format string, v ...interface{}) {
	addLogToLogBufchan(DebugLvl, "", format, v...)
}
func WarningLogNoInScene(format string, v ...interface{}) {
	addLogToLogBufchan(WarningLvl, "", format, v...)
}

func ErrorLogNoInScene( format string, v ...interface{}) {
	addLogToLogBufchan(ErrorLvl, "", format, v...)
}


func DebugLog(scenename string, format string, v ...interface{}) {
	addLogToLogBufchan(DebugLvl, scenename, format, v...)
}
func WarningLog(scenename string, format string, v ...interface{}) {
	addLogToLogBufchan(WarningLvl, scenename, format, v...)
}

func ErrorLog(scenename string, format string, v ...interface{}) {
	addLogToLogBufchan(ErrorLvl, scenename, format, v...)
}

// 向log日志队列中写日志信息
func addLogToLogBufchan(loglvl uint16, scenename string, format string, v ...interface{}) {
	if _xlog == nil {
		fmt.Println("addLogToLogBufchan xlog is nil 日志为: ",fmt.Sprintf(format, v...))
		return
	}
	if _xlog.closeFlag.IsFalse() {
		return
	}
	// 未设置对应的日志等级就不能打印
	if !canLogBylvl(loglvl) {
		return
	}
	tem := _xlog.logmodelPool.Get()
	if  tem == nil {
		return
	}
	lm,ok := tem.(*LogModel)
	if !ok {
		return
	}
	lm.OutStr = fmt.Sprintf(format, v...)
	lm.LogGenerateTime = timeutil.GetTimeNow()
	lm.LogLvel = loglvl
	lm.SceneName = scenename
	_xlog.logBufchan <- lm
}

func (xl *Xlog) writeLogToFile(lm *LogModel) {
	if lm == nil {
		return
	}
	isOk := xl.newLogsDir(lm.LogGenerateTime) // 查看目录是否存在
	if !isOk {
		return
	}
	lm.WriteFile = xl.newLogFile(lm) // 创建文件
	if lm.WriteFile == nil {  //文件对象未创建成功
		return
	}
	xl.setOutFile(lm.WriteFile)
	xl.setOutPrefix(lm.LogLvel, lm.LogGenerateTime)
	xl.baseLog.Println(lm.OutStr) // 向输出流输出字符串
	lm.WriteFile.Close()          // 最后关闭文件
	lm.WriteFile = nil
	_xlog.logmodelPool.Put(lm)    // 放回池子
}

// 创建日志日期路径
func (xl *Xlog) newLogsDir(time time.Time) bool {
	dirs := path.Join(xl.initInfo.LogsPath, timeutil.GetYearMonthDayFromatStr(time))
	return osutil.MakeDirAll(dirs)
}

func (xl *Xlog) newLogFile(lm *LogModel) *os.File {
	// 每两个小时一个文件
	filename := timeutil.GetYearMonthDayHourFromatStrBySpan(lm.LogGenerateTime, xl.initInfo.Volatile.FileTimeSpan) + "_" + xl.initInfo.ServerName + "_" + lm.SceneName + ".log"
	str := path.Join(xl.initInfo.LogsPath, timeutil.GetYearMonthDayFromatStr(lm.LogGenerateTime), filename)
	tempfile, err := os.OpenFile(str, os.O_CREATE | os.O_WRONLY |os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("打开日志文件错误 = %v ", err)
		tempfile.Close()
		return nil
	}
	return tempfile
}

func (xl *Xlog) setOutFile(writeFile *os.File) {
	if writeFile == nil {
		xl.baseLog.SetOutput(os.Stdout)
	} else if xl.initInfo.Volatile.IsOutStd && writeFile != nil {
		xl.baseLog.SetOutput(io.MultiWriter(os.Stdout, writeFile))
	} else {
		xl.baseLog.SetOutput(writeFile)
	}
}

func (xl *Xlog) setOutPrefix(reqlvl uint16, t time.Time) {
	// 清除日志时间
	if prefixStr, ok := loglvlStrMap[reqlvl]; ok {
		// 日志等级与生成时间
		xl.baseLog.SetPrefix(fmt.Sprintf("%s%s\t", prefixStr, timeutil.GetTimeALLStr(t)))
	} else {
		xl.baseLog.SetPrefix(timeutil.GetTimeALLStr(t))
		
	}
}

// 日志执行逻辑线程
func (xl *Xlog) run() {
	
	defer _wg.Done()
	// 拉起宕机
	defer func() {
		if rec := recover(); rec != nil {
			buf := make([]byte, LenStackBuf)
			l := runtime.Stack(buf, false)
			fmt.Printf("%v\n%s \n", rec, buf[:l])
			//重开运行
			_wg.Add(1)
			go _xlog.run()
		}
	}()
	
	// 这里接收写的通道，没有数据会一直阻塞，直到通道关闭
	for logmodel := range xl.logBufchan {
		if logmodel == nil {
			break
		}
		 xl.writeLogToFile(logmodel)
	}

}

// 如果不能输出日志都在标准中输出
func canLogBylvl(loglvl uint16) bool {
	if _xlog == nil {
		return false
	}
	return (loglvl & _xlog.initInfo.Volatile.ShowLvl) != 0
}

// 关闭日志
func CloseLog() {
	_xlog.closeFlag.SetFalse()
	_xlog.logBufchan <- nil
	_xlog.onClose()
}

func (xl *Xlog) onClose() {
	_wg.Wait()
	close(xl.logBufchan)
	loglvlStrMap = nil
	fmt.Println("Log doClose")
}
