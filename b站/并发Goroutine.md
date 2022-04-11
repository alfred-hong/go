基础概念

Go语言的并发

> 描述； Go语言并发是通过goroutine (go [ruːˈtiːn])实现的，goroutine类似于线程，属于用户态的线程（比内核态线程更轻量级）, 我们可以根据需要创建成千上万个 goroutine 并发工作。它是由Go语言的运行时（runtime）调度完成（自己编写的线程调度），而线程是由操作系统调度完成。

除此之外Go语言还提供channel在多个goroutine间进行通信。所以goroutine和channel 是 Go 语言秉承的CSP（Communicating Sequential Process-通信顺序过程）并发模式的重要实现基础。

**并发 (concurrency):** 逻辑上具有处理多个同时性任务的能力。（某段时间里你在用微信和两个女朋友聊天）

**并行 (parallesim):** 物理上同一时刻执行多个并发任务。（这一刻你和你朋友都在用微信和女朋友聊天）

**进程 (Process):** 是操作系统中某数据集合操作的一个程序运行的资源，进程是资源分配的最小单位。(车辆生产的某一整条流水线作业，例如从零件装配到整车完备)

**线程 (Thread):** 是运行的进程任务中的某一条执行流程，在进程中的更小的运行单位，线程是CPU调度的最小单位。(流水线中某一个流水作业岗位，比如外观包装岗位工)

**用户态:** 由程序员自己定义的线程调度。

**内核态:** 由操作系统预定义的线程调度。



Goroutine 入门

> 描述: 在java/c++中我们要实现并发编程的时候，我们通常需要自己维护一个线程池，并且需要自己去包装一个又一个的任务，同时需要自己去调度线程执行任务并维护上下文切换，这一切通常会耗费程序员大量的心智。

> Q: 那么能不能有一种机制，程序员只需要定义很多个任务，让系统去帮助我们把这些任务分配到CPU上实现并发执行呢？
>
> 答: 那当然是有的，即Go语言中的 goroutine 就是这样一种机制，其概念类似于线程,但 goroutine是由Go的运行时（runtime）调度和管理的, Go程序会智能地将 goroutine 中的任务合理地分配给每个CPU。
> 这就是为何Go语言之被称为现代化的编程语言，就是因为它在语言层面已经内置了调度和上下文切换的机制。

如何使用 goroutine（并发）?

> 描述: Go语言中使用goroutine非常简单,将需要并发执行的任务包装成为一个函数,并在调用函数的时候前面加上go关键字,此时便可以开启一个Goroutine去执行该函数的任务。
> 非常注意，一个goroutine必定对应一个函数，但可以创建多个goroutine去执行相同的函数。



示例1.启动单个Goroutine

描述：启动goroutine的方式非常简单，只需要在调用的函数（普通函数和匿名函数）前面加上一个go关键字，即可。

```go
func hello() {
	fmt.Println("Hello Goroutine!")
}
func main() {
  	go	hello()               // 在调用hello函数前面加上关键字go，即启动一个goroutine去执行hello这个函数。
	fmt.Println("main goroutine done!")
 	time.Sleep(time.Second) // 如果不加此延时，则执行结果只打印了main goroutine done!，并不会打印Hello Goroutine!
}
```

> Q: 为什么在未加延时函数的情况下执行结果只打印了main goroutine done!？
>
> 答: 因为在程序启动时，Go程序就会为 main() 函数创建一个默认的goroutine。
>
> 当 main() 函数返回的时候该goroutine就结束了，所有在 main() 函数中启动的 goroutine 会一同结束。
> 例如 main函数所在的goroutine就像是权利的游戏中的夜王，其他的goroutine都是异鬼，夜王一死它转化的那些异鬼也就全部GG了。
> 所以为了让main函数等一等hello函数，最简单粗暴的方式就是使用time.Sleep()，不过后续还要更好的解决办法那就是采用 sync.WaitGroup来等待goroutine任务执行完毕。



示例2.启动多个Goroutine

> 描述: 在Go语言中实现并发就是这样简单，我们还可以启动多个goroutine。

> 此处我将解决示例1中的延时阻塞问题，如何监测 Goroutine 什么时候结束?
>
> 1.Goroutine 在对应调用函数运行完成时结束，其次是main函数执行完成时，由main函数创建的那些goroutine都结束。
> 2.此处我们这里使用sync.WaitGroup来实现goroutine的同步执行结束检测。
> wg.Add(1) : goroutine 任务执行计数器+1
> wg.Done() : goroutine 任务执行完毕计数器-1
> wg.Wait() : 当前 goroutine 任务执行全部完毕后且计数器为0,等待结束退出Main函数

代码示例:

```go
var wg sync.WaitGroup
func hello(i int) {
	defer wg.Done() // goroutine结束就登记-1
	fmt.Println("Hello Goroutine!", i)
}
func main() {

	for i := 0; i < 10; i++ {
		wg.Add(1) // 启动一个goroutine就登记+1
		go hello(i)
	}
	wg.Wait() // 等待所有登记的goroutine都结束
}
```


Tips: 执行上述代码您会发现每次打印的数字的顺序都不一致，这是由于10个goroutine是并发执行的，而goroutine的调度是随机的。

实践演示(1).常规方式以及匿名函数方式并发调度

