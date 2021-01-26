package common

const (
	JobSaveDir = "/cron/jobs/" //保存目录

	JobKillDir = "/cron/killer/" //删除目录

	JobLockDir = "/cron/lock/" //任务锁目录

	//服务注册目录
	JobWorkerDir = "/cron/workers/"

	JobSaveEvent = 1 //保存任务事件

	JobDeleteEvent = 2 //删除任务事件

	JobKillerEvent = 3 //杀死任务
)
