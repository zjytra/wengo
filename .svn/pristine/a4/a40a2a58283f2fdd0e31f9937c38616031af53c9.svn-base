/*
创建时间: 2020/5/1
作者: zjy
功能介绍:

*/

package strutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func StrToInt8(str string) int8 {
	return int8(StrToInt64(str))
}

func StrToUint8(str string) uint8 {
	return uint8(StrToUint64(str))
}

func StrToInt16(str string) int16 {
	return int16(StrToInt64(str))
}

func StrToUint16(str string) uint16 {
	return uint16(StrToUint64(str))
}

func StrToInt32(str string) int32 {
	return int32(StrToInt64(str))
}

func StrToUint32(str string) uint32 {
	return uint32(StrToUint64(str))
}

func StrToInt(str string) int {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return i
}

func StrToInt64(str string) int64 {
	i, e := strconv.ParseInt(str,10,64)
	if e != nil {
		return 0
	}
	return i
}

func StrToUint64(str string) uint64 {
	i, e := strconv.ParseUint(str,10,64)
	if e != nil {
		return 0
	}
	return i
}

//查看字符串是否包含空格或特殊字符
func StringHasSpaceOrSpecialChar(str string) bool {
	isMatch, erro := regexp.MatchString(`[\s]+|[\W]+`,str) //匹配到一个或多个空格或者非单词字符
	if erro != nil {
		fmt.Println("StringHasSpaceOrSpecialChar %v,str = %v", erro,str)
		return isMatch
	}
	return isMatch
}


//sql注入匹配
func StringHasSqlKey(str string) bool {
	patternstr := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	isMatch, erro := regexp.MatchString(patternstr,str) //
	if erro != nil {
		fmt.Println("StringHasSpaceOrSpecialChar %v,str = %v", erro,str)
		return isMatch
	}
	return isMatch
}
// 判断字符串是否有数据  无数据返回true
func StringIsNil(str string) bool {
	return len(str) == 0 ||  strings.Compare(str,"") == 0
}

