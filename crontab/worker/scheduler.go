package worker

import (
	"crontabInit/common"
	"fmt"
	"time"
)

//任务调度
type Scheduler struct {
	jobEventChan      chan *common.JobEvent               //etcd中的任务事件
	jobPlanTable      map[string]*common.JobSchedulerPlan //任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo   //任务执行表
	jobResultChan     chan *common.JobExecuteResult       //任务结果队列
}

var (
	Gscheduler *Scheduler
)

//处理任务事件
func (s *Scheduler) handlerJobEvent(jobEv *common.JobEvent) {
	var (
		jobSchedulerPlan *common.JobSchedulerPlan
		jobExisted       bool
		err              error
	)

	switch jobEv.EventType {
	case common.JobSaveEvent: //保存事件
		if jobSchedulerPlan, err = common.BuildJobSchedulerPlan(jobEv.Job); nil != err {
			return
		}
		s.jobPlanTable[jobEv.Job.Name] = jobSchedulerPlan
	case common.JobDeleteEvent: //删除事件
		//如果存在就删除
		if jobSchedulerPlan, jobExisted = s.jobPlanTable[jobEv.Job.Name]; jobExisted {
			delete(s.jobPlanTable, jobEv.Job.Name)
		}
	}
}

//重新计算任务调度状态
func (s *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		now      time.Time
		nearTime *time.Time
		jobPlan  *common.JobSchedulerPlan
	)

	//如果任务表为空
	if len(s.jobPlanTable) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}
	//当前时间
	now = time.Now()

	//1.遍历所有任务
	for _, jobPlan = range s.jobPlanTable {
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//尝试执行任务，可能上一个任务还没执行完
			s.TryStartJob(jobPlan)
			fmt.Println("执行任务", jobPlan.Job.Name)
			jobPlan.NextTime = jobPlan.Expr.Next(now)

		}
		//统计一个最近要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	//下次调度间隔 最近要执行的任务调度时间-当前时间，就是下次要调度的时间
	if nearTime != nil {
		scheduleAfter = (*nearTime).Sub(now)
	}
	return
}

func (s *Scheduler) TryStartJob(jobPlan *common.JobSchedulerPlan) {
	//调度和执行是两个任务
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobexecuting   bool
	)

	//执行的任务可能会执行很久，1分钟60次，可能上个
	if jobExecuteInfo, jobexecuting = s.jobExecutingTable[jobPlan.Job.Name]; jobexecuting {
		fmt.Println("任务还未退出", jobPlan.Job.Name)
		return
	}

	//构建执行状态信息
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	//保存执行状态
	s.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	//执行任务
	GExecutor.ExecuteJob(jobExecuteInfo)
	//todo
	fmt.Println("执行任务", jobPlan.Job.Name)
	fmt.Println("计划时间", jobExecuteInfo.PlanTime, jobExecuteInfo.RealTime)
}

//调度协程
func (s *Scheduler) schedulerLoop() {
	var (
		jobEvent       *common.JobEvent
		schedulerAfter time.Duration
		schedulerTimer *time.Timer
		jobResult      *common.JobExecuteResult
	)

	//初始化，查看最近的调度时间是多久
	schedulerAfter = s.TrySchedule()

	//创建定时器
	schedulerTimer = time.NewTimer(schedulerAfter)

	//定时任务common.Job
	for {
		select {
		case jobEvent = <-s.jobEventChan: //监听任务变化事件
			//对内存中的任务列表做增删改查
			s.handlerJobEvent(jobEvent)
		case <-schedulerTimer.C: //最近的任务到期了
			fmt.Println(456)
		case jobResult = <-s.jobResultChan: //监听任务执行结果
			s.handJobResult(jobResult)
		}
		//调度一次任务
		schedulerAfter = s.TrySchedule()
		//重置调度间隔
		schedulerTimer.Reset(schedulerAfter)

	}
}

//处理任务执行结果
func (s *Scheduler) handJobResult(jobResult *common.JobExecuteResult) {
	//删除正在执行的任务
	delete(s.jobExecutingTable, jobResult.Executeinfo.Job.Name)
	fmt.Println(string(jobResult.OutPut), jobResult.Err)
}

//推送任务变化事件
func (s *Scheduler) PushJobEvent(job *common.JobEvent) {
	s.jobEventChan <- job
}

//初始化调度器
func InitScheduler() (err error) {
	Gscheduler = &Scheduler{
		jobEventChan:      make(chan *common.JobEvent, 1000),
		jobPlanTable:      make(map[string]*common.JobSchedulerPlan),
		jobExecutingTable: make(map[string]*common.JobExecuteInfo),
		jobResultChan:     make(chan *common.JobExecuteResult, 1000),
	}
	go Gscheduler.schedulerLoop()
	return
}

//回传任务结果
func (s *Scheduler) PushJobResult(jobRes *common.JobExecuteResult) {
	s.jobResultChan <- jobRes
}
