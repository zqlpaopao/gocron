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
		//putOp clientv3.Op
		//getop clientv3.Op
		//opRes clientv3.OpResponse
		lease         clientv3.Lease
		leaseGrantRes *clientv3.LeaseGrantResponse
		leaseId       clientv3.LeaseID
		keepResChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepRes       *clientv3.LeaseKeepAliveResponse
		ctx           context.Context
		canelFunc     context.CancelFunc
		txn           clientv3.Txn
		tsnResp       *clientv3.TxnResponse
	)

	conf = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	//建立连接
	if client, err = clientv3.New(conf); nil != err {
		log.Fatal(err)
	}

	//lease实现锁自动过期
	//op操作
	//txn 事务

	//1、上锁--创建租约，自动续租，拿着租约取抢占一个key
	//申请租约
	lease = clientv3.NewLease(client)
	//申请5s的租约
	if leaseGrantRes, err = lease.Grant(context.TODO(), 5); nil != err {
		log.Fatal(err)
	}

	//拿到租约id
	leaseId = leaseGrantRes.ID

	//准备取消续租的context
	ctx, canelFunc = context.WithCancel(context.TODO())

	//确保函数推出后，自动停止续租
	defer canelFunc()

	//立即释放租约
	defer func() {
		if _, err = lease.Revoke(context.TODO(), leaseId); nil != err {
			fmt.Println(err)
		}

	}()

	//自动续租
	if keepResChan, err = lease.KeepAlive(ctx, leaseId); nil != err {
		log.Fatal(err)
	}

	//处理应答续租的协程
	go func() {
		for {
			select {
			case keepRes = <-keepResChan:
				if keepResChan == nil {
					fmt.Println("租约失效")
					goto END
				} else { //每1s一次
					fmt.Println("收到续租应答")
				}
			}
		}
	END:
	}()

	//如果锁设置不成功，then 设置它，else 设置失败
	kv = clientv3.NewKV(client)

	txn = kv.Txn(context.TODO())

	//原子性
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9")) //否则抢锁失败

	// 提交事务
	if tsnResp, err = txn.Commit(); nil != err {
		log.Fatal(err)
	}

	//是否抢到锁
	if !tsnResp.Succeeded {
		fmt.Println("锁被抢占", string(tsnResp.Responses[0].GetResponseRange().Kvs[0].Value))
	}

	//2、处理任务
	fmt.Println("处理事务")
	time.Sleep(5 * time.Second)

	//3、defer 会释放租约和取消续约

}
