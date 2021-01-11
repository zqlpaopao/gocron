package master

import (
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
func handleJobSave(w http.ResponseWriter, r *http.Request) {

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
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      mux,
	}

	GapiServer = &ApiServer{httpServer: httpServer}

	//启动服务端
	go httpServer.Serve(listen)

	return
}
