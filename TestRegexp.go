package main

import (
	"fmt"
	"regexp"
)

func main() {
	patternstr := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, erro := regexp.Compile(patternstr)//匹配到一个或多个空格或者非单词字符
	if erro != nil {
		fmt.Println("%v",erro)
	}
	isok := re.MatchString("update")
	fmt.Println(isok)
	
}
