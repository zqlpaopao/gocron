package worker

import (
	"context"
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
			cmd    *exec.Cmd
			err    error
			output []byte
			result *common.JobExecuteResult
		)

		result = &common.JobExecuteResult{
			Executeinfo: info,
			OutPut:      make([]byte, 0),
			StartTime:   time.Now(), //任务开时间
		}
		//执行shell命令
		cmd = exec.CommandContext(context.TODO(), "/bin/bash", "-c", info.Job.Command)

		//执行并捕获输出
		output, err = cmd.CombinedOutput()

		//任务结束时间
		result.EndTime = time.Now()
		result.OutPut = output
		result.Err = err

		//任务执行完毕，把执行结果给scheduler，scheduler会从executing中删除执行记录
		Gscheduler.PushJobResult(result)

	}()
}

//初始化执行器
func InitExecutor() {
	GExecutor = &Executor{}
	return
}
