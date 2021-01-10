package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

//建立连接

func main() {
	var (
		conf   clientv3.Config
		client *clientv3.Client
		err    error
	)

	conf = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	//建立连接
	if client, err = clientv3.New(conf); nil != err {
		log.Fatal(err)
	}
	fmt.Println(client)
}
