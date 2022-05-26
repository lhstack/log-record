package server

import (
	bytes2 "bytes"
	"fmt"
	"log-record/utils"
	"net"
)

//Conn 包装链接
type Conn struct {
	Conn   net.Conn
	Attach map[string]any
}

var server *Server

type Server struct {
	listenPort        string
	receiveBufferSize int
	listener          net.Listener
}

func init() {
	//获取监听端口
	listenPort := utils.EnvGetOrDefaultStringValue("LISTEN_PORT", "6330")
	receiveBufferSizeInt := utils.EnvGetOrDefaultIntValue("RECEIVE_BUFFER_SIZE", 128)
	server = &Server{receiveBufferSize: receiveBufferSizeInt, listenPort: listenPort}
	fmt.Printf("init server ,LISTEN_PORT %v,RECEIVE_BUFFER_SIZE %v \n", listenPort, receiveBufferSizeInt)
}

func Stop() {
	err := server.listener.Close()
	if err != nil {
		return
	}
	fmt.Println("stop server successful")
}

func Start(hand func(buffer *bytes2.Buffer, conn *Conn) error) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", server.listenPort))
	if err != nil {
		fmt.Println("listener port failure ", err)
	}
	server.listener = listen
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go accept(conn, hand)
		}
	}()
	fmt.Println("start server successful")
}

func accept(conn net.Conn, hand func(buffer *bytes2.Buffer, conn *Conn) error) {
	buffer := bytes2.NewBuffer([]byte{})
	packConn := &Conn{Conn: conn, Attach: make(map[string]any, 1)}
	for {
		bytes := make([]byte, server.receiveBufferSize)
		readCount, err := conn.Read(bytes)
		if err != nil {
			fmt.Println(err)
			_ = conn.Close()
			break
		}
		if readCount == -1 {
			_ = conn.Close()
			break
		}
		if readCount > 0 {
			buffer.Write(bytes[0:readCount])
			err := hand(buffer, packConn)
			if err != nil {
				fmt.Println(err)
				_ = conn.Close()
				break
			}
		}
	}
}
