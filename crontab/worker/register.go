package worker

import (
	"context"
	"crontabInit/common"
	"go.etcd.io/etcd/clientv3"
	"net"
	"time"
)

//注册戒掉到etcd： /cron/worker/IP地址
type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP string //本机IP
}

var Gregister *Register

//注册到/cron/workers/IP,并自动续租
func (r *Register) keepLine() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		err            error
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse

		cancelCtx  context.Context
		cancelFunc context.CancelFunc
	)
	for {
		//注册地址
		regKey = common.JobWorkerDir + r.localIP
		cancelFunc = nil
		if leaseGrantResp, err = r.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}

		//自动续租
		if keepAliveChan, err = r.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
			goto RETRY
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		//注册到etcd
		if _, err = r.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}

		//处理续租应答
		for {
			select {
			case keepAliveResp = <-keepAliveChan:
				if keepAliveResp == nil { //续租失败
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}

}

//获取本季ip
func getLocalIp() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)

	//获取所有网卡
	if addrs, err = net.InterfaceAddrs(); nil != err {
		return
	}
	//取第一个非Lo的网卡IP
	for _, addr = range addrs {
		//这个网路地址是ipv4 或者ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			//跳过ipv6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() //127.0.0.1
				return
			}
		}
	}
	err = common.ErrNoLocalIpFound
	return
}

func InitRegister() (err error) {
	var (
		client  *clientv3.Client
		config  clientv3.Config
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
	)

	//初始化配置
	config = clientv3.Config{
		Endpoints:   GConfig.EtcdEndpoints,
		DialTimeout: time.Duration(GConfig.EtcdDialTimeout) * time.Millisecond, //5s
	}

	//建立连接
	if client, err = clientv3.New(config); nil != err {
		return
	}

	//KV和lease的api子集
	kv = clientv3.KV(client)
	lease = clientv3.Lease(client)

	//获取本机ip
	if localIp, err = getLocalIp(); nil != err {
		return
	}
	Gregister = &Register{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIp,
	}

	Gregister.keepLine()
	return

}
