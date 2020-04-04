package utils

import (
	"encoding/json"
	"io/ioutil"
	"llvvlv00.org/zinx/ziface"
	"os"
)

// 存储一切有关zinx框架的全局参数，供其他模块使用
// 一些参数是可以通过zinx.json 由用户进行配置

type GlobalObj struct {
	// server
	TcpServer ziface.IServer	// 当前Zinx全局的Server对象
	Host string					// 当前服务器珠玑监听的IP
	TcpPort int					// 当前服务器主机监听的端口号
	Name string					// 当前服务器名称


	// zinx
	Version string	//当前Zinx的版本好
	MaxConn int 	//当前服务器主机允许的最大连接数
	MaxPackageSize uint32 //当前Zinx框架数据包的最大值
	WorkerPoolSize uint32 //当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //Zinx框架允许用户最多开辟多少个worker
}

var GlobalObject *GlobalObj

// 提供一个init方法，初始化当前的GlobalObject
func init() {
	// 如果配置文件没有加载,默认的值
	GlobalObject = &GlobalObj{
		Name:"ZinxServerApp",
		Version:"V0.7",
		TcpPort:8999,
		Host:"0.0.0.0",
		MaxConn:1000,
		MaxPackageSize:4096,
		WorkerPoolSize:10,	// worker 工作池的队列的个数
		MaxWorkerTaskLen:1024, // 每个worker对应的消息队列的任务的数量最大值
	}

	//应该尝试从config/zinx.json去加载一些用户自定义的参数
	GlobalObject.Reload()
}



// 从zinx.json去加载用于自定义的参数
func (g *GlobalObj) Reload() {
	path, err :=os.Getwd()
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadFile(path + "/conf/zinx.json")
	if err != nil {
		panic(err)
	}

	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}
