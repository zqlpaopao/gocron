package main

import (
	"crontabInit/master"
	"flag"
	"fmt"
	"runtime"
)

var (
	confFile string //配置文件路径
)

func initENV() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func initArgs() {
	//master-config ./master.config
	flag.StringVar(&confFile, "config", "./master.json", "传入的master.json")
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
	if err = master.InitConfig(); err != nil {
		goto ERR
	}
	//启动Api服务
	if err = master.InitApiServer(); nil != err {
		goto ERR
	}

	return

ERR:
	fmt.Println(err)
}
