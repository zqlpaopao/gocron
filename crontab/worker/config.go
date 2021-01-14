package worker

import (
	"encoding/json"
	"io/ioutil"
)

//解析config
type Config struct {
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
}

var (
	GConfig *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	//读取配置文件
	//fmt.Println(filename)
	if content, err = ioutil.ReadFile(filename); nil != err {
		return
	}

	//2.反序列化
	if err = json.Unmarshal(content, &conf); nil != err {
		return
	}
	GConfig = &conf
	return
}
