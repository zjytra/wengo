/*
创建时间: 2020/2/21
作者: zjy
功能介绍:

*/

package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/zjytra/wengo/app/appdata"
	"github.com/zjytra/wengo/csvdata"
)

//读取控制台命令
func (this *App)ReadConsle()  {
	// scanner := bufio.NewScanner(os.Stdin)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Shell")
	fmt.Println("---------------------")
	
	defer this.AppWG.Done()
	for this.AppIsRun(){
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		switch text {
		case "reloadcsv":
			csvdata.ReLoadCommonCsvData(appdata.WorkPool)
		case "close":
			this.CloseApp()
			break
		}
	}
}
