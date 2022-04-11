Go语言基础之Socket网络编程

现在的我们几乎每天都在使用互联网，但是你知道程序是如果通过网络互相通信吗？

> 描述: 相信大部分人通常是一知半解的，作为一个程序员👨‍💻‍，对于网络模型你应该了解，知道网络到底是怎么进行通信的，进行工作的，为什么服务器能够接收到请求，做出响应。这里面的原理应该是每个 Web 程序员应该了解的。

本章我们就一起来学习下Go语言中的网络编程，关于网络编程其实是一个很庞大的领域，本文只是简单的演示了如何使用net包进行TCP和UDP通信。

基础概念介绍

> 描述: 互联网的核心是一系列协议，总称为互联网协议(Internet Protocol Suite)，正是这一些协议规定了电脑如何连接和组网，并通过各种协议实现不同的功能, 下面简单介绍一些协议涉及的基础知识概念。

此处不得不老生重谈OSI七层网络模型 ，OSI是国际标准化组织 1984 提出的模型标准，简称 OSI(Open Systems Interconnection Model)，主要统一标准来规范网络协议让各个硬件厂商可以协同工作，相信如果你学习过网络工程或者操作系统方面的课程至少你是了解它的。每一层都有自己的功能，就像建筑物一样，每一层都靠下一层支持，通常用户接触到的只是最上面的那一层。

