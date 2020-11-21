/*
创建时间: 2019/11/6
作者: zjy
功能介绍:
工具包
*/

package xutil

import (
	"errors"
	"fmt"
	"github.com/zjytra/wengo/xutil/strutil"
	"math"
	"path"
	"runtime"
	"strings"
	"unsafe"
)



// 是否错误，有错返回 true无错返回false
func IsError(err error) bool {
	if err != nil {
		buf := make([]byte, 4096)
		l := runtime.Stack(buf, false)
		fmt.Printf("%v \n%s", err, buf[:l])
		return true
	}
	return false
}

// 是否错误，有错返回 true无错返回false
func IsErrorNoPrintf(err error) bool {
	return err != nil
}

func SprintfAssertObjErro(objtype string) error{
	_, file, line, _ := runtime.Caller(0)
	return errors.New(fmt.Sprintf("文件%s 第%d行 断言数据 %s 异常",file,line,objtype))
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	if strutil.StringIsNil(str) {
		return str
	}
	var upperStr string
	vv := []rune(str)
	if vv[0] >= 97 && vv[0] <= 122 { // 后文有介绍
		vv[0] -= 32 // string的码表相差32位
		upperStr = string(vv[0]) + string(vv[1:len(vv)])
	} else {
		fmt.Println("Not begins with lowercase letter,")
		return str
	}
	
	return upperStr
}

// 是否是xlsx 文件
func IsXlsx(fileName string) bool {
	return path.Ext(fileName) == ".xlsx" && !strings.HasPrefix(fileName, "~$")
}

// 验证csv行数据是否有效
// 除第三行外,行没有注释 str首字符 != #  ASCII表 35
// 并且id不为nil
func ValidCsvRow(str string, rownum int) bool {
	if strutil.StringIsNil(str) {
		return false
	}
	if rownum != 2 && str[0] == 35 {
		return false
	}
	return true
}



// 获取包的字符串名称
func GetPackageStr(pkgname string) string {
	return fmt.Sprintf("\"%s\"", pkgname)
}

// 主机是否是小端序列编码
func IsLittleEndian() bool {
	n := 0x1234
	// 转换获取小端的数值
	f := *((*byte)(unsafe.Pointer(&n)))
	return (f ^ 0x34) == 0
}
func MaxInt64(one, two int64) int64 {
	if one >= two {
		return one
	}
	return two
}

func MinInt64(one, two int64) int64 {
	if one <= two {
		return one
	}
	return two
}

func MaxUint32(one, two uint32) uint32 {
	if one >= two {
		return one
	}
	return two
}

func MinUint32(one, two uint32) uint32 {
	if one <= two {
		return one
	}
	return two
}

func MaxInt(one, two int) int {
	if one >= two {
		return one
	}
	return two
}

func MinInt(one, two int) int {
	if one <= two {
		return one
	}
	return two
}

// 将两个16位命令组合在一起
func MakeUint32(main, sub uint16) uint32 {
	return uint32(main)<<16 | uint32(sub)
}

// 将32位拆为两个 16位
func UnUint32(cmd uint32) (mn, sub uint16) {
	mn = uint16(cmd >> 16)
	// 去掉高位
	sub = uint16(cmd & math.MaxUint16)
	return
}


func IntDataSize(data interface{}) int {
	switch data := data.(type) {
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 1
	case []bool:
		return len(data)
	case []int8:
		return len(data)
	case []uint8:
		return len(data)
	case int16, uint16, *int16, *uint16:
		return 2
	case []int16:
		return 2 * len(data)
	case []uint16:
		return 2 * len(data)
	case int32, uint32, *int32, *uint32:
		return 4
	case []int32:
		return 4 * len(data)
	case []uint32:
		return 4 * len(data)
	case int64, uint64, *int64, *uint64:
		return 8
	case []int64:
		return 8 * len(data)
	case []uint64:
		return 8 * len(data)
	}
	return 0
}

//移除切片中某个元素
func RemoveSliceByElem(slice []interface{}, elem interface{}) []interface{}{
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveSliceByElem(slice,elem)
		}
	}
	return slice
}

//移除切片中某个元素
func RemoveUint32SliceByElem(slice []uint32, elem uint32) []uint32{
	if len(slice) == 0 {
		return slice
	}
	for i, v := range slice {
		if v == elem {
			slice = append(slice[:i], slice[i+1:]...)
			return RemoveUint32SliceByElem(slice,elem)
		}
	}
	return slice
}


