package master

import (
	"context"
	"crontabInit/common"
	"encoding/json"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
	"time"
)

//任务管理器
type JobMge struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	GJobMgr *JobMge
)

//初始化管理器
func InitJobMge() (err error) {
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

	//赋值
	GJobMgr = &JobMge{
		client: client,
		kv:     kv,
		lease:  lease,
	}
	return
}

//保存任务
func (jobMge *JobMge) SaveJob(job *common.Job) (old *common.Job, err error) {
	//把任务保存到/cron/jobs/任务名
	var (
		jobKey   string
		jobValue []byte
		putRes   *clientv3.PutResponse
		oldJob   common.Job
	)

	//etcd保存key
	jobKey = common.JobSaveDir + job.Name

	//任务信息json
	if jobValue, err = json.Marshal(job); nil != err {
		return
	}

	if putRes, err = jobMge.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); nil != err {
		return
	}

	//如果是更新，返回旧值
	if putRes.PrevKv != nil {
		if err = json.Unmarshal(putRes.PrevKv.Value, &oldJob); nil != err {
			err = nil
			return
		}
		old = &oldJob
	}

	return
}

//删除任务
func (jobMge *JobMge) DeleteJob(name string) (old *common.Job, err error) {
	var (
		jobKey    string
		deleteRes *clientv3.DeleteResponse
		deletePre common.Job
	)

	//etcd删除key
	jobKey = common.JobSaveDir + name

	//删除
	if deleteRes, err = GJobMgr.kv.Delete(context.TODO(), jobKey, clientv3.WithPrevKV()); nil != err {
		return
	}

	if len(deleteRes.PrevKvs) > 0 {

		if err = json.Unmarshal(deleteRes.PrevKvs[0].Value, &deletePre); nil != err {
			err = nil
			return
		}
		old = &deletePre
	}
	return
}

//任务列表
func (jobMge *JobMge) ListJob() (old []*common.Job, err error) {
	var (
		dirKey  string
		getResp *clientv3.GetResponse
		kvPair  *mvccpb.KeyValue
		job     *common.Job
	)

	//任务目录
	dirKey = common.JobSaveDir

	//获取目录下所有任务
	if getResp, err = GJobMgr.kv.Get(context.TODO(), dirKey, clientv3.WithPrefix()); nil != err {
		return
	}

	//遍历所有任务
	for _, kvPair = range getResp.Kvs {
		job = &common.Job{}
		if err = json.Unmarshal(kvPair.Value, job); nil != err {
			continue
		}
		old = append(old, job)
	}
	return
}
