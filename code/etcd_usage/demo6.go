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
		conf          clientv3.Config
		client        *clientv3.Client
		err           error
		leads         clientv3.Lease
		leadsGrantres *clientv3.LeaseGrantResponse
		leadsId       clientv3.LeaseID
		kv            clientv3.KV
		putRes        *clientv3.PutResponse
		getRes        *clientv3.GetResponse
		keepResChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepRes       *clientv3.LeaseKeepAliveResponse
	)

	conf = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	//建立连接
	if client, err = clientv3.New(conf); nil != err {
		log.Fatal(err)
	}

	//申请租约
	leads = clientv3.NewLease(client)

	//申请一个10s的租约
	if leadsGrantres, err = leads.Grant(context.TODO(), 10); nil != err {
		log.Fatal(err)
	}
	//拿到租约id
	leadsId = leadsGrantres.ID

	//自动续租,启动协程，定期的取续租
	if keepResChan, err = leads.KeepAlive(context.TODO(), leadsId); nil != err {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case keepRes = <-keepResChan:
				if keepRes == nil {
					fmt.Println("租约失效")
					goto END
				} else {
					fmt.Println("收到续租应答", keepRes.ID)
				}
			}
		}
	END:
	}()

	//获取kv PI子集
	kv = clientv3.NewKV(client)

	//put一个key，让其与租约关联起来，实现10s自动过期
	if putRes, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leadsId)); nil != err {
		log.Fatal(err)
	}
	fmt.Println("写入成功", putRes.Header.Revision)

	//查看是否过期
	for {
		if getRes, err = kv.Get(context.TODO(), "/cron/lock/job1"); nil != err {
			log.Fatal(err)
		}

		if getRes.Count == 0 {
			fmt.Println("过期了")
			break
		}

		fmt.Println("还没过期", getRes.Kvs)
		time.Sleep(2 * time.Second)
	}
}
