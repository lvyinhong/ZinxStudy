package main

import (
	"fmt"
	"io"
	"llvvlv00.org/zinx/znet"
	"net"
	"time"
)

/**
 模拟客户端
 */

func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)

	// 1、直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println(" client start err, ", err, " exit!")
		return
	}

	for {
		// 发送封包的msg消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(0, []byte("ZinxV0.5 client Test Message")))
		if err != nil {
			fmt.Println("Pack error: ", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("send msg error: ", err)
			return
		}

		// 服务器应该给我们回复一个message数据，MsgID:1 ping... ping... ping...
		// 先读取流中的head部分，得到ID和dataLen

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error", err)
			break
		}

		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error", err)
			break
		}

		if msgHead.GetMsgLen() > 0 {
			// 再根据DataLen进行第二次读取，将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error ", err)
				return
			}
			fmt.Println("======>Recv server Msg: ID = ", msg.GetMsgId(), ", len = ", msg.GetMsgLen(), ", data = ", string(msg.GetMsgData()))
		}
		// cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
