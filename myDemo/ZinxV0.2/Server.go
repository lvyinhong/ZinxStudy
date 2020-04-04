package main

import (
	"fmt"
	"llvvlv00.org/zinx/ziface"
	"llvvlv00.org/zinx/znet"
)
// 基于zinx框架开发的服务器端应用程序

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test PreRouter
func (this *PingRouter)PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ...\n"))
	if err != nil {
		fmt.Printf("call back before ping error")
	}
}

// Test Handle
func (this *PingRouter)Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("ping ... ping ... ping ...\n"))
	if err != nil {
		fmt.Printf("call back ping error")
	}
}

// Test PostHandler
func (this *PingRouter)PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_,err := request.GetConnection().GetTCPConnection().Write([]byte("after ping ...\n"))
	if err != nil {
		fmt.Printf("call back after ping error")
	}
}

func main() {
	//1、创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.2]")
	s.AddRouter(&PingRouter{})
	//2、启动server
	s.Serve()
}
