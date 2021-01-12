package master

import (
	"crontabInit/common"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

//任务http接口
type ApiServer struct {
	httpServer *http.Server
}

//定义单利对象
var (
	GapiServer *ApiServer
)

//保存任务接口
//job = {"name":"job1","command":"echo hello","cronExpr":"* * * * * *"}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	//保存到etcd中
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	//1.解析表单数据
	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	//2.读取表单数据
	postJob = r.PostForm.Get("job")

	//3.反序列化
	if err = json.Unmarshal([]byte(postJob), &job); nil != err {
		goto ERR
	}

	//4.保存到etcd
	if oldJob, err = GJobMgr.SaveJob(&job); nil != err {
		fmt.Println(err)
		goto ERR
	}

	//5.返回正常应答 {"errno":0,"msg":"ok","data":{...}}
	if bytes, err = common.BuildResponse(0, "success", oldJob); nil == err {
		w.Write(bytes)
	}

ERR:
	//6.返回异常应答
	if bytes, err = common.BuildResponse(-1, "fail", nil); nil != err {
		w.Write(bytes)
	}
	return
}

//初始化服务
func InitApiServer() (err error) {
	var (
		mux        *http.ServeMux
		listen     net.Listener
		httpServer *http.Server
	)

	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	//监听端口
	if listen, err = net.Listen("tcp", ":"+strconv.Itoa(GConfig.ApiPort)); nil != err {
		return
	}

	//创建HTTP服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(GConfig.ApiReadTimeout) * time.Second,
		WriteTimeout: time.Duration(GConfig.ApiWriteTimeout) * time.Second,
		Handler:      mux,
	}

	GapiServer = &ApiServer{httpServer: httpServer}

	//启动服务端
	go httpServer.Serve(listen)

	return
}