```go
package main
import (
	"fmt"
	"time"
)
// 方式1，常规函数方式
func hello(count int) {
	fmt.Printf("欢迎你第 %d 次\n", count)
}
func demo1() {
	for i := 0; i < 5; i++ {
		go hello(i) // 开启一个单独的goroutine去执行hello函数(任务)
	}
}

// 方式2.匿名函数方式
func demo2() {
	for i := 0; i < 5; i++ {
		go func(count int) {
			fmt.Printf("第 %d 次欢迎你\n", count)
		}(i)
	}
}

// 程序启动之后会创建一个主Goroutine去执行。
func main() {
	fmt.Println("[*] main start")
	demo1()
	demo2()
	time.Sleep(time.Second)  // 最暴力简单的延时函数
	fmt.Println("[-] main end")
}
// 如果main函数结束了，则由main启动的goroutine也结束了。
```


执行结果:

```sh
[*] main start
第 4 次欢迎你
欢迎你第 0 次
欢迎你第 1 次
欢迎你第 2 次
欢迎你第 3 次
欢迎你第 4 次
第 0 次欢迎你
第 1 次欢迎你
第 2 次欢迎你
第 3 次欢迎你  # 你会发现输出完此行时，会有卡顿。这是由于未到延时时间代码便执行完毕，所以终端处于阻塞模式，等待延时时间得到来。
[-] main end
```


上面的方式在实际的开发工作不会这样使用，一般会利用sync.WaitGroup对象，在函数中当goroutine全部执行完毕后，将会自动结束运行, 可以参考下述实践示例2。



实践演示(2).使用sync.WaitGroup对线下等待并发线程执行完毕

```go
// 优雅的使用 goroutine 并行调度。
package main
import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)
// 实例化结构体得到对象
var wg sync.WaitGroup
// goroutine 将调用的函数
func f1() {
	defer wg.Done() // goroutine结束就登记-1
	fmt.Println("输出的随机数：", rand.Intn(10))
}
// 当main函数任务中所有的goroutine都结束了，才结束main函数
func main() {
	fmt.Println("Start Main")
	// goroutine 调用
	for i := 0; i < 10; i++ {
		rand.Seed(time.Now().Unix()) // 随机数种子传入的是一个int64类型
		wg.Add(1)  // 启动一个goroutine就登记+1
		go f1()
	}
	wg.Wait() 	// 监听并等待main中的goroutine是否都运行结束(等待wg组为0时)。
	fmt.Println("End Main")  // 输出再也没有卡顿了
}
```


执行结果:

```sh
Start Main
输出的随机数： 2
输出的随机数： 2
输出的随机数： 8
输出的随机数： 2
输出的随机数： 1
输出的随机数： 7
输出的随机数： 6
输出的随机数： 2
输出的随机数： 7
输出的随机数： 0
End Main   // 此处输出再也没有卡顿了。
```

Goroutine 特性

(1) Goroutine 可增长的栈

> 描述: OS线程（操作系统线程）一般都有固定的栈内存（通常为2MB）,一个goroutine的栈在其生命周期开始时只有很小的栈（典型情况下2KB），goroutine的栈不是固定的，他可以按需增大和缩小，goroutine的栈大小限制可以达到1GB，虽然极少会用到这么大。所以在Go语言中一次性创建十万左右的goroutine也是可以的。

Tips : 总结 goroutine 初始栈大小为2k，其大小可以按照需要扩容缩。

(2) Goroutine 调度模型

> 描述:GPM是Go语言运行时（runtime）层面的实现，是go语言自己实现的一套调度系统。区别于操作系统调度OS线程，它比起OS的调度更轻量级些。

goroutine GMP模型:

​	G 其就是个goroutine里面除了存放本goroutine信息外，还存放与所在P的绑定等信息。
​	P 管理着一组goroutine队列，P里面会存储当前goroutine运行的上下文环境（函数指针，堆栈地址及地址边界），P会对自己管理的goroutine队列做一些调度（比如把占用CPU时间较长的goroutine暂停、运行后续的goroutine等等）当自己的队列消费完了就去全局队列里取，如果全局队列里也消费完了会去其他P的队列里抢任务。
​	M (machine) 是Go运行时（runtime）对操作系统内核线程的虚拟，M与内核线程(Kernel Thread)一般是一一映射的关系，一个groutine最终是要放到M上执行的；

> Q: P与M有何关系?
> 描述: P与M通常是一一对应的,他们关系是P管理着一组G挂载在M上运行。当一个G长久阻塞在一个M上时，runtime会新建一个M，阻塞G所在的P会把其他的G挂载在新建的M上，当旧的G阻塞完成或者认为其已经死掉时则回收旧的M。

P的个数是通过 runtime.GOMAXPROCS设定（最大256 核），Go1.5版本之后默认为物理线程数。在并发量大的时候会增加一些P和M，但不会太多，切换上下文太频繁的话得不偿失。

单从线程调度讲，Go语言相比起其他语言的优势在于OS线程是由OS内核来调度的，goroutine则是由Go运行时（runtime）自己的调度器调度的，这个调度器使用一个称为m:n调度的技术（复用/调度m个goroutine到n个OS线程）。 其一大特点是goroutine的调度是在用户态下完成的， 不涉及内核态与用户态之间的频繁切换，包括内存的分配与释放，都是在用户态维护着一块大的内存池， 不直接调用系统的malloc函数（除非内存池需要改变），成本比调度OS线程低很多。 另一方面充分利用了多核的硬件资源，近似的把若干goroutine均分在物理线程上， 再加上本身goroutine的超轻量，以上种种保证了go调度方面的性能。

