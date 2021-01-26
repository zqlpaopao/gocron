package main

import (
	"crontabInit/master"
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
	flag.StringVar(&confFile, "config", "./master/main/master.json", "传入的master.json")
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
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	//初始化集群管理器
	if err = master.InitWorkerMgr(); err != nil {
		goto ERR
	}

	//任务管理器
	if err = master.InitJobMge(); nil != err {
		goto ERR
	}

	//启动Api服务
	if err = master.InitApiServer(); nil != err {
		goto ERR
	}
	for {
		time.Sleep(time.Second)
	}

	return

ERR:
	fmt.Println(err)
}
