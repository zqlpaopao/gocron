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
		putRes *clientv3.PutResponse
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

	//查看前一次的值 clientv3.WithPrevKV()
	//查看相同前缀的
	if putRes, err = kv.Put(context.TODO(), "/cron/jobs/jobs1", "bytes", clientv3.WithPrevKV()); nil != err {
		log.Fatal(err)
	} else {
		fmt.Println(putRes.Header)       //cluster_id:14841639068965178418 member_id:10276657743932975437 revision:7 raft_term:3
		fmt.Println(putRes.PrevKv)       //nil
		fmt.Println(putRes.OpResponse()) //{0xc000352270 <nil> <nil> <nil>}
	}
}
