package socket

import (
	. "demogo/tcpdemo/logs"
	"demogo/tcpdemo/protocol"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

var conseq int32
var Upload_chan chan string

type netConn struct {
	Seq  int32
	Des  string
	Conn net.Conn
}

var newConn netConn
var Map_connection map[string]netConn

func CheckError(err error) {
	if err != nil {
		Logger.Critical(err)
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Remove1(slc *[]netConn, item int) {
	s := *slc
	s = append(s[:item], s[item+1:]...)
	*slc = s
}

func init() {

	Upload_chan = make(chan string, 1024)
	Map_connection = make(map[string]netConn)

	go func() {
		netListen, err := net.Listen("tcp", "localhost:6080")
		CheckError(err)
		defer netListen.Close()

		newConn = netConn{}

		Logger.Debug("Waiting for clients")
		var index int
		for {
			conn, err := netListen.Accept()
			if err != nil {
				continue
			}

			Logger.Debug(conn.RemoteAddr().String(), " tcp connect success")
			go handleConnection(conn, newConn, index)
		}

	}()
}

func handleConnection(conn net.Conn, newConn netConn, index int) {

	// 缓冲区，存储被截断的数据
	tmpBuffer := make([]byte, 0)

	//接收解包
	readerChannel := make(chan protocol.Message, 1024)
	fmt.Printf("%d connection connected into server\n", index)
	go reader(conn, newConn, readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {

			if err != io.EOF {

				for _, value := range Map_connection {

					if value.Conn == conn {

						//从MAP中干掉失效连接
						Logger.Criticalf("connector:[%s],ip:[%s],connection lost", value.Des, conn.RemoteAddr().String())
						delete(Map_connection, value.Des)
					}
				}
				break
			}
		}
		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)

	}
	defer conn.Close()
}

func reader(conn net.Conn, newConn netConn, readerChannel chan protocol.Message) {
	for {
		select {

		case data := <-readerChannel:

			switch data.MsgType {

			case 0:

				Logger.Debug(conn.RemoteAddr().String(), "receive regist string: ", data.MsgContent)
				checkserver := strings.Split(data.MsgContent, "@")
				if checkserver != nil {

					newConn.Des = checkserver[0]
					atomic.AddInt32(&conseq, 1)
					newConn.Seq = conseq
					newConn.Conn = conn
					Map_connection[newConn.Des] = newConn
				}
			/*
				新来的注册client，需要Server先发送心跳包，开始双方之间的aliveCheck，同时启动SetDeadline，
				如超时未收到消息，则关闭链接
				客户端先注册，注册后Server向客户端发送友好心跳，客户端收到心跳后需要在5秒内回复Server，否则
				Server认为此链接失效，将断开此链接..
			*/
			case 1:
				conn.Write(protocol.Enpack(&protocol.Message{"heartbeat", 1}))
				conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))

			case 2:
				//fmt.Println(data.MsgContent)
				Decode(data.MsgContent)
				//todo.. write back ok msg..

			default:
				fmt.Println("weird happens")
			}
		}
	}
}
