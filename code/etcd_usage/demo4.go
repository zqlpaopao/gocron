package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

func main() {
	var (
		conf   clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		getRes *clientv3.GetResponse
	)

	conf = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	//建立连接
	if client, err = clientv3.New(conf); nil != err {
		log.Fatal(err)
	}

	//put key
	kv = clientv3.NewKV(client)

	//读取前缀key的所有选项
	if getRes, err = kv.Get(context.TODO(), "/cron/jobs/", clientv3.WithPrefix()); nil != err {
		log.Fatal(err)
	} else {
		fmt.Println(getRes.Kvs)
	}
}
