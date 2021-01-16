package main

import (
	"crontabInit/worker"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string //配置文件路径
)

func initENV() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	//master-config ./master.config
	flag.StringVar(&confFile, "config", "./worker/main/worker.json", "传入的master.json")
	flag.Parse()

}

func main() {
	var (
		err error
	)

	//初始化线程
	initENV()

	//初始化命令行参数
	initArgs()

	//加载配置
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	//启动调度器
	if err = worker.InitScheduler(); nil != err {
		goto ERR
	}

	//执行期
	worker.InitExecutor()

	//任务管理器
	if err = worker.InitJobMge(); nil != err {
		goto ERR
	}

	//启动监听
	//worker.WatchJobs()
	for {
		time.Sleep(time.Second)
	}

	return

ERR:
	fmt.Println(err)
}