![WeiyiGeek.OSI七层网络与TCP/IP四层网络模型](https://i0.hdslb.com/bfs/article/6cb224ce36f3496d54218c62dbbb1f8fb821d9ca.png@942w_654h_progressive.png)


如上图所示按照不同的模型划分会有不用的分层，但是不论按照什么模型去划分，越往上的层越靠近用户，越往下的层越靠近硬件，在软件开发中我们使用最多的是上图中将互联网划分为五个分层的模型，即应用层、传输层、网络层、数据链路层、物理层。

下面简单介绍各层: (如有兴趣想深入学习的可以自行Google)

(1) 物理层 : 即通过我们计算机或其它设备通过网络硬件与外界互联网通信，它主要规定了网络的一些电气特性，作用是负责传送0和1的电信号(比特流)。

> 例如： 以太网、无线Lan、PPP、双绞线、光纤、无线。

(2) 数据链路层 : 确定了物理层传输的0和1的分组方式及代表的意义, 通过以太网(Ethernet)的协议规定一组电信号构成一个数据包，叫做帧(Frame)。

> 每一帧分成两个部分：标头(Head)和数据(Data)
> Head : 包含数据包的一些说明项，比如发送者(MAC地址)、接受者(MAC地址)、数据类型等等；(其长度固定为18字节)
> Data : 则是数据包的具体内容。(其长度，最短为46字节，最长为1500字节)
> Tips: 因此整个帧最短为64字节，最长为1518字节, 所以如果数据很长就必须分割成多个帧进行发送。

(3) 网络层 : 使得我们能够区分不同的计算机是否属于同一个子网络(子网)，该地址就叫做网络地址(即IP地址)，此时每台计算机有了两种地址，一种是MAC地址(硬件网卡唯一标识)，另一种是IP地址。

> IP地址: 则是网络管理员分配的，它可以帮助我们确定计算机所在的子网络，MAC地址则将数据包送到该子网络中的目标网卡。
> 根据IP协议发送的数据就叫做IP数据包，IP数据包也分为标头和数据两个部分：
> 标头部分主要包括版本、长度、IP地址等信息(长度为20到60字节)
> 数据部分则是IP数据包的具体内容，整个数据包的总长度最大为65535字节。

(4) 传输层 : 有了上述的MAC地址和IP地址可以就可以在互联网上任意两台主机上建立通信。但是如果是要与主机上某一程序进行通信我们还需要一个端口(Port)，从而便可以让两个程序通过网络进行收发数据。

> IP和端口我们就能实现唯一确定互联网上一个程序，进而实现网络间的程序通信, 此外可以选择常用的TCP或者UDP协议进行通信。
> 端口: 是0到65535之间的一个整数，正好16个二进制位。0到1023的端口被系统占用，用户只能选用大于1023的端口
> TCP (Transmission Control Protocol): 面向连接的、可靠的、基于字节流的传输层通信协议(经过三次握手四层挥手)，由IETF的RFC 793定义，TCP数据包没有长度限制，理论上可以无限长，但是为了保证网络的效率，通常TCP数据包的长度不会超过IP数据包的长度，以确保单个TCP数据包不必再分割。 (主要用于对通信的信息比较重要的场景)
> UDP (User Datagram Protocol) : 为应用程序提供了一种无需建立连接就可以发送封装的 IP 数据包的方法,我们可以持续的发送信息但并不关心其是否正常到达, 由IETF的 RFC 768 定义, 其数据包非常简单，”标头”部分一共只有8个字节，总长度不超过65,535字节，正好放进一个IP数据包。(主要用于视屏直播流、非实时性的控制指令发送的场景)

(5) 应用层 : ”应用层”的作用就是规定应用程序使用的数据格式，由于互联网是开放架构，数据来源五花八门，必须事先规定好通信的数据格式，否则接收方根本无法获得真正发送的数据内容。

> 例如: 我们TCP协议之上常见的Email、HTTP、FTP等协议，这些协议就组成了互联网协议的应用层。

如下图所示，发送方的HTTP数据经过互联网的传输过程中会依次添加各层协议的标头信息，接收方收到数据包之后再依次根据协议解包得到数据。

![WeiyiGeek.一张解释互联网的传输过程的图](https://i0.hdslb.com/bfs/article/c37c6278e244555f4f3c8906a383e838d49ec7af.png@942w_650h_progressive.png)


知识扩展

> Q:那么，发送者和接受者是如何标识呢？
> 答: 以太网规定，连入网络的所有设备都必须具有网卡接口。数据包必须是从一块网卡，传送到另一块网卡。网卡的地址，就是数据包的发送地址和接收地址，这叫做MAC地址。每块网卡出厂的时候，都有一个全世界独一无二的MAC地址，长度是48个二进制位，通常用12个十六进制数表示例如00-FF-81-D5-15-F8。前6个十六进制数是厂商编号，后6个是该厂商的网卡流水号, 此时有了MAC地址就可以定位网卡和数据包的路径了。

> Q: 有了MAC地址之后，如何把数据准确的发送给接收方呢？
> 描述: 首先通过ARP协议来获取接受方的MAC地址, 然后通过广播(broadcasting)的方式，向本网络内所有计算机都发送，让每台计算机读取这个包的标头，找到接收方的MAC地址，然后与自身的MAC地址相比较，如果两者相同就接受这个包，做进一步处理，否则就丢弃这个包。

> Q: 网络地址的协议?
> 描述: 规定网络地址的协议叫做IP协议,目前，广泛采用的是IP协议第四版，简称IPv4。IPv4这个版本规定，网络地址由32个二进制位组成，我们通常习惯用分成四段的十进制数表示IP地址，从0.0.0.0~255.255.255.255，当然里面包含三个私有网段，以及保留的网段，此处不细讲知道即可。



Socket 基础介绍

其实学习其它开发语言你将会发现, 基本高级语言中都有专门进行网络通信的包来提供两个程序的网络通信, 通常针对程序网络通信的开发都描述为Socket编程。

> Q: 什么是Socket编程?
> 描述: Socket(也称作”套接字”)是BSD UNIX的进程通信机制，用于描述IP地址和端口是一个通信链的句柄,Socket可理解为TCP/IP网络的API，它定义了许多函数或例程，程序员可以用它们来开发TCP/IP网络上的应用程序，所以在我们计算机运行的应用程序通常通过"套接字"向网络发出请求或者应答网络请求。

Socket 是应用层与传输层(TCP/IP协议族)通信的中间软件抽象层，在设计模式中Socket其实就是一个门面模式，它把复杂的TCP/IP协议族隐藏在Socket后面，对用户来说只需要调用Socket规定的相关函数，让Socket去组织符合指定的协议数据然后进行通信。

![WeiyiGeek.Socket图解](https://i0.hdslb.com/bfs/article/1feb05ddbbc8b3e3f7fcc78f48e3cafa07998c1f.png@942w_672h_progressive.png)


Go实现C/S端TCP通信

> 描述: TCP/IP(Transmission Control Protocol/Internet Protocol)即传输控制协议/网间协议，是一种面向连接(连接导向)的、可靠的、基于字节流的传输层(Transport layer)通信协议，因为是面向连接的协议，数据像水流一样传输，但它会存在黏包问题(将会在后续演示解决方案)。

1) TCP 服务端

> 描述: 通常一个TCP服务端可以同时连接很多个客户端，例如世界各地的用户使用自己电脑上的浏览器访问淘宝网、京东商城。

由于Go语言中创建多个goroutine实现并发非常方便和高效，所以我们可以每建立一次链接就创建一个goroutine去处理。

TCP服务端程序常规流程：
1.设置监听网络地址与端口
2.接收客户端请求建立链接
3.创建goroutine处理链接

在Go语言中的可以采用net包来实现的TCP服务端的创建,下述罗列出net包中使用的相关方法原型:

​	func net.Listen(network string, address string) (net.Listener, error) : 指定通信协议类型版本、本地网络地址和通信端口进行监听,注意 network参数值必须是“tcp”、“tcp4”、“tcp6”、“unix”或“unixpacket”之一。
​	func (net.Listener).Accept() (net.Conn, error) : 等待Listener对象并返回到侦听器的下一个连接。
​	func (net.Conn).Read(b []byte) (n int, err error) : 从服务端或者客户端连接中读取数据。


2) TCP 客户端

TCP客户端进行TCP通信的流程大致如下：
1.与服务端的建立链接
2.与服务端进行数据收发
3.关闭与服务端的链接

下述罗列出net包中Tcp客户端创建使用的相关方法原型:

​	func net.Dial(network string, address string) (net.Conn, error) : 客户端连接到指定网络上的地址, 同样 networks 参数可选值如下"tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".
​	func (net.Conn).Write(b []byte) (n int, err error) : 写入将数据写入连接即发送给连接的网络对象，并返回发送的字节数。（注意中文字符占3字节）
​	func (net.Conn).LocalAddr() net.Addr : 获取本地客户端连接到服务端的网络地址信息。

简单示例1:TCP Server端与Client端一次通信。

服务端: Server.go

```go
func main() {
	// 1.设置服务端监听端口
	address := "10.20.172.108:22022"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Start Tcp Server on %v Failed!\nerr：%v\n", address, err)
		return
	} else {
		fmt.Printf("Server Listen : %v\n", address)
	}

	// 2.等待客户端建立连接
	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("Accept failed,err: %v\n", err)
		return
	}
	defer conn.Close() // 程序结束时,关闭与客户端打开的TCP连接通道
	
	// 3.与客户端进行消息通信(读取客户端法过来的信息)
	var msg [1024]byte
	n, err := conn.Read(msg[:]) // 注意读取的类型是byte的slice
	if err != nil {
		fmt.Printf("Read from Client conn failed, err:%v\n", err)
		return
	}
	
	// 4.打印客户端发送的信息,注意此需要将[]byte类型的切片,转为字符串类型进行输出.
	fmt.Println(string(msg[:n]))

}
```


客户端: Client.go

```go
func main() {
	// 1.与Server端建立TCP通信链接
	address := "10.20.172.108:22022"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Connect Server failed!\nerr:%v\n", err)
		return
	}

	// 2.发送字符串数据到Server端
	sendMsg := "Hello World! Server, I'm client"
	conn.Write([]byte(sendMsg))
	
	// 3.关闭建立的TCP通信链接
	defer conn.Close()
}
```

将上面的代码保存之后分别编译成server和client可执行文件,具体操作如下:

```sh
➜  Server go build
➜  Server ./Server
Server Listen : 10.20.172.108:22022

➜  Client go build
➜  Client ./Client
```

进阶示例2: TCP Server端与多个Client端持续通信。
服务端 Server.go 

```go
package main
import (
	"bufio"
	"fmt"
	"io"
	"net"
)
func SendReceiveProccess(conn net.Conn) {
	// 3.与客户端进行消息通信(循环读取客户端法过来的信息)
	defer conn.Close() // 关闭当前链接通信对象
	reader := bufio.NewReader(conn)
	var msg [1024]byte // 每次读取1024B
	for {
		n, err := reader.Read(msg[:]) // 注意读取的类型是byte的slice
		// 末尾标识
		if err == io.EOF {
			fmt.Printf("Close conn %v\n", conn.RemoteAddr())
			break
		}
		// 异常时break
		if err != nil {
			fmt.Printf("Read from Client conn failed, Close conn %v\n", conn.RemoteAddr())
			break
		}
		fmt.Println(conn.RemoteAddr(), "->", string(msg[:n]))

		// 将客户端发送的信息又转发给客户端(返回写入的字节数)
		_, err = conn.Write([]byte(msg[:n]))
		if err != nil {
			fmt.Printf("Send failed, Close Client conn %v\n", conn.RemoteAddr())
			break
		}
	}

}

func main() {
	// 1.设置监听端口
	address := "10.20.172.108:22022"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Start Tcp Server on %v Failed!\nerr：%v\n", address, err)
		return
	} else {
		fmt.Printf("Server Listen %v Start......\n", address)
	}
	defer listener.Close() // 关闭服务端监听

	// 2.等待客户端建立连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept failed,err: %v\n", err)
			return
		}
	// 不同的客户端利用Goroutine分配不同的线程进行响应。
		go SendReceiveProccess(conn)
	}

}
```


客户端 Client.go

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// 1.与Server端建立TCP链接
	address := "10.20.172.108:22022"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("Connect Server failed! \n [err]: %v\n", err)
		return
	} else {
		fmt.Printf("Connect Server %v successful!\n", address)
	}
	// 退出通信连接
	defer conn.Close()

	// 2.发送初始连接信息到Server端
	sendMsg := fmt.Sprintf("Hello Server, I'm  %v client.", conn.LocalAddr())
	inputReader := bufio.NewReader(os.Stdin) // 复习点
	conn.Write([]byte(sendMsg))
	
	// 3.循环从服务端接收以及从终端输入发送信息到服务端
	for {
		// 服务端回复的信息
		reply := [1024]byte{}
		n, err := conn.Read(reply[:])
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		fmt.Printf("Server > %v\n", string(reply[:n]))
	
		// 客户端输入字符串信息
		fmt.Printf("请输入消息:")
		sendMsg, _ = inputReader.ReadString('\n') // 复习点 (以\n截至读取)
		sendMsg = strings.TrimSpace(sendMsg)      // 复习点 处理输入字符串的前后空格
		sendMsg = strings.Trim(sendMsg, "\n")     // 复习点 处理输入字符串的最后的换行符
	
		// 但客户端输入quit则退出与服务端建立的TCP通信连接.
		if strings.ToUpper(sendMsg) == "QUIT" {
			fmt.Printf("exit conn.......")
			break
		}
	
		// 发送已经处理过后的信息到服务端
		_, err = conn.Write([]byte(sendMsg))
		if err != nil {
			return
		}
	}

}
```


Server.go与Client.go 的编译&运行&执行结果:

```sh
go build && ./Server
go build && ./Client
```

Go实现C/S端UDP通信

> 描述: UDP协议（User Datagram Protocol）中文名称是用户数据报协议，是OSI（Open System Interconnection，开放式系统互联）参考模型中一种无连接的传输层协议，不需要建立连接就能直接进行数据发送和接收，属于不可靠的、没有时序的通信，但是UDP协议的实时性比较好，通常用于视频直播相关领域。

Tips : UDP 通信相比较于 TCP 通信使用简单, 下述我将UDP服务端和客户端实现的相关方法原型进行展示。

​	func net.ListenUDP(network string, laddr *net.UDPAddr) (*net.UDPConn, error): 设置监听UDP相关的网络地址与端口信息、网络必须是UDP网络名称。
​	func (\*net.UDPConn).ReadFromUDP(b []byte) (int, \*net.UDPAddr, error): 接收conn对象里发送的信息, ReadFromUDP与ReadFrom类似，但返回一个UDPADD。
​	func (\*net.UDPConn).WriteToUDP(b []byte, addr \*net.UDPAddr) (int, error): 发送信息到conn对象里, WriteToUDP的行为类似于WriteTo，但使用UDPADD。
func net.DialUDP(network string, laddr *net.UDPAddr, raddr *net.UDPAddr) (*net.UDPConn, error) : 作用于与UDP服务端建立连接, DialUDP的作用类似于UDP网络的拨号。

进阶示例1.实现UDP服务端与多客户端进行连接通信!

Server.go

```go
package main
import (
	"fmt"
	"net"
	"strings"
	"time"
)
func main() {
	// 1.服务端开启监听 UDP 通信的相关设置
	server_ip := [4]byte{10, 20, 172, 108}
	server_port := 30000
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(server_ip[0], server_ip[1], server_ip[2], server_ip[3]),
		Port: server_port,
	})

	if err != nil {
		fmt.Printf("Listen UDP Server (%v:%v) Failed, err: %v\n", server_ip, server_port, err)
		return
	} else {
		fmt.Printf("[%v] Listening UDP Server %v:%v is successful!\n", time.Now().Format("2006-01-02 15:04:05"), server_ip, server_port)
	}
	// 程序结束时关闭conn资源
	defer conn.Close()
	
	// 2.循环接收和响应数据给客户端，非常主要此处不需要建立连接，直接收发数据。
	for {
		// 获得客户端通信对象以及返回读取的字节数
		var recvMsg [1024]byte
		count, addr, err := conn.ReadFromUDP(recvMsg[:]) // 接收数据
		if err != nil {
			fmt.Println("Read from UDP Client failed. Err:", err)
			return
		}
		// 打印客户端发送的信息到终端
		fmt.Printf("[%v] %v - %v\n", time.Now().Format("2006-01-02 15:04:05"), addr.String(), string(recvMsg[:count-1]))
	
		// 并将接收到的信息更改为大写，再返还给Client。
		reply := strings.ToUpper(string(recvMsg[:count]))
		conn.WriteToUDP([]byte(reply), addr) // 发送数据
	}

}
```


Client.go

```go
package main
import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)
func main() {
	// (1) 与服务端建立UDP通信链接
	server_ip := [4]byte{10, 20, 172, 108}
	server_port := 30000
	socket, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(server_ip[0], server_ip[1], server_ip[2], server_ip[3]),
		Port: server_port,
	})
	if err != nil {
		fmt.Printf("Connect UDP Server (%v:%v) Failed! err: %v\n", server_ip, server_port, err)
		return
	} else {
		fmt.Printf("[%v] - Connect UDP Server %v:%v successful!\n", time.Now().Format("2006-01-02 15:04:05"), server_ip, server_port)
	}
	// 同样关闭建立的通信连接
	defer socket.Close()

	// (2) 发送与接收服务端返回的信息
	var reply [1024]byte
	inputData := bufio.NewReader(os.Stdin)
	for {
		// 终端接收输入要发送的给服务端的内容
		fmt.Print("<- 请输入将要发送的内容:")
		sendMsg, _ := inputData.ReadString('\n')
		sendMsg = strings.TrimSpace(sendMsg) // 取消字符串前后的空格
		socket.Write([]byte(sendMsg))        // 发送数据
		if err != nil {
			fmt.Printf("发送数据失败，err: %v\n", err)
			return
		}
		// 接收来自服务端的反馈的内容
		count, _, err := socket.ReadFromUDP(reply[:]) // 接收数据
		if err != nil {
			fmt.Printf("接收数据失败, err: %v\n", err)
			return
		}
		fmt.Printf("Server -> [%v Bytes] %v \n", count, string(reply[:count]))
	}

}
```


执行结果:

![WeiyiGeek.实现UDP服务端与多客户端进行连接通信!](https://i0.hdslb.com/bfs/article/e6167e40f4d1877f01bd3de5eebb40382821cc93.png@942w_236h_progressive.png)


5.课外知识扩展

1) TCP 黏包

在讲解TCP黏包前我们先来看看TCP黏包会导致什么问题?

黏包示例服务端代码如下：

```go
func process(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var buf [1024]byte
	for {
		n, err := reader.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client发来的数据：", recvStr)
	}
}

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}
}
```

黏包示例客户端代码如下:

```go
func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		conn.Write([]byte(msg))
	}
}
```


将上面的代码保存后分别编译, 先启动服务端然后再启动客户端，可以看到服务端输出结果如下：

```sh
收到client发来的数据： Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?
收到client发来的数据： Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?
收到client发来的数据： Hello, Hello. How are you?Hello, Hello. How are you?
收到client发来的数据： Hello, Hello. How are you?Hello, Hello. How are you?Hello, Hello. How are you?
收到client发来的数据： Hello, Hello. How are you?Hello, Hello. How are you?
```


从客户端示例代码中可以看见客户端分20次发送的数据，但是在服务端并没有成功的输出20次，而是多条数据“粘”到了一起，这就是TCP黏包带来的问题。



> Q: 此时可能会问为什么会出现粘包?
>
> 答：主要原因就是tcp数据传递模式是流模式，在保持长连接的时候可以进行多次的收和发, 而TCP黏包可以发生在在发送端也可发生在接收端, 主要是由于Nagle算法导致的。
> Nagle算法是一种改善网络传输效率的算法， 通常在发送端由于Nagle算法导致的黏包问题,而接收端接收不及时也会造成的接收端粘包。
> 发送端: 简单来说就是当我们提交一段数据给TCP发送时，TCP并不立刻发送此段数据，而是等待一小段时间看看在等待期间是否还有要发送的数据，若有则会一次把这两段数据发送出去，所以说Nagle算法特性其并不适用于某些场景。
> 接收端: TCP会把接收到的数据存在自己的缓冲区中，然后通知应用层取数据。当应用层由于某些原因不能及时的把TCP的数据取出来，就会造成TCP缓冲区中存放了几段数据。

> Q: 有何解决办法?
>
> 答: 出现”粘包”的关键在于接收方不确定将要传输的数据包的大小，因此我们可以对数据包进行封包和拆包的操作。
> 此时需要自己定义一个协议,比如数据包的前4个字节为包头里面存储的是发送的数据的长度，然后通过发送端进行封包、接收端进行拆包的操作来解决此问题。

> Q: 什么是封包?
>
> 答: 封包就是给一段数据加上包头，这样一来数据包就分为包头和包体两部分内容了(过滤非法包时封包会加入”包尾”内容)。包头部分的长度是固定的，并且它存储了包体的长度，根据包头长度固定以及包头中含有包体长度的变量就能正确的拆分出一个完整的数据包。

实践示例：

自定义协议: proto.go

```go
package proto
import (
	"bufio"
	"bytes"
	"encoding/binary"
)
// (1) Encode 将消息编码
func Encode(message string) ([]byte, error) {
	// 1.读取消息的长度，转换成int32类型（占4个字节）以后可以按照需要进行自定义
	var length = int32(len(message))
	var pkg = new(bytes.Buffer) // 向系统为具有读写方法的字节大小可变的缓冲区申请内存。
	// 2.写入消息头
	err := binary.Write(pkg, binary.LittleEndian, length) //注意此处以小端的方式写入.在后续解包时也必须采用小端方式读取
	if err != nil {
		return nil, err
	}
	// 3.写入消息实体
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	// 4.返回封包完毕的缓冲区中数据
	return pkg.Bytes(), nil
}

