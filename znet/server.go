package znet

import (
	"fmt"
	"llvvlv00.org/zinx/utils"
	"llvvlv00.org/zinx/ziface"
	"net"
)

// iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	// 服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int

	//当前Server的消息管理模块，用来绑定msgId和对应的业务处理API关系
	MsgHandler ziface.IMessageHandler

}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s listenner at Ip: %s, Port: %d is srarting",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	fmt.Printf("[Zinx] Version %s, MaxConn:%d, MaxPackageSize: %d\n",
		utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)

	fmt.Printf("[Start] Server Listenner at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 0 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

		// 1、获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error: ", err)
			return
		}

		// 2、监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err", err)
			return
		}
		fmt.Println( "start Zinx server success ", s.Name, ", Listenning ...")

		var cid uint32
		cid = 0
		// 3、阻塞的等待客户端链接，处理客户端链接业务(读写)
		for {
			//如果有客户端链接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 将处理新链接的业务方法和conn进行绑定，得到我们的链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			// 启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

//停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源，状态、一些已经开辟的链接信息进行停止或者回收
}

//运行服务器
func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {

	}
}

// 注册路由
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success !")
}

//初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:utils.GlobalObject.Name,
		IPVersion:"tcp4",
		IP: utils.GlobalObject.Host,
		Port:utils.GlobalObject.TcpPort,
		MsgHandler:NewMsgHandle(),
	}
	return s
}