//通讯协议处理
package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Message struct {
	MsgContent string
	MsgType    int //0:register;1:heartbeat;2:message
}

//注册消息
const (
	Register = "register"
	RgLen    = len(Register)
)

//心跳消息
const (
	HeartBeat = "heartbeat"
	HBLen     = len(HeartBeat)
)

//通讯消息
const (
	MsgHeader  = "message"
	HeaderLen  = len(MsgHeader)
	SaveMsgLen = 4
)

//封包
func Enpack(msg *Message) []byte {

	var rmsg []byte

	switch msg.MsgType {

	case 0:
		rmsg = append(append([]byte(Register), IntToBytes(len(msg.MsgContent))...), []byte(msg.MsgContent)...)
	case 1:
		rmsg = []byte(HeartBeat)
	case 2:
		rmsg = append(append([]byte(MsgHeader), IntToBytes(len(msg.MsgContent))...), []byte(msg.MsgContent)...)
	default:
		fmt.Println("weird happens")
	}

	return rmsg

}

//解包
func Depack(buffer []byte, readerChannel chan Message) []byte {

	length := len(buffer)
	var i int

	for i = 0; i < length; i = i + 1 {

		//解压各种包..心跳包、通信消息包、注册包
		if length < i+len(HeartBeat) {
			//fmt.Printf("i is %d ,remain %s \n", i, buffer[i:len(buffer)])
			break
		}
		if string(buffer[i:i+HeaderLen]) == MsgHeader {

			//如果正好碰到message，但是解不出来后4位。。出去。。
			if length < i+HeaderLen+SaveMsgLen {
				break
			}

			messageLength := BytesToInt(buffer[i+HeaderLen : i+HeaderLen+SaveMsgLen])
			if length < i+HeaderLen+SaveMsgLen+messageLength {
				break
			}

			data := buffer[i+HeaderLen+SaveMsgLen : i+HeaderLen+SaveMsgLen+messageLength]
			i = i + HeaderLen + SaveMsgLen + messageLength - 1
			storemsg := Message{MsgType: 2, MsgContent: string(data)}
			readerChannel <- storemsg

		} else if string(buffer[i:i+HBLen]) == HeartBeat {

			i = i + HBLen - 1
			//fmt.Println(string(buffer[i : i+HBLen]))
			storemsg := Message{MsgType: 1, MsgContent: HeartBeat}
			readerChannel <- storemsg

		} else if string(buffer[i:i+RgLen]) == Register {

			if length < i+RgLen+SaveMsgLen {
				break
			}
			messageLength := BytesToInt(buffer[i+RgLen : i+RgLen+SaveMsgLen])
			if length < i+RgLen+SaveMsgLen+messageLength {
				break
			}

			data := buffer[i+RgLen+SaveMsgLen : i+RgLen+SaveMsgLen+messageLength]
			i = i + RgLen + SaveMsgLen + messageLength - 1

			storemsg := Message{MsgType: 0, MsgContent: string(data)}
			readerChannel <- storemsg

		} else {

			fmt.Printf("Unknown Msg!\n")
			break
		}

		//fmt.Println("in depack msg is ", storemsg)
	}

	if i == length {
		return make([]byte, 0)
	}
	return buffer[i:]
}

//整形转换成字节
func IntToBytes(n int) []byte {
	x := int32(n)

	buf := new(bytes.Buffer)
	//bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, x)
	return buf.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}
