package znet

import (
	"fmt"
	"llvvlv00.org/zinx/utils"
	"llvvlv00.org/zinx/ziface"
	"strconv"
)

// 消息处理


//消息处理模块的实现
type MsgHandler struct {
	// 存放每个MsgID所对应的处理方法
	Apis map[uint32] ziface.IRouter

	ziface.IMessageHandler

	//负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的worker数量
	WorkerPoolSize uint32

}

func NewMsgHandle() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度/执行对应的Router消息处理方法
func(mh *MsgHandler) DoMsgHandler(request ziface.IRequest) {
	//1、从request中找到msgId  根据msgId 调度对应router业务即可
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND, need register!")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func(mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	//1、判断当前 msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id 已经注册了
		panic("repeat register api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//2、添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Printf("Add api msgId = ", msgID, "success !")
}

// 启动一个worker工作池
func (mh *MsgHandler)StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker， 每一个Worker用一个Go来承载
	for i:=0; i< int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1、当前的worker对应的channel消息队列 开辟空间,第0个worker就用第0个channel
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// 2、启动当前的Worker，阻塞等待消息从channel传递进来
		go mh.startOneWorker(i,mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandler)startOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started ...")
	// 不断的阻塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
			case request:= <- taskQueue:
				mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给taskQueue 由Worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1、将消息平均分配给对应的worker
	// 根据客户端建立的ConnID来进行分配, 更好的是给句RequestId来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize

	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID =", request.GetMsgID(),
		"to WorkerID = ", workerID)

	// 2、将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
