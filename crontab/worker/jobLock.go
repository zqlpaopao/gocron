package worker

import (
	"context"
	"crontabInit/common"
	"go.etcd.io/etcd/clientv3"
)

//分布式锁（TXN事务）
type JobLock struct {
	kv    clientv3.KV
	lease clientv3.Lease

	jobName    string             //任务名
	cancelFunc context.CancelFunc //用于自动终止续租
	leaseId    clientv3.LeaseID
	isLocked   bool //是否上锁成功
}

//初始化任务锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (JobLocker *JobLock) {
	JobLocker = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

//尝试上锁
func (jl *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnRes         *clientv3.TxnResponse
	)

	//1.创建5秒租约
	if leaseGrantResp, err = jl.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	//租约id
	leaseId = leaseGrantResp.ID

	//取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	//2.自动续租
	if keepRespChan, err = jl.lease.KeepAlive(cancelCtx, leaseId); nil != err {
		goto FAIL
	}
	//3.处理续租应答的协程
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)

		for {
			select {
			case keepResp = <-keepRespChan: //自动续租的应答
				if keepResp == nil {
					goto END
				}

			}
		}
	END:
	}()

	//4.创建事务txn
	txn = jl.kv.Txn(context.TODO())

	//锁路径 /cron/lock
	lockKey = common.JobLockDir + jl.jobName

	//5.事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	//提交事务
	if txnRes, err = txn.Commit(); nil != err {
		goto FAIL
	}

	//6.成功返回，失败释放租约
	if !txnRes.Succeeded { //锁被占用
		err = common.ErrLockAlreadyRequired
		goto FAIL
	}

	//抢锁成功
	jl.leaseId = leaseId
	jl.cancelFunc = cancelFunc
	jl.isLocked = true
	return

FAIL:
	cancelFunc()                                      //取消自动续租
	_, err = jl.lease.Revoke(context.TODO(), leaseId) //主动取租约
	return
}

//释放锁
func (jl *JobLock) Unlock() {
	if jl.isLocked {
		jl.cancelFunc()                             //取消锁
		jl.lease.Revoke(context.TODO(), jl.leaseId) //释放租约
	}
}
