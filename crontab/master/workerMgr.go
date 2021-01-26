package master

import (
	"context"
	"crontabInit/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"time"
)

///cron/workers/
type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	GWorkerMage *WorkerMgr
)

//获取在县wprker列表
func (w *WorkerMgr) ListWorkers() (workerArr []string, err error) {
	var (
		getResp  *clientv3.GetResponse
		kv       *mvccpb.KeyValue
		workerIp string
	)

	workerArr = make([]string, 0)

	//获取目录下所有kv
	if getResp, err = w.kv.Get(context.TODO(), common.JobWorkerDir, clientv3.WithPrefix()); err != nil {
		return nil, err
	}

	//解析每个节点
	for _, kv = range getResp.Kvs {
		//kv.Key : cron/workers/127.0.0.1
		workerIp = common.ExtractWorkerIp(string(kv.Key))
		workerArr = append(workerArr, workerIp)
	}

	return
}

func InitWorkerMgr() (err error) {
	var (
		client *clientv3.Client
		config clientv3.Config
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	//初始化配置
	config = clientv3.Config{
		Endpoints:   GConfig.EtcdEndpoints,
		DialTimeout: time.Duration(GConfig.EtcdDialTimeout) * time.Millisecond, //5s
	}

	//建立连接
	if client, err = clientv3.New(config); nil != err {
		return
	}

	//KV和lease的api子集
	kv = clientv3.KV(client)
	lease = clientv3.Lease(client)
	GWorkerMage = &WorkerMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}
}
