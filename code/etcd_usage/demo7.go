package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	//"go.etcd.io/etcd/mvcc/mvccpb"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"log"
	"time"
)

func main() {
	var (
		conf   clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		//putRes *clientv3.PutResponse
		getRes         *clientv3.GetResponse
		watchStartRevi int64
		watcher        clientv3.Watcher
		watchChan      clientv3.WatchChan
		watchRes       clientv3.WatchResponse
		event          *clientv3.Event
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

	//模拟key的变化
	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/jobs7", "i am jobs7")

			kv.Delete(context.TODO(), "/cron/jobs/jobs7")
			time.Sleep(time.Second)
		}

	}()

	//先get当前值，并监听后续变化
	if getRes, err = kv.Get(context.TODO(), "/cron/jobs/jobs7"); nil != err {
		log.Fatal(err)
	}

	//现在key是存在的
	if len(getRes.Kvs) > 0 {
		fmt.Println("current", getRes.Kvs[0])
	}

	//创建监听版本,当前这个key的版本加1是监控开始版本
	watchStartRevi = getRes.Header.Revision + 1

	//创建监听器
	watcher = clientv3.NewWatcher(client)

	//启动监听
	fmt.Println("从该版本向后监听", watchStartRevi)
	watchChan = watcher.Watch(context.TODO(), "/cron/jobs/jobs7", clientv3.WithRev(watchStartRevi))

	//处理kv变化事件
	for watchRes = range watchChan {
		for _, event = range watchRes.Events {
			switch event.Type {
			case mvccpb.PUT:
				fmt.Println("修改为", string(event.Kv.Value), "revision", event.Kv.CreateRevision, event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了", "revision", event.Kv.ModRevision)

			}
		}
	}
}