// (2) Decode 解码消息
func Decode(reader *bufio.Reader) (string, error) {
	// 1.读取消息的长度
	lengthByte, _ := reader.Peek(4)           // 读取前4个字节的数据
	lengthBuff := bytes.NewBuffer(lengthByte) // NewBuffer使用buf作为初始内容创建并初始化一个新缓冲区,此处指定要读取数据的长度.
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// 2.Buffered返回缓冲中现有的可读取的字节数,如果获取的字节数小于消息的长度则说明数据包有误.
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// 3.读取真正的消息数据
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	
	// 4.利用slice切片返回四个字节后的消息数据
	return string(pack[4:]), nil

}
```


然后在服务端和客户端分别使用上面定义的proto包的Decode 和 Encode函数处理数据。

服务端代码如下：

```go
package main
import (
	"bufio"
	"fmt"
	"io"
	"net"
	"weiyigeek.top/studygo/Day08/03StickBag/proto"
)
func process(conn net.Conn) {
	// 4.退出时关闭conn资源
	defer conn.Close()
	// 5.NewReader返回其缓冲区具有默认大小的新读取器。
	reader := bufio.NewReader(conn)
	for {
		// 6.解包: 将通过conn对象中获取的缓冲区数据进行解包.
		msg, err := proto.Decode(reader)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println("decode msg failed, err:", err)
			return
		}
		// 7.打印解包后的数据
		fmt.Println("收到client发来的数据：", msg)
	}
}
func main() {
	// 1.设置TCP Server端监听地址和端口
	listen, err := net.Listen("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("listen failed, err:", err)
		return
	}
	// 2.函数结束后关闭监听资源
	defer listen.Close()

	// 3.循环接收客户端发送过来的数据,利用gorontine执行process任务
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:", err)
			continue
		}
		go process(conn)
	}

}
```

客户端代码如下:

```go
package main
import (
	"fmt"
	"net"
	"weiyigeek.top/studygo/Day08/03StickBag/proto"
)
func main() {
	// 1.连接到server端.
	conn, err := net.Dial("tcp", "127.0.0.1:30000")
	if err != nil {
		fmt.Println("dial failed, err", err)
		return
	}

	// 2.程序结束时关闭conn网络连接资源.
	defer conn.Close()
	
	// 3.循环发送20次msg给server端
	for i := 0; i < 20; i++ {
		msg := `Hello, Hello. How are you?`
		// 4.将要发送的信息进行封包处理
		data, err := proto.Encode(msg)
		if err != nil {
			fmt.Println("encode msg failed, err:", err)
			return
		}
		// 5.将处理过的封包进行发送
		conn.Write(data)
	}

}
```

执行结果: 此时发现它不会将多个hello...字符串放在一行了，并且输出也是我们预定的20次。

![WeiyiGeek.黏包解决办法](https://i0.hdslb.com/bfs/article/584083aa59ce90a7599c34f088c212e87bc66bf2.png@942w_563h_progressive.png)


补充知识:字节序列的存储格式之大端(Big-endian)和小端(LittleEndian)存储

​	Big-endian：将高序字节存储在起始地址（高位编址），此种方式便于人类理解。
​	Little-endian：将低序字节存储在起始地址（低位编址）一般在x64/x32的系统中都是小端存储。
举个例子: 如果我们将0x1234abcd写入到以0x0000开始的内存中，则结果为；


注：每个地址存1个字节，2位16进制数是1个字节（0xFF=11111111）；

![WeiyiGeek.大端和小端](https://i0.hdslb.com/bfs/article/45399474323d5b8a2136a3818b701d87a12db84e.png@942w_495h_progressive.png)


补充知识:CPU存储一个字节的数据时其字节内的8个比特之间的顺序是否也有big endian和little endian之分？或者说是否有比特序的不同？

实际上，这个比特序是同样存在的只是多个两个表示名称（MSB 和 LSB）。
MSB的意思是：全称为Most Significant Bit，在二进制数中属于最高有效位，MSB是最高加权位，与十进制数字中最左边的一位类似。
LSB的意思是：全称为Least Significant Bit，在二进制数中意为最低有效位，一般来说，MSB位于二进制数的最左侧，LSB位于二进制数的最右侧。

下面以数字0xB4（10110100）用图加以说明。

```sh
# Big Endian

msb------------------------>lsb
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|   1  |   0  |   1  |   1  |   0  |   1  |   0  |   0  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+


# Little Endian

lsb-------------------------->msb
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|   0  |   0  |   1  |   0  |   1  |   1  |   0  |   1  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

实际上，由于CPU存储数据操作的最小单位是一个字节，其内部的比特序是什么样对我们的程序来说是一个黑盒子。也就是说，你给我一个指向0xB4这个数的指针，对于big endian方式的CPU来说，它是从左往右依次读取这个数的8个比特；而对于little endian方式的CPU来说，则正好相反，是从右往左依次读取这个数的8个比特。而我们的程序通过这个指针访问后得到的数就是0xB4，字节内部的比特序对于程序来说是不可见的，其实这点对于单机上的字节序来说也是一样的。

至此，Go语言中Socket网络编程学习完毕！