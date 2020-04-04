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

type HelloZinxRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *PingRouter)Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	// 先读取客户端的数据，再回写 ping, ping, ping...
	fmt.Println("recv from client: msgID=", request.GetMsgID(),
		", data=", string(request.GetData()))
	err:=request.GetConnection().SendMsg(200, []byte("ping... ping... ping..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn ziface.IConnection)  {
	fmt.Println("====>DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnection BEGIN")); err != nil {
		fmt.Println(err)
	}
}

//链接断开之前需要执行的函数
func DoConnectionLost(conn ziface.IConnection)  {
	fmt.Println("====>DoConnectionLost is Called ...")
	fmt.Println(" conn ID = ", conn.GetConnID())
}


// Test Handle
func (this *HelloZinxRouter)Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	// 先读取客户端的数据，再回写 ping, ping, ping...
	fmt.Println("recv from client: msgID=", request.GetMsgID(),
		", data=", string(request.GetData()))
	err:=request.GetConnection().SendMsg(201, []byte("hello zinx router"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1、创建一个server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.8]")

	//2、 注册链接Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//3、增加路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	//4、启动server
	s.Serve()
}
