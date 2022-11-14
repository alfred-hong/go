package main

//改进echo

import (
	"io"
	"log"
	"net"
)

//仅回显接受到到数据
func echo(conn net.Conn) {
	defer conn.Close()
	//使用io.Copy()将数据从io.Rader复制到io.Writer
	if _, err := io.Copy(conn, conn); err != nil {
		log.Fatalln("Unable to read/write data")
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
