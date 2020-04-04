package znet

import (
	"errors"
	"fmt"
	"io"
	"llvvlv00.org/zinx/utils"
	"llvvlv00.org/zinx/ziface"
	"net"
)

// 链接模块
type Connection struct {
	// 当前Conn隶属于哪个Server
	TcpServer ziface.IServer

	//当前链接的socket TCP套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的链接状态
	isClosed bool

	//告知当前链接已经退出/停止的 channel (由Reader告知Writer由异常)
	ExitChan chan bool

	// 消息的管理MsgID 和对应的业务API关系
	MsgHandler ziface.IMessageHandler

	//无缓冲的管道用于读写Goroutine直接的消息通信
	msgChan chan []byte
}

// 初始化链接模块的方法
func NewConnection(server ziface.IServer,conn *net.TCPConn, connID uint32,msgHandler ziface.IMessageHandler) *Connection {
	c:=&Connection{
		TcpServer:server,
		Conn: conn,
		ConnID:connID,
		MsgHandler:msgHandler,
		isClosed:false,
		msgChan: make(chan []byte),
		ExitChan:make(chan bool, 1),
	}

	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}


// 写消息的Goroutine，专门发送给客户端消息的模块
func (c *Connection)StartWriter() {
	fmt.Println("[Writer] Goroutine is running")
	defer fmt.Println("[Writer is exit] connID = ", c.GetConnID(), " , remote addr is", c.RemoteAddr().String())

	// 不断地阻塞等待channel的消息，进行写给客户端
	for {
		select {
			case data := <-c.msgChan:
				//有数据要写给客户端
				if _,err := c.Conn.Write(data); err != nil {
					fmt.Println("Send data error, ", err)
					return
				}
			case <-c.ExitChan:
				//代表reader已经退出，此时Writer也要退出
				return

		}
	}
}

func (c *Connection)StartReader(){
	fmt.Println("[Reader] Goroutine is running ...")
	defer fmt.Println("[Reader is exit] connID = ", c.GetConnID(), " , remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		//创建一个拆包解包的对象
		dp :=NewDataPack()

		//读取客户端的msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err :=io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error, ", err)
			return
		}
		//拆包，得到msgID 和msgDatalen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			break
		}

		//根据datalen 再次读取Data 放在msg.Data 中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
		   if _, err =io.ReadFull(c.GetTCPConnection(), data); err != nil {
		   		fmt.Println("read msg data error", err)
		   		break
		   }
		}
		msg.SetMsgData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn:c,
			msg:msg,
		}
		// 判断是否已经开启工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		}else {
			// 根据绑定好的msgID，找到对应处理api的msgHandler 并执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}



//启动链接，让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println(" Conn Start() .. ConnID = ", c.ConnID)

	// 启动从当前链接的读数据的业务
	go c.StartReader()

	//TODO 启动从当前链接写数据的业务
	go c.StartWriter()

	// 按照开发中传递进来的 创建链接之后需要调用的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStart(c)
}

//停止链接，结束当前链接的工作
func (c *Connection) Stop()  {
	fmt.Println(" Conn Stop() .. ConnID = ", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 调用开发者注册的销毁链接之前需要执行的处理业务，执行对应的Hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	if err := c.Conn.Close();err != nil {
		fmt.Println(" Conn Conn.Close() err ", err, " .. ConnID = ", c.ConnID)
	}

	// 告知Writer,该链接已经关闭
	c.ExitChan <- true

	// 将当前链接从connMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)

	return
}

// 获取当前链接的绑定socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return  c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection)GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态 IP port
func (c *Connection)RemoteAddr() net.Addr{
	return c.Conn.RemoteAddr()
}

//发送数据，将数据发送给远程的客户端
// 提供一个SendMsg 方法， 将我们要发送给客户端的数据先进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data 进行封包 MsgDataLen|MsgID|Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMessage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("pack error msg")
	}

	// 将数据发送给msgChan Writer读到 msgChan再发送给客户端
	c.msgChan <- binaryMsg
	return nil
}