**Tips: goroutine 组最终是要放在M(内核态)中执行，不过在此之前goroutine已经将任务进行排好队列（底层实现线程池），然后等待分别到操作系统之中。**

课外扩展: https://www.cnblogs.com/sunsky303/p/9705727.html

(3) Goroutine 线程数

> 描述: Go运行时的调度器使用GOMAXPROCS参数来设置使用多少个OS线程来同时执行Go代码，其默认值是机器上的CPU核心数。

Go语言中可以通过runtime.GOMAXPROCS(NUMBER)函数, 设置当前程序并发时占用的CPU逻辑核心数。并可以通过 runtime.NumCPU() 与 runtime.NumGoroutine() 分别查看机器中的逻辑CPU数和当前存在的goroutine数。

Tips: Go1.5版本之前，默认使用的是单核心执行。
Tips: Go1.5版本之后，默认使用全部的CPU逻辑核心数。

例如，在一个8核心的机器上，调度器会把Go代码同时调度到8个OS线程上（GOMAXPROCS是m:n调度中的n）。

> Q: 什么是M:N?
> 答: M:N 即把m个goroutine任务分配给n个操作系统线程去执行。

Go语言中的操作系统线程和goroutine的关系：

​	一个操作系统线程对应用户态的多个goroutine。
​	go语言程序可以同时使用多个操作系统线程。
​	go语言中的goroutine与OS线程是多对多的关系，即 m:n 的关系。


实践案例:
例如: 我们可以通过将任务分配到不同的CPU逻辑核心上实现并行的效果：

```go
package main
import (
	"fmt"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func a() {
	defer wg.Done()
	for i := 0; i < 6; i++ {
		println("Func A() :", i)
	}
}

func b() {
	defer wg.Done()
	for i := 0; i < 6; i++ {
		println("Func B() :", i)
	}
}

func main() {
	fmt.Println("[*] Main Start")
	fmt.Println("当前机器的 CPU 核心数:", runtime.NumCPU()) // NumCPU返回当前进程可用的逻辑CPU数量。
	runtime.GOMAXPROCS(2)  // 占用cpu的两个核
	wg.Add(1)
	go a() // 并发调用a函数（后输出）
	wg.Add(1)
	go b() // 并发调用a函数（先输出）
	fmt.Println("当前机器的 goroutine 数:", runtime.NumGoroutine()) // NumGoroutine返回当前存在的goroutine数。
	wg.Wait()
	fmt.Println("[*] Main End")
}
```

```sh
[*] Main Start
当前机器的 CPU 核心数: 4
当前机器的 goroutine 数: 3
[*] Main End
Func B() : 0
Func B() : 1
Func B() : 2
Func B() : 3
Func B() : 4
Func B() : 5
Func A() : 0
Func A() : 1
Func A() : 2
Func A() : 3
Func A() : 4
Func A() : 5
// 结果说明:  Mac 可以复现成功，而Linux与Windows还是一个goroutine任务做完后再做另外一个goroutine任务,其结果如上。
// 当两个任务只有一个逻辑核心，此时是做完一个任务再做另一个任务。
// 当将逻辑核心数设为2，此时两个任务并行执行，即理想的状态是做一个任务a再做一个任务b，进行交叉执行。
```


Channel 通道

> 描述: 上面章节介绍了Goroutine的基本使用，但是您会发现单纯地将函数并发执行是没有意义的，函数与函数间需要交换数据才能体现并发执行函数的意义。

背景说明
虽然可以使用共享内存进行数据交换，但是共享内存在不同的goroutine中容易发生竞态问题，所以为了保证数据交换的正确性，就必须使用互斥量对内存进行加锁，但是这种做法势必造成性能方面的问题。

解决办法
为了解决上述问题,Go语言的并发模型采用得是(CSP-Communicating Sequential Processes), 它提倡通过通信共享内存而不是通过共享内存而实现通信,其引入了channel的概念。

基础介绍
如果说 goroutine是Go程序并发的执行体，而 channel (英 [ˈtʃænl]) 就是它们之间的连接通道, channel 是可以让一个goroutine发送特定值到另一个goroutine的通信机制。

简单的说: 即通过Channel实现多个goroutine之间的通信。

Go 语言中的通道（channel）是一种特殊的类型, 通道像一个传送带或者队列，总是遵循先入先出（First In First Out）的规则，保证收发数据的顺序。每一个通道都是一个具体类型的导管，也就是声明channel的时候需要为其指定元素类型。

1.channel 类型

> 描述: channel是特殊类型(一种引用类型), 其声明通道类型的格式如下：var 变量 chan 元素类型

示例说明:

```go
var ch1 chan int   // 声明一个传递整型的通道
var ch2 chan bool  // 声明一个传递布尔型的通道
var ch3 chan []int // 声明一个传递int切片的通道
```

2.channel 创建

> 描述: 创建 channel 的格式如下：make(chan 元素类型, [缓冲大小]), channel的缓冲大小是可选的。

Tips: 由于通道是引用类型, 声明的通道后需要使用make函数初始化之后才能使用，并且需要注意通道类型的空值是nil。
Tips: Slice、Map、Channel 需要make函数初始化后方能使用。

