package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

type Job struct {
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//任务调度计划
type JobSchedulerPlan struct {
	Job      *Job                 //要调度的任务信息
	Expr     *cronexpr.Expression //解析好的cronexpr
	NextTime time.Time            //下次调度时间

}

//任务执行时间
type JobExecuteInfo struct {
	Job        *Job
	PlanTime   time.Time          //理论上的调度时间
	RealTime   time.Time          //实际的执行时间
	CancelCtx  context.Context    //任务command的context
	CancelFunc context.CancelFunc //用于取消command执行的cancel函数
}

//HTTP返回
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

//变化事件
type JobEvent struct {
	EventType int //save delete
	Job       *Job
}

//任务执行过结果
type JobExecuteResult struct {
	Executeinfo *JobExecuteInfo //执行状态
	OutPut      []byte          //脚本输出
	Err         error           //脚本错误原因
	StartTime   time.Time       //启动时间
	EndTime     time.Time       //结束时间
}

//任务执行日志
type JobLog struct {
	JobName      string `bson:"jobName"`      //任务名字
	Command      string `bson:"command"`      //脚本命令
	Err          string `bson:"err"`          //错误原因
	Output       string `bson:"output"`       //脚本输出
	PlanTime     int64  `bson:"planTime"`     //计划开始时间
	ScheduleTime int64  `bson:"scheduleTime"` //任务调度时间
	StartTime    int64  `bson:"startTime"`    //任务开始时间
	EndTime      int64  `bson:"endTime"`      //任务结束时间
}

//日志批次
type LogBatch struct {
	Logs []interface{} //多条
}

//应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {

	var response Response

	response.Data = data
	response.Errno = errno
	response.Msg = msg

	resp, err = json.Marshal(response)
	return
}

//反序列化
func UnpackJob(v []byte) (ret *Job, err error) {
	var job *Job
	job = &Job{}

	if err = json.Unmarshal(v, job); nil != err {
		return
	}
	ret = job
	return
}

//提取任务名
///cron/jobs/job10,job10
func ExtractJobName(jobK string) string {
	return strings.TrimPrefix(jobK, JobSaveDir)
}

//任务变化事件有两种： 更新事件，删除事件
func BuildJobEvent(ev int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: ev,
		Job:       job,
	}
}

//构造执行计划
func BuildJobSchedulerPlan(job *Job) (jobSchedulerPlan *JobSchedulerPlan, err error) {

	var (
		expr *cronexpr.Expression
	)

	//解析JOb的cron表达式
	if expr, err = cronexpr.Parse(job.CronExpr); nil != err {
		return
	}

	//生成任务调度计划对象
	jobSchedulerPlan = &JobSchedulerPlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()),
	}

	return
}

//构造执行状态信息
func BuildJobExecuteInfo(jobScheduler *JobSchedulerPlan) (JobExecuteInfos *JobExecuteInfo) {
	JobExecuteInfos = &JobExecuteInfo{
		Job:      jobScheduler.Job,
		PlanTime: jobScheduler.NextTime,
		RealTime: time.Now(), //真正执行时间
	}
	JobExecuteInfos.CancelCtx, JobExecuteInfos.CancelFunc = context.WithCancel(context.TODO())
	return

}

//从/cron/killer/job10 提取job10
func ExtractKillerName(killKey string) string {
	return strings.TrimPrefix(killKey, JobKillDir)
}

//获取workerip
func ExtractWorkerIp(killKey string) string {
	return strings.TrimPrefix(killKey, JobWorkerDir)
}
