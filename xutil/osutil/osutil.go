/*
创建时间: 2020/5/1
作者: zjy
功能介绍:

*/

package osutil

import (
	"fmt"
	"wengo/xutil/strutil"
	"net"
	"os"
	"runtime"
)

func MakeDirAll(dir string) bool {
	if strutil.StringIsNil(dir) { // 路径为nil不能创建
		return false
	}
	exists, err := PathExists(dir)
	if !exists {
		if err != nil {
			fmt.Println(dir, " 不存在需要创建", err)
		}
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return false
		}
	}
	return true
}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 获取目录
func ReadDir(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY, os.ModeDir)
}


//获取正在使用的Mac地址
func GetUpMacAddr()  string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v\n", err)
		return ""
	}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		hasflag := uint(netInterface.Flags) & uint(net.FlagUp) == 1
		if hasflag  {
			return macAddr
		}
	}
	return ""
}
//获取Ip地址
func GetIPs() (ips []string) {
	
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}
	
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

//获取运行时文件及行数
//为0时，打印当前调用文件及行数。为1时，打印上级调用的文件及行数
func GetRuntimeFileAndLineStr(skip int) string{
	//其中calldepth 指的调用的深度，为0时，打印当前调用文件及行数。
	//为1时，打印上级调用的文件及行数，依次类推。
	_, file, line, _ := runtime.Caller(skip + 1)
	return fmt.Sprintf("文件%v第%v行",file,line)
}

//获取运行时文件及行数
//为0时，打印当前调用文件及行数。为1时，打印上级调用的文件及行数
func GetRuntimeFileAndLine(skip int) (file string, line int) {
	//其中calldepth 指的调用的深度，为0时，打印当前调用文件及行数。
	//为1时，打印上级调用的文件及行数，依次类推。由于在外面一层调用所以要加1
	_, file, line, _ = runtime.Caller(skip + 1)
	return
}

// 计算内存
func MemConsumed() uint64 {
	runtime.GC() // GC，排除对象影响
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	return memStat.Sys
}