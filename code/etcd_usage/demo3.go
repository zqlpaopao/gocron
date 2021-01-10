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

	if getRes, err = kv.Get(context.TODO(), "/cron/jobs/jobs1"); nil != err {
		log.Fatal(err)
	} else {
		fmt.Println(getRes.Kvs)    //[key:"/cron/jobs/jobs1" create_revision:7 mod_revision:10 version:4 value:"bytes" ]
		fmt.Println(getRes.Header) //cluster_id:14841639068965178418 member_id:10276657743932975437 revision:10 raft_term:3
		fmt.Println(getRes.Count)  //1
	}
}
