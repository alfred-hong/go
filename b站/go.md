# go笔记

## 环境配置

1. 下载go
```sh
$ sudo pacman -S go
```
下载后系统会将go安装在/usr/lib/go目录下

2. 配置一些环境变量#
>一共需要三个环境变量，分别为：

>GOROOT -> go语言安装目录

>GOPATH -> go语言工作区

>GOBIN -> 存放go语言可执行文件目录

>先创建一个目录用作go语言的工作区

```sh
$ cd ~/指定目录
$ mkdir go
```
然后创建一个目录存放可执行文件

```sh
$ cd go
$ mkdir bin
```
为了随地调用go语言命令和go编译后的可执行文件，可以将\$GOROOT/bin和$GOBIN加入到PATH

将第二部所有操作添加到到.xprofile中

```sh
export GOROOT=/usr/lib/go			#第三方安装包路径
export GOPATH=~/指定目录/go			 #项目路径一般指向src
export GOBIN=~/指定目录/go/bin
export PATH=$PATH:$GOROOT/bin:$GOBIN
```
使.bash_profile生效

```sh
$ source .bash_profile
```
3. 其他

在GOPATH下创建src目录用存放源代码
在GOPATH下创建pkg目录用存放编译后的库文件

GOROOT : 指定安装Go语言开发包的解压路径
GOPATH : 指定外部Go语言代码开发工作区目录

GOPROXY : 指定代理Go语言从公共代理仓库中快速拉取您所需的依赖代码

终端输入go env检查是否安装成功

```sh
在vscode中手动安装环境使用 go install
调试程序vscode
go get github.com/derekparker/delve/cmd/dlv
```

Goroutine 并行设计
描述:

透过Goroutine能够让程序以异步的方式运行，而不需要担心一个函数导致程序中断，因此Go也非常地适合网络服务。

假设有个程序，里面有两个函数：

```go
func main() {
  // 假設 loop 是一個會重複執行十次的迴圈函式。
  // 迴圈執行完畢才會往下執行。
  loop()
  // 執行另一個迴圈。
  loop()
}
```

如此就不需要等待该函数运行完后才能运行下一个函数。

```go
func main() {
  // 透過 `go`，我們可以把這個函式同步執行，
  // 如此一來這個函式就不會阻塞主程式的執行。
  go loop()
  // 執行另一個迴圈。
  loop()
}
```

```go
package main

import (
	"fmt"
)

func add(a int, b int) int{	//强制要求花括号在后面
	var sum int
	sum = a + b
	return sum
}

func main() {

	var c int
	c = add(100,200)
	fmt.Println("add(100,200)=")
}
```

编译: go Build
描述: go build 命令表示将源代码编译成可执行文件。

在hello目录下执行go build(指定.go文件)或者在其他目录执行以下命令go build helloworld(项目需要在GOROOT路径的src目录之中),因为go编译器会去 GOPATH 的src目录下查找你要编译的hello项目

编译: go Build
描述: go build 命令表示将源代码编译成可执行文件。

在hello目录下执行go build(指定.go文件)或者在其他目录执行以下命令go build helloworld(项目需要在GOROOT路径的src目录之中),因为go编译器会去 GOPATH 的src目录下查找你要编译的hello项目

```sh
# - 目录下执行
$ pwd
/root/Go/Day01
$ go build
$ ./Day01
Hello World.

# - 指定main包所在的.go文件
$ go build HelloWorld.go
$ ./HelloWorld
Hello World.

# - 使用-o参数来指定编译后得到的可执行文件的名字
$ go build -o ahelloworld.
$ ./ahelloworld
Hello World.
```

Tips : 如上述编译得到的可执行文件会保存在执行编译命令的当前目录下会有 HelloWorld 可执行文件。



编译&运行: go Run
描述: 我们也可以直接执行程序，该命令本质上也是先编译再执行。

```sh
$ go run HelloWorld.go
Hello World.
```



编译&安装软件包&依赖项: go Install
描述: go install 表示安装的意思，它先编译源代码得到可执行文件，然后将可执行文件移动到GOPATH的bin目录下。因为我们的环境变量中配置了GOPATH下的bin目录，所以我们就可以在任意地方直接执行可执行文件了。

```sh
$ go install                # 生成 Day01 可执行文件
$ go install HelloWorld.go  # 生成 HelloWorld 可执行文件
$ ls ${GOROOT}
$ ls ${GOPATH}/bin/
# Day01  dlv  dlv-dap  gomodifytags  go-outline  gopkgs  goplay  gopls  gotests  HelloWorld  impl  staticcheck
$ /home/weiyigeek/app/program/project/go/bin/Day01
Hello World.
```



跨平台编译: CGO_ENABLED / GOOS / GOARCH
描述: 默认我们go build的可执行文件都是当前操作系统可执行的文件，如果我想在windows下编译一个linux下可执行文件，那需要怎么做呢？

只需要指定目标操作系统的平台和处理器架构即可，例如Windows平台cmd下按如下方式指定环境变量编译出的可以执行文件则可以在Linux 操作系统 amd64 处理器中执行,然后再执行go build命令，得到的就是能够在Linux平台运行的可执行文件了。

```sh
SET CGO_ENABLED=0  # 禁用CGO
SET GOOS=linux     # 目标平台是linux
SET GOARCH=amd64   # 目标处理器架构是amd64
```

注意：如果你使用的是PowerShell终端，那么设置环境变量的语法为 $ENV:CGO_ENABLED=0。



不同平台快速交叉编译:

```sh
# 目标平台是linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
# 目标平台Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build
# 目标平台Mac
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build

# 和go build用法一样
```

简单实践: 在Liunx平台上编译出在Windows上运行的helloWorld.exe可执行文件。

```sh
# Linux
$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o Win-HelloWorld.exe
$ ls
HelloWorld.go  Win-HelloWorld.exe

# Windows: 拷贝后执行
PS D:\Temp> .\Win-HelloWorld.exe
Hello World.
```

Tips : 对比不同平台交叉编译后的可执行文件大小。

```sh
$ ls -la --ignore HelloWorld.go
-rwxrwxr-x 1 weiyigeek weiyigeek 1937799 7月  30 03:23 helloworld         # ~ 1.9 MB
-rwxrwxr-x 1 weiyigeek weiyigeek 2027936 7月  30 03:24 Mac-HelloWorld     # ~ 2.0 MB
-rwxrwxr-x 1 weiyigeek weiyigeek 2098688 7月  30 02:58 Win-HelloWorld.exe # ~ 2.1 MB
```

Http Web Server

描述: 透过Go仅需几行代码就完成HTTP网页服务器的实现。

```go
package main

import (
  "io"
  "net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
  io.WriteString(w, "Hello world!")
}

func main() {
  http.HandleFunc("/", hello)
  http.ListenAndServe(":8000", nil)
}
```



echo 类似命令程序

描述: 用Go去实现一个像Unix中的Echo命令程序。

```go
package main

import (
  "os"
  "flag"
)

var omitNewline = flag.Bool("n", false, "don't print final newline")

const (
  Space   = " "
  Newline = "\n"
)

func main() {
  flag.Parse() // Scans the arg list and sets up flags
  var s string = ""
  for i = 0; i < flag.NArg(); i++ {
      if i > 0 {
          s += Space
      }
      s += flag.Arg(i)
  }
  if !*omitNewline {
      s += Newline
  }
  os.Stdout.WriteString(s)
}
```

## 项目结构

在进行Go语言开发的时候，我们的代码总是会保存在\$GOPATH/src目录下。在工程经过go build、go install或go get等指令后，会将下载的第三方包源代码文件放在 \$GOPATH/src 目录下，产生的二进制可执行文件放在 \$GOPATH/bin目录下，生成的中间缓存文件会被保存在 $GOPATH/pkg 下。

Tips : 如果我们使用版本管理工具（Version Control System，VCS。常用如Git/Svn）来管理我们的项目代码时，我们只需要添加$GOPATH/src目录的源代码即可, bin 和 pkg 目录的内容无需版本控制。



通常来讲GOPATH目标下文件目录组织架构的设置常常有以下三种:

(1)适合个人开发者

