NSQ分布式消息队列实践

描述: 目前比较流行的一个分布式的消息队列是RabbitMQ、ZeroMQ、Kafka(大项目中推荐)、NSQ，本章将着重介绍 NSQ 基础概念、安装部署及Go语言如何操作NSQ分布式消息队列，因为NSQ是采用Go语言进行开发使用的。



NSQ 快速了解

Q: 什么是NSQ?
答: NSQ 是一个基于Go语言的分布式实时消息平台, 具有分布式，易于水平扩展，易于安装，易于集成（主流语言都有对应的客户端库）的特点。

其主要核心概念

> Topic： 在生产者publish时会创建topic，一个topic就是程序发布消息的一个逻辑键。
> Channels： 通道组与消费者相关，是消费者之间的负载均衡，channel在某种意义上来说是一个“队列”。每当一个发布者发送一条消息到一个topic，消息会被复制到所有消费者连接的channel上，消费者通过这个特殊的channel读取消息，实际上，在消费者第一次订阅时就会创建channel。（Channel会将消息进行排列，如果没有消费者读取消息，消息首先会在内存中排队，当量太大时就会被保存到磁盘中）
> Messages: 消息构成了我们数据流的中坚力量，消费者可以选择结束消息，表明它们正在被正常处理，或者重新将他们排队待到后面再进行处理。每个消息包含传递尝试的次数，当消息传递超过一定的阀值次数时，我们应该放弃这些消息，或者作为额外消息进行处理。
> NSQ 的优势：

安装运行简单: 易于配置和部署，并且内置了管理界面。
协议简单: NSQ 支持多种语言客户端接入，其有一个快速的二进制协议，通过短短的几天工作量就可以很简单地实现这些协议，我们还自己创建了我们的纯JS驱动（当时只存在coffeescript驱动）
在线扩容：NSQ 支持横向扩展，没有任何集中式代理。
分布式 : 提倡分布式和分散的拓扑，没有单点故障，支持容错和高可用性，并提供可靠的消息交付保证
NSQ 的特性：

> 持久化模式方案: 采用的方式时内存+硬盘的模式，当内存到达一定程度时就会将数据持久化到硬盘, 如果设置了--mem-queue-size=0则所有的消息将会存储到磁盘。
> 队列中的每条消息至少传递一次。
> 队列中消息不保证有序的。

NSQ 四个重要组件构成：
(1) nsqd：一个负责接收、排队、转发消息到客户端的守护进程，它可以独立运行，不过通常它是由 nsqlookupd 实例所在集群配置的, 其默认监听端口4150和4151。

```sh
# nsqd 简单启动示例
./nsqd -broadcast-address=192.168.1.2:4160
	# -broadcast-address 配置广播地址

# 如果是在搭配nsqlookupd使用的模式下需要还指定nsqlookupd地址，如果是部署了多个nsqlookupd节点的集群，那还可以指定多个-lookupd-tcp-address。
./nsqd -broadcast-address=192.168.1.2 -lookupd-tcp-address=192.168.1.2:4160
```

> nsqlookupd：管理拓扑信息并提供最终一致性的发现服务的守护进程,值得注意其数据并不是持久化保存,也不需要与任何其他nsqlookupd实例协调以满足查询, 因此根据你系统的冗余要求尽可能多地部署nsqlookupd节点(通常三个), 其默认监听端口4160和4161。
> nsqadmin：它是一套实时监控集群状态、执行各种管理任务的Web管理平台, 默认监听端口4171。

```sh
# 启动 nsqadmin 示例
./nsqadmin -lookupd-http-address=192.168.1.2:4161
```

> utilities：常见基础功能、数据流处理工具，如 nsq_stat、nsq_tail、nsq_to_file、nsq_to_http、nsq_to_nsq、to_nsq

NSQ 架构及其工作模式
我们首先开看NSQ工作模式图，值得非常注意的一点就是在非集群模式下可以直接连接指定的nsqd, 而集群模式下则通过nsqlookup查询到nsqd地址再连接。

