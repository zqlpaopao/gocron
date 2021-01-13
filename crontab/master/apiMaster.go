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

//删除任务接口
// POST /job/delete name=job1
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		name  string
		old   *common.Job
		bytes []byte
	)

	//1、解析参数
	if err = r.ParseForm(); nil != err {
		goto ERR
	}

	name = r.PostForm.Get("name")

	//2.删除job
	if old, err = GJobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	//3.正常应答
	if bytes, err = common.BuildResponse(0, "success", old); nil == err {
		w.Write(bytes)
	}
ERR:
	//4.异常应答
	if bytes, err = common.BuildResponse(0, "fail", old); nil != err {
		w.Write(bytes)
	}
}

//列举所有任务
// POST /job/delete name=job1
func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		jobs  []*common.Job
		err   error
		bytes []byte
	)

	//1.获取
	if jobs, err = GJobMgr.ListJob(); nil != err {
		goto ERR
	}

	//2.返回
	if bytes, err = common.BuildResponse(0, "success", jobs); nil == err {
		w.Write(bytes)
	}

ERR:
	if bytes, err = common.BuildResponse(-1, "fail", jobs); nil != err {
		w.Write(bytes)
	}
}

//杀死任务接口
// POST /job/delete name=job1
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)
	if err = r.ParseForm(); nil != err {
		goto ERR
	}

	//要杀死的任务名
	name = r.PostForm.Get("name")

	//杀死任务
	if err = GJobMgr.KillJob(name); nil != err {
		goto ERR
	}

	//2.返回
	if bytes, err = common.BuildResponse(0, "success", nil); nil == err {
		w.Write(bytes)
	}

ERR:
	if bytes, err = common.BuildResponse(-1, "fail", nil); nil != err {
		w.Write(bytes)
	}
	//ETCDCTL_API=3 ./etcdctl watch "/cron/killer/" --prefix
	//监控目录变化
	/*
		PUT
		/cron/killer/job1

		DELETE
		/cron/killer/job1

	*/

}

//初始化服务
func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listen        net.Listener
		httpServer    *http.Server
		staticDir     http.Dir
		staticHandler http.Handler
	)

	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	//静态文件 index.html 匹配最长的
	staticDir = http.Dir(GConfig.Webroot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler)) // ./webroot/index.html

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
