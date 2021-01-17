package worker

import (
	"crontabInit/common"
	"os/exec"
	"time"
)

//任务执行器
type Executor struct {
}

var (
	GExecutor *Executor
)

//执行一个任务
func (e *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *common.JobExecuteResult
			jobLock *JobLock
		)

		result = &common.JobExecuteResult{
			Executeinfo: info,
			OutPut:      make([]byte, 0),
		}

		///初始化分布式锁
		jobLock = GJobMgr.CreateJobLock(info.Job.Name)

		//上锁
		//解决分布式集群节点间时间同步不准确问题,节点间误差最多1ms
		//time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		err = jobLock.TryLock()
		defer jobLock.Unlock() //释放锁

		if nil != err {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			result.StartTime = time.Now()
			//执行shell命令
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)

			//执行并捕获输出
			output, err = cmd.CombinedOutput()

			//任务结束时间
			result.EndTime = time.Now()
			result.OutPut = output
			result.Err = err

			//任务执行完毕，把执行结果给scheduler，scheduler会从executing中删除执行记录
			Gscheduler.PushJobResult(result)
		}

	}()
}

//初始化执行器
func InitExecutor() {
	GExecutor = &Executor{}
	return
}
