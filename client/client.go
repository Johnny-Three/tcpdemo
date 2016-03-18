package main

import (
	. "demogo/tcpdemo/logs"
	"demogo/tcpdemo/protocol"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Walkday struct {
	Walkdate  int64  `json:"walkdate"`
	Walkhour  string `json:"walkhour"`
	Walktotal int    `json:"walktotal"`
	Recipe    string `json:"recipe"`
}

type Walkdata struct {
	Userid    int       `json:"userid"`
	Timestamp int64     `json:"timestamp"`
	Walkdays  []Walkday `json:"walkdays"`
}

func send(conn net.Conn) {

	var total int

	n, err := conn.Write(protocol.Enpack(&protocol.Message{"javaserver@client", 0}))

	if err != nil {
		total += n
		fmt.Printf("write %d bytes, error:%s\n", n, err)
		os.Exit(1)
	}
	total += n
	fmt.Printf("write regist %d bytes this time, %d bytes in total\n", n, total)

	var total0 int

	n0, err0 := conn.Write(protocol.Enpack(&protocol.Message{"heartbeat", 1}))

	if err0 != nil {
		total0 += n0
		fmt.Printf("write %d bytes, error:%s\n", n0, err0)
		os.Exit(1)
	}
	total0 += n0
	fmt.Printf("write heartbeat %d bytes this time, %d bytes in total\n", n0, total0)

}

func HandleRead(conn net.Conn) {

	// buffer ..
	tmpBuffer := make([]byte, 0)

	//get packet
	readerChannel := make(chan protocol.Message, 1024)
	go reader(conn, readerChannel)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {

			if err != io.EOF {

				Logger.Debug(conn.RemoteAddr().String(), " connection error: ", err)
				conn.Close()
				return
			}
		}

		tmpBuffer = protocol.Depack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}
	defer conn.Close()
}

func reader(conn net.Conn, readerChannel chan protocol.Message) {
	for {
		select {

		case data := <-readerChannel:

			switch data.MsgType {

			case 0:
				fmt.Println("zero")
			case 1:

				n, err := conn.Write(protocol.Enpack(&protocol.Message{"heartbeat", 1}))
				conn.SetReadDeadline(time.Now().Add(time.Duration(2) * time.Second))

				if err != nil {
					fmt.Printf("write %d bytes, error:%s\n", n, err)
					os.Exit(1)
				}
				time.Sleep(1 * time.Second)

			case 2:
				fmt.Println("two")
			default:
				fmt.Println("weird happens")
			}

		}
	}
}

func main() {

	server := "localhost:6080"

	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	go HandleRead(conn)
	send(conn)

	var s Walkdata

	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1452873600,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1452960000,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})

	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1453046400,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1453132800,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1455897600,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1455984000,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})
	s.Walkdays = append(s.Walkdays, Walkday{Walkdate: 1456070400,
		Walktotal: 13000,
		Walkhour:  "32,0,0,0,0,0,3000,544,0,696,492,673,1219,15,0,0,938,4000,359,0,1148,6321,3941,67",
		Recipe:    "3790,3,3"})

	s.Timestamp = 1455724804

	for i := 0; i < 100000; i++ {

		total := 0

		for i := 454080; i < 454393; i++ {

			s.Userid = i

			b, err := json.Marshal(s)
			if err != nil {
				fmt.Println("json err:", err)
			}

			_, err1 := conn.Write(protocol.Enpack(&protocol.Message{string(b), 2}))
			if err1 != nil {
				fmt.Println("in for run ", err1)
				os.Exit(1)
			}
			time.Sleep(time.Duration(2) * time.Second)

			total += 1
			fmt.Println("total send msg is ", total)

		}
	}
}