```sh
var ch chan int
fmt.Println(ch)   // 未初始返回 <nil>

ch4 := make(chan int) // 必须用make函数初始化之后才能使用
ch5 := make(chan bool)
ch6 := make(chan []int,10)  // 带缓冲区的
```


Tips: 声明时指定通道中元素类型,定义使用时通道必须使用make函数初始化才能使用。

3.channel 操作

> 描述: 通道有发送(send) 、接收 (receive) 和 关闭(close) 三种操作, 但值得注意的是发送和接收都使用<-符号。

send ：将一个值发送到通道中。
receive : 从一个通道中接收值。
close : 我们通过调用内置的close函数来关闭通道,也可以不用关闭通道,因为通道是引用类型的变量,所以它会在程序结束后自动给GC(垃圾回收)。


示例说明

```sh
// 定义一个通道(初始化通道)
ch := make(chan int)

// 发送 send
ch <- 10 // 把10发送到ch中

// 接收 receive
x := <-ch  // 从ch中接收值并赋值给变量x
<-ch       // 从ch中接收值，忽略结果

// 关闭 close
close(x)
```


温馨提示:

关于关闭通道需要注意的事情是，只有在通知接收方goroutine所有的数据都发送完毕的时候才需要关闭通道。
通道是可以被垃圾回收机制回收的，它和关闭文件是不一样的，在结束操作之后关闭文件是必须要做的，但关闭通道不是必须的。
通道关闭后的特点如下所描述:

```sh
* 对一个关闭的通道再发送值就会导致panic。
* 对一个关闭的通道进行接收会一直获取值直到通道为空。
* 对一个关闭的并且没有值的通道执行接收操作会得到对应类型的零值。
* 关闭一个已经关闭的通道会导致panic。
```

4.channel 缓冲

> 描述: 我们可以为 channel 设置缓冲或者不设置缓冲区，其两者概念和区别如下。

无缓冲的通道: 又称为阻塞的通道(必须有接收才能发送)。
有缓冲的通道: 为解决无缓冲的通道存在的问题孕育而生。


无缓冲的通道

描述: 使用无缓冲通道进行通信将导致发送和接收的goroutine同步化,因此无缓冲通道也被称为同步通道。

废话不多说,先看引入的示例代码。

```go
func main() {
	ch := make(chan int)
	ch <- 10  // 代码会阻塞在 ch <- 10这一行代码形成死锁
	fmt.Println("发送成功")
}
```


上面这段代码能够通过编译，但是执行的时候会出现以下错误：

```go
fatal error: all goroutines are asleep - deadlock!
goroutine 1 [chan send]:
main.main()
.../src/github.com/main.go:8 +0x50
```

> Q: 为什么会出现deadlock错误呢？
>
> 答: 因为我们使用ch := make(chan int)创建的是无缓冲的通道，无缓冲的通道只有在有人接收值的时候才能发送值。
> 例如: 就像你住的小区如果没有快递柜和代收点，快递员给你打电话必须要把这个物品送到你的手中，简单来说就是无缓冲的通道必须有接收才能发送。

> Q: 那如何解决这个问题呢？
>
> 答: 一种方法是启用一个goroutine去接收值，另外一种方式就是采用带缓冲的通道（后续介绍）。

```go
func recv(c chan int) {
	ret := <-c
	fmt.Println("接收成功", ret)
}
func main() {
	ch := make(chan int)
	go recv(ch) // 启用goroutine从通道接收值
	ch <- 10
	fmt.Println("发送成功")
}
```

Tips: 在无缓冲通道上的发送操作会阻塞，直到另一个goroutine在该通道上执行接收操作，这时值才能发送成功，两个goroutine将继续执行。
Tips: 在无缓冲通道上的如果接收操作先执行，接收方的goroutine将阻塞，直到另一个goroutine在该通道上发送一个值。

有缓冲的通道

> 描述: 使用有缓冲区的通道可以解决无缓冲的通道阻塞问题, 我们可以在使用make函数初始化通道的时候为其指定通道的容量，例如：

```go
func main() {
	ch := make(chan int, 1) // 创建一个容量为1的有缓冲区通道
	ch <- 10
	fmt.Println("发送成功")
}
```

只要通道的容量大于零，那么该通道就是有缓冲的通道，通道的容量表示通道中能存放元素的数量, 就像你小区的快递柜只有那么个多格子，格子满了就装不下了就阻塞了，等到别人取走一个快递员就能往里面放一个。

同时我们可以使用内置的len函数获取通道内元素的数量，使用cap函数获取通道的容量，虽然我们很少会这么做。

实践案例:

