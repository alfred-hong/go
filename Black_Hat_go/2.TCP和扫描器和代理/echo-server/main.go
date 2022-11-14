package main

import (
	"io"
	"log"
	"net"
)

//仅回显接受到到数据
func echo(conn net.Conn) {
	defer conn.Close()
	//创建一个缓冲区来存储接收到到数据
	b := make([]byte, 512)
	for {
		//通过conn.Read接收到缓冲区
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client disconnected")
			break
		}
		if err != nil {
			log.Println("Unexpected err")
			break
		}
		log.Printf("Received %d bytes: %s\n", size, string(b))

		//通过conn.Write发送数据
		log.Println("Writing data")
		if _, err := conn.Write(b[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}

	}
}

func main() {
	// 在所有接口上绑定TCP端口20080
	listener, err := net.Listen("tcp", ":20080")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}
	log.Println("Listening on 0.0.0.0:20080")
	for {
		//等待连接。在已建立到连接上创建net.Conn
		conn, err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		//处理连接
		go echo(conn)
	}
}