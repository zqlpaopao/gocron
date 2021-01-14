package worker

import (
	"crontabInit/common"
	"fmt"
	"time"
)

//任务调度
type Scheduler struct {
	jobEventChan chan *common.JobEvent               //etcd中的任务事件
	jobPlanTable map[string]*common.JobSchedulerPlan //任务调度计划表
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
			//todo 尝试执行任务，可能上一个任务还没执行完
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

//调度协程
func (s *Scheduler) schedulerLoop() {
	var (
		jobEvent       *common.JobEvent
		schedulerAfter time.Duration
		schedulerTimer *time.Timer
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

		}
		//调度一次任务
		schedulerAfter = s.TrySchedule()
		//重置调度间隔
		schedulerTimer.Reset(schedulerAfter)

	}
}

//推送任务变化事件
func (s *Scheduler) PushJobEvent(job *common.JobEvent) {
	s.jobEventChan <- job
}

//初始化调度器
func InitScheduler() (err error) {
	Gscheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),
		jobPlanTable: make(map[string]*common.JobSchedulerPlan),
	}
	go Gscheduler.schedulerLoop()
	return
}