```go
package main

import (
	"fmt"
	"sync"
)

// make 函数申请内存空间的传入对象（实例化三种类型）
var s []int          // slice 切片
var m map[string]int // map 字典映射
var c chan int       // 指定通道中元素的类型

// 定义全局的waitGroup
var wg sync.WaitGroup

// 无缓冲的通道示例
func noBuffer() {
	fmt.Println(c)      // 未初始化的通道返回 nil （未向内存中申请空间）
	c := make(chan int) // 不带缓冲区通道的初始化 （但必须有对应的接收）
	fmt.Println("将 10 发生到 channel c 之中")
	wg.Add(1)
	go func() { // 并行任务的顺序非常重要，此处不能放在 c <- 10 后否则终端将会一直处于阻塞状态
		defer wg.Done()
		x := <-c
		fmt.Printf("Backgroup Goroutine 从 channel c 中取得 %v \n\n", x)
	}()
	c <- 10 // 将 10 发生到 channel c 之中（注意此行放的顺序）
	wg.Wait()
	defer close(c) // 关闭通道
}

// 有缓冲的通道示例
func useBuffer() {
	fmt.Println(c)        // 未初始化的通道返回 nil （未向内存中申请空间）
	c = make(chan int, 2) // 带缓冲区通道的初始化
	fmt.Println("通道缓冲数量（发送前）:", len(c))
	c <- 10                                // 将 10 发生到 channel c 之中
	fmt.Println("同样将 10 发生到 channel c 之中") // 此处将不会阻塞
	c <- 20                                // 将 10 发生到 channel c 之中
	fmt.Println("然后将 20 发生到 channel c 之中") // 如何缓冲区通道初始化为1，则此处将阻塞,如果初始化通道缓冲区大于等于2将会不阻塞
	fmt.Println("通道缓冲数量（发送后）:", len(c))
	x := <-c
	fmt.Printf("第一次，从channel c中取到了 %v\n", x)
	x = <-c
	fmt.Printf("第二次，从channel c中取到了 %v\n", x)
	fmt.Printf("channel c ptr = %p \n", c)
	defer close(c) // 关闭通道
}

func main() {
	noBuffer()
	useBuffer()
}
```

执行结果:

```go
<nil>
将 10 发生到 channel c 之中
Backgroup Goroutine 从 channel c 中取得 10

<nil>
通道缓冲数量: 0
同样将 10 发生到 channel c 之中
然后将 20 发生到 channel c 之中
通道缓冲数量: 2
第一次，从channel c中取到了 10
第二次，从channel c中取到了 20
channel c ptr = 0xc0000240e0
```




5.channel 遍历

> 描述: 当一个通道发送到通道队列里有多个值时, 此时我们想取出通道队列的所有值时，我们可以使用for range 遍历通道，并且当通道被关闭的时候就会退出for range遍历。

当向通道中发送完数据时，我们可以通过close函数来关闭通道，如果此时再往该通道发送值会引发panic，从该通道取值的操作会先取完通道中的值，再然后取到的值一直都是对应类型的零值。

> Q: 那如何判断一个通道是否被关闭了呢？
>
> 第一种方法<-ch, 第二种方法for range遍历通道。

实践案例:

```go
// channel 遍历实践操作
var wg sync.WaitGroup

// 方式1
func method1() {
	ch1 := make(chan int) // 不带缓冲区
	ch2 := make(chan int)
	// 开启goroutine将1~9的数发送到ch1中
	go func() {
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
		close(ch1) // 关闭通道 ch2 （此时只能读不能写）
	}()
	// 开启goroutine从ch1中接收值，并将该值的平方发送到ch2中
	go func() {
		for {
			i, ok := <-ch1 // 通道关闭后再取值到末尾时，ok=false 【关键点值得学习】
			if !ok {
				break
			}
			ch2 <- i * i // 同样求取通道的平方
		}
		close(ch2) // 关闭通道 ch2 （此时只能读不能写）
	}()
	// 在主goroutine中从ch2中接收值打印
	fmt.Println("方式1:")
	for i := range ch2 { // 通道关闭后会退出for range循环
		fmt.Printf("%d ", i)
	}
}

// 负责将10～19发送到ch1中
func f1(ch1 chan int) {
	defer wg.Done()
	for i := 10; i < 20; i++ {
		ch1 <- i
	}
	close(ch1) // 关闭通道 ch1 （此时只能读不能写）
}

// 负责将接收ch1值的值进行平方运算
func f2(ch1, ch2 chan int) {
	defer wg.Done()
	for num := range ch1 {   //【关键点】
		ch2 <- num * num
	}
	close(ch2) // 关闭通道 ch2 （此时只能读不能写）
}

// 方式2
func method2() {
	ch1 := make(chan int, 10) // 带缓冲区
	ch2 := make(chan int, 10)
	// goroutine 等待组数量设置
	wg.Add(2)
	// 开启 goroutine
	go f1(ch1)
	go f2(ch1, ch2)
	// 等待全部 goroutine 任务执行完毕
	wg.Wait()
	fmt.Println("方式2:")
	// 通道关闭后会退出for range循环
	for ret := range ch2 {
		fmt.Printf("%d ", ret)
	}
}

func main() {
	// 匿名函数
	method1()
	fmt.Println()
	// 常规函数
	method2()
}
```


Tips: 从上面的例子中我们看到有两种方式在接收值的时候判断该通道是否被关闭.

单向通道

> 描述: 有的时候我们会将通道作为参数在多个任务函数间传递，很多时候我们在不同的任务函数中使用通道都会对其进行限制，比如限制通道在函数中只能发送或只能接收。

所以在这样场景下Go语言中提供了单向通道来处理这种情况。

> out chan<- int 是一个只写单向通道（只能对其写入int类型值），可以对其执行发送操作但是不能执行接收操作;
> in <-chan int 是一个只读单向通道（只能从其读取int类型值），可以对其执行接收操作但是不能执行发送操作;
> Tips: 在函数传参及任何赋值操作中可以将双向通道转换为单向通道，但反过来是不可以的。

例如，我们把上面的例子改造如下：

