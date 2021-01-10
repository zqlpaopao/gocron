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
		delRes *clientv3.DeleteResponse
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

	//del
	if delRes, err = kv.Delete(context.TODO(), "/cron/jobs/jobs1", clientv3.WithPrevKV()); nil != err {
		log.Fatal(err)
	}

	if len(delRes.PrevKvs) > 0 {
		fmt.Println(delRes.PrevKvs)
	}
}
