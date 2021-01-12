package common

import "encoding/json"

type Job struct {
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//HTTP返回
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
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