```go
// 函数参数 out： 一个只写单向通道。
func counter(out chan<- int) {
	for i := 0; i < 100; i++ {
		out <- i
	}
	close(out)
}

// 函数参数 out： 一个只写单向通道,函数参数 in: 一个只读单向通道。
func squarer(out chan<- int, in <-chan int) {
	for i := range in {
		out <- i * i
	}
	close(out)
}

// 函数参数 in : 一个只读单向通道。
func printer(in <-chan int) {
	for i := range in {
		fmt.Println(i)
	}
}

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go counter(ch1) // ch1 此时只写
	go squarer(ch2, ch1) // ch2 此时只写，ch1 此时只读
	printer(ch2) // ch2 此时只读
}



```


实践案例
描述: 请传入通道中浮点数类型的三次方。

```go
// Unidirectional channel
var wg sync.WaitGroup
var once sync.Once

// 通道只写操作
func f1(ch1 chan<- float64) {
	defer wg.Done()
	for i := 0; i < 10; i++ {
		ch1 <- float64(i)
	}
	close(ch1)
}

// 通道只写与只读操作
func f2(ch2 chan<- float64, ch1 <-chan float64) {
	defer wg.Done()
	for i := range ch1 {
		ch2 <- math.Pow(i, 3)
	}
	// 确保某个操作只执行一次，防止因为通道关闭而其他并发(进程)不能读取
  anonymousefun := (func() {
		close(ch2)
	})
	once.Do(anonymousefun)
}

func main() {
	ch1 := make(chan float64, 10) // 带缓冲区的通道
	ch2 := make(chan float64, 10) // 带缓冲区的通道
	wg.Add(3)
	go f1(ch1)
	go f2(ch2, ch1) // 多次 goroutine 同一个函数
	go f2(ch2, ch1)
	wg.Wait()
  // 循环遍历，直到读取到末尾
	for {
		x, ok := <-ch2
		if !ok {
			break
		}
		fmt.Printf("%.0f ", x)
	}
	fmt.Println()
}
```


执行结果: 0 1 8 27 64 125 216 343 512 729

通道总结

> 描述: channel 常见的异常总结,非常注意关闭已经关闭的channel也会引发panic，如下表：

| Channel       | Nil   | 非空                         | 空的                         | 满了                         | 没满               |
| ------------- | ----- | ---------------------------- | ---------------------------- | ---------------------------- | ------------------ |
| 发送(Send)    | 阻塞  | 发送值                       | 发送值                       | 阻塞                         | 发送值             |
| 接受(Receive) | 阻塞  | 接受值                       | 阻塞                         | 接受值                       | 接受值             |
| 关闭(Close)   | panic | 关闭成功，读完数据后返回零值 | 关闭成功，读完数据后返回零值 | 关闭成功，读完数据后返回零值 | 关闭成功，返回零值 |

Goroutine 池

> 描述: 在工作中我们通常会使用可以指定启动的goroutine数量以worker pool模式,并且利用控制 goroutine的数量，来防止goroutine泄漏和暴涨等问题。

简单描述：就是预定义一定数量的Goroutine去执行任务(函数)。

实践案例1:

```go
// goroutine 同步等待组
var wg sync.WaitGroup
var one sync.Once

func worker(id int, jobs <-chan int, results chan<- int) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("worker:%d start job:%d\n", id, j)  // start
		time.Sleep(time.Millisecond * 500)            // 延时500s查看效果比较明显
		fmt.Printf("worker:%d end job:%d\n", id, j)    // end
		results <- j * 2
	}
	// 保证调用的函数只执行一次 (关键点)
	one.Do(func() {
		close(results)
	})
}

// 五个任务给3给goroutine池执行
func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	// goroutine 池(此处3个goroutine组成)
	wg.Add(3)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results) // 当执行到第五的一个任务后将介绍（所以一共会打印10次）
	}
	// 设置5个任务
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	// 保证调用的函数只执行一次 (关键点)
	one.Do(func() {
		close(jobs)
	})
	wg.Wait()
	// 输出结果
	fmt.Println("results channle length: ", len(results)) // 长度为 5
	// 方式1,打印后阻塞(会一直从channel中取数)
	// for {
	// 	x, ok := <-results  // results 通道被关闭时ok=False。
	// 	if !ok {
	// 		fmt.Println(ok)
	// 		break
	// 	}
	// 	fmt.Println(ok, x)
	// }

	// 方式2,打印后阻塞(会一直从channel中取数)
	// for i := range results {
	// 	fmt.Println(i)
	// }
	
	// 此种方式不会阻塞
	for a := 1; a <= 5; a++ {
		<-results
	}

}
```


执行结果:

```go
worker:3 start job:1
worker:2 start job:3
worker:1 start job:2
worker:1 end job:2
worker:1 start job:4
worker:2 end job:3
worker:2 start job:5
worker:3 end job:1
worker:1 end job:4
worker:2 end job:5
results channle length:  5
```


实践案例2:针对上述示例进行优化，优化并发流程,利用信号变量通知goroutine任务结束关闭通道.

