// 创建时间: 2019-10-2019/10/17
// 作者: zjy
// 功能介绍:
// 1.主要入口
// 2.
// 3.
package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"wengo/app"
	"runtime"
	"time"
)
// main 初始化工作
func init() {
}

func TestRedis() {
	rediscli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	
	_, erro := rediscli.Ping().Result()
	if erro != nil {
		fmt.Print(erro)
		return
	}
	erro = rediscli.Set("woshishui", "jiaopige", time.Minute*2).Err()
	if erro != nil {
		fmt.Print(erro)
		return
	}
	getval := rediscli.Get("woshishui")
	if geterro := getval.Err(); geterro != nil {
		if geterro == redis.Nil {
			fmt.Println("key does not exists")
			return
		}
		fmt.Printf("geterro %v ", geterro)
		return
	}
	fmt.Printf(getval.Val()+"\n")
}
// 各服务器主入口
func main() {
	TestRedis()
	// pro := profile.Start(profile.MemProfile,profile.ProfilePath("./profiles"))
	// 设置最大运行核数
	runtime.GOMAXPROCS(runtime.NumCPU())
	app.NewApp().AppStart()
	// 等待退出 在app 退出后整个程序退出
	// pro.Stop()
}