描述: 我们知道源代码都是存放在GOPATH的src目录下，那我们可以按照下图来组织我们的代码。
![WeiyiGeek.适合个人开发者](https://i0.hdslb.com/bfs/article/e850cd44113c98ed6b73060c17c415810c6fa1f1.png@930w_558h_progressive.png)


(2)适合企业开发场景

描述: 此种目录结构设置更适合企业开发环境,以代码仓库为前缀并以公司内部组织架构为基准,其次是项目名称，最后是各个模块开发的名称。
![WeiyiGeek.适合企业开发场景](https://i0.hdslb.com/bfs/article/8e2305d5dc1de204fa366cac4abac9dd9d674016.png@942w_410h_progressive.png)


(3)目前流行的项目结构

描述: Go语言中也是通过包来组织代码文件，我们可以引用别人的包也可以发布自己的包，但是为了防止不同包的项目名冲突，我们通常使用顶级域名来作为包名的前缀，这样就不担心项目名冲突的问题了。

因为不是每个个人开发者都拥有自己的顶级域名，所以目前流行的方式是使用个人的github用户名来区分不同的包。
![WeiyiGeek.目前流行的项目结构](https://i0.hdslb.com/bfs/article/6430f2bbc96a68ddeb83b826ab9fc5fd41f368d7.png@942w_452h_progressive.png)
目前流行的项目结构
举例说明: 张三和李四都有一个名叫studygo的项目，那么这两个包的路径就会是：

```go
import "github.com/zhangsan/studygo"
import "github.com/lisi/studygo"
```

举例说明: 同样如果我们需要从githuab上下载别人包的时候如：

```go
go get github.com/jmoiron/sqlx, 那么这个包会下载到我们本地GOPATH目录下的src/github.com/jmoiron/sqlx。
```

总结说明: 我们的开发学习示例基本按照第三种项目结构进行。 

而Go语言推荐使用驼峰法式命名。

```go
# 下划线连接
student_name

# 小驼峰法式 (推荐方式)
studentName

# 大驼峰法式
StudentName 
```

25个关键字

```go
var const ：     变量和常量的声明
var varName type  或者 varName : = value
package and import: 导入
func：   用于定义函数和方法
return ：用于从函数返回
defer someCode ：在函数退出之前执行
go :      用于并行
select    用于选择不同类型的通讯
interface 用于定义接口
struct    用于定义抽象数据类型
break、case、continue、for、fallthrough、else、if、switch、goto、default 流程控制
chan  用于channel通讯
type  用于声明自定义类型
map   用于声明map类型数据
range 用于读取slice、map、channel数据
```

37个保留字

```go
# Constants:
true  false  iota  nil

# Types:
int  int8  int16  int32  int64
uint  uint8  uint16  uint32  uint64  uintptr
float32  float64  complex128  complex64
bool  byte  rune  string  error

# Functions:
make  len  cap  new  append  copy  close  delete
complex  real  imag
panic  recover
```

变量声明 ：var 变量名 变量类型 

```go
# 单一声明: 变量声明以关键字var开头，变量类型放在变量的后面，行尾无需分号。
var name string
var age int
var isOk bool

# 批量声明: 每声明一个变量就需要写var关键字会比较繁琐，go语言中还支持批量变量声明。
var (
  a string
  b int
  c bool
  d float32
) 
```

变量初始化：var 变量名 类型 = 表达式 

声明会默认初始化

1.整型和浮点型变量的默认值为0。
2.字符串变量的默认值为空字符串。
3.布尔型变量默认为false。
4.切片、函数、指针变量的默认为nil。

```go
//# 单一变量初始化
var name string = "WeiyiGeek"
var age int = 18

//# 批量变量初始化
var name, age = "WeiyiGeek", 20 

//类型推导初始化
var name = "WeiyiGeek"
var age = 18 
```

短变量声明:= 

```go
func main() {
  count := 10
  username := "WeiyiGeek"
} 
```

匿名变量 _

在使用多重赋值时，如果想要忽略某个值，可以使用匿名变量（anonymous variable）- 特殊变量。 

```go
func foo() (int, string) {
  return 10, "Q1mi"
}
func main() {
  x, _ := foo()
  _, y := foo()
  fmt.Println("x=", x)
  fmt.Println("y=", y)
} 
```

**示例演示：**

```go
package main

import "fmt"

// 变量声明(单一-全局)
var singleName string
var notUseVar bool

// 变量声明(批量-全局)
var (
  multiName string
  multiAge  int8
)

func main() {
  // 对声明后的变量赋值
  singleName = "Weiyi_"
  multiName = "Geek"
  multiAge = 18

  // 变量初始化（局部）
  var name string = "WeiyiGeek"
  var sex, addr = "boy", "China"

  // 类型推导变量
  var flag = true
  var count = 1024

  // 简短变量声明（此种类型只能在函数中使用）
  briefCount := 65535

  fmt.Printf("My Name is %s, Sex is %s , Address: %s\n", name, sex, addr)
  fmt.Println("Alias Name :", singleName, multiName, " Age is :", multiAge)
  fmt.Print("类型推导 ：", flag, count)
  fmt.Println(", 简短变量 ：", briefCount)
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
My Name is WeiyiGeek, Sex is boy , Address: China
Alias Name : Weiyi_ Geek  Age is : 18
类型推导 ：true 1024, 简短变量 ： 65535
```

> Tips : Go语言中变量必须先声明后使用，而且声明变量(非全局变量)后必须使用，如有不使用的变量编译时报错。
>
> 函数外的每个语句都必须以关键字开始 (var、const、func) 等
>
> :=不能使用在函数外。
>
> 匿名变量或者叫哑元变量(\_)多用于占位，表示忽略值，即当有些数据必须用变量接收但又不使用它时，可以采用\_来接收改值。
>
> 变量在同一个作用域中代码块({})中不能重复声明同名的变量。 



常量：在定义的时候必须赋值。 

```go
// 单一声明: 声明了pi和e这两个常量之后，在整个程序运行期间它们的值都不能再发生变化了。
const pi = 3.1415
const e = 2.7182

// 批量声明
const (
  pi = 3.1415
  e = 2.7182
)

// 批量声明（如果省略了值则表示和上面一行的值相同）
// 常量n1、n2、n3的值都是100。
const (
  n1 = 100
  n2
  n3
) 
```

常量计数器：iota是go语言的常量计数器，只能在常量的表达式中使用。

Tips : iota在const关键字出现时将被重置为0, const中每新增一行常量声明将使iota计数一次 (iota可理解为const语句块中的行索引)。

应用场景: 使用iota能简化定义，在定义枚举时很有用。

```go
//1.使用_跳过某些值
const (
  n1 = iota //0
  n2        //1
  _
  n4        //3
)

//2.iota声明中间插队
const (
  n1 = iota //0
  n2 = 100  //100
  n3 = iota //2
  n4        //3
)
const n5 = iota //0

//3.多个iota定义在一行
const (
  a, b = iota + 1, iota + 2 //1,2
  c, d                      //2,3
  e, f                      //3,4
)

//4.定义数量级   （这里的<<表示左移操作，1<<10表示将1的二进制表示向左移10位，也就是由1变成了10000000000，也就是十进制的1024。同理2<<2表示将2的二进制表示向左移2位，也就是由10变成了1000，也就是十进制的8。）
const (
  _  = iota
  KB = 1 << (10 * iota)
  MB = 1 << (10 * iota)
  GB = 1 << (10 * iota)
  TB = 1 << (10 * iota)
  PB = 1 << (10 * iota)
)

```

**示例演示:**

```go
package main

import "fmt"

// 单一常量声明
const pi = 3.1415926535898

// 批量常量声明
const (
  e    = 2.7182
  flag = false
)

// 特殊批量常量声明
const (
  a = 1
  b
  _
  c
)

// iota 常量计数器
const (
  _     = iota               // 0
  d, e1 = iota + 1, iota + 2 // 2,3 常量名称不能重复
  f, g  = iota + 1, iota + 2 // 3,4
)

const (
  _  = iota             // 0
  KB = 1 << (10 * iota) // 1024
  MB = 1 << (10 * iota)
  GB = 1 << (10 * iota)
  TB = 1 << (10 * iota)
  PB = 1 << (10 * iota)
)

func main() {
  fmt.Println("pi :", pi)
  fmt.Println("e :", e, " , flag:", false)
  fmt.Println("特殊批量常量声明:", a, b, c)
  fmt.Println("iota 常量计数器 :", d, e1, f, g)
  fmt.Println("文件体积大小 :", KB, MB, GB, TB, PB)
}
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
pi : 3.1415926535898
e : 2.7182  , flag: false
特殊批量常量声明: 1 1 1
iota 常量计数器 : 2 3 3 4
文件体积大小 : 1024 1048576 1073741824 1099511627776 1125899906842624
```

Tips : 常量声明后不能在程序中进行重新赋值更改。



基本数据类型

基本的整型、浮点型、布尔型、字符串外，还有数组、切片、结构体、函数、map、通道（channel）等。

整型

| 类型   | 描述                                                         |
| ------ | ------------------------------------------------------------ |
| uint8  | 无符号8位整型（0-255）byte                                   |
| uint16 | 无符号16位整型（0-65535）                                    |
| uint32 | 无符号32位整型（0-4294967295）                               |
| uint64 | 无符号64位整型（0-18446744073709551615）                     |
| int8   | 有符号8位整型（-128到127 ）                                  |
| int16  | 有符号16位整型（-32768到32767 ）short                        |
| int32  | 有符号32位整型（-2147483648到2147483647 ）                   |
| int64  | 有符号64位整型（-9223372036854775808到9223372036854775807）long |

特殊整型

| 类型    | 描述                               |
| ------- | ---------------------------------- |
| uint    | 32位系统是uint32，64位系统是uint64 |
| int     | 32位系统是int32，64位系统是int64   |
| uintptr | 无符号整型，用于存放一个指针       |

> 注意： 在使用int和 uint类型时，不能假定它是32位或64位的整型，而是考虑int和uint可能在不同平台上的差异。
>
> 获取对象的长度的内建len()函数返回的长度可以根据不同平台的字节长度进行变化。实际使用中，切片或 map 的元素数量等都可以用int来表示。在涉及到二进制传输、读写文件的结构描述时，为了保持文件的结构不会受到不同编译目标平台字节长度的影响，不要使用int和 uint。 

数字字面量语法（Number literals syntax） 

v := 0b00101101， 代表二进制的 101101，相当于十进制的 45。
v := 0o377，代表八进制的 377，相当于十进制的 255。
\- v := 0x1p\-2，代表十六进制的 1 除以 2²，也就是 0.25。
而且还允许我们用 _ 来分隔数字，比如说： v := 123_456 表示 v 的值等于 123456。 

示例：

```go
package main

import "fmt"

func main(){
  // 十进制以不同的进制展示
  var a int = 10
  fmt.Printf("%b \n", a)   // 1010  占位符%b表示二进制
  fmt.Printf("%o \n", a)   // 12    占位符%o表示八进制
  fmt.Printf("%d \n", a)   // 10    占位符%d表示十进制
  fmt.Printf("0x%x \n", a) // 0xa  占位符%x表示十六进制

  // 八进制(以0开头)
  var b int = 077
  fmt.Printf("%b \n", b)   // 111111
  fmt.Printf("%o \n", b)   // 77
  fmt.Printf("%d \n", b)   // 63
  fmt.Printf("0x%x \n", b) // 0x3f

  // 十六进制(以0x开头)
  var c int = 0xff
  fmt.Printf("0x%x \n", c) // 0xff
  fmt.Printf("0X%X \n", c) // 0xFF

  // 数字字面量语法（Number literals syntax）
  binary := 0b1111
  octal := 0o17
  digital := 15
  hexadecimal := 0xf
  specialhexa := 0x8p-2    // 8 / 2^2 = 2
  underline := 10_24

  fmt.Printf("binary : %b , digital ： %d\n", binary, binary)
  fmt.Printf("octal : %o , digital ： %d\n", octal, octal)
  fmt.Printf("digital type (变量类型): %T,digital ： %d\n", digital, digital)
  fmt.Printf("hexadecimal : %x, digital ： %d, specialhexa : %f\n", hexadecimal, hexadecimal, specialhexa)
  fmt.Printf("underline : %d \n", underline)
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
1010 
12 
10 
0xa 
111111 
77 
63 
0x3f 
0xff 
0XFF 
binary : 1111 , digital ： 15
octal : 17 , digital ： 15
digital type (变量类型): int,digital ： 15
hexadecimal : f, digital ： 15, specialhexa : 2.000000
underline : 1024 
```



浮点型

float32 的浮点数的最大范围约为 3.4e38，其常量定义：math.MaxFloat32。
float64 的浮点数的最大范围约为 1.8e308，其常量定义：math.MaxFloat64。 

```go
//打印浮点数时，可以使用fmt包配合动词%f 
package main
import (
    "fmt"
    "math"
)
func main() {
  	var floatnumber float64 = 1024.00
  	fmt.Printf("数据类型: %T , floatnumber: %.1f\n", floatnumber, floatnumber)
  	fmt.Printf("%f,%.2f\n", math.Pi, math.Pi) // 保留小数点后两位
 	fmt.Printf("float32的浮点数的最大范围 :%d ~ %f\n", 0, math.MaxFloat32)
  	fmt.Printf("float64的浮点数的最大范围 :%d ~ %f\n", 0, math.MaxFloat64)
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
数据类型: float64 , floatnumber: 1024.0
3.141593,3.14
float32的浮点数的最大范围 :0 ~ 340282346638528859811704183484516925440.000000
float64的浮点数的最大范围 :0 ~ 179769313486231570814527423731704356798070567525844996598917476803157260780028538760589558632766878171540458953514382464234321326889464182768467546703537516986049910576551282076245490090389328944075868508455133942304583236903222948165808559332123348274797826204144723168738177180919299881250404026184124858368.000000
```

复数

复数有实部和虚部，complex64的实部和虚部为32位，complex128的实部和虚部为64位。 

```go
var c1 complex64
c1 = 1 + 2i
var c2 complex128
c2 = 2 + 3i
fmt.Println(c1) // (1+2i)
fmt.Println(c2) // (2+3i) 
```

布尔值

```go
fmt.Println("布尔型示例:")
var flag bool = true
fmt.Printf("数据类型: %T ,任意类型输出: %v", flag, flag)  // 数据类型: bool ,任意类型输出: true 
```

>注意：
>
>布尔类型变量的默认值为false。Go 语言中不允许将整型强制转换为布尔型。布尔型无法参与数值运算，也无法与其他类型进行转换。

字符串

Go语言中的字符串以原生数据类型出现，使用字符串就像使用其他原生数据类型（int、bool、float32、float64 等）一样。字符串的内部实现使用UTF-8编码。 字符串的值为双引号(")中的内容，可以在Go语言的源码中直接添加非ASCII码字符 。

```go
s1 := "hello"
s2 := "你好"
c1 := 'g'
c2 := 'o'
```

Tips : Go 语言中用双引号包裹的是字符串，而单引号包裹的是字符。 

转义符

| 转义符 | 含义                     |
| ------ | ------------------------ |
| \r     | 回车符（返回行首）       |
| \n     | 换行符（下一行同列位置） |
| \t     | 制表符                   |
| \\'    | 单引号                   |
| \\"    | 双引号                   |
| \\\    | 反斜杠                   |

```go
package main
import (
    "fmt"
)
func main() {
  	s1 := "'c:\\weiyigeek\\go\\hello'"
    fmt.Println("str :=",s1)
    fmt.Println("str := \"c:\\Code\\weiyigeek\\go.exe\"")
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
str := 'c:\weiyigeek\go\hello'
str := "c:\Code\weiyigeek\go.exe"
```

多行字符串

```go
//定义一个多行字符串时，就必须使用反引号字符
s1 := `第一行
第二行
第三行
`
s2 := `c:\weiyigeek\go\hello`  // 注意点此处没用转义符(\) 也能输出路径
fmt.Println(s1,s2) 
```

```bash
[root@MyArch go]# go run "/root/go/src/test.go"
第一行
第二行
第三行
 c:\weiyigeek\go\hello
```



Tips: 反引号间换行将被作为字符串中的换行，但是所有的转义字符均无效，文本将会原样输出 

字符串常用操作

| 方法                                 | 介绍           |
| ------------------------------------ | -------------- |
| len(str)                             | 求长度         |
| +或fmt.Sprintf                       | 拼接字符串     |
| string.Split                         | 分割           |
| strings.Contains                     | 判断是否包含   |
| strings.HasPrefix,strings.HasSuffix  | 前缀/后缀判断  |
| strings.Index(),string.LastIndex()   | 子串出现的位置 |
| strings.Join(a[] string, sep string) | join操作       |

```go
// 字符串型示例
func stringdemo() {
  // 字符
  c1 := 'a'
  c2 := 'A'

  // 字符串 (单行与多行以及转义)
  s1 := "Name"
  s2 := "姓名"
  s3 := `
  这是一个
        多行字符串案例！
  This is mutlilineString Example！
  Let's Go   // 特点：原样输出
  `
  // 转义演示
  s4 := "'c:\\weiyigeek\\go\\hello'"
  s5 := `c:\weiyigeek\go\hello`

  fmt.Printf("c1 char : %c,\t c2 char %c -> digital : %d\n", c1, c2, c2)
  fmt.Println(s1, s2)
  fmt.Println(s3)
  fmt.Println(s4, s5)

  // 字符串常用函数
  fmt.Println("s1 String length:", len(s1), "s2 string length:", len(s2))

  info := fmt.Sprintf("%s (%s): %s", s1, s2, "WeiyiGeek")
  fmt.Println("Infomation : "+"个人信息", info)

  fmt.Println("字符串分割 :", strings.Split(s5, "\\"))

  fmt.Println("判断字符串是否包含go", strings.Contains(s3, "go"))

  fmt.Println(strings.HasPrefix(s1, "N"), strings.HasSuffix(s1, "e"))

  fmt.Println(strings.Index(s4, "weiyigeek"), strings.LastIndex(s4, "weiyigeek"))

  s6 := strings.Split(s5, "\\")
  fmt.Println("字符串间隔符 : ", strings.Join(s6, "-"))
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
c1 char : a,     c2 char A -> digital : 65
Name 姓名

  这是一个
        多行字符串案例！
  This is mutlilineString Example！
  Let's Go   // 特点：原样输出
  
'c:\weiyigeek\go\hello' c:\weiyigeek\go\hello
s1 String length: 4 s2 string length: 6
Infomation : 个人信息 Name (姓名): WeiyiGeek
字符串分割 : [c: weiyigeek go hello]
判断字符串是否包含go false
true true
4 4
字符串间隔符 :  c:-weiyigeek-go-hello
```

byte和rune类型 与 字符类型

组成每个字符串的元素叫做“字符”，可以通过遍历或者单个获取字符串元素获得字符。 字符用单引号（’）包裹起来 。

```go
var a = '中'
var b = 'x'
c := 'a' 
```

两种字符类型

> uint8类型，或者叫 byte 型，代表了ASCII码的一个字符（1B）。
> rune类型，代表一个 UTF-8字符, 并且一个rune字符由一个或多个byte组成（3B~4B）。 
>
> Tips : 当需要处理中文、日文或者其他复合字符时，则需要用到rune类型。rune类型实际是一个int32。

```go
// 遍历字符串
func traversalString() {
  	s := "hello沙河"
  
    // byte 类型
  	for i := 0; i < len(s); i++ {
    	fmt.Printf("%v(%c) ", s[i], s[i])
  	}
  	fmt.Println()
  
  	// rune 类型
  	for _, r := range s {
    	fmt.Printf("%v(%c) ", r, r)
  	}
  	fmt.Println()
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
104(h) 101(e) 108(l) 108(l) 111(o) 230(æ) 178(²) 153() 230(æ) 178(²) 179(³) 
104(h) 101(e) 108(l) 108(l) 111(o) 27801(沙) 27827(河) 
```

>Q: 为什么出现上述情况?
>
>答: 因为UTF8编码下一个中文汉字由3~4个字节（4*8bit）组成，所以我们不能简单的按照字节去遍历一个包含中文的字符串，否则就会出现上面输出中第一行的结果。

Tips : 字符串底层是一个byte数组，所以可以和[]byte类型相互转换。字符串是不能修改的字符串是由byte字节组成，所以字符串的长度是byte字节的长度。



类型转换

>Go语言中只有强制类型转换，没有隐式类型转换。该语法只能在两个类型之间支持相互转换的时候使用。
强制类型转换的基本语法如下：
T(表达式)  # 其中，T表示要转换的类型。表达式包括变量、复杂算子和函数返回值等.
Tips : Boolen 类型不能强制转换为整型。 

```sh
//计算直角三角形的斜边长时使用math包的Sqrt()函数，该函数接收的是float64类型的参数
func sqrtDemo() {
  var a, b = 3, 4
  var c int
  // math.Sqrt() 接收的参数是float64类型，需要强制转换
  c = int(math.Sqrt(float64(a*a + b*b)))
  fmt.Println(c)
} 
```

Tips : 在Go语言中不同类型的值不能直接赋值，例如float32类型变量a的值不能直接赋值给floa64类型变量b的值。 



字符串类型转换

>如果修改字符串，需要先将其转换成[]rune或[]byte，完成后再转换为string。无论哪种转换，都会重新分配内存，并复制字节数组。
>
>在一个字符串中如果既有中文也存在英文，我们则可以使用byte[]类型(1B)来存放ASCII码表示的字符(0~255)，如果是中文则使用rune\[\](4B)类型来存放或者周转。 

```go
func changeString() {
  // 强制类型转换
  s1 := "big" 
  byteS1 := []byte(s1)
  byteS1[0] = 'p'
  fmt.Println(string(byteS1))

  s2 := "白萝卜"
  runeS2 := []rune(s2)
  runeS2[0] = '红'
  fmt.Println(string(runeS2))
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
pig
红萝卜
```

示例：

```go
// Byte与Rune类型示例
func brdemo() {
  var c1 = 'a' // int32 类型
  var c2 = 'A' // int32 类型
  z1 := '中'    // int32 类型
  z2 := '文'    // int32 类型
  z3 := "中"    // string 类型 (双引号)

  // 字符不同格式输出
  fmt.Printf("字符 ：%d (%c) , %d (%c) \n", c1, c1, c2, c2)
  fmt.Printf("中文字符 ：%d (%v) = %c , %d (%v) = %c \n", z1, z1, z1, z2, z2, z2)
  fmt.Printf("单双引号不同类型 : c1 = %c (%T) , z2 = %c (%T) ,  z3 = %s (%T) \n", c1, c1, z2, z2, z3, z3)

  // 中英文字符串修改
  s1 := "a和我都爱中国"
  s2 := "为 Hello 中国 World,Go 语言 学习"

  // 将字符类型转化为byte类型
  c3 := byte(c2)
  fmt.Printf("强制转化类型 : c2 = %c (%T) , byte(c2) = %c (%T) \n", c2, c2, c3, c3)

  // 将字符串类型转化为string类型
  r1 := []rune(s1) // 强制转化字符串为一个rune切片
  r1[0] = '您'      // 注意此处需传入为字符
  fmt.Println("修改后中文字符串输出(未类型转换)：", r1)
  fmt.Println("修改后中文字符串输出(已类型转换)：", s1, string(r1)) // 强制转化rune切片为字符串

  // 将整型转化成为浮点数类型
  // 计算直角三角形的斜边长
  var a, b = 3, 4
  var c int = int(math.Sqrt(float64(a*a + b*b)))
  fmt.Println("计算直角三角形的斜边长 (a=3,b=4) c =", c)

  // 统计字符串中中文个数
  res := []rune(s2)
  reslen := len(res)
  count := 0
  for i := 0; i < reslen; i++ {
    if res[i] > 255 {
      count++
    }
  }
  fmt.Printf("字符串:%s (Length = %d),一共有 %d 个中文字符", s2, reslen, count)//无换行符
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
字符 ：97 (a) , 65 (A) 
中文字符 ：20013 (20013) = 中 , 25991 (25991) = 文 
单双引号不同类型 : c1 = a (int32) , z2 = 文 (int32) ,  z3 = 中 (string) 
强制转化类型 : c2 = A (int32) , byte(c2) = A (uint8) 
修改后中文字符串输出(未类型转换)： [24744 21644 25105 37117 29233 20013 22269]
修改后中文字符串输出(已类型转换)： a和我都爱中国 您和我都爱中国
计算直角三角形的斜边长 (a=3,b=4) c = 5
字符串:为 Hello 中国 World,Go 语言 学习 (Length = 25),一共有 7 个中文字符(base)
```

算数运算符

\+ - \* / % 加减乘除，求余

> 注意：++ --  在go语言中是单独的语句，并不是运算符。

关系运算符

== != > >= < <= 

> Tips : Go 语言是强类型的所以必须相同类型变量才能进行比较。 

逻辑运算符

&& || ! 	：and or not

位运算符

& | ^ ：与 或 异或 

<< >> ：左移，高位丢弃，低位补0。右移同理。2的n次方。

赋值运算符

= += -= \*= /= %= <<= >>= &= |= ^=

例子：

```go
a += 1  // a = a + 1
a %= 3  // a = a % 3
a <<= 4 // a = a << 4
a ^= 5  // a = a ^ 5 
```



流程控制

> Go语言中常用的是if和for，switch和goto是为了简化代码，降低重复率而生的扩展类流程控制



if else(分支结构) 

```go
if 表达式1 {
  分支1
} else if 表达式2 {
  分支2
} else{
  分支3
} 
```

> Go语言规定与if匹配的左括号{必须与if和表达式放在同一行，{放在其他位置会触发编译错误。 同理，与else匹配的{也必须与else写在同一行，else也必须与上一个if或else if右边的大括号在同一行。

```go
func ifDemo1() {
  score := 65
  if score >= 90 {
    fmt.Println("A")
  } else if score > 75 {
    fmt.Println("B")
  } else {
    fmt.Println("C") // 输出结果
  }
}
```

if条件判断特殊写法 ：可以在 if 表达式之前添加一个执行语句 

```go
func ifDemo2() {
  score := 88 // 注意变量作用域的影响
  if score := 65; score >= 90 {
    fmt.Println("A", score)
  } else if score > 75 {
    fmt.Println("B", score)
  } else {
    fmt.Println("C", score) // 输出结果
  }
  fmt.Println("score : ", score)	//score变量作用域只在if...else代码块中有效。 
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
C 65
score :  88
```



for(循环结构)

```go
for 初始语句;条件表达式;结束语句{
   循环体语句
} 
```

条件表达式返回true时循环体不停地进行循环，直到条件表达式返回false时自动退出循环。

```go
func forDemo() {  for i := 0; i < 10; i++ {    fmt.Println(i)  }}

//for循环的初始语句可以被忽略，但是初始语句后的分号必须要写，例如：
func forDemo2() {  i := 0  for ; i < 10; i++ {    fmt.Println(i)  }}

//for循环的初始语句和结束语句都可以省略，例如：
func forDemo3() {  i := 0  for i < 10 {    fmt.Println(i)    i++  }}

//for无限循环，这种写法类似于其他编程语言中的while，在while后添加一个条件表达式，满足条件表达式时持续循环，否则结束循环。
//例如: for循环可以通过break、goto、return、panic语句强制退出循环。
for {
  循环体语句
} 
```

for range(键值循环) 

描述: Go语言中可以使用for range遍历数组、切片、字符串、map 及通道（channel）。

通过for range遍历的返回值有以下规律：

 1. 数组、切片、字符串返回索引和值。

 2. map返回键和值。

 3. 通道（channel）只返回通道内的值。

```go
s1 := "Hello,Go 输出的是中文"
for i, v := range s1 {
  fmt.Printf("Index : %d ,Value : %s , Number : %v \n", i, string(v), v)
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
Index : 0 ,Value : H , Number : 72 
Index : 1 ,Value : e , Number : 101 
Index : 2 ,Value : l , Number : 108 
Index : 3 ,Value : l , Number : 108 
Index : 4 ,Value : o , Number : 111 
Index : 5 ,Value : , , Number : 44 
Index : 6 ,Value : G , Number : 71 
Index : 7 ,Value : o , Number : 111 
Index : 8 ,Value :   , Number : 32 
Index : 9 ,Value : 输 , Number : 36755 
Index : 12 ,Value : 出 , Number : 20986 
Index : 15 ,Value : 的 , Number : 30340 
Index : 18 ,Value : 是 , Number : 26159 
Index : 21 ,Value : 中 , Number : 20013 
Index : 24 ,Value : 文 , Number : 25991 
```

switch case(选择语句) 

```go
func switchDemo1() {
  finger := 3
  switch finger {
  case 1:
    fmt.Println("大拇指")
  case 2:
    fmt.Println("食指")
  case 3:
    fmt.Println("中指")
  case 4:
    fmt.Println("无名指")
  case 5:
    fmt.Println("小拇指")
  default:
    fmt.Println("无效的输入！")
  }
} 
```

> Go语言规定每个switch只能有一个default分支, 但一个分支可以有多个值，多个case值中间使用英文逗号分隔。 

```go
func testSwitch3() {
  switch n := 7; n {
  case 1, 3, 5, 7, 9:
    fmt.Println("奇数")
  case 2, 4, 6, 8:
    fmt.Println("偶数")
  default:
    fmt.Println(n)
  }
} 
```

分支还可以使用表达式，这时候switch语句后面不需要再跟判断变量。 

```go
func switchDemo4() {
  age := 30
  switch {
  case age < 25:
    fmt.Println("好好学习吧")
  case age > 25 && age < 35:
    fmt.Println("好好工作吧")
  case age > 60:
    fmt.Println("好好享受吧")
  default:
    fmt.Println("活着真好")
  }
} 
```

fallthrough语法: 可以执行满足条件的case的下一个case，是为了兼容C语言中的case设计的（值得学习）。

```go
func switchDemo5() {
  s := "a"
  switch {
  case s == "a":
    fmt.Println("a")
    fallthrough
  case s == "b":
    fmt.Println("b")
  case s == "c":
    fmt.Println("c")
  default:
    fmt.Println("...")
  }
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
a
b
```

goto(跳转到指定标签) 

例如：双层嵌套的for循环要退出时

```go
func gotoDemo1() {
  var breakFlag bool
  for i := 0; i < 10; i++ {
    for j := 0; j < 10; j++ {
      if j == 2 {
        // 设置退出标签
        breakFlag = true
        break
      }
      fmt.Printf("%v-%v\n", i, j)
    }
    // 外层for循环判断
    if breakFlag {
      break
    }
  }
} 
```

简化代码

```go
func gotoDemo2() {
  for i := 0; i < 10; i++ {
    for j := 0; j < 10; j++ {
      if j == 2 {
        // 设置退出标签
        goto breakTag
      }
      fmt.Printf("%v-%v\n", i, j)
    }
  }
  return
  // 标签
  breakTag:
    fmt.Println("正结束for循环")
    fmt.Println("已结束for循环")
} 
```

break(跳出循环)：break语句还可以在语句后面添加标签，表示退出某个标签对应的代码块，标签要求必须定义在对应的for、switch和 select的代码块上。

```go
func breakDemo1() {
BREAKDEMO1:
  for i := 0; i < 10; i++ {
    for j := 0; j < 10; j++ {
      if j == 2 {
        break BREAKDEMO1
      }
      fmt.Printf("%v-%v\n", i, j)
    }
  }
  fmt.Println("...")	//正常执行
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
0-0
0-1
...
```

continue(继续下次循环) ：在 continue语句后添加标签时，表示开始标签对应的循环。 

```go
func continueDemo() {
forloop1:
  for i := 0; i < 5; i++ {
    // forloop2:
    for j := 0; j < 5; j++ {
      if i == 2 && j == 2 {
        continue forloop1
      }
      fmt.Printf("%v-%v\n", i, j)
    }
  }
} 
```



数组

定义： 数组的长度必须是常量，并且长度是数组类型的一部分。一旦定义，长度不能变。注意: [5]int和[10]int是不同的类型。

```go
var 数组变量名 [元素数量]T

// 定义一个长度为3元素类型为int的数组a
var a [3]int 
```

> Tips ：数组可以通过下标进行访问，下标是从0开始，最后一个元素下标是：len-1，访问越界（下标在合法范围之外），则触发访问越界，会panic。 

初始化

初始化数组时可以使用初始化列表来设置数组元素的值。 

```go
func main() {
	var testArray [3]int                        //数组会初始化为int类型的零值
	var numArray = [3]int{1, 2}                 //使用指定的初始值完成初始化
	var cityArray = [3]string{"北京", "上海", "深圳"} //使用指定的初始值完成初始化
	fmt.Println(testArray)                      //[0 0 0]
	fmt.Println(numArray)                       //[1 2 0]
	fmt.Println(cityArray)                      //[北京 上海 深圳]
} 
```

让编译器根据初始值的个数自行推断数组的长度。

```go
func main() {
	var testArray [3]int
	var numArray = [...]int{1, 2}
	var cityArray = [...]string{"北京", "上海", "深圳"}
	fmt.Println(testArray)                          //[0 0 0]
	fmt.Println(numArray)                           //[1 2]
	fmt.Printf("type of numArray:%T\n", numArray)   //type of numArray:[2]int
	fmt.Println(cityArray)                          //[北京 上海 深圳]
	fmt.Printf("type of cityArray:%T\n", cityArray) //type of cityArray:[3]string
} 
```

指定索引值的方式来初始化数组。

```go
func main() {
	a := [...]int{1: 1, 3: 5}
  b := [...]int{1:100,9:200}      // [0 100 0 0 0 0 0 0 200 ]
	fmt.Println(a)                  // [0 1 0 5]
	fmt.Printf("type of a:%T\n", a) // type of a:[4]int
} 
```

数组的遍历

```go
func main() {
	var a = [...]string{"北京", "上海", "深圳"}
	// 方法1：for循环遍历
	for i := 0; i < len(a); i++ {
		fmt.Println(a[i])
	}

	// 方法2：for range遍历
	for index, value := range a {
		fmt.Println(index, value)
	}
} 
```

多维数组(嵌套数组)

```go
func main() {
	a := [3][2]string{
		{"北京", "上海"},
		{"广州", "深圳"},
		{"成都", "重庆"},
	}
	fmt.Println(a) //[[北京 上海] [广州 深圳] [成都 重庆]]
	fmt.Println(a[2][1]) //支持索引取值:重庆
} 
```

注意： 多维数组只有第一层可以使用...来让编译器推导数组长度。

```go
//支持的写法
a := [...][2]string{
	{"北京", "上海"},
	{"广州", "深圳"},
	{"成都", "重庆"},
}
//不支持多维数组的内层使用...
b := [3][...]string{
	{"北京", "上海"},
	{"广州", "深圳"},
	{"成都", "重庆"},
}
```

二维数组的遍历

```go
func main() {
	a := [3][2]string{
		{"北京", "上海"},
		{"广州", "深圳"},
		{"成都", "重庆"},
	}

	// 方式1. for range 方式
	for _, v1 := range a {
		for _, v2 := range v1 {
			fmt.Printf("%s\t", v2)
		}
		fmt.Println()
	}
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
北京    上海
广州    深圳
成都    重庆
```

数组是值类型：赋值和传参会复制整个数组。因此改变副本的值，不会改变本身的值。

```go
// go 语言中默认传参是值传递（拷贝原变量的值即 Ctrl+c 、Ctrl+v ）
func modifyArray(x [3]int) {
	x[0] = 100
}

func modifyArray2(x [3][2]int) {
	x[2][0] = 100
}
func main() {
	a := [3]int{10, 20, 30}
	modifyArray(a) //在modify中修改的是a的副本x，不会更改数组a的元素
	fmt.Println(a) //[10 20 30]

	b := [3][2]int{
		{1, 1},
		{1, 1},
		{1, 1},
	}
	modifyArray2(b) //在modify中修改的是b的副本x，不会更改数组b的元素
	fmt.Println(b)  //[[1 1] [1 1] [1 1]]
} 
```

> 注意：数组支持 “==“、”!=” 操作符，因为内存总是被初始化过的。[n]\*T表示指针数组，\*[n]T表示数组指针 。 
>

```go
package main

import "fmt"

func main() {
	// 定义一个长度为3元素类型为int的数组a
	var a [2]int      // 默认为0
	var a1 [2]string  // 默认为空
	var a2 [2]bool    // 默认为false
	var a3 [2]float64 // 默认为0

	fmt.Printf("a 数组类型 %T , 元素: %v\n", a, a)
	fmt.Printf("a1 数组类型 %T , 元素: %v\n", a1, a1)
	fmt.Printf("a2 数组类型 %T , 元素: %v\n", a2, a2)
	fmt.Printf("a3 数组类型 %T , 元素: %v\n", a3, a3)

	// 数组初始化
	// 方式1.使用初始化列表来设置数组元素的值
	var b = [3]int{1, 2} // 三个元素，未指定下标元素的其值为 0
	var c = [3]string{"Let's", "Go", "语言"}
	// 方式2.根据初始值的个数自行推断数组的长度
	var d = [...]float32{1.0, 2.0}
	e := [...]bool{true, false, false}
	// 方式3.使用指定索引值的方式来初始化数组
	var f = [...]int{1: 1, 3: 8} // 只有 下标为1的其值为1，下标为3的其值为8，初开之外都为0
	g := [...]string{"Weiyi", "Geek"}

	fmt.Printf("b 数组类型 %T , 元素: %v\n", b, b)
	fmt.Printf("c 数组类型 %T , 元素: %v\n", c, c)
	fmt.Printf("d 数组类型 %T , 元素: %v\n", d, d)
	fmt.Printf("e 数组类型 %T , 元素: %v\n", e, e)
	fmt.Printf("f 数组类型 %T , 元素: %v\n", f, f)
	fmt.Printf("f 数组类型 %T , 元素: %v\n", g, g)

	// 数组指定元素获取
	fmt.Println("c[1] 元素获取 : ", c[1])
	// 数组遍历
	// 方式1
	alen := len(c)
	for i := 0; i < alen; i++ {
		fmt.Printf("c[%d]: %s ", i, c[i])
	}
	fmt.Println()
	// 方式2
	for i, v := range c {
		fmt.Printf("c[%d]: %s ", i, v) // 注意如果是切片类型需要强转为string
	}
	fmt.Println()

	// 多维数组
	// 方式1
	s1 := [3][2]string{
		{"北京", "上海"},
		{"广州", "深圳"},
		{"成都", "重庆"},
	}

	// 方式2
	s2 := [...][2]string{
		{"Go", "C"},
		{"PHP", "Python"},
		{"Shell", "Groovy"},
	}
	fmt.Println(s1[2][1])    //支持索引取值:重庆
	fmt.Println(len(s1), s1) //[[北京 上海] [广州 深圳] [成都 重庆]]
	fmt.Println(len(s2), s2)

	// 多维数组遍历
	// 方式1
	s1len := len(s1)
	for i := 0; i < s1len; i++ {
		s1length := len(s1[i])
		for j := 0; j < s1length; j++ {
			fmt.Printf("s1[%d][%d] = %v ", i, j, s1[i][j])
		}
	}
	fmt.Println()

	// 方式2 （推荐方式）
	for i, v1 := range s2 {
		for j, v2 := range v1 {
			fmt.Printf("s2[%d][%d] = %v ", i, j, v2)
		}
	}
	fmt.Println()

	// 多维数组元素更改
	s1[1][0] = "Test"
	s1[1][1] = "Change"
	fmt.Println(s1)
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
a 数组类型 [2]int , 元素: [0 0]
a1 数组类型 [2]string , 元素: [ ]
a2 数组类型 [2]bool , 元素: [false false]
a3 数组类型 [2]float64 , 元素: [0 0]
b 数组类型 [3]int , 元素: [1 2 0]
c 数组类型 [3]string , 元素: [Let's Go 语言]
d 数组类型 [2]float32 , 元素: [1 2]
e 数组类型 [3]bool , 元素: [true false false]
f 数组类型 [4]int , 元素: [0 1 0 8]
f 数组类型 [2]string , 元素: [Weiyi Geek]
c[1] 元素获取 :  Go
c[0]: Let's c[1]: Go c[2]: 语言 
c[0]: Let's c[1]: Go c[2]: 语言 
重庆
3 [[北京 上海] [广州 深圳] [成都 重庆]]
3 [[Go C] [PHP Python] [Shell Groovy]]
s1[0][0] = 北京 s1[0][1] = 上海 s1[1][0] = 广州 s1[1][1] = 深圳 s1[2][0] = 成都 s1[2][1] = 重庆 
s2[0][0] = Go s2[0][1] = C s2[1][0] = PHP s2[1][1] = Python s2[2][0] = Shell s2[2][1] = Groovy 
[[北京 上海] [Test Change] [成都 重庆]]
```

切片：slice

因为数组的长度是固定的并且数组长度属于类型的一部分，所以数组有很多的局限性。 

```go
func arraySum(x [3]int) int{
    sum := 0
    for _, v := range x{
        sum = sum + v
    }
    return sum
} 
```

这个求和函数只能接受[3]int类型，其他的都不支持。 再比如，`a := [3]int{1, 2, 3}`

数组a中已经有三个元素了，我们不能再继续往数组a中添加新元素了, 所以为了解决上述问题我们引入了Python一样切片的编程语言特性。 

> 描述: 切片（Slice）是一个拥有相同类型元素的可变长度的序列。它是基于数组类型做的一层封装。
>
> 特点:
>
> 切片它非常灵活，支持自动扩容。
> 切片是一个引用类型，它的内部结构包含地址、长度和容量。切片一般用于快速地操作一块数据集合。

```go
var name []T

// 关键字解析
- name:表示变量名
- T:表示切片中的元素类型 
```

Tips : 在定义时可看出与数组定义var array [number]T间的区别，其不需要设置元素个数 

```go
func main() {
	// 声明切片类型
	var a []string              //声明一个字符串切片
	var b = []int{}             //声明一个整型切片并初始化
	var c = []bool{false, true} //声明一个布尔切片并初始化
	var d = []bool{false, true} //声明一个布尔切片并初始化
	fmt.Println(a)              //[]
	fmt.Println(b)              //[]
	fmt.Println(c)              //[false true]
	fmt.Println(a == nil)       //true
	fmt.Println(b == nil)       //false
	fmt.Println(c == nil)       //false
	// fmt.Println(c == d)      //切片是引用类型，不支持直接比较，只能和nil比较
} 
```

切片长度与容量 

> 描述: 切片拥有自己的长度和容量，我们可以通过使用内置的len()函数求长度，使用内置的cap()函数求切片的容量。 

```go
// 切片长度与容量
var lth = []int{}
var lth64 = []float64{1, 2, 3}
fmt.Println("切片长度", len(lth), ",切片容量", cap(lth))      // 切片长度 0 ,切片容量 0
fmt.Println("切片长度", len(lth64), ",切片容量", cap(lth64))  // 切片长度 3 ,切片容量 3 
```

切片表达式

> 描述: 切片表达式从字符串、数组、指向数组或切片的指针构造子字符串或切片。

它有两种变体：一种指定low和high两个索引界限值的简单的形式，另一种是除了low和high索引界限值外还指定容量的完整的形式。

简单切片表达式

> 描述: 切片的底层就是一个数组，所以我们可以基于数组通过切片表达式得到切片。 切片表达式中的low和high表示一个索引范围（左包含，右不包含），也就是下面代码中从数组a中选出1<=索引值<4的元素组成切片s，得到的切片长度=high-low，容量等于得到的切片的底层数组的容量。

```go
func main() {
	a := [5]int{1, 2, 3, 4, 5}
	s := a[1:3]  // s := a[low:high]
	fmt.Printf("s:%v len(s):%v cap(s):%v\n", s, len(s), cap(s)) // 5 - 1
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
s:[2 3] len(s):2 cap(s):4
```

为了方便起见，可以省略切片表达式中的任何索引。省略了low则默认为0；省略了high则默认为切片操作数的长度:

```go
a[2:]  // 等同于 a[2:len(a)]
a[:3]  // 等同于 a[0:3]
a[:]   // 等同于 a[0:len(a)] 
```

> 注意：对于数组或字符串，如果0 <= low <= high <= len(a)，则索引合法，否则就会索引越界（out of range）。
>
> Tips : 对切片再执行切片表达式时（切片再切片），high的上限边界是切片的容量cap(a)，而不是长度。
>
> 常量索引必须是非负的，并且可以用int类型的值表示;对于数组或常量字符串，常量索引也必须在有效范围内。如果low和high两个指标都是常数，它们必须满足low <= high。如果索引在运行时超出范围，就会发生运行时panic。 

```go
func main() {
	a := [5]int{1, 2, 3, 4, 5}
	s1 := a[1:3]  // s1 := a[low:high]
	fmt.Printf("s1:%v len(s1):%v cap(s1):%v\n", s1, len(s1), cap(s1))
	s2 := s1[3:4]  // 索引的上限是cap(s)而不是len(s)
	fmt.Printf("s2:%v len(s2):%v cap(s2):%v\n", s2, len(s2), cap(s2))
} 
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
s1:[2 3] len(s1):2 cap(s1):4
s2:[5] len(s2):1 cap(s2):1
```

完整切片表达式

> 描述: 对于数组，指向数组的指针，或切片a(注意不能是字符串)支持完整切片表达式：

`a[low : high : max]`

> 描述: 上面的代码会构造与简单切片表达式a[low: high]相同类型、相同长度和元素的切片。另外它会将得到的结果切片的容量设置为max-low。在完整切片表达式中只有第一个索引值（low）可以省，它默认为0。 

```go
func main() {
	a := [5]int{1, 2, 3, 4, 5}
	t := a[1:3:5]
	fmt.Printf("t:%v len(t):%v cap(t):%v\n", t, len(t), cap(t))
}
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
t:[2 3] len(t):2 cap(t):4
```

> Tips : 完整切片表达式需要满足的条件是0 <= low <= high <= max <= cap(a)，其他条件和简单切片表达式相同。 

切片遍历 ：和数组一样

```go
func main() {
	s := []int{1, 3, 5}
	for i := 0; i < len(s); i++ {
		fmt.Println(i, s[i])
	}
	for index, value := range s {
		fmt.Println(index, value)
	}
} 
```

切片的本质

描述: 切片的本质就是对底层数组的封装，它包含了三个信息：底层数组的指针、切片的长度（len）和切片的容量（cap）。

举个例子，现在有一个数组a := [8]int{0, 1, 2, 3, 4, 5, 6, 7}，切片s1 := a[:5]，相应示意图如下。

![WeiyiGeek.slice_01](https://i0.hdslb.com/bfs/article/9df3291e8cba1df403b00f4dd67ca100e0b52897.png@942w_447h_progressive.png)

切片s2 := a[3:6]，相应示意图如下：

![WeiyiGeek.slice_02](https://i0.hdslb.com/bfs/article/cd6e549aff74fe8aa8bc262566706e0b8f39b342.png@942w_447h_progressive.png)


Tips ： 由上面两图可知切片的容量是数组长度 - 切片数组起始索引下标，例如 a[1:] = 8 - 1 其容量为7 

make() 方法构造切片

> 描述: 我们上面都是基于数组来创建的切片，如果需要动态的创建一个切片，我们就需要使用内置的make()函数 

```go
make([]T, size, cap)

# 参数说明
- T:切片的元素类型
- size:切片中元素的数量
- cap:切片的容量 
```

```go
func main() {
	a := make([]int, 2, 10)
	fmt.Println(a)      //[0 0]
	fmt.Println(len(a)) //2
	fmt.Println(cap(a)) //10
}
```

上面代码中a的内部存储空间已经分配了10个，但实际上只用了2个。 容量并不会影响当前元素的个数，所以len(a)返回2，cap(a)则返回该切片的容量。 

append() 方法切片添加元素

> 描述: Go语言的内建函数append()可以为切片动态添加元素。 可以一次添加一个元素，可以添加多个元素，也可以添加另一个切片中的元素（后面加…）。 

```go
func main(){
	var s []int
	s = append(s, 1)        // [1]
	s = append(s, 2, 3, 4)  // [1 2 3 4]
	s2 := []int{5, 6, 7}
	s = append(s, s2...)    // [1 2 3 4 5 6 7]
}
```


注意： 通过var声明的零值切片可以在append()函数直接使用，无需初始化。

```go
var s []int
s = append(s, 1, 2, 3)
```


注意： 没有必要像下面的代码一样初始化一个切片再传入append()函数使用，

```go
s := []int{}  // 没有必要初始化
s = append(s, 1, 2, 3)

var s = make([]int)  // 没有必要初始化
s = append(s, 1, 2, 3)
```


描述: 每个切片会指向一个底层数组，这个数组的容量够用就添加新增元素。当底层数组不能容纳新增的元素时，切片就会自动按照一定的策略进行“扩容”，此时该切片指向的底层数组就会更换。“扩容”操作往往发生在append()函数调用时，所以我们通常都需要用原变量接收append函数的返回值。

```go
func main() {
	//append()添加元素和切片扩容
	var numSlice []int
	for i := 0; i < 10; i++ {
		numSlice = append(numSlice, i)
		fmt.Printf("%v  len:%d  cap:%d  ptr:%p\n", numSlice, len(numSlice), cap(numSlice), numSlice)
	}
}
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
[0]  len:1  cap:1  ptr:0xc0000ba000
[0 1]  len:2  cap:2  ptr:0xc0000ba030
[0 1 2]  len:3  cap:4  ptr:0xc0000be020
[0 1 2 3]  len:4  cap:4  ptr:0xc0000be020
[0 1 2 3 4]  len:5  cap:8  ptr:0xc0000c0040
[0 1 2 3 4 5]  len:6  cap:8  ptr:0xc0000c0040
[0 1 2 3 4 5 6]  len:7  cap:8  ptr:0xc0000c0040
[0 1 2 3 4 5 6 7]  len:8  cap:8  ptr:0xc0000c0040
[0 1 2 3 4 5 6 7 8]  len:9  cap:16  ptr:0xc0000c2000
[0 1 2 3 4 5 6 7 8 9]  len:10  cap:16  ptr:0xc0000c2000
```

>注意：append()函数将元素追加到切片的最后并返回该切片。切片numSlice的容量按照1，2，4，8，16这样的规则自动进行扩容，每次扩容后都是扩容前的2倍。

append() 函数还支持一次性追加多个元素 

```sh
var citySlice []string
// 追加一个元素
citySlice = append(citySlice, "北京")
// 追加多个元素
citySlice = append(citySlice, "上海", "广州", "深圳")
// 追加切片
a := []string{"成都", "重庆"}
citySlice = append(citySlice, a...)
fmt.Println(citySlice)       //  [北京 上海 广州 深圳 成都 重庆] 
```

copy()方法复制切片 

> Tips : 由于切片是引用类型，所以a和b其实都指向了同一块内存地址。修改b的同时a的值也会发生变化。 

```go
func main() {
	a := []int{1, 2, 3, 4, 5}
	b := a
	fmt.Println(a) //[1 2 3 4 5]
	fmt.Println(b) //[1 2 3 4 5]
	b[0] = 1000
	fmt.Println(a) //[1000 2 3 4 5]
	fmt.Println(b) //[1000 2 3 4 5]
} 
```

Go语言内建的copy()函数可以迅速地将一个切片的数据复制到另外一个切片空间中 

```go
copy(destSlice, srcSlice []T)

# 参数:
- srcSlice: 数据来源切片
- destSlice: 目标切片
```

```go
func main() {
	// copy()复制切片
	a := []int{1, 2, 3, 4, 5}
	c := make([]int, 5, 5)
	copy(c, a)     //使用copy()函数将切片a中的元素复制到切片c
	fmt.Println(a) //[1 2 3 4 5]
	fmt.Println(c) //[1 2 3 4 5]
	c[0] = 1000
	fmt.Println(a) //[1 2 3 4 5]
	fmt.Println(c) //[1000 2 3 4 5]
} 
```

从切片中删除元素

描述: Go语言中并没有删除切片元素的专用方法，我们可以使用切片本身的特性来删除元素。

```go
func main() {
	// 从切片中删除元素
	a := []int{30, 31, 32, 33, 34, 35, 36, 37}
	// 要删除索引为2的元素
	a = append(a[:2], a[3:]...)
	fmt.Println(a) // [30 31 33 34 35 36 37]
} 
```

总结一下就是：要从切片a中删除索引为index的元素，操作方法是a = append(a[:index], a[index+1:]...) 

切片相关操作

判断切片是否为空

> 描述: 要检查切片是否为空，请始终使用len(s) == 0来判断，而不应该使用s == nil来判断。

```go
d := [5]int{1, 2, 3, 4, 5}
// 判断切片是否为空
if len(d) != 0 {
  fmt.Println("变量 d 切片不为空: ", d)
}
```

切片不能直接比较

> 描述: 切片之间是不能比较的，我们不能使用==操作符来判断两个切片是否含有全部相等元素。 切片唯一合法的比较操作是和nil比较。 一个nil值的切片并没有底层数组，一个nil值的切片的长度和容量都是0。

但是我们不能说一个长度和容量都是0的切片一定是nil

```go
var s1 []int         //len(s1)=0;cap(s1)=0;s1==nil
s2 := []int{}        //len(s2)=0;cap(s2)=0;s2!=nil
s3 := make([]int, 0) //len(s3)=0;cap(s3)=0;s3!=nil
```

所以要判断一个切片是否是空的，要是用len(s) == 0来判断，不应该使用s == nil来判断。

切片的赋值拷贝

> 描述: 下面的代码中演示了拷贝前后两个变量共享底层数组，对一个切片的修改会影响另一个切片的内容，这点需要特别注意。

```go
func main() {
	s1 := make([]int, 3) //[0 0 0]
	s2 := s1             //将s1直接赋值给s2，s1和s2共用一个底层数组
	s2[0] = 100
	fmt.Println(s1) //[100 0 0]
	fmt.Println(s2) //[100 0 0]
}
```


切片的扩容策略

> 描述: 可以通过查看$GOROOT/src/runtime/slice.go源码，其中扩容相关代码如下：

```go
newcap := old.cap
doublecap := newcap + newcap
if cap > doublecap {
	newcap = cap
} else {
	if old.len < 1024 {
		newcap = doublecap
	} else {
		// Check 0 < newcap to detect overflow
		// and prevent an infinite loop.
		for 0 < newcap && newcap < cap {
			newcap += newcap / 4
		}
		// Set newcap to the requested cap when
		// the newcap calculation overflowed.
		if newcap <= 0 {
			newcap = cap
		}
	}
}
```


从上面的代码可以看出以下内容：

首先判断，如果新申请容量（cap）大于2倍的旧容量（old.cap），最终容量（newcap）就是新申请的容量（cap）。
否则判断，如果旧切片的长度小于1024，则最终容量(newcap)就是旧容量(old.cap)的两倍，即（newcap=doublecap），
否则判断，如果旧切片的长度大于等于1024，则最终容量（newcap）从旧容量（old.cap）开始循环增加原来的1/4，即（newcap=old.cap,for {newcap += newcap/4}）直到最终容量（newcap）大于等于新申请的容量(cap)，即（newcap >= cap）
如果最终容量（cap）计算值溢出，则最终容量（cap）就是新申请容量（cap）。

> Tips : 需要注意的是，切片扩容还会根据切片中元素的类型不同而做不同的处理，比如int和string类型的处理方式就不一样。

```go
package main

import "fmt"

func main() {
	// 切片声明与定义
	var a []string              //声明一个字符串切片
	var b = []int{}             //声明一个整型切片并初始化
	var c = []bool{false, true} //声明一个布尔切片并初始化

	// - 切片 a 变量值为空/零值。
	if a == nil {
		fmt.Println("a 切片元素:", a)
	}
	fmt.Println("b 切片元素:", b)
	fmt.Println("c 切片元素:", c)

	// 切片长度与容量
	var lth = []int{}
	var lth64 = []float64{1, 2, 3}
	fmt.Println("切片长度", len(lth), ",切片容量", cap(lth))
	fmt.Println("切片长度", len(lth64), ",切片容量", cap(lth64))

	// 切片表达式
	d := [5]int{1, 2, 3, 4, 5}
	s := [5]string{"Let", "'s", "Go", "语言", "学习"}
	s1 := d[1:3]   // s := d[low(包含):high(不包含)] == d[1] d[2]
	s2 := d[2:]    // 等同于 a[2:5]  == d[2] d[3] d[4]
	s3 := d[:3]    // 等同于 a[0:3]  == d[0] d[1] d[2]
	s4 := d[:]     // 等同于 a[0:5]  == d[0] d[1] d[2] d[3] d[4]
	s5 := s[1:4:5] // 等同于 s[1:4] == s[1] s[2] s[3]

	fmt.Printf("s1:%v len(s1):%v cap(s1):%v\n", s1, len(s1), cap(s1)) // 注意此种情况 { 2 .. 5 容量为 4 }
	fmt.Printf("s2:%v len(s2):%v cap(s2):%v\n", s2, len(s2), cap(s2)) // { 3 .. 5 容量为 3 }
	fmt.Printf("s3:%v len(s3):%v cap(s3):%v\n", s3, len(s3), cap(s3)) // 注意此种情况 { 1 .. 5 容量为 5 }
	fmt.Printf("s4:%v len(s4):%v cap(s4):%v\n", s4, len(s4), cap(s4)) // { 1 .. 5 容量为 5}
	fmt.Printf("s5:%v len(s5):%v cap(s5):%v\n", s5, len(s5), cap(s5)) // s5:['s Go 语言] len(s5):3 cap(s5):4

	// 判断切片是否为空
	if len(d) != 0 {
		fmt.Println("变量 d 切片不为空: ", d)
	}

	// 切片遍历
	for i, v := range s {
		fmt.Printf("i: %d, v: %v , 切片指针地址: %p \n", i, v, &v)
	}
	fmt.Println()

	// make() 构造切片
	e := make([]int, 2, 10)
	fmt.Printf("e:%v len(e):%d cap(e):%d \n", e, len(e), cap(e)) // 长度 2，容量为 10

	// append() 添加元素 {7,8,9}
	f := append(e, 7, 8, 9)                                      // f:[0 0 7 8 9] len(f):5 cap(f):10
	f = append(f, e...)                                          // 追加切片
	fmt.Printf("f:%v len(f):%d cap(f):%d \n", f, len(f), cap(f)) // 长度 7，容量为 10

	// copy() 复制切片
	slice1 := []int{1, 2, 3, 4, 5}
	slice2 := make([]int, 7, 7)
	copy(slice2, slice1)
	slice2[6] = 2048
	fmt.Println("slice1 : ", slice1, "\t slice2 :", slice2)

	// 切片赋值拷贝
	slice3 := make([]int, 3)
	slice4 := slice3
	slice4[0] = 1024
	slice4[2] = 4096
	fmt.Printf("slice3 : %v, ptr : %p \n", slice3, slice3)
	fmt.Printf("slice4 : %v, ptr : %p \n", slice4, slice4)
} 
```
```sh
(base) [root@MyArch runtime]# go run "/root/go/src/test.go"
a 切片元素: []
b 切片元素: []
c 切片元素: [false true]
切片长度 0 ,切片容量 0
切片长度 3 ,切片容量 3
s1:[2 3] len(s1):2 cap(s1):4
s2:[3 4 5] len(s2):3 cap(s2):3
s3:[1 2 3] len(s3):3 cap(s3):5
s4:[1 2 3 4 5] len(s4):5 cap(s4):5
s5:['s Go 语言] len(s5):3 cap(s5):4
变量 d 切片不为空:  [1 2 3 4 5]
i: 0, v: Let , 切片指针地址: 0xc000010260 
i: 1, v: 's , 切片指针地址: 0xc000010260 
i: 2, v: Go , 切片指针地址: 0xc000010260 
i: 3, v: 语言 , 切片指针地址: 0xc000010260 
i: 4, v: 学习 , 切片指针地址: 0xc000010260 

e:[0 0] len(e):2 cap(e):10 
f:[0 0 7 8 9 0 0] len(f):7 cap(f):10 
slice1 :  [1 2 3 4 5]    slice2 : [1 2 3 4 5 0 2048]
slice3 : [1024 0 4096], ptr : 0xc0000160f0 
slice4 : [1024 0 4096], ptr : 0xc0000160f0 
```

Tips 总结: 数组是值类型，且包含元素的类型和元素个数，需注意元素的个数(数组长度)属于数组类型的一部分。 

Map 映射

> 描述: Map 是一种无序的基于key-value的数据结构, 并且它是引用类型，所以必须初始化值周才能进行使用。

```go
map[KeyType]ValueType

// # 参数说明:
- KeyType:表示键的类型。
- ValueType:表示键对应的值的类型。
```

>Tips : map类型的变量默认初始值为nil，需要使用make()函数来分配内存。语法为：make(map[KeyType]ValueType, [cap]), 其中cap表示map的容量，该参数虽然不是必须的，但是我们应该在初始化map的时候就为其指定一个合适的容量。

Map 基础使用

```go
// 1.采用Make初始化Map类型的变量。
scoreMap := make(map[string]int, 8)
scoreMap["小明"] = 100
fmt.Println(scoreMap["小明"])
fmt.Printf("type of a:%T\n", scoreMap)

// 2.在声明时填充元素。
userInfo := map[string]string{
  "username": "WeiyiGeek",
  "password": "123456",
}
fmt.Println(userInfo) 
```

Map 键值遍历

> 描述: 在进行Map类型的变量遍历之前，我们先学习判断map中键是否存在。

(1) 键值判断 

> 描述: 判断Map中某个键是否存在可以采用如下特殊写法: value, ok := map[key]

```go
scoreMap := make(map[string]int)
scoreMap["小明"] = 100
value, ok := scoreMap["张三"]
if ok {
  fmt.Println("scoreMap 存在该 '张三' 键")
} else {
  fmt.Println("scoreMap 不存在该键值")
}
```


(2) 键值遍历
描述: Go 语言中不像Python语言一样有多种方式进行遍历, 大道至简就 for...range 遍历 Map 就可以搞定。

```go
scoreMap := make(map[string]int)
scoreMap["Go"] = 90
scoreMap["Python"] = 100
scoreMap["C++"] = 60
// 遍历 k-v 写法
for k, v := range scoreMap {
  fmt.Println(k, v)
}

// 遍历 k 写法
for k := range scoreMap {
  fmt.Println(k)
}

// 遍历 v 写法
for _, v := range scoreMap {
  fmt.Println(v)
}
```


Tips ：遍历map时的元素顺序与添加键值对的顺序无关。

4.Map 键值删除

> 描述: 我们可使用 delete() 内建函数 从map中删除一组键值对, delete() 函数的格式如下: delete(map, key)其中 map:表示要删除键值对的map, key: 表示要删除的键值对的键。

```go
scoreMap := make(map[string]int)
scoreMap["张三"] = 90
scoreMap["小明"] = 100
delete(scoreMap, "小明" )  // 将`小明:100`从map中删除
for k,v := range scoreMap{
  fmt.Println(k, v)
}
```


5.值为map类型的切片

> 描述: 第一次看到时可能比较绕，其实可以看做在切片中存放Map类型变量。

```go
func demo3() {
	var mapSlice = make([]map[string]string, 3)
	for index, value := range mapSlice {
		fmt.Printf("index:%d value:%v\n", index, value)
	}
	fmt.Println()
	// 对切片中的map元素进行初始化
	mapSlice[0] = make(map[string]string, 10)
	mapSlice[1] = make(map[string]string, 10)
	mapSlice[2] = make(map[string]string, 10)
	mapSlice[0]["name"] = "WeiyiGeek"
	mapSlice[0]["sex"] = "Man"
	mapSlice[1]["姓名"] = "极客"
	mapSlice[1]["性别"] = "男"
	mapSlice[2]["hobby"] = "Computer"
	mapSlice[2]["爱好"] = "电脑技术"
	for i, v := range mapSlice {
		//fmt.Printf("index:%d value:%v\n", i, v)
		for _, value := range v {
			fmt.Printf("index:%d value:%v\n", i, value)
		}
	}
}
```

```sh
[root@MyArch runtime]# go run "/root/go/src/test.go"
index:0 value:map[]
index:1 value:map[]
index:2 value:map[]

index:0 value:WeiyiGeek
index:0 value:Man
index:1 value:男
index:1 value:极客
index:2 value:电脑技术
index:2 value:Computer
```


6.值为切片类型的map

> 描述: 同样在Map中存放切片类型的数据。

```go
// 值为切片类型的map
func demo4() {
	var sliceMap = make(map[string][]string, 3)
	var key = [2]string{"Country", "City"}
	fmt.Println("初始化 sliceMap 其值 : ", sliceMap)

	for _, v := range key {
  // 判断键值是否存在如果不存在则初始化一个容量为2的切片
  value, ok := sliceMap[v]
  if !ok {
    value = make([]string, 0, 2)
  }
  if v == "Country" {
    value = append(value, "中国")
  } else {
    value = append(value, "北京", "上海", "台湾")
  }
  // 将切片值赋值给Map类型的变量
  sliceMap[v] = value
} 
```
> 初始化 sliceMap 其值 :  map[]
> map[City:[北京 上海 台湾] Country:[中国]]


Tips : 非常重要、重要 Slice切片与Map 在使用时一定要做初始化操作(在内存空间申请地址)。

7.示例演示

1.Map类型的基础示例

```go
func demo1() {
	// 1.Map 定义
	var a1 map[string]int8  // (未分配内存)
	fmt.Println("Map 类型 的 a1 变量 :", a1)
	if a1 == nil {
		fmt.Println("默认初始化的Map类型的a1变量值: nil")
	}

	// 2.基本使用利用Make进行分配内存空间存储Map。
	b1 := make(map[string]string, 8)
	b1["姓名"] = "WeiyiGeek"
	b1["性别"] = "男|man"
	b1["爱好"] = "计算机技术"
	b1["出生日期"] = "2021-08-08"
	// 指定输出
	fmt.Printf("b1['姓名'] = %v \n", b1["姓名"])
	// 整体输出
	fmt.Printf("Map b1 Type: %T , Map b1 Value: %v \n", b1, b1)

	// 3.在声明时填充元素。
	c1 := map[string]string{
		"username": "WeiyiGeek",
		"sex":      "Man",
		"hobby":    "Computer",
	}
	// 指定输出
	fmt.Printf("c1['username'] = %v \n", c1["username"])
	// 整体输出
	fmt.Printf("Map c1 Type: %T , Length : %d , Map c1 Value: %v \n", c1, len(c1), c1)

	// 4.判断c1中的键值时候是否存在 sex Key.
	value, ok := c1["sex"]
	if ok {
		fmt.Println("c1 Map 变量中存在 'sex' 键 = ", value)
	} else {
		fmt.Println("c1 Map 变量中不存在 sex 键")
	}

	// 5.遍历Map
	for k, v := range b1 {
		fmt.Println(k, "=", v)
	}

	// 6.删除指定键值对，例如删除c1中的hobby键值。
	delete(c1, "hobby")
	fmt.Printf("Map 现存在的键 : ")
	for k := range c1 {
		fmt.Print(k, " ")
	}
} 
```
```sh
[root@MyArch runtime]# go run "/root/go/src/test.go"
Map 类型 的 a1 变量 : map[]
默认初始化的Map类型的a1变量值: nil
b1['姓名'] = WeiyiGeek 
Map b1 Type: map[string]string , Map b1 Value: map[出生日期:2021-08-08 姓名:WeiyiGeek 性别:男|man 爱好:计算机技术] 
c1['username'] = WeiyiGeek 
Map c1 Type: map[string]string , Length : 3 , Map c1 Value: map[hobby:Computer sex:Man username:WeiyiGeek] 
c1 Map 变量中存在 'sex' 键 =  Man
爱好 = 计算机技术
出生日期 = 2021-08-08
姓名 = WeiyiGeek
性别 = 男|man
Map 现存在的键 : username sex
```

2.按照指定顺序遍历map

```go
func demo2() {
	rand.Seed(time.Now().UnixNano()) //初始化随机数种子

	// 申请并初始化一个长度为 200 的 Map
	var scoreMap = make(map[string]int, 200)
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("stu%02d", i) //生成stu开头的字符串
		value := rand.Intn(100)          //生成0~99的随机整数
		scoreMap[key] = value
	}

	//取出map中的所有key存入切片keys
	var keys = make([]string, 0, 200)
	for key := range scoreMap {
		keys = append(keys, key)
	}

	//对切片进行排序
	sort.Strings(keys)

	//按照排序后的key遍历map
	for _, key := range keys {
		fmt.Println(key, scoreMap[key])
	}
} 
```
```sh
[root@MyArch runtime]# go run "/root/go/src/test.go"
stu00 86
stu01 87
stu02 4
stu03 53
stu04 26
stu05 71
stu06 92
stu07 63
stu08 77
stu09 19
stu10 49
stu11 14
stu12 24
stu13 53
stu14 46
stu15 1
stu16 33
stu17 4
stu18 39
stu19 19
```


Tips : 探究上述示例中Array 数组、Slice 切片、Map 映射有序与无序输出演示。

```go
func demo5() {
	// Array
	var Arr = [...]int{1, 2, 6, 4, 5}
	// Slice
	var Sli = []int{1, 2, 6, 4, 5}
	// Map
	var Map = map[string]int{
		"a1": 1,
		"b2": 2,
		"c3": 3,
		"d6": 6,
		"e5": 5,
	}

	fmt.Printf("Type : %T, Value : %v \n", Arr, Arr)
	for _, A := range Arr {
		fmt.Printf("%v ", A)
	}
	fmt.Println()
	fmt.Printf("Type : %T, Value : %v \n", Sli, Sli)
	for _, S := range Sli {
		fmt.Printf("%v ", S)
	}
	fmt.Println()
	fmt.Printf("Type : %T, Value : %v \n", Map, Map)
	for _, M := range Map {
		fmt.Printf("%v ", M)
	}
} 
```
```sh
[root@MyArch runtime]# go run "/root/go/src/test.go"
Type : [5]int, Value : [1 2 6 4 5] 
1 2 6 4 5 
Type : []int, Value : [1 2 6 4 5] 
1 2 6 4 5 
Type : map[string]int, Value : map[a1:1 b2:2 c3:3 d6:6 e5:5] 
1 2 3 6 5 
```

指针

> 描述: Go 语言中的指针区别于C/C++中的指针，Go语言中的指针不能进行偏移和运算是安全指针。

Go 语言中三个重要概念: 指针地址、指针类型以及指针取值。

简单回顾: 任何程序数据载入内存后，在内存都有他们的地址这就是指针。而为了保存一个数据在内存中的地址，我们就需要指针变量。

比如，“永远不要高估自己”这句话是我的座右铭，我想把它写入程序中，程序一启动这句话是要加载到内存（假设内存地址0x123456），我在程序中把这段话赋值给变量A，把内存地址赋值给变量B。这时候变量B就是一个指针变量, 通过变量A和变量B都能找到我的座右铭。

Go语言中的指针操作非常简单，我们只需要记住两个符号：&（取地址） 和 *（根据地址取值）。 



指针地址

> 描述: 每个变量在运行时都拥有一个地址，该地址代表变量在内存中的位置。

Go语言中使用&字符放在变量前面对变量进行“取地址”操作。 

```go
ptr := &v    // v的类型为T
// # 参数
// v:代表被取地址的变量，类型为T
// ptr:用于接收地址的变量，ptr的类型就为*T，称做T的指针类型。*代表指针。 
```

指针类型

> 描述: Go语言中的值类型（int、float、bool、string、array、struct）都有对应的指针类型，如：*int、*int64、*string等。

```go
func main() {
    a := 10
    b := &a
    fmt.Printf("a:%d ptr:%p\n", a, &a)             // a:10 ptr:0xc00001a078 (指针地址)
    fmt.Printf("*b:%d ptr:%p type:%T\n",*b, b, b)  // b:10 ptr:0xc00001a078 type:*int (指针类型)
    fmt.Printf("&b ptr:%p ",&b)                     // &b ptr:0xc00000e018
} 
```

指针取值

> 描述: 在对普通变量使用&操作符取地址后会获得这个变量的指针，然后可以对指针使用*操作，也就是指针取值，代码如下。

```go
func main() {
	//指针取值
	a := 10
	b := &a // 取变量a的地址，将指针保存到b中
	fmt.Printf("type of b:%T\n", b)
	c := *b // 指针取值（根据指针去内存取值）
	fmt.Printf("type of c:%T\n", c)
	fmt.Printf("value of c:%v\n", c)
}
```

指针特性

> 描述: 通过上面的指标变量、类型、取值的学习，我们了解到取地址操作符&和取值操作符*是一对互补操作符,其中&取出地址，*根据地址取出地址指向的值。

> Tips : 变量、指针地址、指针变量、取地址、取值的相互关系和特性如下：

1.对变量进行取地址（&）操作，可以获得这个变量的指针变量（指针地址）。
2.对指针变量进行取值（*）操作，可以获得指针变量指向的原变量的值。

```go
func modify1(x int) {
	x = 100
}

func modify2(x *int) {
	*x = 100
}

func main() {
    a := 10
    modify1(a)
    fmt.Println(a) // 10
    modify2(&a)
    fmt.Println(a) // 100
} 
```

内存地址分配

> 描述: 在Go语言中对于引用类型的变量，我们在使用的时候不仅要声明它，还要为它分配内存空间，否则我们的值就没办法存储。而对于值类型的声明不需要分配内存空间，是因为它们在声明的时候已经默认分配好了内存空间。

> Tips ：Go语言中new和make是内建的两个函数，他主要用来分配内存。

例如:执行下述例子中的代码会引发panic错误

```go
func main() {
  // 声明
  var a *int
   // 定义
  *a = 100

  fmt.Println(*a)

  var b map[string]int
  b["沙河娜扎"] = 100
  fmt.Println(b)
} 
```

New 函数

> 描述: new是Go语言的一置的函数它的函数签名如下：

`func new(Type) *Type` 其中，Type 表示类型，new 函数只接受一个参数，这个参数是一个类型，*Type 表示类型指针，new 函数返回一个指向该类型内存地址的指针。


Tips ：New 函数不太常用但由它可以得到一个类型的指针，并且该指针对应的值应该为该类型的零值。

```go
func main() {
 // 只是声明了一个指针变量a但是没有初始化
  a := new(int)
  b := new(bool)
  fmt.Printf("%T\n", a) // *int
  fmt.Printf("%T\n", b) // *bool
  fmt.Println(*a)       // 0
  fmt.Println(*b)       // false
} 
```

Tips : 指针作为引用类型需要初始化后才会拥有内存空间才可以给它赋值,所以需要按照下述方式使用内置的new函数对a进行初始化之后就可正常对其赋值了。

```go
func main() {
  var a *int
  a = new(int)
  *a = 10
  fmt.Println(*a)
} 
```

make 函数

> 描述: make也是用于内存分配的，区别于new，它只用于slice、map以及chan的内存创建，而且它返回的类型就是这三个类型本身，而不是他们的指针类型，因为这三种类型就是引用类型，所以就没有必要返回他们的指针了。

函数签名如下：

`func make(t Type, size ...IntegerType) Type`

Tips : Type 主要是 slice、map 以及channel类型，并且必须使用make进行初始化后，才能对它进行操作。

```go
func main() {
 // 只是声明变量b是一个map类型的变量
  var b map[string]int
   //使用make函数进行初始化操作之后
  b = make(map[string]int, 10)
   //才能对其进行键值对赋值：
  b["WeiyiGeek"] = 100
  fmt.Println(b)
} 
```

总结:new 函数与 make函数的区别

> 二者都是用来做内存分配的。
> make只用于slice、map以及channel的初始化，返回的还是这三个引用类型本身；
> new用于类型的内存分配，并且内存对应的值为类型零值，返回的是指向类型的指针。

```go
// 转入int类型的参数
func normal(x int) {
	x = 65535
	fmt.Printf("Func Param &x ptr : %p \n", &x)
}

// 传入的参数为指针类型
func pointer(x *int) {
	*x = 65535
	fmt.Printf("Func Param x ptr : %p \n", x)
}

func demo1() {
	// 1.2获得变量a的内存地址
	a := 1024
	b := &a
	fmt.Printf("a : %d , a ptr: %p, b ptr : %v , *b = %d \n", a, &a, b, *b)
	fmt.Printf("b type: %T, &b ptr : %p \n", b, &b)
	fmt.Println()
	// 2.针对变量a的内存地址进行重赋值（此时会覆盖变量a的原值）
	*b = 2048
	fmt.Printf("Change -> a : %d , a ptr: %p, b ptr : %v , *b = %d \n", a, &a, b, *b)
	fmt.Printf("b type: %T, &b ptr : %p \n\n", b, &b)

	// 3.指针传值
	c := 4096
	normal(c)
	fmt.Println("After Normal Function c : ", c)
	pointer(&c)
	fmt.Printf("After Pointer Function c : %v, c ptr: %p \n\n", c, &c)

	// 4.new 内存地址申请
	var a4 *int
	//*a4 = 100 // 此行会报 _panic 错误，因为未分配内存空间
	fmt.Println("a4 ptr : ", a4) // 空指针 （<nil>）还没有内存地址

	d := new(int)                         // 申请一块内存空间 （内存地址）
	fmt.Printf("%T ，%p, %v \n", d, d, *d) // 其指针类型默认值为 0 与其类型相关联。
	*d = 8192                             // 对该内存地址赋值
	fmt.Printf("%T ，%p, %v \n\n", d, d, *d)

	// 5.make 内存地址申请
	var b5 map[string]string
	//b5["Name"] = "WeiyGeek" //此行会报 _panic 错误，因为未分配内存空间
	fmt.Printf("%T , %p , %v\n", b5, &b5, *&b5)

	b5 = make(map[string]string, 10) // 申请一块内存空间 （内存地址）
	b5["Name"] = "WeiyGeek"          // 此时便可对该Map类型进行赋值了
	b5["Address"] = "ChongQIng China"
	fmt.Printf("%T , %p , %v\n\n", b5, &b5, b5)
}
```

```sh
[root@MyArch go]# go run "/root/go/src/test.go"
a : 1024 , a ptr: 0xc0000140c0, b ptr : 0xc0000140c0 , *b = 1024 
b type: *int, &b ptr : 0xc00000e028 

Change -> a : 2048 , a ptr: 0xc0000140c0, b ptr : 0xc0000140c0 , *b = 2048 
b type: *int, &b ptr : 0xc00000e028 

Func Param &x ptr : 0xc0000140f8 
After Normal Function c :  4096
Func Param x ptr : 0xc0000140f0 
After Pointer Function c : 65535, c ptr: 0xc0000140f0 

a4 ptr :  <nil>
*int ，0xc000014110, 0 
*int ，0xc000014110, 8192 

map[string]string , 0xc00000e038 , map[]
map[string]string , 0xc00000e038 , map[Address:ChongQIng China Name:WeiyGeek]
```

函数定义

描述: Go语言中支持函数、匿名函数和闭包，并且函数在Go语言中属于一等公民。

```go
func 函数名(参数)(返回值){
 	函数体
} 
```

其中：

> 函数名：由字母、数字、下划线组成。但函数名的第一个字母不能是数字。注意在同一个包内函数名也称不能重名（包的概念详见后文）。
> 参数：参数由参数变量和参数变量的类型组成，多个参数之间使用,分隔。
> 返回值：返回值由返回值变量和其变量类型组成，也可以只写返回值的类型，多个返回值必须用()包裹，并用,分隔。
> 函数体：实现指定功能的代码块。

```go
// 方式1
func sayHello() {
  fmt.Println("Hello World, Let's Go")
} 
```

函数调用

> 描述: 定义了函数之后，我们可以通过函数名()的方式调用函数。

```go
func main() {
  fmt.Println("Start")
  sayHello()
  fmt.Println("End")
} 
```

函数参数

描述: 通常我们需要为函数传递参数进行相应的处理以达到我们最终需要的产物。

函数参数类型

固定参数
可变参数



固定参数

常规参数类型
针对固定函数的参数我们需要制定其类型

```go
func intSum(x int, y int) {
 	fmt.Println("x + y =",x+y)
}
```

参数类型简写
函数的参数中如果相邻变量的类型相同，则可以省略类型

```go
func intSum(x , y int) {
 	fmt.Println("x + y =",x+y)
}
```

Tips : 上面的代码中，intSum函数有两个参数，这两个参数的类型均为int，因此可以省略x的类型，因为y后面有类型说明，x参数也是该类型。



可变参数

> 描述: 可变参数是指函数的参数数量不固定。Go语言中的可变参数通过在参数名后加 ... 来标识(是否似曾相识，我们在数组那章节时使用过它，表示自动判断数组中元素个数进行初始化操作)。

```go
func intSum2(x ...int) int {
  fmt.Println(x) //x是一个切片
  sum := 0
  for _, v := range x {
    sum = sum + v
   }
  return sum
}
```

```sh
调用上面的函数：

ret1 := intSum2()
ret2 := intSum2(10)
ret3 := intSum2(10, 20)
ret4 := intSum2(10, 20, 30)
fmt.Println(ret1, ret2, ret3, ret4) //0 10 30 60
```

注意：可变参数通常要作为函数的最后一个参数。



固定参数搭配可变参数使用时，可变参数要放在固定参数的后面

```go
func intSum3(x int, y ...int) int {
  fmt.Println(x, y)
  sum := x
  for _, v := range y {
  	sum = sum + v
   }
  return sum
}
```

```sh
调用上述函数：

ret5 := intSum3(100)
ret6 := intSum3(100, 10)
ret7 := intSum3(100, 10, 20)
ret8 := intSum3(100, 10, 20, 30)
fmt.Println(ret5, ret6, ret7, ret8) //100 110 130 160
```

Tips : 本质上，函数的可变参数是通过切片来实现的。

函数返回

> 描述: 与其他编程语言一样，Go语言中通过return关键字向外输出返回值。

单返回值

```go
func sum(x, y int)(res int) {
 return x + y
}
```

多返回值

描述: Go语言中函数支持多返回值，函数如果有多个返回值时必须用()将所有返回值包裹起来。

```go
func calc(x, y int) (int, int) {
  sum := x + y
  sub := x - y
  return sum, sub
}
```


函数调用并接收返回值: sum,sub := calc(5, 3) // 8 , 2

返回值命名

描述: 函数定义时可以给返回值命名，并在函数体中直接使用这些变量，最后通过return关键字返回。

```go
func calc(x, y int) (sum, sub int) {
  sum = x + y
  sub = x - y
  return
}
```


函数调用并接收返回值:

sum,sub := calc(5, 3) // 8 , 2

Tips ：如果使用返回值命令时，只要其中一个返回值命名则另外一个返回值也必须命名。

返回值补充

> 描述: 当我们的一个函数返回值类型为slice时，nil可以看做是一个有效的slice，没必要显示返回一个长度为0的切片。

```go
func someFunc(x string) []int {
  if x == "" {
  return nil // 没必要返回[]int{}
  }
  ...
}
```


函数中变量作用域

全局变量

> 描述: 全局变量是定义在函数外部的变量，它在程序整个运行周期内都有效, 在函数中可以访问到全局变量。

```go
package main
import "fmt"

//定义全局变量num
var num int64 = 10
func testGlobalVar() {
 fmt.Printf("num=%d\n", num) //函数中可以访问全局变量num
}
func main() {
 testGlobalVar()         // 10
 fmt.Printf("num=%d\n", num) // 10
}
```


局部变量

> 描述: 局部变量由分为两类一种是在函数内部定义的局部变量，另外一种则是在函数内部代码块中定义的局部变量

类1.在函数内定义的变量无法在该函数外使用。

例如: 下面的示例代码main函数中无法使用testLocalVar函数中定义的变量x

```go
func testLocalVar() {
  //定义一个函数局部变量x,仅在该函数内生效
  var x int64 = 100
  fmt.Printf("x=%d\n", x)
}

func main() {
  testLocalVar()
  fmt.Println(x) // 此时无法使用变量x,并且此时会报错undefine x。
}
```

类2.在函数内的语句块定义的变量

> 描述: 通常我们会在if条件判断、for循环、switch语句上使用这种定义变量的方式。

```go
// if 代码块
func testLocalVar2(x, y int) {
  fmt.Println(x, y) //函数的参数也是只在本函数中生效
  if x > 0 {
  z := 100 //变量z只在if语句块生效
  fmt.Println(z)
  }
  //fmt.Println(z) //此处无法使用变量z
}


// for 代码块
func testLocalVar3() {
  for i := 0; i < 10; i++ {
  fmt.Println(i) //变量i只在当前for语句块中生效
  }
  // fmt.Println(i) //此处无法使用变量i
}
```

Tips : 如果局部变量和全局变量重名，则优先访问局部变量。

Tips : 函数中查找变量的顺序步骤,(1) 现在函数内部查找，(2) 在函数上层或者外层查找, (3) 最后在全局中查找(此时如果找不到则会报错)

函数类型与变量

定义函数类型

描述: 我们可以使用type关键字来定义一个函数类型，具体格式如下：`type calculation func(int, int) int`

上面语句定义了一个calculation类型，它是一种函数类型，这种函数接收两个int类型的参数并且返回一个int类型的返回值。

简单来说，凡是满足这个条件的函数都是calculation类型的函数，例如下面的add和sub都是calculation类型的函数。

```go
type calculation func(int, int) int
func add(x, y int) int {
  return x + y
}
func sub(x, y int) int {
  return x - y
}
// add和sub都能赋值给calculation类型的变量。
var c calculation
c = add
fmt.Println(c(1,2)) // 函数类型变量传递
```


函数类型变量

> 描述: 我们可以声明函数类型的变量并且为该变量赋值：

```go
func main() {
  var c calculation               // 声明一个calculation类型的变量c
  c = add                         // 把add赋值给calculation类型的变量c
  fmt.Printf("type of c:%T\n", c) // type of c:main.calculation  (区别点)
  fmt.Println(c(1, 2))            // 像调用add一样调用c 1 + 2 = 3

  f := sub                        // 将函数sub赋值给变量f1
  fmt.Printf("type of f:%T\n", f) // type of f:func(int, int) int (区别点) 非函数类型的变量
  fmt.Println(f(30, 20))          // 像调用add一样调用f 30 - 20 = 10
}
```

7.高阶函数 (huidiaohanshu)

> 描述: 高阶函数分为函数作为参数和函数作为返回值两部分。

函数作为参数

```go
func add(x, y int) int {
  return x + y
}
func calc(x, y int, op func(int, int) int) int {
  return op(x, y)
}
func main() {
  ret2 := calc(10, 20, add)
  fmt.Println(ret2) //30
}
```

函数作为返回值

函数也可以作为返回值示例(此种方式非常值得学习):

```go
func do(s string) (func(int, int) int, error) {
  switch s {
  case "+":
  return add, nil
  case "-":
  return sub, nil
  default:
  err := errors.New("无法识别的操作符")
  return nil, err
  }
}
```


函数补充

递归函数

```go
func recursion() {
 recursion() /* 函数调用自身 */
}

func main() {
 recursion()
}
```

Tips : 递归函数对于解决数学上的问题是非常有用的，就像计算阶乘，生成斐波那契数列等。

示例1.n的阶乘计算

```go
func factorial(n uint64) (ret uint64) {
  if n <= 1 {
  return 1
  }
  return n * factorial(n-1)
 }

func demo1() {
  fmt.Println("5 的阶乘 : ", factorial(5))
}
```

示例2.利用递归求斐波那契数列

```go
// 方式1
func Fibonacci(count uint64) (ret uint64) {
  if count == 0 {
  return 0
  }
  if count == 1 || count == 2 {
  return 1
  }
  ret = Fibonacci(count-1) + Fibonacci(count-2)
  return
}

func demo3() {
  count := 10
  fmt.Printf("%v 个斐波那契数列:", count)
  for i := 1; i < count; i++ {
  fmt.Printf("%v ", Fibonacci(uint64(i)))
  }
}
```

```go
// 方式2.值得学习
// fib returns a function that returns successive Fibonacci numbers.
func fib() func() int {
  a, b := 0, 1
  return func() int {
  a, b = b, a+b
  return a
  }
}
func main() {
  f := fib()
  fmt.Println(f(), f(), f(), f(), f())
}
```

匿名函数

> 描述: 函数当然还可以作为返回值，但是在Go语言中函数内部不能再像之前那样定义函数了，只能定义匿名函数。
> 匿名函数多用于实现回调函数和闭包。

Tips : 匿名函数因为没有函数名，所以没办法像普通函数那样调用，所以匿名函数需要保存到某个变量或者作为立即执行函数:

```go
func main() {
  // 方式1.将匿名函数保存到变量
  add := func(x, y int) {
   fmt.Println(x + y)
  }
  add(10, 20) // 通过变量调用匿名函数

  //方式2.自执行函数：匿名函数定义完加()直接执行
  func(x, y int) {
    fmt.Println(x + y)
  }(10, 20)
}
```


闭包

> 描述: 闭包指的是一个函数和与其相关的引用环境组合而成的实体。简单来说，闭包=函数+外遍变量的引用, 例如在第三方包里只能传递一个不带参数的函数，此时我们可以通过闭包的方式创建一个带参数处理的流程，并返回一个不带参数的函数。

Tips : 非常注意引用的外部外部变量在其生命周期内都是存在的（即下次调用还能使用该变量值）。

闭包基础示例1：

```go
func adder() func(int) int {
  var x int  // 在f的生命周期内，变量x也一直有效
  return func(y int) int {
    x += y
    return x
  }
}
func main() {
  var f = adder()
  fmt.Println(f(10)) //x=0,y=10 ->  x = 10
  fmt.Println(f(20)) //x=10,y=20 -> x = 30
  fmt.Println(f(30)) //x=30,y=30 -> x = 60

  f1 := adder()
  fmt.Println(f1(40)) //40
  fmt.Println(f1(50)) //90
}
```


闭包基础示例2:

```go
package main

import (
"fmt"
"math"
)

// 1.假设这是个第三方包
func f1(f func()) {
fmt.Printf("# This is f1 func , Param is f func() : %T \n", f)
f() // 调用传入的函数
}

// 2.自己实现的函数
func f2(x, y int) {
fmt.Printf("# This is f2 func , Param is x,y: %v %v\n", x, y)
fmt.Printf("x ^ y = %v \n", math.Pow(float64(x), float64(y)))
}

// 要求 f1(f2) 可以执行，此时由于f1 中的传递的函数参数并无参数，所以默认调用执行一定会报错。
// 此时我们需要一个中间商利用闭包和匿名函数来实现,返回一个不带参数的函数。

func f3(f func(int, int), x, y int) func() {
tmp := func() {
f(x, y) // 此处实际为了执行f2函数
}
return tmp // 返回一个不带参数的函数，为返回给f1函数
}

func main() {
ret := f3(f2, 2, 10) // 此时函数并为执行只是将匿名函数进行返回。先执行 f3(fun,x,y int)
f1(ret)              // 当传入f1中时ret()函数便会进行执行。再执行 f1() ,最后执行 f2(x,y int)
}
```

执行结果:
```sh
# This is f1 func , Param is f func() : func()

# This is f2 func , Param is x,y: 2 10

x ^ y = 1024
```

Tips : 变量f是一个函数并且它引用了其外部作用域中的x变量，此时f就是一个闭包。并且在f的生命周期内，变量x也一直有效。



闭包进阶示例1：相比较于上面这种方式该种是将x变量放入函数参数之中，在进行函数调用时赋值。

```go
func adder2(x int) func(int) int {
  return func(y int) int {
  x += y
  return x
  }
}
func main() {
  var f = adder2(10) // `在f的生命周期内，变量x也一直有效。`
  fmt.Println(f(10)) //20
  fmt.Println(f(20)) //40
  fmt.Println(f(30)) //70

  f1 := adder2(20)
  fmt.Println(f1(40)) //60
  fmt.Println(f1(50)) //110
}
```


闭包进阶示例2：判断文件名称是否以指定的后缀结尾，是则返回原文件名称，否则返回文件名称+指定后缀的文件。

```go
func makeSuffixFunc(suffix string) func(string) string {
  return func(name string) string {
     // 判断name变量中的字符串是否已suffix结尾
  if !strings.HasSuffix(name, suffix) {
  return name + suffix
  }
  return name
  }
}

func main() {
  jpgFunc := makeSuffixFunc(".jpg")
  txtFunc := makeSuffixFunc(".txt")
  fmt.Println(jpgFunc("test")) //test.jpg
  fmt.Println(txtFunc("test")) //test.txt
}
```


闭包进阶示例3：该示例中函数同时返回add,sub两个函数.

```go
func calc(base int) (func(int) int, func(int) int) {
  add := func(i int) int {
  base += i
  return base
  }

  sub := func(i int) int {
  base -= i
  return base
  }
  return add, sub
}

func main() {
  f1, f2 := calc(10)
  fmt.Println(f1(1), f2(2)) //11 9
  fmt.Println(f1(3), f2(4)) //12 8
  fmt.Println(f1(5), f2(6)) //13 7
}
```

Important : 闭包其实并不复杂，只要牢记闭包=函数+引用环境。

defer 语句

> 描述: Go语言中的defer语句会将其后面跟随的语句进行延迟处理。在defer归属的函数即将返回时，将延迟处理的语句按defer定义的逆序进行执行(压栈-后进先出)，也就是说，先被defer的语句最后被执行，最后被defer的语句，最先被执行。

(1) defer 执行时机

> 描述: 在Go语言的函数中return语句在底层并不是原子操作，它分为给返回值赋值和RET指令两步。而defer语句执行的时机就在返回值赋值操作后，RET指令执行前。

![WeiyiGeek.defer执行时机](https://i0.hdslb.com/bfs/article/6ea5df085af746954c716a7838370ceded6dab69.png@942w_354h_progressive.png)


(2) 简单例子

示例1: defer 延迟特性演示

```go
func main() {
  fmt.Println("start")
  defer fmt.Println(1)
  defer fmt.Println(2)
  defer fmt.Println(3)
  fmt.Println("end")
}
```

```sh
// 输出结果：
start
end
3
2
1
```


示例2.探究程序执行开始时间以及最后函数返回前的时间

```go
func funcTime() int {
  fmt.Println("函数开始时间: ", time.Now().Local())
  var x = 0
  defer fmt.Println("init x = ", x) // 注意点: 此处已经将x=0值赋值了，只是没有被输出。// 最终输出
  for i := 0; i <= 100; i++ {
  x += i
  }
  defer fmt.Println("函数返回前时间: ", time.Now().Local()) // 再输出
  defer fmt.Println("ret x = ", x)                   // 后进先出 -> 先输出
  return x
}
```

```sh
// 输出结果:
函数开始时间:  2021-08-15 18:28:58.37787611 +0800 CST
ret x =  5050
函数返回前时间:  2021-08-15 18:28:58.377991344 +0800 CST
init x =  0
```

Tips : 由于defer语句延迟调用的特性，所以defer语句能非常方便的处理资源释放问题。比如：资源清理、文件关闭、解锁及记录时间等。

(3) 经典面试案例
示例1:

```go
package main
import "fmt"
// 函数返回值无命名
func f1() int {
	x := 5 // 局部变量
	defer func() {
		x++
	}()
	return x // 1.返回值 x = 5, 2.defer 语句执行后修改的是 x = 6，3.RET指令最后返回的值是 5 (由于无返回值命令则就是return已赋予的值5)
}

// 函数返回值命名 y 进行返回
func f2() (x int) {
	defer func() {
		x++
	}()
	return 5 // 1.返回值 x = 5, 2.defer 语句执行后修改的是 x = 6，3.RET指令最后返回的x值是 6 (由于存在返回值命名x则就是return x 值6)
}

// 函数返回值命名 y 进行返回
func f3() (y int) {
	x := 5 // 局部变量
	defer func() {
		x++ // 修改 x 变量的值 x + 1
	}()
	return x // 1.返回值 x = y = 5, 2.defer 语句执行后修改的是 x ，3.RET指令最后返回的y值还是 5
}

// 匿名函数无返回值
func f4() (x int) {
	defer func(x int) {
		x++ // 改变得是函数中局部变量x，非外部x变量。
	}(x)
	return 5 // 1.返回值 x = 5, 2.defer 语句执行后 x 副本 = 6 , 3.RET指令最后返回的值还是 5
}

// 匿名函数中返回值
func f5() (x int) {
	defer func(x int) int {
		x++ // 改变得是函数中局部变量x，非外部x变量。
		return x
	}(x)
	return 5 // 1.返回值 x = 5, 2.defer 语句执行后 x 副本 = 6 , 3.RET指令最后返回的值还是 5
}

// 传入一个指针到匿名函数中(方式1)
func f6() (x int) {
	defer func(x *int) {
		*x++
	}(&x)
	return 5 // 1.返回值 x = 5, 2.由于defer语句，传入x指针地址到匿名函数中 x = 6, 3.RET指令最后返回的值 6
}

// 传入一个指针到匿名函数中(方式2)
func f7() (x int) {
	defer func(x *int) int {
		(*x)++
		return *x
	}(&x)
	return 5 // 1.返回值x = 5, 2.由于defer语句，传入x指针地址到匿名函数中 x = 6, 3.RET指令最后返回值 6
}

func main() {
	fmt.Println("f1() = ", f1())
	fmt.Println("f2() = ", f2())
	fmt.Println("f3() = ", f3())
	fmt.Println("f4() = ", f4())
	fmt.Println("f5() = ", f5())
	fmt.Println("f6() = ", f6())
	fmt.Println("f7() = ", f7())
}
```

```sh
// 执行结果:
f1() =  5
f2() =  6
f3() =  5
f4() =  5
f5() =  5
f6() =  6
f7() =  6
```


示例2.问下面代码的输出结果是？（提示：defer注册要延迟执行的函数时该函数所有的参数都需要确定其值）

```go
func calc(index string, a, b int) int {
     ret := a + b
     fmt.Println(index, a, b, ret)
     return ret
}

func main() {
     x := 1
     y := 2
     defer calc("AA", x, calc("A", x, y))
     // calc("A", x, y) =>calc("A", 1, 2) = 3  {"A" , 1, 2, 3}
     // defer calc("AA", 1, 3) = 4 {"AA", 1, 3, 4}
     x = 10
     defer calc("BB", x, calc("B", x, y))
     // calc("B", x, y) = calc("B", 10, 2) = 12  {"B" , 10, 2, 12}
     // defer calc("BB", 10, 12) = 22 {"BB",10,12,22}
     y = 20
}
```

```sh
// 执行结果:
{"A" , 1, 2, 3}
{"B" , 10, 2, 12}
{"BB", 10, 12, 22}
{"AA", 1, 3, 4}
```

Tips : 当遇到defer语句时其中的函数中调用的变量值是外部变量时，是离该defer语句最近的外部变量其赋予的值(存在于一个变量多次赋值的场景)。

函数总结示例

```go
package main
import (
	"errors"
	"fmt"
	"strings"
	"time"
)
// 函数：将一段代码封装到代码块之中
// 1.无参函数
func f1() {
	fmt.Println("Hello World, Let's Go")
}

// 2.有参函数
func f2(name string) {
	fmt.Println("Hello", name)
}

// 3.函数返回值
func f3(i int, j int) int {
	sum := i + j
	return sum
}

// 4.函数多命名返回值与参数类型简写
func f4(x, y int) (sum, sub int) {
	sum = x + y
	sub = x - y
	return
}

// 5.可变参数
func f5(title string, value ...int) string {
	return fmt.Sprintf("Title : %v , Value : %v \n", title, value)
}

// 6.变量作用域之全局变量

const PATH = "/home/weiyigeek"

var author = "WeiyiGeek"

func f6() {
	fmt.Println("author:", author, ",Home PATH:", PATH)
}

// 7.变量作用域之局部变量
func f7(x, y int) {
	localAuthor := "WeiyiGeek" // 局部变量外部无法引用
	fmt.Println("localAuthor = ", localAuthor, ",x = ", x, ",y = ", y)
	// 语句块定义的变量
	if x > 0 {
		z := 1024
		fmt.Println(z)
	}
	for i := 0; i < 10; i++ {
		fmt.Print(i, " ")
	}
	// fmt.Println(z,i)  //此处无法使用变量z 和 i
	fmt.Println()
}

// 8.函数类型与变量
type calc func(int, int) int

func sum(x, y int) int {
	return x + y
}
func sub(x, y int) int {
	return x - y
}
func f8() {
	// 方式1
	var c calc
	c = sum
	fmt.Printf("type of c:%T , c(1,2) ： %v \n", c, c(1, 2)) // type of c:main.calculation  (区别点)

	// 方式2
	d := sub
	fmt.Printf("type of d:%T , d(1,2) ： %v \n", d, d(1, 2)) // type of d:func(int, int) int (区别点)

}

// 9.函数作为参数值或者作为返回值
func mul(x, y int) int {
	return x * y
}
func div(x, y int) int {
	return x / y
}

// 函数作为参数值
func calculation(x, y int, op func(int, int) int) int {
	return op(x, y)
}

// 函数作为返回值
func ops(s string) (func(int, int) int, error) {
	switch s {
	case "+":
		return sum, nil
	case "-":
		return sub, nil
	case "*":
		return mul, nil
	case "/":
		return div, nil
	default:
		err := errors.New("无法识别的操作符")
		return nil, err
	}
}

func f9() {
	// 演示1
	fmt.Printf("Type : %T , calculation (10 , 20, mul) = %v \n", calculation(10, 20, mul), calculation(10, 20, mul))

	// 演示2
	value, _ := ops("/")
	fmt.Printf("Type : %T , ops('/') ->  div(100,10) = %v \n\n", value(100, 10), value(100, 10))

}

// 10.匿名函数
func f10() {
	// 方式1
	muls := func(x, y int) int {
		fmt.Println("匿名函数1 之 x , y =", x, y)
		return x * y
	}
	ret := muls(3, 2)
	fmt.Println("匿名函数1 返回结果: ", ret)

	// 方式2
	func(x, y int) {
		fmt.Println("匿名函数2 之 x , y =", x*y)
	}(3, 2)

}

// 11.闭包
func adder1() func(int) int {
	var x int
	return func(y int) int {
		x += y
		return x
	}
}

func adder2(x int) func(int) int {
	return func(y int) int {
		x += y
		return x
	}
}

func makeSuffixFunc(suffix string) func(string) string {
	return func(name string) string {
		if !strings.HasSuffix(name, suffix) {
			return name + suffix
		}
		return name
	}
}

func f11() {
	// 方式1
	var f = adder1()
	fmt.Printf("\n闭包 adder1: %v\n", f(10)) //x=0,y=10  -> x = 10
	fmt.Println("闭包 adder1:", f(20))       //x=10,y=20 -> x = 30
	fmt.Println("闭包 adder1:", f(30))       //x=30,y=30 -> x = 60

	// 方式2
	g := adder2(10)
	fmt.Printf("闭包 adder2: %v\n", g(10)) //x=10,y=10 -> x = 20
	fmt.Println("闭包 adder2:", g(20))     //x=20,y=20 -> x = 40
	fmt.Println("闭包 adder2:", g(30))     //x=40,y=30 -> x = 70
	
	// 示例3
	testJPG := makeSuffixFunc("jpg")
	fmt.Printf("闭包 makeSuffixFunc : file test = %v , file test.jpg = %v \n\n", testJPG("test"), testJPG("test.jpg"))

}

// 12.defer 语句使用演示
func funcTime() int {
	fmt.Println("函数开始时间: ", time.Now().Local())
	var x = 0
	defer fmt.Println("init x = ", x) // 注意点: 此处已经将x=0值赋值了，只是没有被输出。 // 最终输出
	for i := 0; i <= 100; i++ {
		x += i
	}
	defer fmt.Println("函数返回前时间: ", time.Now().Local()) // 再输出
	defer fmt.Println("ret x = ", x)                   // 后进先出 -> 先输出
	return x
}

func f12() {
	ret := funcTime()
	fmt.Println("defer 示例1： 1+2+3+....+99+100 =", ret)
}

func main() {
	f1()

	f2("WeiyiGeek")
	
	fmt.Println(f3(1, 1))
	
	x, y := f4(1, 3)
	fmt.Printf("x = %d ,y = %d \n", x, y)
	
	fmt.Println(f5("我是一串数字:", 1, 2, 3, 4))
	
	f6()
	
	f7(1, 2)
	
	f8()
	
	f9()
	
	f10()
	
	f11()
	
	f12()

}
```

执行结果:

```sh
Hello World, Let''s Go

Hello WeiyiGeek

2

x = 4 ,y = -2

Title : 我是一串数字: , Value : [1 2 3 4]


author: WeiyiGeek ,Home PATH: /home/weiyigeek

localAuthor =  WeiyiGeek ,x =  1 ,y =  2

1024

0 1 2 3 4 5 6 7 8 9

type of c:main.calc , c(1,2) ： 3

type of d:func(int, int) int , d(1,2) ： -1

Type : int , calculation (10 , 20, mul) = 200

Type : int , ops('/') ->  div(100,10) = 10


匿名函数1 之 x , y = 3 2

匿名函数1 返回结果:  6

匿名函数2 之 x , y = 6


闭包 adder1: 10

闭包 adder1: 30

闭包 adder1: 60

闭包 adder2: 20

闭包 adder2: 40

闭包 adder2: 70

闭包 makeSuffixFunc : file test = testjpg , file test.jpg = test.jpg


函数开始时间:  2021-08-15 19:35:19.159014152 +0800 CST

ret x =  5050

函数返回前时间:  2021-08-15 19:35:19.159208306 +0800 CST

init x =  0

defer 示例1： 1+2+3+....+99+100 = 5050
```

Go语言基础之错误处理

> 描述: Go语言中目前(1.16 版本中)是没有异常处理机制(Tips ：说是在2.x版本中将会加入异常处理机制)，但我们可以使用error接口定义以及panic/recover函数来进行异常错误处理。

error 接口定义

描述: 在Golang中利用error类型实现了error接口，并且可以通过errors.New或者fmt.Errorf来快速创建错误实例。

主要应用场景: 在 Go 语言中，错误是可以预期的，并且不是非常严重，不会影响程序的运行。对于这类问题可以用返回错误给调用者的方法，让调用者自己决定如何处理，通常采用 error 接口进行实现。

```go
type error interface {
  Error() string
}
```


Go语言的标准库代码包errors方法：

```go
// 方式1.在errors包中的New方法（Go 1.13 版本）。
package errors
// go提供了errorString结构体，其则实现了error接口
type errorString struct {
  text string
}
func (e *errorString) Error() string {
  return e.text
}
// 在errors包中，还提供了New函数，来实例化errorString，如下：
func New(text string) error {
  return &errorString{text}
}

// 方式2.另一个可以生成error类型值的方法是调用fmt包中的Errorf函数(Go 1.13 版本以后)
package fmt
import "errors"
func Errorf(format string, args ...interface{}) error{
	return errors.New(Sprintf(format,args...))
}
```

采用 errors 包中装饰一个错误;

```go
errors.Unwrap(err error)	//通过 errors.Unwrap 函数得到被嵌套的 error。
errors.Is(err, target error)	//用来判断两个 error 是否是同一个
errors.As(err error, target interface{})	//error 断言
```


实际示例1:

```go
package main

import (
	"errors"
	"fmt"
	"math"
)

// 错误处理
// 1.Error
func demo1() {
	// 1.声明并初始化为error类型
	var errNew error = errors.New("# 错误信息来自 errors.New 方法。")
	fmt.Println(errNew)

	// 2.调用标准库中Errorf方法
	errorfFun := fmt.Errorf("- %s", "错误信息来自 fmt.Errorf 方法。")
	fmt.Println(errorfFun)
	
	// 3.实际案例
	result, err := func(a, b float64) (ret float64, err error) {
		err = nil
		if b == 0 {
			err = errors.New("此处幂指数不能为0值,其结果都为1")
			ret = 1
		} else {
			ret = math.Pow(a, b)
		}
		return
	}(5, 0)
	
	if err != nil {
		fmt.Println("# 输出错误信息:", err)
		fmt.Printf("5 ^ 0 = %v", result)
	} else {
		fmt.Printf("5 ^ 2 = %v", result)
	}

}

func main() {
	demo1()
}
```

执行结果:

```sh
# 错误信息来自 errors.New 方法。

- 错误信息来自 fmt.Errorf 方法。

# 输出错误信息: 此处幂指数不能为0值,其结果都为1

5 ^ 0 = 1
```


实际示例2:

```go
package main

import (
    "fmt"
)

// 定义一个 DivideError 结构 (值得学习)
type DivideError struct {
  dividee int
  divider int
}
// 实现 `error` 接口 (值得学习)
func (de *DivideError) Error() string {
  strFormat := `
  Cannot proceed, the divider is zero.
  dividee: %d
  divider: 0
`
  return fmt.Sprintf(strFormat, de.dividee)
}

// 定义 `int` 类型除法运算的函数
func Divide(varDividee int, varDivider int) (result int, errorMsg string) {
  if varDivider == 0 {
    dData := DivideError{
            dividee: varDividee,
            divider: varDivider,
    }
    errorMsg = dData.Error()
    return
  } else {
    return varDividee / varDivider, ""
  }
}

func main() {
  // 正常情况
  if result, errorMsg := Divide(100, 10); errorMsg == "" {
    fmt.Println("100/10 = ", result)
  }
  // 当除数为零的时候会返回错误信息
  if _, errorMsg := Divide(100, 0); errorMsg != "" {
    fmt.Println("errorMsg is: ", errorMsg)
  }
}
```

执行结果:

```sh
100/10 =  10
errorMsg is:
  Cannot proceed, the divider is zero.
  dividee: 100
  divider: 0
```


panic 函数

> 描述: 当遇到某种严重的问题时需要直接退出程序时，应该调用panic函数从而引发的panic异常, 所以panic用于不可恢复的错误类似于Java的Error。

具体流程：是当panic异常发生时，程序会中断运行，并立即执行在该goroutine，随后程序崩溃并输出日志信息。日志信息包括panic、以及value的函数调用的堆栈跟踪信息。

panic 函数语法定义:

`func panic(v interface{})`

Tips : panic函数接受任何值作为参数



示例1.数组越界会自动调用panic

```go
func TestA() {
  fmt.Println("func TestA{}")
}

func TestB(x int) {
  var a [10]int
  a[x] = 111
}

func TestC() {
  fmt.Println("func TestC()")
}

func main() {
TestA()
TestB(20) //发生异常,中断程序
TestC()
}
```

```sh
>>> func TestA{}
panic: runtime error: index out of rang
```

示例2.调用panic函数引发的panic异常

```go
func A() {
	fmt.Println("我是A函数 - 正常执行")
}

func B() {
	fmt.Println("我是B函数 - 正在执行")
	panic("func B():panic")
	fmt.Println("我是B函数 - 结束执行")
}

func C() {
	fmt.Println("我是c函数 - 正在执行")
}

func demo2() {
	A()
	B() //发生异常,中断程序
	C()
}
```

```sh
我是A函数 - 正常执行
我是B函数 - 正在执行
发生异常: panic
"func B():panic"
Stack:
	2  0x00000000004b69a5 in main.B
	    at /home/weiyigeek/app/project/go/src/weiyigeek.top/studygo/Day02/05error.go:47
	3  0x00000000004b6a8a in main.demo2
	    at /home/weiyigeek/app/project/go/src/weiyigeek.top/studygo/Day02/05error.go:57
	4  0x00000000004b6ac5 in main.main
	    at /home/weiyigeek/app/project/go/src/weiyigeek.top/studygo/Day02/05error.go:63
```

> 什么时候使用Error，什么时候使用Panic?
>
> 对于真正意外的情况，那些表示不可恢复的程序错误，例如索引越界、不可恢复的环境问题、栈溢出、数据库连接后需操作，我们才使用 panic。
> 对于其他的错误情况，我们应该是期望使用 error 来进行判定。


recover 函数

> 描述: panic异常会导致程序崩溃,而recover函数专门用于“捕获”运行时的panic异常,它可以是当前程序从运行时panic的状态中恢复并重新获得流程控制权。

通常我们会使用 Recover 捕获 Panic 异常，例如Java中利用Catch Throwable来进行捕获异常。

```go
// Java
try {
  ...
} catch (Throwable t) {
  ...
}

// C++
try {
  ...
} catch() {

}
```

panic 函数语法定义:

`func recover() interface{}`

Tips: 在未发生panic时调用recover会返回nil。

**流程说明**: 如果调用了内置函数recover,并且定义该defer语句的函数发生了panic异常,recover会使程序从panic中恢复,并返回panic value。导致panic异常的函数不会继续运行，但能正常返回。

示例1:panic与recover联合使用，此处采用 panic 演示的代码中的B函数进行继续修改
描述: 在Go语言中可以通过defer定义的函数去执行一些错误恢复的行为

```go
func recoverB() (err error) {
	fmt.Println("我是recoverB 函数 - 正在执行")
	// 必须是 defer 语句中以及在panic函数前
	defer func() {
		x := recover()
		if x != nil {
			err = fmt.Errorf("# 1.进行 recover（恢复） Panic 导致的程序异常,从此之后将会继续执行后续代码：\n%v", x)
		}
	}() // 此处利用匿名函数
	//panic("# 2.recoverB 函数中捕获 Panic")
	panic(errors.New("# 2.recoverB 函数中出现 Panic"))
	fmt.Println("我是recoverB 函数 - 结束执行") // 无法访问的代码
	return
}
func demo3() {
	A()
	err := recoverB()
	if err != nil {
		fmt.Println("#recoverB 输出的信息：", err)
	}
	C()
}
```

```sh
我是A函数 - 正常执行
我是recoverB 函数 - 正在执行

# recoverB 输出的信息： # 1.进行 recover（恢复） Panic 导致的程序异常,从此之后将会继续执行后续代码：

# 2.recoverB 函数中出现 Panic

我是c函数 - 正在执行
```


示例 2.recover捕获异常后的异常，不能再次被recover捕获。

```go
func demo4() {
	// 采用匿名函数进行立即执行该函数
	defer func() { //   声明defer，
		fmt.Println("----调用 defer func1 start----")
		err := recover() // 此处输出为 nil ，因为panic只能被 recover 捕获一次
		fmt.Printf("# 第二次 捕获 : %#v \n", err)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("----调用 defer func1 end----")
	}()

	defer func() { //   声明defer，压栈操作后进先出。
		fmt.Println("----调用 defer func2 start----")
		if err := recover(); err != nil {
			fmt.Println("# 第一次 捕获:", err) // 这里的err其实就是panic传入的内容
		}
		fmt.Println("----调用 defer func2 end----")
	}()
	
	panic("panic 异常 抛出 测试！")

}
```

```sh
----调用 defer func2 start----

# 第一次 捕获: panic 异常 抛出 测试！

----调用 defer func2 end----
----调用 defer func1 start----

# 第二次 捕获 : <nil>

----调用 defer func1 end----
```

> Q: panic() 与 recover() 位置区别?
> 答: panic函数可以在任何地方引发(但panic退出前会执行defer指定的内容)，但recover函数只有在defer调用的函数中有效并且一定要位于panic语句之前。

TIPS : 非常注意下面这种“错误方式”, 他可能会形成僵尸服务进程，导致Health Check失效。

```go
defer func() {
  if err := recover(); err != nil {
    Log.Error("Recovered Panic", err)
  }
}()
```

Q: panic 和 os.Exit 联用时对recover的影响

os.Exit 退出时不会调用defer指定的函数.
os.Exit 退出时不会输出当前调用栈信息.



错误处理最佳实践

1、预定义错误，code里判断
2、及早失败，避免嵌套

Go语言基础之结构体

> 描述: Go语言中没有“类”的概念，也不支持“类”的继承等面向对象的概念。但 Go语言中通过结构体的内嵌再配合接口比面向对象具有更高的扩展性和灵活性。

* Go语言中的基础数据类型可以表示一些事物的基本属性，但是当我们想表达一个事物的全部或部分属性时，这时候再用单一的基本数据类型明显就无法满足需求了（局限性）。
* Go语言提供了一种自定义数据类型，可以封装多个基本数据类型，这种数据类型叫结构体(英文名称struct), 我们可以通过struct来定义自己的类型了。

简单得说: 结构体时一种数据类型，一种我们自己可以保持多个维度数据的类型。 所以与其他高级编程语言一样，Go语言也可以采用结构体的特性, 并且Go语言通过struct来实现面向对象。

类型定义

> 描述: 在Go语言中有一些基本的数据类型，如string、int{}整型、float{}浮点型、boolean布尔等数据类型， Go语言中可以使用type关键字来定义自定义类型(实际上定义了一个全新的类型)。

Tips : 我们可以基于内置的基本类型定义，也可以通过struct定义。

```go
//将MyInt定义为int类型
type MyInt int
```

通过type关键字的定义，MyInt就是一种新的类型，它具有int的特性。

类型别名

> 描述: 类型别名从字面意义上都很好理解，即类型别名本章上与原类型一样, 就比如像一个孩子小时候有小名、乳名，上学后用学名，英语老师又会给他起英文名，但这些名字都指的是他本人。

示例演示:

```go
// TypeAlias只是Type的别名，本质上TypeAlias与Type是同一个类型
type TypeAlias = Type
```

我们之前见过的rune和byte就是类型别名，他们的定义如下：

```go
type byte = uint8
type rune = int32
```

Tips: 采用int32别名创建一个变量的几种方式。

```go
type MyInt32 = int32
// 方式1
var i MyInt32
i = 1024
// 方式2
var j MyInt32 = 1024
// 方式3
var k  = MyInt32(1024)
// 方式4
l := MyInt32(1024)  // 此处并非是函数，而是一个强制类型转换而已
```

> Q: 类型定义和类型别名有何区别?
>
> 答: 类型别名与类型定义表面上看只有一个等号的差异，我们通过下面的这段代码来理解它们之间的区别。

示例演示1:

```go
//1.类型定义
type NewInt int

//2.类型别名
type MyInt = int

// 类型定义 与 类型别名 区别演示
func demo1() {
	// 类型定义的使用
	var i NewInt
	i = 1024
	fmt.Printf("Type of i: %T, Value:%v \n", i, i)

	// 类型别名的使用
	var j MyInt
	j = 2048
	fmt.Printf("Type of j: %T, Value:%v \n", j, j)
	
	// rune 也是类型别名底层还是int32类型
	var k rune
	k = '中'
	fmt.Printf("Type of j: %T, Value:%c \n", k, k)

}
```

```sh
Type of i: main.NewInt, Value:1024
Type of j: int, Value:2048
Type of j: int32, Value:中
```

结果显示说明:

i 变量的类型是main.NewInt，表示main包下定义的NewInt类型。
j 变量的类型是int，因MyInt类型只会在代码中存在，编译完成时并不会有MyInt类型。

结构体的定义

> 描述: 语言内置的基础数据类型是用来描述一个值的，而结构体是用来描述一组值的。比如一个人有名字、年龄和居住城市等，本质上是一种聚合型的数据类型。

使用type和struct关键字来定义结构体，具体代码格式如下：

```go
type 类型名 struct {
  字段名 字段类型
  字段名 字段类型
  …
}
```

其中:

类型名：标识自定义结构体的名称，在同一个包内不能重复。
字段名：表示结构体字段名。结构体中的字段名必须唯一。
字段类型：表示结构体字段的具体类型。
举例说明: 以定义一个Person（人）结构体为例:

```go
// 方式(0)
var v struct{}

// 方式(1)
type person struct {
	name string
	city string
	age  int8
}

// 方式(2): 同样类型的字段也可以写在一行
type person1 struct {
	name, city string
	age   int8
}
```

Tips : 上面创建了结构体一个person的自定义类型，它有name、city、age三个字段，分别表示姓名、城市和年龄。这样我们使用这个person结构体就能够很方便的在程序中表示和存储人信息了。

结构体实例化

> 描述: 只有当结构体实例化时，才会真正地分配内存。也就是必须实例化后才能使用结构体的字段。

Tips ：结构体本身也是一种类型，我们可以像声明内置类型一样使用var关键字声明结构体类型。例如:var 结构体实例 结构体类型。

> 描述: 结构体初始化是非常必要，因为没有初始化的结构体，其成员变量都是对应其类型的零值。

结构体示例化的三种语法格式:

```go
type demo struct {
  username string
  city string
}

// 1.方式1.利用`.`进行调用指定属性
var m1 demo
demo.username = "WeiyiGeek"

// 2.方式2.使用键值对初始化
m2 := demo {username: "WeiyiGeek",city:"重庆",}
m2 := &demo {username: "WeiyiGeek",city:"重庆",} // ==> new(demo) 此种方式会在结构体指针里面实践。

// 3.方式3.使用值的列表初始化
m3 := demo {
  "WeiyiGeek",
  "重庆"
}
m3 := &demo {
  "WeiyiGeek",
  "重庆"
}
```

Tips : 特别注意在使用值的列表初始化这种格式初始化时, (1)必须初始化结构体的所有字段,(2)初始值的填充顺序必须与字段在结构体中的声明顺序一致,(3) 该方式不能和键值初始化方式混用。



示例演示: 下述演示三种基础方式进行结构体的实例化。

```go
// 1.结构体初识还是老示例采用结构体描述人员信息并进行赋值使用
type Person struct {
	name  string
	age   uint8
	sex   bool
	hobby []string
}

func demo1() {
	// 方式1.声明一个Persin类型的变量x
	var x Person
	// 通过结构体中的属性进行赋值
	x.name = "WeiyiGeek"
	x.age = 20
	x.sex = true // {Boy,Girl)
	x.hobby = []string{"Basketball", "乒乓球", "羽毛球"}
	// 输出变量x的类型以及其字段的值
	fmt.Printf("Type of x : %T, Value : %v \n", x, x)
	x.name = "WeiyiGeeker"
  // 我们通过.来访问结构体的字段（成员变量）, 例如x.name和x.age等。
	fmt.Printf("My Name is %v \n", x.name)

	// 方式2.在声明是进行赋值(key：value，或者 value)的值格式

  // 使用键值对初始化
	var y = Person{
		name:  "Go",
		age:   16,
		sex:   false,
		hobby: []string{"Computer", "ProgramDevelopment"},
	}
	fmt.Printf("Type of y : %T, Value : %v \n", y, y)
	// 非常注意此种方式是按照结构体中属性顺序进行赋值,同样未赋值的为该类型的零值
  // 使用值的列表初始化
	z := Person{
		"WeiyiGeek",
		10,
		true,
		[]string{},
	}
	fmt.Printf("Type of z : %T, Value : %v \n", z, z)
}
```

```sh
Type of x : main.Person, Value : {WeiyiGeek 20 true [Basketball 乒乓球 羽毛球]}
My Name is WeiyiGeeker
Type of y : main.Person, Value : {Go 16 false [Computer ProgramDevelopment]}
Type of z : main.Person, Value : {WeiyiGeek 10 true []}
```

Tips : 如果没有给结构体中的属性赋值，则默认采用该类型的零值。



结构体内存布局

描述: 结构体占用一块连续的内存，但是需要注意空结构体是不占用空间的。

连续内存空间

示例演示:

```go
// 示例1.空结构体是不占用空间的
var v struct{}
fmt.Println(unsafe.Sizeof(v))  // 0


// 示例2.结构体占用一块连续的内存
type test struct {
	a int8
	b int8
	c int8
	d int8
}
n := test{
	1, 2, 3, 4,
}
fmt.Printf("n.a %p, int8 size: %d\n", &n.a, unsafe.Sizeof(bool(true)))
fmt.Printf("n.b %p\n", &n.b)
fmt.Printf("n.c %p\n", &n.c)
fmt.Printf("n.d %p\n", &n.d)

// 执行结果:
n.a 0xc0000a0060
n.b 0xc0000a0061
n.c 0xc0000a0062
n.d 0xc0000a0063
```


内存对齐分析

[进阶知识点] 关于在 Go 语言中恰到好处的内存对齐
描述: 在讲解前内存对齐前, 我们先丢出两个struct结构体引发思考:

示例1. 注意两个结构体中声明不同元素类型的顺序。

```go
type Part1 struct {
	a bool
	b int32
	c int8
	d int64
	e byte
}

type Part2 struct {
	e byte
	c int8
	a bool
	b int32
	d int64
}
```


在开始之前，希望你计算一下 Part1 与 Part2 两个结构体分别占用的大小是多少呢?

```go
func typeSize() {
	fmt.Printf("bool size: %d\n", unsafe.Sizeof(bool(true)))
	fmt.Printf("int32 size: %d\n", unsafe.Sizeof(int32(0)))
	fmt.Printf("int8 size: %d\n", unsafe.Sizeof(int8(0)))
	fmt.Printf("int64 size: %d\n", unsafe.Sizeof(int64(0)))
	fmt.Printf("byte size: %d\n", unsafe.Sizeof(byte(0)))
	fmt.Printf("string size: %d\n", unsafe.Sizeof("WeiyiGeek"))  // 注意上面声明的结构体中没有该类型。
}

// 输出结果
bool size: 1
int32 size: 4
int8 size: 1
int64 size: 8
byte size: 1
string size: 16
```


这么一算 Part1/Part2 结构体的占用内存大小为 1+4+1+8+1 = 15 个字节。相信有的小伙伴是这么算的，看上去也没什么毛病

真实情况是怎么样的呢？我们实际调用看看，如下：

```go
func main() {
	part1 := Part1{}
	fmt.Printf("part1 size: %d, align: %d\n", unsafe.Sizeof(part1), unsafe.Alignof(part1))
	fmt.Println()
	part2 := Part2{}
	fmt.Printf("part2 size: %d, align: %d\n", unsafe.Sizeof(part2), unsafe.Alignof(part2))
}
```

执行结果:

```sh
part1 size: 32, align: 8
part2 size: 16, align: 8

Tips : `unsafe.Sizeof` 来返回相应类型的空间占用大小
Tips : `unsafe.Alignof` 来返回相应类型的对齐系数
```


从上述结果中可以看见 part1 占用32个字节而 part2 占用16字节,此时 part1 比我们上面计算结构体占用字节数多了16 Byte, 并且相同的元素类型但顺序不同的 part2 是正确的只占用了 16 Byte, 那为什么会出现这样的情况呢？同时这充分地说明了先前的计算方式是错误的。

在这里要提到 “内存对齐” 这一概念，才能够用正确的姿势去计算，接下来我们详细的讲讲它是什么

> Q: What 什么是内存对齐?
> 答:有的小伙伴可能会认为内存读取，就是一个简单的字节数组摆放(例图1) 表示一个坑一个萝卜的内存读取方式。但实际上 CPU 并不会以一个一个字节去读取和写入内存, 相反 CPU 读取内存是一块一块读取的，块的大小可以为 2、4、6、8、16 字节等大小, 块大小我们称其为内存访问粒度(例图2)：
>
> ![WeiyiGeek.内存对齐](https://i0.hdslb.com/bfs/article/cbdb6330853dcdebe52b45d73111247a1a7aca30.png@894w_182h_progressive.png)


在样例中，假设访问粒度为 4。 CPU 是以每 4 个字节大小的访问粒度去读取和写入内存的。这才是正确的姿势

> Q: Why 为什么要关心对齐?
>
> 你正在编写的代码在性能（CPU、Memory）方面有一定的要求
> 你正在处理向量方面的指令
> 某些硬件平台（ARM）体系不支持未对齐的内存访问

> Q: Why 为什么要做对齐?
>
> 平台（移植性）原因：不是所有的硬件平台都能够访问任意地址上的任意数据。例如：特定的硬件平台只允许在特定地址获取特定类型的数据，否则会导致异常情况
> 性能原因：若访问未对齐的内存，将会导致 CPU 进行两次内存访问，并且要花费额外的时钟周期来处理对齐及运算。而本身就对齐的内存仅需要一次访问就可以完成读取动作
>
> ![WeiyiGeek.内存申请](https://i0.hdslb.com/bfs/article/7a1fa5379fa2bf8421544374b46a5bd39cfb6657.png@573w_330h_progressive.png)

在上图中，假设从 Index 1 开始读取，将会出现很崩溃的问题, 因为它的内存访问边界是不对齐的。因此 CPU 会做一些额外的处理工作。如下：

1.CPU 首次读取未对齐地址的第一个内存块，读取 0-3 字节。并移除不需要的字节 0
2.CPU 再次读取未对齐地址的第二个内存块，读取 4-7 字节。并移除不需要的字节 5、6、7 字节
3.合并 1-4 字节的数据
4.合并后放入寄存器
从上述流程可得出，不做 “内存对齐” 是一件有点 "麻烦" 的事。因为它会增加许多耗费时间的动作, 而假设做了内存对齐，从 Index 0 开始读取 4 个字节，只需要读取一次，也不需要额外的运算。这显然高效很多，是标准的空间换时间做法



默认系数
描述: 在不同平台上的编译器都有自己默认的 “对齐系数”，可通过预编译命令 #pragma pack(n) 进行变更，n 就是代指 “对齐系数”。一般来讲，我们常用的平台的系数如下：32 位：4, 64 位：8, 例如, 前面示例中的对齐系数是8验证了我们系统是64位的。

另外要注意不同硬件平台占用的大小和对齐值都可能是不一样的。因此本文的值不是唯一的，调试的时候需按本机的实际情况考虑

不同数据类型的对齐系数

```go
func main() {
  fmt.Printf("bool align: %d\n", unsafe.Alignof(bool(true)))
  fmt.Printf("byte align: %d\n", unsafe.Alignof(byte(0)))
  fmt.Printf("int8 align: %d\n", unsafe.Alignof(int8(0)))
  fmt.Printf("int32 align: %d\n", unsafe.Alignof(int32(0)))
  fmt.Printf("int64 align: %d\n", unsafe.Alignof(int64(0)))
  fmt.Printf("string align: %d\n", unsafe.Alignof("WeiyiGeek"))
  fmt.Printf("map align: %d\n", unsafe.Alignof(map[string]string{}))
}
```

执行结果:

```sh
bool align: 1
byte align: 1
int8 align: 1
int32 align: 4
int64 align: 8
string align: 8
map align: 8
```

通过观察输出结果，可得知基本都是 2^n，最大也不会超过 8。这是因为我手提（64 位）编译器默认对齐系数是 8，因此最大值不会超过这个数。

Tips: 在上小节中提到了结构体中的成员变量要做字节对齐。那么想当然身为最终结果的结构体，也是需要做字节对齐的



对齐规则

1.结构体的成员变量，第一个成员变量的偏移量为 0。往后的每个成员变量的对齐值必须为编译器默认对齐长度（#pragma pack(n)）或当前成员变量类型的长度（unsafe.Sizeof），取最小值作为当前类型的对齐值。其偏移量必须为对齐值的整数倍
2.结构体本身，对齐值必须为编译器默认对齐长度（#pragma pack(n)）或结构体的所有成员变量类型中的最大长度，取最大数的最小整数倍作为对齐值
3.结合以上两点，可得知若编译器默认对齐长度（#pragma pack(n)）超过结构体内成员变量的类型最大长度时，默认对齐长度是没有任何意义的


分析流程

**Step 1.首先我们先来分析 part1 结构体 到底经历了些什么，影响了 “预期” 结果**

![img](https://i0.hdslb.com/bfs/article/861892fefa98d4bb566d1776b724995c14f7d0f9.png@669w_737h_progressive.png)


成员对齐步骤

第一个成员 a
	类型为 bool
	大小/对齐值为 1 字节
	初始地址，偏移量为 0。占用了第 1 位
第二个成员 b
	类型为 int32
	大小/对齐值为 4 字节
	根据规则 1，其偏移量必须为 4 的整数倍。确定偏移量为 4，因此 2-4 位为 Padding(理解点)。而当前数值从第 5 位开始填充，到第 8 位。如下：axxx|bbbb
第三个成员 c
	类型为 int8
	大小/对齐值为 1 字节
	根据规则1，其偏移量必须为 1 的整数倍。当前偏移量为 8。不需要额外对齐，填充 1 个字节到第 9 位。如下：	axxx|bbbb|c...
第四个成员 d
	类型为 int64
	大小/对齐值为 8 字节
	根据规则 1，其偏移量必须为 8 的整数倍。确定偏移量为 16，因此 9-16 位为 Padding。而当前数值从第 17 位开始写入，到第 24 位。如下：axxx|bbbb|cxxx|xxxx|dddd|dddd
第五个成员 e
	类型为 byte
	大小/对齐值为 1 字节
	根据规则 1，其偏移量必须为 1 的整数倍。当前偏移量为 24。不需要额外对齐，填充 1 个字节到第 25 位。如下：axxx|bbbb|cxxx|xxxx|dddd|dddd|e...


整体对齐步骤

​	在每个成员变量进行对齐后，根据规则 2，整个结构体本身也要进行字节对齐，因为可发现它可能并不是 2^n，不是偶数倍。显然不符合对齐的规则
​	根据规则 2，可得出对齐值为 8。现在的偏移量为 25，不是 8 的整倍数。因此确定偏移量为 32。对结构体进行对齐


结果说明：

最终 Part1 内存布局 axxx|bbbb|cxxx|xxxx|dddd|dddd|exxx|xxxx

通过本节的分析，可得知先前的 “推算” 为什么错误？
是因为实际内存管理并非 “一个萝卜一个坑” 的思想。而是一块一块。通过空间换时间（效率）的思想来完成这块读取、写入。另外也需要兼顾不同平台的内存操作情况



**Step 2.通过上述我们可知根据成员变量的类型不同，其结构体的内存会产生对齐等动作。而像 part2 结构体一样，按照变量类型对齐值从小到大，进行依次排序进行占用内存空间的结果分析。**

通过开头的示例我们可知，只是 “简单” 对成员变量的字段顺序(类型占用字节数从小到大排序)进行改变，就改变了结构体占用大小。

![img](https://i0.hdslb.com/bfs/article/4626ccc778fa93883c3802e0849befd5360f1362.png@671w_615h_progressive.png)


成员对齐

第一个成员 e
	类型为 byte
	大小/对齐值为 1 字节
	初始地址，偏移量为 0。占用了第 1 位
第二个成员 c
	类型为 int8
	大小/对齐值为 1 字节
	根据规则1，其偏移量必须为 1 的整数倍。当前偏移量为 2。不需要额外对齐
第三个成员 a
	类型为 bool
	大小/对齐值为 1 字节
	根据规则1，其偏移量必须为 1 的整数倍。当前偏移量为 3。不需要额外对齐
第四个成员 b
	类型为 int32
	大小/对齐值为 4 字节
	根据规则1，其偏移量必须为 4 的整数倍。确定偏移量为 4，因此第 3 位为 Padding(理解点)。而当前数值从第 4 位开始填充，到第 8 位。如下：ecax|bbbb
第五个成员 d
	类型为 int64
	大小/对齐值为 8 字节
	根据规则1，其偏移量必须为 8 的整数倍。当前偏移量为 8。不需要额外对齐，从 9-16 位填充 8 个字节。如下：ecax|bbbb|dddd|dddd
整体对齐: 由于符合规则 2，则不需要额外对齐。

结果说明:

Part2 内存布局：ecax|bbbb|dddd|dddd



总结

通过对比 Part1 和 Part2 的内存布局，你会发现两者有很大的不同。如下：

Part1：axxx|bbbb|cxxx|xxxx|dddd|dddd|exxx|xxxx
Part2：ecax|bbbb|dddd|dddd
仔细一看，Part1 存在许多 Padding。显然它占据了不少空间，那么 Padding 是怎么出现的呢？

通过本文的介绍，可得知是由于不同类型导致需要进行字节对齐，以此保证内存的访问边界

那么也不难理解，为什么调整结构体内成员变量的字段顺序就能达到缩小结构体占用大小的疑问了，是因为巧妙地减少了 Padding 的存在。让它们更 “紧凑” 了。这一点对于加深 Go 的内存布局印象和大对象的优化非常有帮

当然了，没什么特殊问题，你可以不关注这一块。但你要知道这块知识点 😄



指针类型结构体

结构体指针实例化

> 描述: 我们还可以通过使用new关键字(对基础类型进行实例化)对结构体进行实例化，得到的是结构体的地址。

创建一个结构体指针格式:

```go
// 方式1.New 实例化
var p2 = new(person)
fmt.Printf("%T\n", p2)     // *main.person
fmt.Printf("p2=%#v\n", p2) // p2=&main.person{name:"", city:"", age:0}
// 在Go语言中支持对结构体指针直接使用.来访问结构体的成员。
p2.name = "WeiyiGeek"
p2.age = 22
p2.city = "重庆"
fmt.Printf("p2=%#v\n", p2)  //显示出其结构体结构: p2=&main.person{name:"WeiyiGeek", city:"重庆", age:22}

// 方式2.使用&对结构体进行取地址操作相当于对该结构体类型进行了一次new实例化操作。
p3 := &person{}
fmt.Printf("%T\n", p3)     //*main.person
fmt.Printf("p3=%#v\n", p3) //p3=&main.person{name:"", city:"", age:0}
p3.name = "WeiyiGeek"
p3.age = 30
p3.city = "重庆"
fmt.Printf("p3=%#v\n", p3) //p3=&main.person{name:"WeiyiGeek", city:"重庆", age:30}
```

Tips ：p3.name = "WeiyiGeek"其实在底层是(*p3).name = "Geeker"，这是Go语言帮我们实现的语法糖。

示例演示:

```go
type Person struct {
	name  string
	age   uint8
	sex   bool
	hobby []string
}

// 3.结构体指针
func demo3() {
	// 方式1.结构体利用new实例化在内存中申请一块空间
	var p1 = new(Person)
	(*p1).name = "WeiyiGeek" // 取得地址存放的值并将其进行覆盖
	p1.age = 20              // Go语言的语法糖自动根据指针找到对应地址的值并将其值覆盖。
	fmt.Printf("Type of p1 : %T, Struct 实例化结果: %#v\n", p1, p1)

	// 方式2.采用取地址&符号进行实例化结构体(效果与new差不多)
	p2 := &Person{}
	(*p2).name = "Golang" // 取得地址存放的值并将其进行覆盖
	p2.age = 12           // Go语言的语法糖自动根据指针找到对应地址的值并将其值覆盖。
	p2.sex = true
	fmt.Printf("Type of p2 : %T, Struct 实例化结果: %#v\n", p2, p2)
	
	// 5.使用键值对初始化(也可以对结构体指针进行键值对初始化)
	// 当某些字段没有初始值的时候，该字段可以不写。此时没有指定初始值的字段的值就是该字段类型的零值。
	p3 := &Person{
		name: "北京",
	}
	fmt.Printf("p3 Value = %#v \n", p3)
	
	// 6.使用值的列表初始化
	// 初始化结构体的时候可以简写，也就是初始化的时候不写键，直接写值：
	p4 := &Person{
		"WeiyiGeek",
		20,
		false,
		[]string{},
	}
	fmt.Printf("p4 Value = %#v \n", p4)
	
	// 4.探究Struct结构体开辟的是连续的内存空间(内存对齐效果)
	fmt.Printf("*p2 size of = %d, p2 align of = %d \n", unsafe.Sizeof(*p2), unsafe.Alignof(p2))
	fmt.Printf("Pointer p2 = %p, \name = %p,p2.name size of = %d \nnage = %p, p2.age size of = %d\nsex = %p, p2.sex size of = %d\nhobby = %p,p2.hobby size of = %d \n", p2, &p2.name, unsafe.Sizeof((*p2).name), &p2.age, unsafe.Sizeof(p2.age), &p2.sex, unsafe.Sizeof(p2.sex), &p2.hobby, unsafe.Sizeof(p2.hobby))

}
```

执行结构:

```sh
Type of p1 : *main.Person, Struct 实例化结果: &main.Person{name:"WeiyiGeek", age:0x14, sex:false, hobby:[]string(nil)}
Type of p2 : *main.Person, Struct 实例化结果: &main.Person{name:"Golang", age:0xc, sex:true, hobby:[]string(nil)}
p3 Value = &main.Person{name:"北京", age:0x0, sex:false, hobby:[]string(nil)}
p4 Value = &main.Person{name:"WeiyiGeek", age:0x14, sex:false, hobby:[]string{}}
// 结构体占用一块连续的内存地址。
*p2 size of = 48, p2 align of = 8
Pointer p2 = 0xc0001181b0,
name = 0xc0001181b0,p2.name size of = 16
age = 0xc0001181c0, p2.age size of = 1
sex = 0xc0001181c1, p2.sex size of = 1
hobby = 0xc0001181c8,p2.hobby size of = 24
```

从上述Person 结构体指针 p2 内存对齐结果中可知，元素类型占用的大小 16 + 1 + 1 + 24 = 42 Byte, 但是收到整体对齐的规则约束，该 p2 指针类型的结构体占用的内存空间大小为 48 Byte。



结构体指针函数传递

描述: 我们可以将指针类型的结构体进行地址传递在函数中修改其元素属性内容。

```go
func personChange(p Person) {
	p.name = "Change"   // 拷贝的是 p4 指针类型的结构的副本(值引用)
}

func personPointerChange(p *Person) {
	p.name = "PointerChange"  // 传递的是 p4 的地址，所以修改的是 p4.name 的属性值
}

func demo4() {
	p4 := &Person{
		name: "WeiyiGeek",
	}
	personChange(*p4)  // 值传递
	fmt.Printf("personChange(*p4) ->	name = %v \n", p4.name)

	personPointerChange(p4) // 地址传递
	fmt.Printf("personPointerChange(*p4) ->	name = %v", p4.name)

}
```

```sh
personChange(*p4) ->	name = WeiyiGeek
personPointerChange(*p4) ->	name = PointerChange
```

Tips : Go 语言中函数传的参数永远传的是拷贝, 如果要修改原数据必须进行取地址传递并修改。



结构体指针构造函数

> 描述: Go语言的结构体没有构造函数，但我们可以自己实现一个。

Tips: Go语言构造函数约定俗成用new进行开头，例如 newDog()。

例如: 下方的代码就实现了一个person的构造函数。

```go
// (1) 结构体构造函数
type Person struct {
	name, city string
	age        uint8
}

// 方式1.值传递(拷贝副本) 返回的是结构体
func newPerson(name, city string, age uint8) Person {
	return Person{
		name: name,
		city: city,
		age:  age,
	}
}

// 方式2.地址(指针类型变量)传递返回的是结构体指针
func newPointerPerson(name, city string, age uint8) *Person {
	return &Person{
		name: name,
		city: city,
		age:  age,
	}
}

func demo1() {
	// (1) 通过定义的函数直接进行结构体的初始化(值拷贝的方式)
	var person = newPerson("WeiyiGeek", "重庆", 20)
	fmt.Printf("newPerson Type : %T, Value : %v\n", person, person)
	// (2) 通过定义的函数直接传入指针类型的结构体进行初始化(地址拷贝的方式)
	var pointerperson = newPointerPerson("Go", "world", 12)
	fmt.Printf("newPointerPerson Type : %T, Value : %v\n", pointerperson, pointerperson)
}
```

执行结果:

```sh
newPerson Type : main.Person, Value : {WeiyiGeek 重庆 20}
newPointerPerson Type : *main.Person, Value : &{Go world 12}
```

Tips ：因为struct是值类型，如果结构体比较复杂的话，值拷贝性能开销会比较大，所以该构造函数返回的是结构体指针类型。



结构体方法与接收者

描述: Go语言中的方法（Method）是一种作用于特定类型变量的函数, 这种特定类型变量叫做接收者（Receiver）, 接收者的概念就类似于其他语言中的 this 或者 self。

结构体方法

定义格式：

```sh
func (接收者变量 接收者类型) 方法名(参数列表) (返回参数) {
  函数体
}
```

其中，

接收者变量：接收者中的参数变量名在命名时，官方建议使用接收者类型名称首字母的小写，而不是self、this之类的命名。例如 Person类型 的接收者变量应该命名为 p，Connector类型的接收者变量应该命名为c等。
接收者类型：接收者类型和参数类似，可以是指针类型和非指针类型。
方法名、参数列表、返回参数：具体格式与函数定义相同。
Tips : 结构体方法名称写法约束规定，如果其标识符首字母是大写的就表示对外部包可见(例如 java 中 public 指定的函数或者是类公共的)，如果其标识符首字母是小写的表示对外部不可见(不能直接调用), 当然这是一种开发习惯非强制必须的。



示例演示:

```go
//Person 结构体
type Person struct {
	name string
	age  int8
}

//NewPerson 构造函数
func NewPerson(name string, age int8) *Person {
	return &Person{
		name: name,
		age:  age,
	}
}

//Dream Person做梦的方法
func (p Person) Dream() {
	fmt.Printf("%s的梦想是学好Go语言！\n", p.name)
}

func main() {
	p1 := NewPerson("WeiyiGeek", 25)
	p1.Dream()  // WeiyiGeek的梦想是学好Go语言！
}
```

Tips : 方法与函数的区别是，函数不属于任何类型，方法属于特定的类型。



值类型的接收者

> 描述: 当方法作用于值类型接收者时，Go语言会在代码运行时将接收者的值复制一份。

在值类型接收者的方法中可以获取接收者的成员值，但修改操作只是针对副本，无法修改接收者变量本身。

例如: 我们为 Person 添加一个SetAge方法，来修改实例变量的年龄, 验证是否可被修改。

```go
//  使用值接收者：SetAge2 设置p的年龄
func (p Person) SetAge2(newAge int8) {
	p.age = newAge
}
func main() {
	p1 := NewPerson("WeiyiGeek", 25)
	p1.Dream()
	fmt.Println(p1.age) // 25
	p1.SetAge2(30) // (*p1).SetAge2(30)
	fmt.Println(p1.age) // 25
}
```


指针类型的接收者

> 描述: 指针类型的接收者由一个结构体的指针组成，由于指针的特性，调用方法时修改接收者指针的任意成员变量，在方法结束后，修改都是有效的。此种方式就十分接近于其他语言中面向对象中的this或者self达到的效果。

例如：我们为 Person 添加一个SetAge方法，来修改实例变量的年龄。

```sh
// 使用指针接收者 : SetAge 设置p的年龄: 传入的 Person 实例化后的变量的地址 p ，并通过p.属性进行更改其内容存储的内容。
func (p *Person) SetAge(newAge int8) {
	p.age = newAge
}
//调用
func main() {
	p1 := NewPerson("WeiyiGeek", 25)
	fmt.Println(p1.age) // 25
	p1.SetAge(30)
	fmt.Println(p1.age) // 30
}
```

> Q: 什么时候应该使用指针类型接收者?
>
> 一是、需要修改接收者中的值。
> 二是、接收者是拷贝代价比较大的大对象。
> 三是、保证一致性，如果有某个方法使用了指针接收者，那么其他的方法也应该使用指针接收者。

案例演示:

```go
// 结构体方法和接收者, 只能被Person结构体实例化的对象进行调用，不能像函数那样直接调用。此处还是采用上面声明的结构体
func (p Person) ChangePersonName(name string) {
	p.name = name
	fmt.Printf("# 执行 -> ChangePersonName 方法 -> p Ptr : %p ,value : %v\n", &p, p.name)
}
func (p *Person) ChangePointerPersonName(name string, age uint8) {
	p.name = name
	p.age = age
	fmt.Printf("# 执行 -> ChangePointerPersonName 方法 -> p Ptr : %p (关键点),value : %v\n", p, p.name)
}
func demo2() {
	// 利用构造函数进行初始化
	p1 := newPerson("小黄", "Beijing", 20)
	fmt.Printf("p1 Pointer : %p , Struct : %+v \n", &p1, p1)
	// 调用 ChangePersonName 方法
	p1.ChangePersonName("小黑") // 值类型的接收者(修改的是p1结构体副本的值)
	fmt.Printf("	p1 Pointer : %p , Struct : %+v \n", &p1, p1)
	// 调用 ChangePointerPersonName 方法
	p1.ChangePointerPersonName("小白", 30) //指针类型的接收者 (修改的是p1结构体元素的值)
	fmt.Printf("	p1 Pointer : %p , Struct : %+v \n", &p1, p1)
}
```

执行结果:

```sh
p1 Pointer : 0xc00010c150 , Struct : {name:小黄 city:Beijing age:20}

# 执行 -> ChangePersonName 方法 -> p Ptr : 0xc00010c1b0 ,value : 小黑

  p1 Pointer : 0xc00010c150 , Struct : {name:小黄 city:Beijing age:20}

# 执行 -> ChangePointerPersonName 方法 -> p Ptr : 0xc00010c150 (关键点),value : 小白

  p1 Pointer : 0xc00010c150 , Struct : {name:小白 city:Beijing age:30}
```


任意类型的接收者

> 描述: 在Go语言中接收者的类型可以是任何类型，不仅仅是结构体，任何类型都可以拥有方法。

举个例子，我们基于内置的int类型使用type关键字可以定义新的自定义类型，然后为我们的自定义类型添加方法。

```go
// 3.任意类型的接收者都可以拥有自己的方法
// MyInt 将int定义为自定义MyInt类型
type MyInt int
// SayHello 为MyInt添加一个SayHello的方法
func (m MyInt) SayHello(s string) {
	fmt.Printf("Hello, 我是一个int, %s", s)
}
// ChangeM 为MyInt添加一个ChangeM的方法
func (m *MyInt) ChangeM(newm MyInt) {
	fmt.Printf("# Start old m : %d -> new m : %d \n", *m, newm)
	*m = newm  // 关键点修改m其值，此处非拷贝的副本
	fmt.Printf("# End old m : %d -> new m : %d \n", *m, newm)
}
func demo3() {
	// 声明
	var m1 MyInt
	// 赋值
	m1 = 100
  // 方式2
  m2 := MyInt(255)
	// 调用类型方法
	m1.SayHello("Let'Go")
	fmt.Printf("SayHello -> Type m1 : %T, value : %+v \n", m1, m1)
	// 调用类型方法修改m1其值
	m1.ChangeM(1024)
	fmt.Printf("ChangeM -> Type m1 : %T, value : %+v \n", m1, m1)
}
```

执行结果:

```sh
Hello, 我是一个int, Let'GoSayHello -> Type m1 : main.MyInt, value : 100

# Start old m : 100 -> new m : 1024

# End old m : 1024 -> new m : 1024

ChangeM -> Type m1 : main.MyInt, value : 1024
```

Tips : 非常注意，非本地类型不能定义方法，也就是说我们不能给别的包的类型定义方法。



匿名结构体与匿名字段

> 描述: 在定义一些临时数据结构等场景下还可以使用匿名结构体。

示例演示:

```go
// 匿名结构体(只能使用一次，所以常常使用与临时场景)
// 2.匿名结构体(只能使用一次，所以常常使用与临时场景)
func demo2() {
	var temp struct {title string;address []string}
	temp.title = "地址信息"
	temp.address = []string{"中国", "重庆", "江北区"}
	fmt.Printf("Type of temp : %T\nStruct define: %#v \nValue : %v\n", temp, temp, temp)
}
```

```sh
Type of temp : struct { title string; address []string }
Struct define: struct { title string; address []string }{title:"地址信息", address:[]string{"中国", "重庆", "江北区"}}
Value : {地址信息 [中国 重庆 江北区]}
```

> 描述: 结构体允许其成员字段在声明时没有字段名而只有类型，这种没有名字的字段就称为匿名字段。

Tips: 这里匿名字段的说法并不代表没有字段名，而是默认会采用类型名作为字段名，结构体要求字段名称必须唯一，因此一个结构体中同种类型的匿名字段只能有一个。

```go
type Anonymous struct {
	string
	int
}
func demo4() {
	a1 := Anonymous{"WeiyiGeek", 18}
	fmt.Printf("Struct: %#v ，字段1: %v , 字段2: %v \n", a1, a1.string, a1.int)
}
```

```sh
Struct: main.Anonymous{string:"WeiyiGeek", int:18} ，字段1: WeiyiGeek , 字段2: 18
```


嵌套结构体与匿名字段

> 描述: 结构体中可以嵌套包含另一个结构体或结构体指针, 并且上面user结构体中嵌套的Address结构体也可以采用匿名字段的方式。

并且为了防止嵌套结构体的相同的字段名冲突，所以在这种情况下为了避免歧义需要通过指定具体的内嵌结构体字段名。

```go
//Address 地址结构体
type Address struct {
	Province string
	City     string
}

//Email 邮箱结构体
type Email struct {
	Account    string
	CreateTime string
}

//User 用户结构体
type User struct {
	Name    string
	Gender  string
	Address Address
}

//AnonUser 用户结构体
type AnonUser struct {
	Name    string
	Gender  string
	Address // 采用结构体的匿名字段来嵌套结构体Address
	Email   // 采用结构体的匿名字段来嵌套结构体Email
}

// 1.嵌套结构体
func demo1() {
	// 结构体初始化
	user := User{
		Name:   "WeiyiGeek",
		Gender: "男",
		Address: Address{
			Province: "重庆",
			City:     "重庆",
		},
	}
	fmt.Printf("Struct : %#v \n", user)
	fmt.Printf("Name = %v, Address City = %v \n", user.Name, user.Address.City)
}

// 2.嵌套匿名字段防止字段名称冲突
func demo2() {
	var anonuser = AnonUser{
		Name:   "WeiyiGeek",
		Gender: "男",
		Address: Address{
			"重庆",
			"重庆",
		},
		Email: Email{
			"Master@weiyigeek.top",
			"2021年8月23日 10:21:36",
		},
	}
	fmt.Printf("Struct : %#v\n", anonuser)
	fmt.Printf("Name = %v,Address Province = %v, Email Account = %v \n", anonuser.Name, anonuser.Address.Province, anonuser.Email.Account)
}
```

```sh
// 嵌套结构体
Struct : main.User{Name:"WeiyiGeek", Gender:"男", Address:main.Address{Province:"重庆", City:"重庆"}}
Name = WeiyiGeek, Address City = 重庆

//嵌套匿名字段
Struct : main.AnonUser{Name:"WeiyiGeek", Gender:"男", Address:main.Address{Province:"重庆", City:"重庆"}, Email:main.Email{Account:"Master@weiyigeek.top", CreateTime:"2021年8月23日 10:21:36"}}
Name = WeiyiGeek,Address Province = 重庆, Email Account = Master@weiyigeek.top
```

Tips : 当访问结构体成员时会先在结构体中查找该字段，找不到再去嵌套的匿名字段中查找。

结构体的“继承”

> 描述: Go语言中使用结构体也可以实现其他编程语言中面向对象的继承。

```go
package main

import "fmt"

// 父
type Animal struct{ name string }

func (a *Animal) voice(v string) {
	fmt.Printf("我是动物，我叫 %v, 我会叫 %s,", a.name, v)
}

// 子
type Dog struct {
	eat string
	*Animal
}

func (d *Dog) love() {
	fmt.Printf("狗狗喜欢吃的食物是 %v.\n", d.eat)
}

type Cat struct {
	eat string
	*Animal
}

func (c *Cat) love() {
	fmt.Printf("猫猫喜欢吃的食物是 %v.\n", c.eat)

}

func main() {
	d1 := &Dog{
		//注意嵌套的是结构体指针
		Animal: &Animal{
			name: "小黄",
		},
		eat: "bone",
	}
	d1.voice("汪汪.汪汪.")
	d1.love()

	c1 := &Cat{
		//注意嵌套的是结构体指针
		Animal: &Animal{
			name: "小白",
		},
		eat: "fish",
	}
	c1.voice("喵喵.喵喵.")
	c1.love()

}
```

```sh
我是动物，我叫 小黄, 我会叫 汪汪.汪汪.,狗狗喜欢吃的食物是 bone. 

我是动物，我叫 小白, 我会叫 喵喵.喵喵.,猫猫喜欢吃的食物是 fish.
```

结构体与“JSON”

> 描述: JSON(JavaScript Object Notation) 是一种轻量级的数据交换格式,其优点是易于人阅读和编写，同时也易于机器解析和生成。

Tips : JSON键值对是用来保存JS对象的一种方式，键/值对组合中的键名写在前面并用双引号""包裹，使用冒号:分隔，然后紧接着值；多个键值之间使用英文,分隔。

在Go中我们可以通过结构体序列号生成json字符串，同时也能通过json字符串反序列化为结构体得实例化对象，在使用json字符串转换时, 我们需要用到"encoding/json"包。

结构体标签（Tag）

> 描述: Tag是结构体的元信息，可以在运行的时候通过反射的机制读取出来，Tag在结构体字段的后方定义，由一对反引号包裹起来，具体的格式如下：key1:"value1" key2:"value2",可以看到它由一个或多个键值对组成。键与值使用冒号分隔，值用双引号括起来。同一个结构体字段可以设置多个键值对tag，不同的键值对之间使用空格分隔。

例如: 我们为Student结构体的每个字段定义json序列化时使用的Tag。

```go
type Student struct {
	ID     int    `json:"id"` //通过指定tag实现json序列化该字段时的key
	Gender string //json序列化是默认使用字段名作为key
	name   string //私有不能被json包访问
}
```

注意事项： 为结构体编写Tag时，必须严格遵守键值对的规则。结构体标签的解析代码的容错能力很差，一旦格式写错，编译和运行时都不会提示任何错误，通过反射也无法正确取值。例如不要在key和value之间添加空格。

```go
package main

import (
	"encoding/json"
	"fmt"
)

// 结构体转json字符串的三种示例
// 结构体中的字段首字母大小写影响的可见性，表示不能对外使用
type Person1 struct{ name, sex string }

// 结构体对象字段可以对外使用
type Person2 struct{ Name, Sex string }

// 但json字符串中键只要小写时可以采用此种方式
type Person3 struct {
	Name string `json:"name"`
	Sex  string `json:"age"`
}

// # 结构体实例化对象转JSON字符串
func serialize() {
	// 示例1.字段首字母大小写影响的可见性
	person1 := &Person1{"weiyigeek", "男孩"}
	person2 := &Person2{"WeiyiGeek", "男生"}
	person3 := &Person3{"WeiyiGeek", "男人"}

	//序列化
	p1, err := json.Marshal(person1)
	p2, err := json.Marshal(person2)
	p3, err := json.Marshal(person3)
	if err != nil {
		fmt.Printf("Marshal Failed ：%v", err)
		return
	}
	
	// 由于返回是一个字节切片，所以需要强转为字符串
	fmt.Printf("person1 -> %v\nperson2 -> %v\nperson3 -> %v\n", string(p1), string(p2), string(p3))

}

// # JSON字符串转结构体实例化对象

type Person4 struct {
	Name string    `json:"name"`
	Sex  string    `json:"sex"`
	Addr [3]string `json:"addr"`
}

func unserialize() {
	jsonStr := `{"name": "WeiyiGeek","sex": "man","addr": ["中国","重庆","渝北"]}`
	p4 := Person4{}

	// 在其内部修改p4的值
	err := json.Unmarshal([]byte(jsonStr), &p4)
	if err != nil {
		fmt.Printf("Unmarhal Failed: %v", err)
		return
	}
	fmt.Printf("jsonStr -> Person4 : %#v\nPerson4.name : %v\n", p4, p4.Name)

}

func main() {
	serialize()
	unserialize()
}
```

```sh
person1 -> {}
person2 -> {"Name":"WeiyiGeek","Sex":"男生"}
person3 -> {"name":"WeiyiGeek","age":"男人"}
jsonStr -> Person4 : main.Person4{Name:"WeiyiGeek", Sex:"man", Addr:[3]string{"中国", "重庆", "渝北"}}
Person4.name : WeiyiGeek
```


结构体和方法补充知识点

> 描述: 因为slice和map这两种数据类型都包含了指向底层数据的指针，因此我们在需要复制它们时要特别注意。

```go
package main

import "fmt"

type Person struct {
	name   string
	age    int8
	dreams []string
}

// 不推荐的方式
func (p *Person) SetDreams(dreams []string) {
	p.dreams = dreams
}

// 正确的做法是在方法中使用传入的slice的拷贝进行结构体赋值。
func (p *Person) NewSetDreams(dreams []string) {
	p.dreams = make([]string, len(dreams))
	copy(p.dreams, dreams)
}

func main() {
	// (1) 不安全的方式
	p1 := Person{name: "小王子", age: 18}
	data := []string{"吃饭", "睡觉", "打豆豆"}
	p1.SetDreams(data)
	// 你真的想要修改 p1.dreams 吗？
	data[1] = "不睡觉"        // 会覆盖更改切片中的值从而影响p1中的dreams字段中的值
	fmt.Println(p1.dreams) // [吃饭 不睡觉 打豆豆]

	// (2) 推荐方式
	p2 := Person{name: "WeiyiGeek", age: 18}
	data2 := []string{"计算机", "网络", "编程"}
	p2.NewSetDreams(data2)
	data2[1] = "NewMethod" // 由于NewSetDreams返回中是将拷贝的副本给p2的dreams字段，所以此处更改不会影响其值，
	fmt.Println(p2.dreams) // [计算机 网络 编程]

}
```

```sh
[吃饭 不睡觉 打豆豆]

[计算机 网络 编程]
```

Tips: 同样的问题也存在于返回值slice和map的情况，在实际编码过程中一定要注意这个问题。 



接口类型

> 描述: 在Go语言中接口（interface）是一种类型，一种抽象的类型, 其定义了一个对象的行为规范，只定义规范不实现，由具体的对象来实现规范的细节。

如 interface 是一组 method 的集合，是duck-type programming的一种体现。接口做的事情就像是定义一个协议（规则），只要一台机器有洗衣服和甩干的功能，我就称它为洗衣机。不关心属性（数据），只关心行为（方法）。

Tips: 为了保护你的Go语言职业生涯，请牢记接口（interface）是一种类型。 

> Q: 为什么要使用接口?
> 在我们编程过程中会经常遇到：
>
> 比如一个网上商城可能使用支付宝、微信、银联等方式去在线支付，我们能不能把它们当成“支付方式”来处理呢？
> 比如三角形，四边形，圆形都能计算周长和面积，我们能不能把它们当成“图形”来处理呢？
> 比如销售、行政、程序员都能计算月薪，我们能不能把他们当成“员工”来处理呢？
> 例如:面的代码中定义了猫和狗，然后它们都会叫，你会发现main函数中明显有重复的代码，如果我们后续再加上猪、青蛙等动物的话，我们的代码还会一直重复下去。那我们能不能把它们当成“能叫的动物”来处理呢？ 

```go
type Cat struct{}
func (c Cat) Say() string { return "喵喵喵" }
type Dog struct{}
func (d Dog) Say() string { return "汪汪汪" }
func main() {
	c := Cat{}
	fmt.Println("猫:", c.Say())  // 猫: 喵喵喵
	d := Dog{}
	fmt.Println("狗:", d.Say())  // 狗: 汪汪汪
} 
```

Go语言中为了解决类似上面的问题，就设计了接口这个概念。接口区别于我们之前所有的具体类型，接口是一种抽象的类型。当你看到一个接口类型的值时，你不知道它是什么，唯一知道的是通过它的方法能做什么。



接口的定义

描述: Go语言提倡面向接口编程,每个接口由数个方法组成，接口的定义格式如下：

```go
type 接口类型名 interface{
    方法名1( 参数列表1 ) 返回值列表1
    方法名2( 参数列表2 ) 返回值列表2
    …
}
```


参数说明:

接口名：使用type将接口定义为自定义的类型名。Go语言的接口在命名时一般会在单词后面添加er，如有写操作的接口叫Writer，有字符串功能的接口叫Stringer等。接口名最好要能突出该接口的类型含义。
方法名：当方法名首字母是大写且这个接口类型名首字母也是大写时，这个方法可以被接口所在的包（package）之外的代码访问。
参数列表、返回值列表：参数列表和返回值列表中的参数变量名可以省略。
基础示例:

```go
type writer interface{
    Write([]byte) error
}
```

Tips: 当你看到这个接口类型的值时，你不知道它是什么，唯一知道的就是可以通过它的Write方法来做一些事情。

Tips :实现接口的条件, 即一个对象只要全部实现了接口中的方法，那么就实现了这个接口。换句话说接口就是一个需要实现的方法列表。



接口类型变量

> Q: 那实现了接口有什么用呢？
>
> 答: 接口类型变量能够存储所有实现了该接口的实例，接口类型变量实际上你可以看做一个是一个合约。

基础示例:

```go
// 定义一个接口类型writer的变量w。
var w writer // 声明一个writer类型的变量w
```

Tips： 观察下面的代码，体味此处_的妙用

```go
// 摘自gin框架routergroup.go
type IRouter interface{ ... }
type RouterGroup struct { ... }
var _ IRouter = &RouterGroup{}  // 确保RouterGroup实现了接口IRouter
```

```go
package main

import "fmt"

// 接口声明定义以及约定必须实现的方法
type speaker interface {
	speak()
	eat(string)
}

// 人结构体
type person struct{ name, language string }
func (p person) speak() {
	fmt.Printf("我是人类，我说的是%v, 我叫%v\n", p.language, p.name)
}
func (p person) eat(food string) { fmt.Printf("喜欢的食物: %v\n", food) }

// 猫结构体
type cat struct{ name, language string }
func (c cat) speak() {
	fmt.Printf("动物猫，说的是%v, 叫%v\n", c.language, c.name)
}
func (c cat) eat(food string) { fmt.Printf("喜欢的食物: %v\n", food) }

// 狗结构体
type dog struct{ name, language string }
func (d dog) speak() {
	fmt.Printf("动物狗，说的是%v, 叫%v\n", d.language, d.name)
}
func (d dog) eat(food string) { fmt.Printf("喜欢的食物: %v\n", food) }

func talk(s speaker) {
	s.speak()
}

// (1) 接口基础使用演示
func demo1() {
	p := person{"WeiyiGeek", "汉语"}
	c := cat{"小白", "喵喵 喵喵..."}
	d := dog{"阿黄", "汪汪 汪汪...."}
	talk(p)
	talk(c)
	talk(d)
}

// (2) 接口类型的使用(可看作一种合约)方法不带参数以及方法带有参数
func demo2() {
	// 定义一个接口类型writer的变量w。
	var s speaker
	fmt.Printf("Type %T\n", s) // 动态类型

	s = person{"接口类型-唯一", "汉语"} // 动态值
	fmt.Printf("\nType %T\n", s) // 动态类型
	s.speak()
	s.eat("瓜果蔬菜")
	
	s = cat{"接口类型-小白", "喵喵..."} // 动态值
	fmt.Printf("\nType %T\n", s) // 动态类型
	s.speak()
	s.eat("fish")
	
	s = dog{"接口类型-阿黄", "汪汪..."} // 动态值
	fmt.Printf("\nType %T\n", s) // 动态类型
	s.speak()
	s.eat("bone")

}

func main() {
	demo1()
	fmt.Println()
	demo2()
}
```

```sh
我是人类，我说的是汉语, 我叫WeiyiGeek
动物猫，说的是喵喵 喵喵..., 叫小白
动物狗，说的是汪汪 汪汪...., 叫阿黄

Type <nil>

Type main.person
我是人类，我说的是汉语, 我叫接口类型-唯一
喜欢的食物: 瓜果蔬菜

Type main.cat
动物猫，说的是喵喵..., 叫接口类型-小白
喜欢的食物: fish

Type main.dog
动物狗，说的是汪汪..., 叫接口类型-阿黄
喜欢的食物: bone
```

注意: 带参数和不带参数的函数,在接口中实现的不是同一个方法，所以当某个结构体中没有完全实现接口中的方法将会报错。



接口实现之值接收者和指针接收者

> Q: 使用值接收者实现接口和使用指针接收者实现接口有什么区别呢?
>
> 值接收者实现接口: 结构体类型和结构体指针类型的变量都可以存储，由于因为Go语言中有对指针类型变量求值的语法糖，结构体指针变量内部会自动求值（取指针地址中存储的值）。
>
> 指针接收者实现接口: 只能存储结构体指针类型的变量。

```go
package main

import (
	"fmt"
)

// 接口类型声明
// (1) 值接收者实现接口
type Mover interface {
	move()
}
type dog struct{}
func (d dog) move() { fmt.Println("值接收者实现接口 -> 狗...移动....")  } // 关键点

// 使用值接收者实现接口之后，不管是dog结构体还是结构体指针*dog类型的变量都可以赋值给该接口变量.
func demo1() {
	var m1 Mover
	var d1 = dog{} // 值类型
	m1 = d1        // m1可以接收dog类型的变量
	fmt.Printf("Type : %#v \n", m1)
	m1.move()

	var d2 = &dog{} // 指针类型
	m1 = d2         // x可以接收指针类型的(*dog)类型的变量
	fmt.Printf("Type : %#v \n", m1)
	m1.move()

}

// (2)指针接收者实现接口
type Runer interface{ run() }
type cat struct{}
func (c *cat) run() { fmt.Println("指针接收者实现接口 -> 猫...跑....") }
// 此时实现run接口的是*cat类型，所以不能给m1传入cat类型的c1，此时x只能存储*dog类型的值。
func demo2() {
	var m1 Runer
	var c1 = cat{}
	//m1不可以接收dog类型的变量
	// m1 = c1 // 报错信息: cannot use c1 (variable of type cat) as Runer value in assignment: missing method run (run has pointer receiver)compilerInvalidIfaceAssign
	fmt.Printf("Type : %#v \n", c1)

	//m1只能接收*dog类型的变量
	var c2 = &cat{}
	m1 = c2
	fmt.Printf("Type : %#v \n", c2)
	m1.run()

}
func main() {
	demo1()
	fmt.Println()
	demo2()
}



```

```sh
Type : main.dog{}
值接收者实现接口 -> 狗...移动....
Type : &main.dog{}
值接收者实现接口 -> 狗...移动....

Type : main.cat{}
Type : &main.cat{}
指针接收者实现接口 -> 猫...跑....
```


面试题: 注意这是一道你需要回答“能”或者“不能”的题！
问: 首先请观察下面的这段代码，然后请回答这段代码能不能通过编译？

```go
package main

import "fmt"

type People interface {
	Speak(string) string
}

type Student struct{}

func (stu *Student) Speak(think string) (talk string) {
	if think == "man" {
		talk = "你好,帅哥"
	} else {
		talk = "您好,美女"
	}
	return
}

func main() {
	var peo People = Student{} // 此处为关键点
	think := "woman"
	fmt.Println(peo.Speak(think))
}
```

答案: 是不行会报 ./interface.go:21:6: cannot use Student{} (type Student) as type People in assignment: Student does not implement People (Speak method has pointer receiver) (exit status 2)错误，由于指针接收者实现接口必须是有指针类型的结构体实例化对象以及其包含的方法。

接口与类型

一个类型实现多个接口

> 描述: 一个结构体类型可以同时实现多个接口，而接口间彼此独立，不知道对方的实现。

例如: 狗可以叫也可以动,我们就分别定义Sayer接口和Mover接口

```go
// Sayer 接口
type Sayer interface { say() }
// Mover 接口
type Mover interface { move() }
// dog既可以实现Sayer接口，也可以实现Mover接口。
type dog struct {	name string }
// 实现Sayer接口
func (d dog) say() { fmt.Printf("%s会叫 汪汪汪\n", d.name) }
// 实现Mover接口
func (d dog) move() { fmt.Printf("%s会动 \n", d.name) }

func main() {
  var a = dog{name: "旺财"}
	var x Sayer = a // 将dog类型赋予给Sayer接口类型的变量x，此时它可以调用say方法
	var y Mover = a // 将dog类型赋予给Mover接口类型的变量y，此时它可以调用move方法
	x.say() // 旺财会叫 汪汪汪
	y.move() // 旺财会动
}
```


多个类型实现同一接口

> 描述: Go语言中不同的类型还可以实现同一接口,比如我们前面Person、Cat、Dog结构体类型中实现的Speak()方法。

例如：我们定义一个Mover接口，它要求结构体类型中必须有一个move方法, 如狗可以动，汽车也可以动。

```go
// Mover 接口
type Mover interface { move() }
type dog struct { name string }
type car struct { brand string }
// dog类型实现Mover接口
func (d dog) move() {	fmt.Printf("%s会跑\n", d.name) }
// car类型实现Mover接口
func (c car) move() { fmt.Printf("%s速度120迈\n", c.brand) }
func main() {
	var x Mover
	var a = dog{name: "旺财"}
  x = a
	x.move() // 旺财会跑
	var b = car{brand: "保时捷"}
	x = b
	x.move() // 保时捷速度120迈
}
```


非常注意: 并且一个接口的方法，不一定需要由一个类型完全实现，接口的方法可以通过在类型中嵌入其他类型或者结构体来实现。

示例演示:

```go
package main

import "fmt"

// 接口类型
type android interface {
	telephone(int64)
	music()
}

// 结构体声明 实现music方法
type mp3 struct{}
// 实现接口中的方法
func (m *mp3) music() { fmt.Println("播放音乐.....")}

// 结构体声明
type mobilephone struct {
	production string
	mp3        // 嵌入mp3结构体并拥有它的方法
}

// 实现接口中的方法
func (mb *mobilephone) telephone(number int64) { fmt.Printf("%v 手机, 正在拨打 %v 电话....\n", mb.production, number)}

func main() {
	// android 接口类型
	var a android
	// 指针类型结构体变量mb
	var mp = &mobilephone{production: "小米"}
	a = mp
	fmt.Printf("Type : %#v\n", a) // android 接口类型变量输出
	a.telephone(10086)
	a.music()
}
```

```sh
Type : &main.mobilephone{production:"小米", mp3:main.mp3{}}
小米 手机, 正在拨打 10086 电话....
播放音乐.....
```

接口嵌套

> 描述: 接口与接口间可以通过嵌套创造出新的接口,嵌套得到的接口的使用与普通接口一样，这里我们让cat实现animal接口。

示例演示:

```go
// Sayer 接口
type Sayer interface {say()}
// Mover 接口
type Mover interface {move()}
// 接口嵌套
type animal interface {
	Sayer
	Mover
}
// cat 结构体
type cat struct {
	name string
}
// 接口方法的实现
func (c cat) say() {fmt.Printf("%v 喵喵喵",c.name)}
func (c cat) move() {fmt.Printf("%v 猫会动",c.name)}
func main() {
	var x animal
	x = cat{name: "花花"}
	x.move() //喵喵喵
	x.say()  //猫会动
}
```


空接口

空接口的定义

> 描述: 空接口是指没有定义任何方法的接口,因此任何类型都实现了空接口, 该类型的变量可以存储任意类型的变量。他会在我们以后GO编程中常常出现。

例如:

```go
// interface 是关键字，并不是类型。
// 方式1.但一般不会采用此种方式
var empty interface{}

// 方式2.我们可以直接忽略接口名称(空接口类型)
interface{}
```


空接口的应用


空接口作为函数的参数: 使用空接口实现可以接收任意类型的函数参数。

空接口作为map的值: 使用空接口实现可以保存任意值的字典。
示例演示:

```go
package main

import "fmt"

// (1) 空接口作为函数参数
func showType(a interface{}) { fmt.Printf("参数类型:%T, 参数值:%v\n", a, a) }
func main() {
	// (2) 空接口作为map的值
	var m1 map[string]interface{}     // 类似于Java中的 Map<String,Object> m1
	m1 = make(map[string]interface{}) // 为Map申请一块内存空间
	// 可以存储任意类型的值
	m1["name"] = "WeiyiGeek"
	m1["age"] = 20
	m1["sex"] = true
	m1["hobby"] = [...]string{"Computer", "NetSecurity", "Go语言编程学习"}

	fmt.Printf("#空接口作为map的值\n%#v", m1)
	fmt.Println(m1)
	
	fmt.Printf("\n#空接口作为函数参数\n")
	showType(nil)
	showType([]byte{'a'})
	showType(true)
	showType(1024)
	showType("我是一串字符串")

}
```

```sh
#空接口作为map的值
map[string]interface {}{"age":20, "hobby":[3]string{"Computer", "NetSecurity", "Go语言编程学习"}, "name":"WeiyiGeek", "sex":true}
map[age:20 hobby:[Computer NetSecurity Go语言编程学习] name:WeiyiGeek sex:true]

#空接口作为函数参数
参数类型:<nil>, 参数值:<nil>
参数类型:[]uint8, 参数值:[97]
参数类型:bool, 参数值:true
参数类型:int, 参数值:1024
参数类型:string, 参数值:我是一串字符串
```


Tips : 因为空接口可以存储任意类型值的特点，所以空接口在Go语言中的使用十分广泛。

接口之类型断言

> 描述: 空接口可以存储任意类型的值，那我们如何获取其存储的具体数据呢？

接口值

> 描述: 一个接口的值（简称接口值）是由一个具体类型和具体类型的值两部分组成的,这两部分分别称为接口的动态类型和动态值。

我们来看一个具体的例子:

```go
var w io.Writer
w = nil
w = os.Stdout
w = new(bytes.Buffer)
```

请看下图分解： 

![WeiyiGeek.动态类型与动态值](https://i0.hdslb.com/bfs/article/df0b4dfdfcd79b3a5c181b7f529c41504dddece8.png@942w_761h_progressive.png)

想要判断空接口中的值这个时候就可以使用类型断言，其语法格式：x.(T)，其中：

> x：表示类型为interface{}的变量
>
> T：表示断言x可能是的类型。

该语法返回两个参数，第一个参数是x转化为T类型后的变量，第二个值是一个布尔值，若为true则表示断言成功，为false则表示断言失败。

示例演示

```go
package main
import "fmt"
// 示例1.采用if进行判断断言
func assert(x interface{}) {
	v, ok := x.(string) // v 接受是string类型
	if ok {
		fmt.Printf("assert successful : %v, typeof %T\n", v, v)
	} else {
		fmt.Printf("assert failed 非 string 类型! : %v, typeof %T\n", x, x)
	}
}
func demo1() {
	var x interface{}
	x = "WeiyiGeek"
	assert(x) // assert successful : WeiyiGeek, typeof string
	x = 1024
	assert(x) // assert failed 非 string 类型! : 1024, typeof int
}

// 示例2.如果要断言多次就需要写多个if判断，这个时候我们可以使用switch语句来实现：
func justifyType(x interface{}) {
	switch v := x.(type) {
	case string:
		fmt.Printf("x is a string，value is %v\n", v)
	case int:
		fmt.Printf("x is a int is %v\n", v)
	case bool:
		fmt.Printf("x is a bool is %v\n", v)
	default:
		fmt.Println("unsupport type！")
	}
}
func demo2() {
	var x interface{}
	x = "i'm string"
	justifyType(x)
	x = 225
	justifyType(x)
	x = true
	justifyType(x)
}

func main() {
	demo1()
	fmt.Println()
	demo2()
}
```


执行结果:

```sh
assert successful : WeiyiGeek, typeof string
assert failed 非 string 类型! : 1024, typeof int

x is a string，value is i'm string
x is a int is 225
x is a bool is true
```

接口总结:
描述: 关于需要注意的是只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要定义接口,不要为了接口而写接口，那样只会增加不必要的抽象，导致不必要的运行时损耗。 
