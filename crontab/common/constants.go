package common

const (
	JobSaveDir = "/cron/jobs/" //保存目录

	JobKillDir = "/cron/killer/" //删除目录

	JobLockDir = "/cron/lock/" //任务锁目录

	JobSaveEvent = 1 //保存任务事件

	JobDeleteEvent = 2 //删除任务事件
)
