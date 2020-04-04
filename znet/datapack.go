package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"llvvlv00.org/zinx/utils"
	"llvvlv00.org/zinx/ziface"
)

// 封包，拆包的具体模块

type DatePack struct {}

//拆包封包实例的一个初始化方法
func NewDataPack() *DatePack {
	return &DatePack{}
}
//获取包的头的长度方法
func(db *DatePack)GetHeadLen() uint32 {
	//Datalen uint32(4个字节） + ID uint32(4个字节)
	return 8

}

//封包方法
func(db *DatePack)Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen 写进databuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 将MsgId写进databuf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将data数据写进databuf中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(), nil
}

//拆包方法（将包的Head信息都读出来）之后再根据Head信息里的data长度，再进行一次读取
func(db *DatePack)Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息，将得到datalen和MsgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	//读取MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断datalen 是否已经超出了我们允许的最大包长度
	if (utils.GlobalObject.MaxPackageSize > 0) && (msg.DataLen > utils.GlobalObject.MaxPackageSize) {
		return nil, errors.New("too Large msg data recv!")
	}

	return msg, nil

}