![img](https://i0.hdslb.com/bfs/article/489696a8bb5b1bc9ea74e55025855bdbd0ce2e50.png@942w_713h_progressive.png)!


上图中每个nsqd实例旨在一次处理多个数据流, 该数据流称为"topics",并且topic与channels是1对多的关系, 每个channel都会收到topic所有消息的副本，实际上下游的服务是通过对应的channel来消费topic消息。

> topic 在首次使用时创建，方法是将其发布到指定topic，或者订阅指定topic上的channel
> channel 是通过订阅指定的channel在第一次使用时创建的。
> topic 和 channel 都相互独立地缓冲数据，防止缓慢的消费者导致其他chennel的积压（同样适用于topic级别），但是channel可以并且通常会连接多个客户端。

假设所有连接的客户端都处于准备接收消息的状态，则每条消息将被传递到随机客户端，如下图所示:

![img](https://i0.hdslb.com/bfs/article/39a7b1d8d71d22d0b55794336093c372bf67132d.gif)


总而言之消息是从topic -> channel（每个channel接收该topic的所有消息的副本）多播的，但是从channel -> consumers均匀分布（每个消费者接收该channel的一部分消息）。



NSQ 消息队列的应用场景
(1) 异步处理: 我们可以利用消息队列把业务流程中的非关键流程异步化，从而显著降低业务请求的响应时间。

![img](https://i0.hdslb.com/bfs/article/a24a86b67abe751fe7faa6ac484e3195b684a15f.png@942w_677h_progressive.png)

(2) 应用解耦: 通过使用消息队列将不同的业务逻辑解耦，降低系统间的耦合，提高系统的健壮性，后续有其他业务要使用订单数据可直接订阅消息队列，提高系统的灵活性。

![img](https://i0.hdslb.com/bfs/article/8eda641a922cabe43f0a37f04ce5a097c2d62883.png@942w_540h_progressive.png)


(3) 流量削峰: 在类似秒杀（大秒）等场景下，某一时间可能会产生大量的请求，使用消息队列能够为后端处理请求提供一定的缓冲区，保证后端服务的稳定性，例如:秒杀请求 --Write--> 消息队列 --根据规则读取--> 请求处理。

(4) 消息通信: 消息队列一般都内置了高效的通信机制，因此也可以用在纯的消息通讯, 例如实现点对点消息队列，或者聊天室进行消息发布和接收等。

// # 点对点客户端A -->> 消息队列 <<-- 客户端B// # 聊天室订阅主题进行消息发布和接收客户端A <<-->> 消息队列 <<-->> 客户端B

例如: NSQ接收和发送消息流程如下图所示。

![img](https://i0.hdslb.com/bfs/article/a975d1a026cfafebe20db91823f0266347b52690.png@870w_425h_progressive.png)


NSQ 官网地址: https://nsq.io/



NSQ 安装配置

从NSQ官方下载页面(https://nsq.io/deployment/installing.html), 根据自己的平台下载并解压到指定目录, 然后设置环境变量即可。

此处使用Docker方式安装部署测试
实践环境说明:

```sh
$ docker --version
Docker version 19.03.15, build 99e3ed8919
$ docker-compose --version
docker-compose version 1.25.0, build unknown

nsq v1.2.1
```


步骤01.首先创建一个 docker-compose.yml 存放了容器运行配置清单。

```sh
# 注意操作系统中是否安装 docker-compose, 如没有安装请执行yum或者apt安装即可。
$ vim docker-compose.yml
version: '2'
services:
  nsqlookupd:
    container_name: nsqlookupd
    image: nsqio/nsq
    command: /nsqlookupd
    ports:
      - "4160:4160"
      - "4161:4161"
  nsqd:
    image: nsqio/nsq
    container_name: nsqd
    command: /nsqd --lookupd-tcp-address=nsqlookupd:4160
    ports:
      - "4150:4150"
      - "4151:4151"
  nsqadmin:
    image: nsqio/nsq
    container_name: nsqadmin
    command: /nsqadmin --lookupd-http-address=nsqlookupd:4161
    ports:
      - "4171:4171"
```


Tips:从 上面的docker-compose.yml文件可以看到 nsqd服务 需要注册到 nsqlookupd 的 4160 端口, 而 nsqadmin服务 需要注册到 nsqlookupd 的 4161 端口

步骤02.在该yml文件同级目录下执行如下命令进行创建并后台运行容器。

```sh
$ docker-compose up -d
  # Creating network "docker_default" with the default driver
  # Creating nsqd       ... done
  # Creating nsqadmin   ... done
  # Creating nsqlookupd ... done
$ docker ps
  # CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                                                            NAMES
  # 6ff1c7b396fe        nsqio/nsq           "/nsqadmin --lookupd…"   8 seconds ago       Up 6 seconds        4150-4151/tcp, 4160-4161/tcp, 4170/tcp, 0.0.0.0:4171->4171/tcp   nsqadmin
  # 33eeda6384e5        nsqio/nsq           "/nsqlookupd"            8 seconds ago       Up 6 seconds        4150-4151/tcp, 4170-4171/tcp, 0.0.0.0:4160-4161->4160-4161/tcp   nsqlookupd
  # ae4830fd10e1        nsqio/nsq           "/nsqd --lookupd-tcp…"   8 seconds ago       Up 6 seconds        4160-4161/tcp, 0.0.0.0:4150-4151->4150-4151/tcp, 4170-4171/tcp   nsqd
```


步骤03.访问nsqadmin提供的消息队列监控的Web管理平台 http://10.10.107.225:4171/lookup。


![img](https://i0.hdslb.com/bfs/article/80d647ce0bd2f61f6fdf94b4775b65439dd15e9d.png@942w_599h_progressive.png)至此安装完毕，通过docker来部署NSQ是非常简单的。





NSQ 实践操作

1.go-nsq 安装

描述: NSQ官方为了开发者提供了Go语言版的客户端go-nsq(https://github.com/nsqio/go-nsq)，更多客户端支持请查看CLIENT LIBRARIES(https://nsq.io/clients/client_libraries.html)。

go-nsq库安装命令如下所示:

```sh
➜  src cd weiyigeek.top
➜  weiyigeek.top go get -u github.com/nsqio/go-nsq
  # go: downloading github.com/nsqio/go-nsq v1.1.0
  # go: downloading github.com/golang/snappy v0.0.1
  # go: downloading github.com/golang/snappy v0.0.4
  # go get: added github.com/golang/snappy v0.0.4
  # go get: added github.com/nsqio/go-nsq v1.1.0
```




2.go-nsq 简单使用

生产者Producer
简单的生产者示例代码如下：
```go
// studygo/Day09/NSQ/demo1/producer.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nsqio/go-nsq"
)

// NSQ Producer(生产者)示例演示
var Producer *nsq.Producer

type NsqProducer struct {
	nsqd_host string
	nsqd_port int
}

// 初始化NSQ生产者
func (NP NsqProducer) InitProducer() (*nsq.Producer, error) {
	// NewConfig返回一个新的默认nsq配置。
	config := nsq.NewConfig()
	// 组合nsqd服务连接地址。
	nsqdAddr := fmt.Sprintf("%s:%d", NP.nsqd_host, NP.nsqd_port)
	Producer, err := nsq.NewProducer(nsqdAddr, config)
	if err != nil {
		fmt.Printf("create producer failed, err:%v\n", err)
		return nil, err
	}
	return Producer, nil
}

func main() {
	// 1.实例化以及初始化
	nsqd := &NsqProducer{
		nsqd_host: "10.10.107.225",
		nsqd_port: 4150,
	}
	Producer, err := nsqd.InitProducer()
	if err != nil {
		fmt.Printf("Init producer failed, err:%v\n", err)
		return
	} else {
		log.Printf("Init Producer success!")
	}

	// 2.从标准输入读取
	fmt.Println("请输入你要向Topic_Demo消息队列传递的消息:")
	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("read string from stdin failed, err:%v\n", err)
			continue
		}
		data = strings.TrimSpace(data)
		// 3.当前输入Q或者q时退出程序
		if strings.ToUpper(data) == "Q" {
			break
		}
		// 4.向 'Topic_Demo' publish 数据
		err = Producer.Publish("Topic_Demo", []byte(data))
		if err != nil {
			fmt.Printf("publish msg to nsq failed, err:%v\n", err)
			continue
		}
	}
}
```

消费者-Consumer
简单的消费者示例代码如下：
```go
package main
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	nsq "github.com/nsqio/go-nsq"
)

// NSQ Consumer(消费者) 示例演示
// NsqConsumer 是一个消费者类型结构体
type NsqConsumer struct {
	Title string
}

// HandleMessage 是需要实现的处理消息的方法
func (m *NsqConsumer) HandleMessage(msg *nsq.Message) (err error) {
	fmt.Printf("%s : recv from %v, msg:%v\n", m.Title, msg.NSQDAddress, string(msg.Body))
	return
}

// 初始化Consumer消费者
func (NC *NsqConsumer) InitConsumer(topic string, channel string, address string) (err error) {
	// 1.NewConfig返回一个新的默认nsq配置
	config := nsq.NewConfig()

	// 2.查找轮询间隔此处设置15s
	config.LookupdPollInterval = 15 * time.Second
	
	// 3.NewConsumer为指定的主题/频道创建新的Consumer实例
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		fmt.Printf("create consumer failed, err:%v\n", err)
		return
	}
	
	// 4.AddHandler为此使用者接收的消息设置处理程序，可以多次调用该函数以添加其他处理程序，注意其必须放在连接NSQD和NSQLOOKUP前。
	c.AddHandler(NC)
	
	// 5.两种方式连接到NSQD
	// if err := c.ConnectToNSQD(address); err != nil {   // 直接连NSQD，tcp长连接
	if err := c.ConnectToNSQLookupd(address); err != nil { // 通过lookupd查询连NSQD，更易于分布式容错和高可用
		return err
	}
	return nil
}

func main() {
	// 实例化消费者
	consumer := &NsqConsumer{
		Title: "NSQ_USE",
	}

	// 初始化连接NSQD进行获取消息队列中的值
	err := consumer.InitConsumer("Topic_Demo", "channel_first", "10.10.107.225:4161")
	if err != nil {
		fmt.Printf("init consumer failed, err:%v\n", err)
		return
	}
	
	// 定义一个信号的通道
	c := make(chan os.Signal)
	// 转发键盘中断信号到c
	signal.Notify(c, syscall.SIGINT)
	// 通道输出消息队列中的值阻塞
	<-c
}
```

执行结果: 从结果可以看到当开启多个消费者时会一个发送一次消息队列中的信息。
```sh
# 生产者
➜  demo1 go run .
2021/12/28 21:09:40 Init Producer success!
请输入你要向Topic_Demo消息队列传递的消息:
Whoami
2021/12/28 21:21:04 INF    1 (10.10.107.225:4150) connecting to nsqd
WeiyiGeek
...
topic测试
2021/12/28 21:21:56 INF    1 (10.10.107.225:4150) connecting to nsqd
channel
test

# 消费者01
➜  demo2 go run .
2021/12/28 21:21:37 INF    1 [Topic_Demo/channel_first] (ae4830fd10e1:4150) connecting to nsqd
NSQ_USE : recv from ae4830fd10e1:4150, msg:Whoami
NSQ_USE : recv from ae4830fd10e1:4150, msg:WeiyiGeek
NSQ_USE : recv from ae4830fd10e1:4150, msg:topic测试
NSQ_USE : recv from ae4830fd10e1:4150, msg:channel

# 消费者02
➜  demo2 go run .
2021/12/28 21:22:19 INF    1 [Topic_Demo/channel_first] querying nsqlookupd http://10.10.107.225:4161/lookup?topic=Topic_Demo
2021/12/28 21:22:19 INF    1 [Topic_Demo/channel_first] (ae4830fd10e1:4150) connecting to nsqd
NSQ_USE : recv from ae4830fd10e1:4150, msg:test
```

![img](https://i0.hdslb.com/bfs/article/e94d88bf640f71698306937409e073f35c41daa3.png@942w_251h_progressive.png)

Tips: 在客户端执行是如果采用ConnectToNSQLookupd方法即通过lookupd查询连NSQD，需要在hosts绑定对应的容器hostname和宿主机地址（粗暴解决）。例如此处

```sh
$ cat /etc/hosts
127.0.0.1       localhost
127.0.1.1       Ubuntu-PC
10.10.107.225 ae4830fd10e1
```

Tips: 此处我们可以通过nsqdadmin提供的Web页面查看到我们Publish的topic，以及我们生产者向队列传递的值，和消费者从通道中接收到的值，点击页面上的Topic_Demo就能进入一个展示更多详细信息的页面, 而在/counter页面显示处理的消息数量。

![img](https://i0.hdslb.com/bfs/article/b33aaefc2cb6332aa7fdb52efb616a3a49f6e06c.png@942w_1020h_progressive.png)

Tips: 在/lookup界面支持创建topic和channel, 这是提供了一种在将服务部署到生产环境之前设置流层次结构的方法, 如果频道名称为空，则只创建主题。



3.go-nsq 直连方式

描述: 上面实践了通过nsqlookupd的http接口查询后长连接到nsqd, 本节将简单演示直连nsqd（tcp长连接）写法。

```go
package main
import (
  "flag"
  "log"
  "time"
  "github.com/nsqio/go-nsq"
)

func main() {
  go startConsumer()
  startProducer()
}

var url string
func init() {
  //nsqd 服务具体ip,端口根据实际情况传入或者修改默认配置, 写入和获取都是采用同一个。
  flag.StringVar(&url, "url", "10.10.107.225:4150", "nsqd")
  flag.Parse()
}

// 生产者
func startProducer() {
    cfg := nsq.NewConfig()
    producer, err := nsq.NewProducer(url, cfg)
    if err != nil {
        log.Fatal(err)
    }
    // 发布指定的消息
    for {
      if err := producer.Publish("DirectConnection", []byte("test message")); err != nil {
          log.Fatal("publish error: " + err.Error())
      }
      time.Sleep(1 * time.Second)
  }
}

// 消费者
func startConsumer() {
    cfg := nsq.NewConfig()
    consumer, err := nsq.NewConsumer("DirectConnection", "first", cfg)
    if err != nil {
        log.Fatal(err)
    }
    // 设置消息处理函数
    consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
        log.Println(string(message.Body))
        return nil
    }))
    // 连接到单例nsqd
    if err := consumer.ConnectToNSQD(url); err != nil {
        log.Fatal(err)
    }
    <-consumer.StopChan
}
```


至此完毕！