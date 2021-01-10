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
		//delRes *clientv3.DeleteResponse
		putOp clientv3.Op
		getop clientv3.Op
		opRes clientv3.OpResponse
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

	//OP
	putOp = clientv3.OpPut("/cron/jobs/job8", "")

	//执行op
	if opRes, err = kv.Do(context.TODO(), putOp); nil != err {
		log.Fatal(err)
	}

	fmt.Println(opRes.Put().Header.Revision)

	//getop
	getop = clientv3.OpGet("/cron/jobs/job8")

	//执行op
	if opRes, err = kv.Do(context.TODO(), getop); nil != err {
		log.Fatal(err)
	}

	fmt.Println("数据revision", opRes.Get().Kvs[0].ModRevision)
	fmt.Println("数据value", string(opRes.Get().Kvs[0].Value))
}
