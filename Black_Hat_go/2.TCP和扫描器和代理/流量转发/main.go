package main

import (
	"io"
	"log"
	"net"
)

func handle(src net.Conn) {
	dst, err := net.Dial("tcp", "joescatcam.website:80")
	if err != nil {
		log.Fatalln("Unable to connect to our unreachable host")
	}
	defer dst.Close()

	//在goroutine中运行防止io.Copy 被阻塞
	go func() {
		//将源到输出复制到目标
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()

	//将目标到输出复制回源
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}

//代理在80端口上监听并连接到joescatcam.website:80收发到所有数据
func main() {
	//在本地80上监听
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalln("Unable to bind to port")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go handle(conn)
	}
}
