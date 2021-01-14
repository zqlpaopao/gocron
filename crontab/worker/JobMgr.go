package worker

import (
	"context"
	"crontabInit/common"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//任务管理器
type JobMge struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	GJobMgr *JobMge
)

//监听任务变化
func (jobMge *JobMge) watcherJobs() (err error) {
	var (
		getRes             *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchRes           clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobName            string
		jobEvent           *common.JobEvent
	)

	//1.get /cron/jobs 目录下的所有任务，并且获知当前集群的revision
	if getRes, err = jobMge.kv.Get(context.TODO(), common.JobSaveDir, clientv3.WithPrefix()); nil != err {
		return
	}

	//for range
	for _, kvpair = range getRes.Kvs {
		//反序列化
		if job, err = common.UnpackJob(kvpair.Value); nil == err {
			//把这个任务同步给scheduler（调度协程）
			jobEvent = common.BuildJobEvent(common.JobSaveEvent, job)
			//启动后将etcd中的事件同步给Scheduler
			Gscheduler.PushJobEvent(jobEvent)
		}
	}

	//2.从该revision向后监听事件
	go func() {
		//从当前版本的的后续版本开始监听
		watchStartRevision = getRes.Header.Revision + 1
		watchChan = jobMge.watcher.Watch(context.TODO(), common.JobSaveDir, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		//处理监听
		for watchRes = range watchChan {
			for _, watchEvent = range watchRes.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: //保存事件
					if job, err = common.UnpackJob(watchEvent.Kv.Value); nil != err {
						continue
					}
					//构造event事件
					jobEvent = common.BuildJobEvent(common.JobSaveEvent, job)
					fmt.Println(jobEvent)

				case mvccpb.DELETE: //删除事件
					//delete /cron/jobs/job10
					jobName = common.ExtractJobName(string(watchEvent.Kv.Key))

					job = &common.Job{Name: jobName}

					//构造喊出event
					jobEvent = common.BuildJobEvent(common.JobDeleteEvent, job)
					fmt.Println(jobEvent)

				}
				//无论删除还是添加都需要推给schdluer
				Gscheduler.PushJobEvent(jobEvent)

			}
		}
	}()

	return
}

//初始化管理器
func InitJobMge() (err error) {
	var (
		client  *clientv3.Client
		config  clientv3.Config
		kv      clientv3.KV
		lease   clientv3.Lease
		watcher clientv3.Watcher
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
	watcher = clientv3.NewWatcher(client)

	//赋值
	GJobMgr = &JobMge{
		client:  client,
		kv:      kv,
		lease:   lease,
		watcher: watcher,
	}

	err = GJobMgr.watcherJobs()
	return
}
