log.日志标准库

描述: 无论是软件开发的调试阶段还是软件上线之后的运行阶段，日志一直都是非常重要的一个环节，我们也应该养成在程序中记录日志的好习惯。

描述: Go语言内置的log包实现了简单的日志服务, 本文介绍了标准库log的基本使用。

1.使用Logger

描述: log包定义了Logger类型，该类型提供了一些格式化输出的方法。本包也提供了一个预定义的“标准”logger，可以通过调用函数Print系列(Print|Printf|Println）、Fatal系列（Fatal|Fatalf|Fatalln）、和Panic系列（Panic|Panicf|Panicln）来使用，比自行创建一个logger对象更容易使用。

例如，我们可以像下面的代码一样直接通过log包来调用上面提到的方法，默认它们会将日志信息打印到终端界面：

```
package main

import (
	"log"
)

func main() {
	log.Println("这是一条很普通的日志。")
	v := "很普通的"
	log.Printf("这是一条%s日志。\n", v)
	log.Fatalln("这是一条会触发fatal的日志。")
	log.Panicln("这是一条会触发panic的日志。")
}
```
编译并执行上面的代码会得到如下输出：
```
2017/06/19 14:04:17 这是一条很普通的日志。
2017/06/19 14:04:17 这是一条很普通的日志。
2017/06/19 14:04:17 这是一条会触发fatal的日志。
```

Tips : logger 会打印每条日志信息的日期、时间，默认输出到系统的标准错误。
Tips : Fatal 系列函数会在写入日志信息后调用 os.Exit(1)。
Tips : Panic 系列函数会在写入日志信息后调用 panic。



2.配置Logger

标准配置
描述: 默认情况下的logger只会提供日志的时间信息，但是很多情况下我们希望得到更多信息，比如记录该日志的文件名和行号等。log标准库中为我们提供了定制这些设置的方法。

Tips : log 标准库中的Flags函数会返回标准logger的输出配置，而SetFlags函数用来设置标准logger的输出配置。
```
func Flags() int
func SetFlags(flag int)
```
flag 选项
描述: log标准库提供了如下的flag选项，它们是一系列定义好的常量。
```
const (
  // 控制输出日志信息的细节，不能控制输出的顺序和格式。
  // 输出的日志在每一项后会有一个冒号分隔：例如 2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
  Ldate         = 1 << iota     // 日期：2009/01/23
  Ltime                         // 时间：01:23:23
  Lmicroseconds                 // 微秒级别的时间：01:23:23.123123（用于增强Ltime位）
  Llongfile                     // 文件全路径名+行号： /a/b/c/d.go:23
  Lshortfile                    // 文件名+行号：d.go:23（会覆盖掉Llongfile）
  LUTC                          // 使用UTC时间
  LstdFlags     = Ldate | Ltime // 标准logger的初始值
)
```

下面我们在记录日志之前先设置一下标准logger的输出选项如下：
```
func main() {
  log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
  log.Println("这是一条很普通的日志。")
}
```

编译执行后得到的输出结果如下：
```
2017/06/19 14:05:17.494943 .../log_demo/main.go:11: 这是一条很普通的日志。
```

配置日志前缀
描述: log标准库中还提供了关于日志信息前缀的两个方法：
```
func Prefix() string
func SetPrefix(prefix string)
```
其中Prefix函数用来查看标准logger的输出前缀，SetPrefix函数用来设置输出前缀。
```
func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println("这是一条很普通的日志。")
	log.SetPrefix("[小王子]")
	log.Println("这是一条很普通的日志。")
}
```
上面的代码输出如下：
```
[小王子]2017/06/19 14:05:57.940542 .../log_demo/main.go:13: 这是一条很普通的日志。
```
Tips : 这样我们就能够在代码中为我们的日志信息添加指定的前缀，方便之后对日志信息进行检索和处理。



配置日志输出位置
描述: SetOutput 函数用来设置标准logger的输出目的地，默认是标准错误输出。func 
```
func SetOutput(w io.Writer)
```
例如，下面的代码会把日志输出到同目录下的xx.log文件中。
```
func main() {
	logFile, err := os.OpenFile("./xx.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println("这是一条很普通的日志。")
	log.SetPrefix("[小王子]")
	log.Println("这是一条很普通的日志。")
}
```
如果你要使用标准的logger，我们通常会把上面的配置操作写到init函数中。
```
func init() {
	logFile, err := os.OpenFile("./xx.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
}
```



3.创建Logger

描述: log标准库中还提供了一个创建新logger对象的构造函数–New，支持我们创建自己的logger示例。

New 函数的签名如下：
```
func New(out io.Writer, prefix string, flag int) *Logger
```
New 创建一个Logger对象。其中参数out设置日志信息写入的目的地。参数prefix会添加到生成的每一条日志前面。参数flag定义日志的属性（时间、文件等）。

举个例子：
```
func main() {
  logger := log.New(os.Stdout, "<New>", log.Lshortfile|log.Ldate|log.Ltime)
  logger.Println("这是自定义的logger记录的日志。")
}
```
将上面的代码编译执行之后，得到结果如下：
```
<New>2017/06/19 14:06:51 main.go:34: 这是自定义的logger记录的日志。
```

4.总结说明

描述: Go内置的log库功能有限，例如无法满足记录不同级别日志的情况，我们在实际的项目中根据自己的需要选择使用第三方的日志库，如logrus、zap等。

案例演示:
```
package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func delay() {
	time.Sleep(time.Duration(1) * time.Second)
}

func main() {
	v := "信息警告提示"

	// 常规信息
	log.Printf("[-] 此处是这一条的日志信息 ：%s 。\n", v)
	delay()
	log.Println("[-] 这是一条换行的日志信息。")
	
	// Flag 选项
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ltime | log.Ldate)
	log.Println("[-] 这是换行的日志信息(Flag 选项)。")
	
	// 配置日志前缀
	log.SetPrefix("[WeiyiGeek] ")
	log.Println("这是换行的日志信息(配置日志前缀)。")
	
	// 创建 logger 使用示例
	logger := log.New(os.Stdout, "[New-WeiGeek] ", log.Lshortfile|log.Ldate|log.Ltime)
	logger.Println("这是自定义的logger记录的日志。")
	
	// 配置日志输出
	logFile, err := os.OpenFile("./Logger.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)               // 每次执行不会覆盖其内容，而会向其追加内容。
	log.Println("[-] 这是要配置日志输出到文件之中（1）") // 会输出到 /Logger.log
	log.Println("[-] 这是要配置日志输出到文件之中（2）") // 会输出到 /Logger.log
	
	// 执行 Panicln 则会exit
	delay()
	log.Fatalln("[-] 这是一条会触发fatal的日志。") // 会输出到 /Logger.log
	delay()
	log.Panicln("[-] 这是一条会触发panic的日志。")
}
```
执行结果:
```
$ cat Logger.log
[WeiyiGeek] 2021/08/03 11:35:29.482279 /home/weiyigeek/app/project/go/src/weiyigeek.top/studygo/package/02logger.go:41: [-] 这是要配置日志输出到文件之中（1）
[WeiyiGeek] 2021/08/03 11:35:29.482338 /home/weiyigeek/app/project/go/src/weiyigeek.top/studygo/package/02logger.go:42: [-] 这是要配置日志输出到文件之中（2）
[WeiyiGeek] 2021/08/03 11:35:30.482714 /home/weiyigeek/app/project/go/src/weiyigeek.top/studygo/package/02logger.go:46: [-] 这是一条会触发fatal的日志。
```
0x04 time.时间标准库

描述: 本文主要介绍了Go语言内置的time包的基本用法，time包提供了时间的显示和测量用的函数。

Time 报预定义的版式，其定义的时间为2006年1月2号 15点04分05秒是Go语言诞生的日子。
```
const (
  ANSIC       = "Mon Jan _2 15:04:05 2006"
  UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
  RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
  RFC822      = "02 Jan 06 15:04 MST"
  RFC822Z     = "02 Jan 06 15:04 -0700"            // 使用数字表示时区的RFC822
  RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
  RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
  RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"  // 使用数字表示时区的RFC1123
  RFC3339     = "2006-01-02T15:04:05Z07:00"
  RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
  Kitchen     = "3:04PM"

  // 方便的时间戳
  Stamp      = "Jan _2 15:04:05"
  StampMilli = "Jan _2 15:04:05.000"
  StampMicro = "Jan _2 15:04:05.000000"
  StampNano  = "Jan _2 15:04:05.000000000"
)

// 可以获取输出其定义常量
fmt.Println(time.ANSIC) # Mon Jan _2 15:04:05 2006
```

Tips: 日历的计算采用的是公历。

1.时间类型

描述: 我们可以通过time.Now()函数获取当前的时间对象，然后获取时间对象的年月日时分秒等信息。

示例代码如下：
```
func demo1() {
	// 获取当前时间
	now := time.Now()

	// 输出当前时间
	fmt.Printf("Current Localtion Time is ：%v, \nUTC Time is: %v\n", now.Local(), now.UTC())
	
	// 分别获取当前时间的年月日/时分秒
	year := now.Year()     //年
	month := now.Month()   //月
	day := now.Day()       //日
	hour := now.Hour()     //小时
	minute := now.Minute() //分钟
	second := now.Second() //秒
	y, m, d := now.Date()  //年月日
	week := now.Weekday()  //周
	fmt.Printf("Current Localtion Time Format: %d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
	fmt.Printf("Year : %d ,Month : %v ,Day : %d ,WeekDay : %v\n", y, m, d, week)
}

// === 执行结果 ===
Current Localtion Time is ：2021-09-27 12:45:41.295602237 +0800 CST,
UTC Time is: 2021-09-27 04:45:41.295602237 +0000 UTC
Current Localtion Time Format: 2021-09-27 12:45:41
Year : 2021 ,Month : September ,Day : 27 ,WeekDay : Monday
```

2.时间戳

描述: 时间戳是自1970年1月1日（00:00:00）至当前时间的总毫秒数与时区无关,它也被称为Unix时间戳（UnixTimestamp）,我们可以使用time.Unix()函数将时间戳转为时间格式。

// # Unix创建一个本地时间，对应sec和nsec表示的Unix时间（从January 1, 1970 UTC至该时间的秒数和纳秒数）。
```
func Unix(sec int64, nsec int64) Time
// nsec的值在[0, 999999999]范围外是合法的。
```
代码演示:
```
//获取当前时间
now := time.Now()
timestamp1 := now.Unix()     //时间戳
timestamp2 := now.UnixNano() //纳秒时间戳
fmt.Printf("current timestamp : %v\n", timestamp1)
fmt.Printf("current timestamp nanosecond: %v\n", timestamp2)

//将时间戳转为时间格式(秒数，纳秒数)
timeObj := time.Unix(timestamp1, 0)
fmt.Println("时间戳转换后的时间 :", timeObj)
year := timeObj.Year()     //年
month := timeObj.Month()   //月
day := timeObj.Day()       //日
hour := timeObj.Hour()     //小时
minute := timeObj.Minute() //分钟
second := timeObj.Second() //秒
fmt.Printf("格式化后 ：%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)

// ==== 执行结果 ====
current timestamp : 1632718244
current timestamp nanosecond: 1632718244062677557
时间戳转换后的时间 : 2021-09-27 12:50:44 +0800 CST
格式化后 ：2021-09-27 12:50:44
```


3.时间间隔

描述: time.Duration是time包定义的一个类型，它代表两个时间点之间经过的时间，以纳秒为单位。
time.Duration表示一段时间间隔，可表示的最长时间段大约290年。

time包中定义的时间间隔类型的常量如下：
```
const (
  Nanosecond  Duration = 1
  Microsecond          = 1000 * Nanosecond
  Millisecond          = 1000 * Microsecond
  Second               = 1000 * Millisecond
  Minute               = 60 * Second
  Hour                 = 60 * Minute
)
```

例如：time.Duration表示1纳秒，time.Second表示1秒。
```
func demo3() {
	fmt.Println(time.Nanosecond)
	fmt.Println(time.Microsecond)
	fmt.Println(time.Millisecond)
	fmt.Println(time.Second)
	fmt.Println(time.Minute)
	fmt.Println(time.Hour)
}

// 执行结果:
1ns
1µs
1ms
1s
1m0s
1h0m0s
```

4.时间时区

描述: 默认输出的时间为 UTC 世界协调时间，我们可以设置CST 中部标准时间 (Central Standard Time) , 而中国属于东八区，我们需要在上述时间+8小时，我们可以利用如下方法。

GMT、UTC、DST、CST时区代表的意义

GMT：Greenwich Mean Time (格林威治标准时间); 英国伦敦格林威治定为0°经线开始的地方，地球每15°经度 被分为一个时区，共分为24个时区，相邻时区相差一小时；例: 中国北京位于东八区，GMT时间比北京时间慢8小时。
UTC: Coordinated Universal Time (世界协调时间)；经严谨计算得到的时间，精确到秒，误差在0.9s以内， 是比GMT更为精确的世界时间
DST: Daylight Saving Time (夏季节约时间) 即夏令时；是为了利用夏天充足的光照而将时间调早一个小时，北美、欧洲的许多国家实行夏令时；
CST: Central Standard Time (中部标准时间) 四个不同时区的缩写：
```
Central Standard Time (USA) UT-6:00 美国标准时间
Central Standard Time (Australia) UT+9:30 澳大利亚标准时间
China Standard Time UT+8:00 中国标准时间
Cuba Standard Time UT-4:00 古巴标准时间
```

代码演示:
```
func demo4() {
	// UTC & CST & 本地时间 并返回与t关联的时区信息。
	now := time.Now()
	fmt.Printf("UTC 世界协调时间 : %v,时区信息: %v\n", now.UTC(), now.UTC().Location())

	var cst = time.FixedZone("CST", 0)
	cstnow := time.Now().In(cst)
	fmt.Printf("CST 中部标准时间 : %v,时区信息: %v\n", cstnow, cstnow.Location())
	
	fmt.Printf("将UTC时间转化为当地时间 : %v,时区信息: %v\n\n", now.Local(), now.Location())
	
	// 中国北京时间东八区
	// 方式1.FixedZone
	var utcZone = time.FixedZone("UTC", 8*3600)
	fmt.Printf("北京时间 : %v\n", now.In(utcZone))
	
	// 方式2.LoadLocation 设置地区
	var cstZone, _ = time.LoadLocation("Asia/Shanghai") //上海
	fmt.Printf("北京时间 : %v\n", now.In(cstZone))
	
	// 输出当前格林威治时间和该时区相对于UTC的时间偏移量（单位秒）
	name, offset := now.In(utcZone).Zone()
	fmt.Println("当前时间时区名称:", name, " 对于UTC的时间偏移量:", offset)

  // 当前操作系统本地时区
  fmt.Println("当前操作系统本地时区",time.Local)
}
```

执行结果:
```
UTC 世界协调时间 : 2021-09-27 04:58:11.866995694 +0000 UTC,时区信息: UTC
CST 中部标准时间 : 2021-09-27 04:58:11.867088566 +0000 CST,时区信息: CST
将UTC时间转化为当地时间 : 2021-09-27 12:58:11.866995694 +0800 CST,时区信息: Local

北京时间 : 2021-09-27 12:58:11.866995694 +0800 UTC
北京时间 : 2021-09-27 12:58:11.866995694 +0800 CST
当前时间时区名称: UTC  对于UTC的时间偏移量: 28800
当前操作系统本地时区: Local
```

5.时间操作

Add
描述: 我们在日常的编码过程中可能会遇到要求时间+时间间隔的需求，Go语言的时间对象有提供Add方法如下：
语法: func (t Time) Add(d Duration) Time



Sub
描述: 求两个时间之间的差值，返回一个时间段t-u。

如果结果超出了Duration可以表示的最大值/最小值，将返回最大值/最小值。要获取时间点t-d（d为Duration），可以使用t.Add(-d)。

语法: func (t Time) Sub(u Time) Duration



Equal
描述: 判断两个时间是否相同，会考虑时区的影响，因此不同时区标准的时间也可以正确比较。本方法和用t==u不同，这种方法还会比较地点和时区信息。

语法: func (t Time) Equal(u Time) bool



Before
如果t代表的时间点在u之前，返回真；否则返回假。

语法: func (t Time) Before(u Time) bool



After
如果t代表的时间点在u之后，返回真；否则返回假。

语法: func (t Time) After(u Time) bool



代码示例
```
func demo5() {
	now := time.Now()
	// 1.求一个小时之后的时间
	later := now.Add(time.Hour) // 当前时间加1小时后的时间
  tomorrow := now.Add(time.Hour * 24) // 当前时间加1天后的时间
	fmt.Println("later :", later, "\ntomorrow: ", tomorrow)
	fmt.Println("later :", later)
	// 2.当前时间与later的差值
	fmt.Println("Sub :", now.Sub(later))
	// 3.当前时间与later是否相等
	fmt.Println("Equal :", now.Equal(later))
	// 3.当前时间是否在later之前
	fmt.Println("Before :", now.Before(later))
	// 3.当前时间是否在later之后
	fmt.Println("After :", now.After(later))
}
```
执行结果:
```
later : 2021-09-27 14:04:53.94642009 +0800 CST m=+3600.000091915
tomorrow:  2021-09-28 13:04:53.94642009 +0800 CST m=+86400.000091915
Sub : -1h0m0s
Equal : false
Before : true
After : false
```



6.定时器

描述: 使用time.Tick(时间间隔)来设置定时器以及使用time.Sleep(Duration)函数来延迟执行，定时器的本质上是一个通道（channel）。

Duration 时间间隔可选参数:
```
time.Nanosecond
time.Microsecond
time.Millisecond
time.Second
time.Minute
time.Hour
```
示例演示:
```
func demo6() {
	ticker := time.Tick(time.Second) //定义一个1秒间隔的定时器
	for i := range ticker {
		fmt.Println(i)              //每秒都会执行的任务
		time.Sleep(time.Second * 5) //休眠5S执行，注意不能直接传递5，除了前面这种方式你还可以利用显示强转整数5 time.Duration(5);
	}
}

执行结果:
➜  Time go run timeDemo.go
2021-09-27 03:35:58.640028069 +0000 UTC m=+1.000158842
2021-09-27 03:35:59.64011495 +0000 UTC m=+2.000245738  # 特殊点(第二执行未经过5S)
2021-09-27 03:36:04.640065081 +0000 UTC m=+7.000195859
2021-09-27 03:36:09.640302389 +0000 UTC m=+12.000433177
2021-09-27 03:36:14.640070051 +0000 UTC m=+17.000200829
```



7. 时间格式化

描述: 时间类型有一个自带的方法Format进行格式化，需要注意的是Go语言中格式化时间模板不是常见的Y-m-d H:M:S而是使用Go的诞生时间2006年1月2号15点04分（记忆口诀为2006 1 2 3 4）, 也许这就是技术人员的浪漫吧。

补充：如果想格式化为12小时方式，需指定PM。

Foramt|格式化
描述: 格式化时间是把Go语言中的时间对象，转换成为字符串类型的时间。

代码演示:
```
func demo7() {
	// 当前UTC时间
	now := time.Now()
	// 设置时区为Asia/Shanghai
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("北京时间 :", now.In(loc), "\n地点时区 :", loc)

	// 1.时间格式化
	fmt.Println("格式1 :", now.Format("2006/01/02"))
	fmt.Println("格式2 :", now.Format("2006/01/02 15:04"))
	fmt.Println("格式3 :", now.Format("15:04 2006/01/02"))
	
	// 24小时制
	fmt.Println("格式4 :", now.Format("2006-01-02 15:04:05.000 Mon Jan"))
	
	// 12小时制
	fmt.Println("格式5 :", now.Format("2006-01-02 03:04:05.000 PM"))
	fmt.Println("格式6 :", now.Format("2006-01-02 03:04:05.000 PM Mon Jan"))
	fmt.Println("Kitchen 格式 :", now.Format(time.Kitchen))
	
	// 时区展示
	fmt.Println("RFC1123 格式 :", now.Format(time.RFC1123))
	fmt.Println("RFC1123 格式 :", now.Format(time.RFC1123Z))
	fmt.Println("RFC3339 格式 :", now.Format(time.RFC3339))
	fmt.Println("RFC3339Nano 格式 :", now.Format(time.RFC3339Nano))
}
```

执行结果:
```
北京时间 : 2021-09-27 14:15:42.716506733 +0800 CST
地点时区 : Asia/Shanghai
格式1 : 2021/09/27
格式2 : 2021/09/27 14:15
格式3 : 14:15 2021/09/27
格式4 : 2021-09-27 14:15:42.716 Mon Sep
格式5 : 2021-09-27 02:15:42.716 PM
格式6 : 2021-09-27 02:15:42.716 PM Mon Sep
Kitchen 格式 : 2:15PM
RFC1123 格式 : Mon, 27 Sep 2021 14:15:42 CST
RFC1123 格式 : Mon, 27 Sep 2021 14:15:42 +0800
RFC3339 格式 : 2021-09-27T14:15:42+08:00
RFC3339Nano 格式 : 2021-09-27T14:15:42.716506733+08:00
```

Parse|解析字符串格式
描述: 将时间字符串解析为时间对象。

通过time.Parse将时间字符串转化为时间类型对象默认是UTC时间, 而通过time.ParseInLocation我们可以指定时区得到CST时间。



代码演示:
```
func demo8() {
	// 1.时间与时区设置
	now := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	// 2.按照指定时区和指定格式解析字符串时间
	timeObj1, _ := time.Parse("2006-01-02 15:04:05", "2021-09-27 14:15:20")
	timeObj2, _ := time.ParseInLocation("2006/01/02 15:04:05", "2021/09/27 14:15:20", time.Local) // 操作系统本地时区
	timeObj3, _ := time.ParseInLocation("2006/01/02 15:04:05", "2021/09/27 14:15:20", loc)        // 指定时区

	fmt.Printf("Now: %v\ntimeObj1: %v\ntimeObj2: %v\ntimeObj3: %v\n", now.Local(), timeObj1, timeObj2, timeObj3)
	
	// 将当地时区转化为UTC时间
	utcLocal := timeObj3.UTC()
	fmt.Println("将当地时区转化为UTC时间:", utcLocal)
	// 将UTC时间转化为当地时间(+8)
	localTime := utcLocal.Local()
	fmt.Println("将UTC时间转化为当地时间:", localTime)
	
	// 3.相互转换后的时间进行对比.
	fmt.Println("相互转换后的时间进行对比:", utcLocal.Equal(localTime))
	
	// 4.输入的时间字符串与当前时间的相差时间.
	d := timeObj3.Sub(now)
	// 可以看到timeObj 时间 与 当前时间 相差 33 分钟 55 秒
	fmt.Println("看到timeObj 时间 与 当前时间 相差:", d.String())
}
```

执行结果：
```
Now: 2021-09-27 14:49:15.392828987 +0800 CST
timeObj1: 2021-09-27 14:15:20 +0000 UTC
timeObj2: 2021-09-27 14:15:20 +0800 CST
timeObj3: 2021-09-27 14:15:20 +0800 CST
将当地时区转化为UTC时间: 2021-09-27 06:15:20 +0000 UTC
将UTC时间转化为当地时间: 2021-09-27 14:15:20 +0800 CST
相互转换后的时间进行对比: true
看到timeObj3 时间 与 当前时间 相差: -33m55.392828987s
```



7. 时间处理常用

(0) 按照str格式化时间(Go诞生之日口诀:6-1-2-3-4-5)
```
//格式化时间格式
fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
//调用结果: 2021-04-30 13:15:02
```
(1) UTC时间互换标准时间
```
//UTC时间转标准时间
func (this *DataSearch) UTCTransLocal(utcTime string) string {
	t, _ := time.Parse("2006-01-02T15:04:05.000+08:00", utcTime)
	return t.Local().Format("2006-01-02 15:04:05")
}

t1 := UTCTransLocal("2021-04-29T14:11:08.000+08:00")
fmt.Println(t1)

// 调用结果: 2021-04-29 22:11:08
```
(2) 标准时间转UTC时间
```
//标准时间转UTC时间
func (this *DataSearch) LocalTransUTC(localTime string) string {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", localTime, time.Local)
	return t.UTC().Format("2006-01-02T15:04:05.000+08:00")
}
t2 := LocalTransUTC("2021-04-29 22:11:08")
fmt.Println(t2)

//调用结果：  2021-04-29T14:11:08.000+08:00
```
(3) str格式化时间转时间戳
```
the_time, err := time.Parse("2006-01-02 15:04:05", "2020-04-29 22:11:08")
if err == nil {
    unix_time := the_time.Unix()
	fmt.Println(unix_time)
}
fmt.Println(the_time)
//调用结果： 1588198268
```
(4) 时间戳转str格式化时间
```
str_time := time.Unix(1588224111, 0).Format("2006-01-02 15:04:05")
fmt.Println(str_time)
//调用结果：2020-04-30 13:21:51 
```