```go
func worker(id int, jobs <-chan int, results chan<- int, notify chan<- struct{}) {
	for j := range jobs {
		time.Sleep(time.Millisecond * 500)          // 延时500s查看效果比较明显
		fmt.Printf("worker:%d end job:%d\n", id, j) // end
		results <- j * 2
		notify <- struct{}{} // 任务执行信号标识,实例化匿名结构体,如此使用非常节省空间.
	}

}

func main() {
	// 1.初始化通道
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	notify := make(chan struct{}, 5) // 作为信号使用,通道类型采用匿名结构体,占用的系统资源较少(常用值得学习),此处作为并发任务执行完毕通知.

	// 2.生成五个任务
	go func() {
		for i := 1; i <= 5; i++ {
			jobs <- i
		}
		close(jobs) // 给通道传递完值后关闭
	}()
	
	// 3.开启三个Goroutine (并发池)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results, notify)
	}
	
	// 4.通过信号变量,验证五个任务是否执行完毕,如果晚报则关闭results通道
	go func() {
		// 循环五次
		for i := 0; i < 5; i++ {
			<-notify
		}
		// 关闭通道，如果不关闭将一直阻塞
		close(results)
	}()
	
	// 5.遍历已关闭的通道值
	for res := range results {
		fmt.Println(res)
	}

}
```


执行结果:

```sh
worker:2 end job:2
worker:3 end job:3
worker:1 end job:1
4
6
2
worker:2 end job:4
8
worker:3 end job:5
10
```


Goroutine 多路复用

> 描述: 在某些场景下我们需要同时从多个通道接收数据, 通道在接收数据时, 如果没有数据可以接收将会发生阻塞。

你也许会写出如下代码使用遍历的方式来实现, 但此种方式虽然可以实现从多个通道接收值的需求，但是运行性能会差很多。所以为了应对这种场景，Go内置了select关键字，可以同时响应多个通道的操作。

```go
for{
  // 尝试从ch1接收值
  data, ok := <-ch1
  // 尝试从ch2接收值
  data, ok := <-ch2
  …
}
```

select 关键字使用类似于switch语句，它有一系列case分支和一个默认的分支。每个case会对应一个通道的通信（接收或发送）过程。select会一直等待，直到某个case的通信操作完成时，就会执行case分支对应的语句。

语法格式：

```go
select{
  case <-ch1:
      ...
  case data := <-ch2:
      ...
  case ch3<-data:
      ...
  default:
      默认操作
}
```

示例演示:
Goroutine Select 多路复用

```go
package main
import (
	"fmt"
	"sync"
)
var wg sync.WaitGroup
func main() {
	// 示例1
	ch1 := make(chan int, 1)
	fmt.Println("// 示例 1	")
	for i := 0; i < 10; i++ {
		select {
      case x := <-ch1:
        fmt.Printf("index : %d , x = %d\n", i, x) // 通道缓冲区为1时,结果时可以预测的.
      case ch1 <- i: // 将会把(变量i)偶数值传给ch1通道
      default:
        fmt.Printf("index : %d , default\n", i)
		}
	}

	// 示例2
	ch2 := make(chan int, 2)
	fmt.Println("// 示例 2")
	wg.Add(1)
	go func() {
		wg.Done()
		for i := 1; i <= 5; i++ {
			ch2 <- i + i
		}
		close(ch2)
	}()
	wg.Wait()
	for i := 1; i <= 5; i++ {
		select {
	  case j := <-ch2:
	    fmt.Println("case 1: ", i, j) // 通道缓冲区为大于1时,结果是不可以预测的.
	  case ch2 <- i:
	    fmt.Println("case 2: ch2 <- i", i)
	  default:
	    fmt.Println("默认执行", i)
		}
	}

}
```

执行结果:

```sh
// 示例 1
index : 1 , x = 0
index : 3 , x = 2
index : 5 , x = 4
index : 7 , x = 6
index : 9 , x = 8
// 示例 2
case 1:  1 2
case 1:  2 4
case 2: ch2 <- i 3
case 1:  4 6
case 1:  5 3  // 注意通道是先进先出类似于队列结构。
```

总结说明
使用select语句能提高代码的可读性。

可处理一个或多个channel的发送/接收操作。
如果多个case同时满足，select会随机选择一个。
对于没有case的select{}会一直等待，可用于阻塞main函数。


Goroutine 并发安全(锁)

> 描述: 有时候在Go代码中可能会存在多个goroutine同时操作一个资源（临界区），这种情况会发生竞态问题（数据竞态）。

类比现实生活中的例子有十字路口被各个方向的的汽车竞争, 还有火车上的卫生间被车厢里的人竞争, 针对十字路口我们可以采用红绿灯进行各个方向依次放行，针对火车上的卫生间我们可以上锁让外部人员知道厕所有人且无法进入。

所以针对于Go语言中并发时资源竞态问题，其采用锁Lock机制进行解决，保证在线程在并发执行时的安全(不被其它线程干扰)。

举个例子:

```go
var x int64
var wg sync.WaitGroup

func add() {
	for i := 0; i < 5000; i++ {
		x = x + 1
	}
	wg.Done()
}
func main() {
	wg.Add(2)
	go add()  // 两个 goroutine 执行add任务
	go add()
	wg.Wait()
	fmt.Println(x) // 由于数据竞争的影响，可能导致最后的结果与期待的不符。输出: 5702 理想为10000
}
```

上面的代码中我们开启了两个goroutine去累加变量x的值，这两个goroutine在访问和修改x变量的时候就会存在数据竞争，导致最后的结果与期待的不符。

互斥锁

> 描述: 互斥锁是一种常用的控制共享资源访问的方法，它能够保证同时只有一个goroutine可以访问共享资源。

Tips: Go语言中使用sync包的Mutex类型来实现互斥锁。

此处我们使用互斥锁来修复上面代码的问题：

