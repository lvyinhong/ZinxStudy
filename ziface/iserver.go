package ziface

// 定义一个服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()

	//路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)

	// 获取当前server 的链接管理器
	GetConnMgr() IConnManager

	// 注册OnConnStart 钩子函数的方法
	SetOnConnStart(func(conn IConnection))
	// 注册OnConnStop 钩子函数的方法
	SetOnConnStop(func(conn IConnection))

	// 调用OnConnStart 钩子函数的方法
	CallOnConnStart(conn IConnection)
	// 调用OnConnStop 钩子函数的方法
	CallOnConnStop(conn IConnection)
}