```go
var x int64
var wg sync.WaitGroup
var lock sync.Mutex

func add() {
	for i := 0; i < 5000; i++ {
		lock.Lock() // 加锁
		x = x + 1
		lock.Unlock() // 解锁
	}
	wg.Done()
}
func main() {
  // 设置两个goroutine池进行执行add任务
	poolCount := 2
	wg.Add(poolCount)
	for i := 0; i < poolCount; i++ {
		go add()
	}
	wg.Wait()
	fmt.Println(x)  // 输出值为预期相同 10000
}
```


使用互斥锁能够保证同一时间有且只有一个goroutine进入临界区，其他的goroutine则在等待锁；当互斥锁释放后，等待的goroutine才可以获取锁进入临界区，多个goroutine同时等待一个锁时，唤醒的策略是随机的。

读写互斥锁

> 描述: 互斥锁是完全互斥的，但是有很多实际的场景下是读多写少的，当我们并发的去读取一个资源不涉及资源修改的时候是没有必要加锁的，所以在这种场景下使用读写锁是更好的一种选择。

Tips: Go语言中使用sync包的RWMutex类型来实现读写互斥锁，所以读写锁分为两种读锁和写锁。

```go
# 加、解 读锁
rwlock.RLock()       # 对读操作进行锁定
rwlock.RUnlock()     # 对读锁定进行解锁

# 加、解 写锁
sync.RWMutex.Lock()   # 对写操作进行锁定
sync.RWMutex.Unlock() # 对写锁定进行解锁
```


Tips: 当一个goroutine获取读锁之后，其他的goroutine如果是获取读锁会继续获得锁，如果是获取写锁就会等待；当一个goroutine获取写锁之后，其他的goroutine无论是获取读锁还是写锁都会等待。

读写锁使用示例:

```go
var (
	x      int64
	wg     sync.WaitGroup
	rwlock sync.RWMutex
)

func write() {
	rwlock.Lock() // 加写锁
	x = x + 1
	time.Sleep(10 * time.Millisecond) // 假设写操作耗时10毫秒
	rwlock.Unlock()                   // 解写锁
	wg.Done()
}

func read() {
	rwlock.RLock()               // 加读锁
	time.Sleep(time.Millisecond) // 假设读操作耗时1毫秒
	rwlock.RUnlock()             // 解读锁
	wg.Done()
}

func main() {
	start := time.Now()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go write()
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go read()
	}
	wg.Wait()
	end := time.Now()
	fmt.Printf("耗费时间 ：%v, x = %d \n", end.Sub(start), x)

  // 执行结果 => 耗费时间 ：104.906297ms, x = 10
}
```

Tips: **需要注意的是读写锁非常适合读多写少的场景，如果读和写的操作差别不大，读写锁的优势就发挥不出来。**
Tips: 在某些场景中也可能使用sync包中Map,它是一个开箱即用(不需要Make初始化)的并发安全的Map。如:Sync.Map.Store(key,value)、Ssync.Map.Load(Key)、sync.Map.LoadOrstore、sync.Map.Delete()、sync.Map.Range()

原子操作

> 描述: 在上面的示例中我们通过锁操作来实现线程(协程)同步,而实际上锁机制的底层是基于原子操作的，其一般直接通过CPU指令实现。

在Go语言中内置了对基本数据类型的一些并发安全操作，例如atomic包中的方法来实现协程同步，它提供了底层的原子级内存操作。

```go
// 修改内存地址中值为delta并返回新值
func AddInt64(addr *int64, delta int64) (new int64)
// 读取内存地址中的值并返回
func LoadInt64(addr *int64) (val int64)
```


实践示例:

```go
// 接口声明
type Counter interface {
	Inc()
	Load() int64
}

// 普通版
type CommonCounter struct {
	counter int64
}
func (c *CommonCounter) Inc() {
	c.counter++ // 写
}
func (c *CommonCounter) Load() int64 {
	return c.counter // 读
}

// 互斥锁版
type MutexCounter struct {
	counter int64
	lock    sync.Mutex
}
func (m *MutexCounter) Inc() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.counter++  // 写
}
func (m *MutexCounter) Load() int64 {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.counter  // 读
}

// 原子操作版
type AtomicCounter struct {
	counter int64
}
func (a *AtomicCounter) Inc() {
	atomic.AddInt64(&a.counter, 1) // 写
}
func (a *AtomicCounter) Load() int64 {
	return atomic.LoadInt64(&a.counter) // 读
}

func test(c Counter) {
	var wg sync.WaitGroup
	start := time.Now()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			c.Inc()
			wg.Done()
		}()
	}
	wg.Wait()
	end := time.Now()
	fmt.Println(c.Load(), end.Sub(start))
}

func main() {
	c1 := CommonCounter{} // 非并发安全
	test(&c1)
	c2 := MutexCounter{}  // 使用互斥锁实现并发安全
	test(&c2)
	c3 := AtomicCounter{} // 并发安全且比互斥锁效率更高
	test(&c3)
}
```

执行结果:

```go
943 1.240583ms
1000 1.514732ms
1000 421.618µs  // 可以看出其性能效率是最高的。
```


Tips: atomic包对于同步算法的实现很有用, 但是这些函数必须谨慎地保证正确使用，除了某些特殊的底层应用，使用通道或者sync包的函数/类型实现同步更好。



至此，Go语言并发介绍到这里就完毕了。 