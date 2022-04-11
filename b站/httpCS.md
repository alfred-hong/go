HTTP C/S实现

> 描述: 在Socket网络编程中我们学习TCP协议与UDP协议的服务端与客户端代码的编写实践, 今天我们来学习在应用层中我们使用最多的HTTP协议，利用Go语言分别实现HTTP的Client端与Server端服务。

Go语言内置的net/http包十分的优秀，其为我们提供了HTTP客户端和服务端的实现。

> Q: 什么是 HTTP 协议?
>
> 答: 超文本传输协议（HTTP，HyperText Transfer Protocol)是互联网上应用最为广泛的一种网络传输协议，所有的WWW文件都必须遵守这个标准，设计HTTP最初的目的是为了提供一种发布和接收HTML页面的方法。
> 例如，我们在浏览器中访问的http的网站，其传输协议便是采用HTTP。



HTTP 服务端

> 描述: 利用Go语言提供的net/http包我们可以非常便利的使用并创建一个服务端, 值得说明的是如果仅仅是实现简单的API接口可以采用原生的http包中提供的方法, 而如果是编写一些Web后端项目通常是采用框架来实现，所以本章节主要对Go语言创建HTTP服务端的基础示例进行说明学习，而Go语言的Web应用开发框框在我后续笔记中将会进行实践讲解。

HTTP服务端实现常用方法原型:

​	func http.ListenAndServe(addr string, handler http.Handler) error: 使用指定的监听地址和处理器启动一个HTTP服务端, 处理器参数通常是nil表示采用包变量DefaultServeMux作为处理器。
​	func http.Handle(pattern string, handler http.Handler){ DefaultServeMux.Handle(pattern, handler) }: Handle在DefaultServeMux中注册给定模式的处理程序函数。
​	func http.HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)){DefaultServeMux.HandleFunc(pattern, handler)}: HandleFunc在DefaultServeMux中注册给定模式的处理程序函数。

```go
// 1.Handle在DefaultServeMux中注册给定模式的处理程序函数
type httpServer struct {}
func (server httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.URL.Path))
}
var fooHandler httpServer
http.Handle("/foo", fooHandler)

// 2.HandleFunc在DefaultServeMux中注册给定模式的处理程序函数
http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
})

// 3.指定的本地网卡的8080断开作为监听地址和处理器启动一个HTTP服务端。
log.Fatal(http.ListenAndServe(":8080", nil))
```


Tips: go http http.Handle 和 http.HandleFunc 区别?

​	http.Handle() 需要自己去定义struct实现这个Handler接口。
​	http.HandleFunc() 则不需要我们自己定义structqi实现,只需要传入访问连接路径以及DefaultServeMux中注册给定模式的处理程序函数。
通常是使用HandleFunc方法,其更加方便简单。

示例1.http.Handle自定义实现Handler接口

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

// 自定义结构体
type httpServer struct{}

// httpServer 自定义实现http请求处理程序函数
func (server httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	switch path {
	case "/":
		w.Write([]byte("根路径 : " + r.URL.Path))
	case "/index":
		w.Write([]byte("首页路径 : " + r.URL.Path))
	case "/hello":
		w.Write([]byte("子网页路径 : " + r.URL.Path))
	default:
		w.Write([]byte("<b>未知路径</b> : https://weiyigeek.top" + r.URL.Path))
	}
}

func main() {
	// 1.声明serve变量的类型为我们自定义结构体
	var server httpServer
	serveraddr := "0.0.0.0:8080"

	// 2.Handle在DefaultServeMux中注册给定模式的处理程序。
	http.Handle("/", server)
	
	// 3.启动httpServer监听并采用包变量DefaultServeMux作为处理器。
	err := http.ListenAndServe(serveraddr, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Http Server %v Started......", serveraddr)
	}

}
```


执行结果如下图所示

![WeiyiGeek.http.Handle结果](https://i0.hdslb.com/bfs/article/0807e700f22827c25e312098ba7fe79b90ceab75.png@615w_417h_progressive.png)


示例2.http.HandleFunc实现用户请求处理

index.html

<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>测试页面2</title>
</head>
<body>
  <p  id="tips" style="color: red;font-size:medium;font-weight: bolder;">Welcome to Visit weiyigeek.top web site</p>
  <button id="msg">点击提示</button>
  <script>
    var tips = document.getElementById("tips").textContent;
    document.getElementById("msg").onclick=function() {
      alert(tips);
    }
 </script>
</body>
</html>


Http Server 端代码:

```go
package main
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)
func main() {
	// 1.声明定义server监听端口地址
	serveraddr := "0.0.0.0:8080"

	// 2.HandleFunc在DefaultServeMux中注册给定模式的处理程序。
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/blog", blog)
	http.HandleFunc("/htmlfile", htmlfile)
	
	// 3.启动httpServer监听并采用包变量DefaultServeMux作为处理器。
	err := http.ListenAndServe(serveraddr, nil)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Http Server %v Started......", serveraddr)
	}

}

// (1) 方式1.此种方式不能直接写入HTML标签代码并返还给客户端.
func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! Go <b>测试页面0</b> Time : %s", time.Now())
}

// (2) 方式2.可以直接写入HTML标签代码并返还给客户端.
func blog(w http.ResponseWriter, r *http.Request) {
	reply := fmt.Sprintf("<b>测试页面1</b><p> 标题.Demo1 Test(Go net/http) </p> <i>I'm WeiyiGeek</i><br/> Time : %s", time.Now().Format("2006-01-02 15:04:05"))
	w.Write([]byte(reply))
}

// (3) 读取本机上的网页为文件返还给客户端.
func htmlfile(w http.ResponseWriter, r *http.Request) {
	res, err := ioutil.ReadFile("./index.html")
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	w.Write([]byte(res))
}
```

将上面的代码编译之后执行，打开电脑上的浏览器在地址栏输入127.0.0.1:8080回车，此时返回如下页面。

![WeiyiGeek.http.HandleFunc实现用户请求处理](https://i0.hdslb.com/bfs/article/022afa36a5ced978a1c6afb75a0dccc3bd859986.png@926w_902h_progressive.png)

示例3.自定义http.Server结构体参数与Handler实现

> 描述: 我们可以创建一个自定义的 Server Handler实现, 通过http.Server指定设置结构体参数来管理服务端的行为。

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// 自定义myHandler结构体，注册Handler处理程序使用
type myHandler struct {
	name string
}
// 自定义myHandler结构体的方法注册到Handler处理程序,【非常注意】、【非常注意】方法名必须为 ServeHTTP
func (handler myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 现在handler基于URL的路径部分（req.URL.Path）来决定执行什么逻辑。
	switch req.URL.Path {
	case "/index":
		fmt.Fprintf(w, "%s\n", "This is Index path")
	case "/weiyigeek":
		fmt.Fprintf(w, "%s -> %s\n", handler.name, "https://weiyigeek.top") // WeiyiGeek -> https://weiyigeek.top
	default:
		// 如果这个handler不能识别这个路径，它会通过调用返回客户端一个HTTP404错误,并响应给客户端表明请求的路径不存在.
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such page: %s\n", req.URL)
	}
}

func main() {
	// 1.实例化一个处理所有请求的Handler接口
	handler := myHandler{name: "WeiyiGeek"}

	// 2.创建一个自定义的Server参数,注意如果Handler为nil则采用http.DefaultServeMux进行处理响应,否则需要我们自己实现结构体的ServeHTTP方法.
	server := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	// 3.启动并监听HTTP服务端
	log.Fatal(server.ListenAndServe())

}
```

执行结果:

![WeiyiGeek.自定义http.Server结构体参数与Handler实现](https://i0.hdslb.com/bfs/article/b6b9b75536674df36ce54990dcbb48663d50ddfa.png@942w_506h_progressive.png)


Tips: 显然我们可以继续向ServeHTTP方法中添加case，但在一个实际的应用中，将每个case中的逻辑定义到一个分开的方法或函数中会很实用。对于更复杂的应用我们可以通过一个ServeMux将一批http.Handler聚集到一个单一的http.Handler中,通过组合来处理更加错综复杂的路由需求。

Tips: **【非常注意】【非常注意】【非常注意】**自己实现的 http.Handler 且必须包含一个 ServeHTTP 方法名, 才能接受和响应客户端。

示例4.ServeMux.HandleFunc实现http服务端
语句http.HandlerFunc(handler.list)是一个转换而非一个函数调用，因为http.HandlerFunc是一个类型, 它有如下的定义：

```go
package http
type HandlerFunc func(w ResponseWriter, r *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
    f(w, r)
}
```


简单示例:

```go
// 自定义类型声明
type dollars float32
func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }
// Map类型的变量声明
type MyHandler map[string]dollars
// 类型方法
func (self MyHandler) list(w http.ResponseWriter, req *http.Request) {
  for item, price := range self {
      fmt.Fprintf(w, "%s: %s\n", item, price)
  }
}
func (self MyHandler) price(w http.ResponseWriter, req *http.Request) {
  item := req.URL.Query().Get("item")
  price, ok := self[item]
  if !ok {
      w.WriteHeader(http.StatusNotFound) // 404
      fmt.Fprintf(w, "no such item: %q\n", item)
      return
  }
  fmt.Fprintf(w, "%s\n", price)
}
func main() {
    handler := MyHandler{"shoes": 50, "socks": 5}
    http.HandleFunc("/list", handler.list)
    http.HandleFunc("/price", handler.price)
    log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
```

Tips: 为了方便net/http包提供了一个全局的ServeMux实例DefaultServerMux和包级别的http.Handle和http.HandleFunc函数, 前面说过我们只需要传入nil值即可。

HTTP 客户端

> 描述: 如果你学过Python爬虫那么你也肯定知道，我们可以采用编程语言提供的相应函数或者方法就要可以对指定的网站进行请求，并可以将请求响应的数据进行清洗然后存储进数据中。

Go语言也为我们提供相应的内置包net/http和net/url以便于我们进行网站API接口(GET、POST、PUT)请求和处理服务端响应的数据。

HTTP客户端请求函数原型:

​	func http.Get(url string) (resp \*http.Response, err error) : 向指定的URL发出Get请求,它将随重定向，最多10个重定向。
​	func http.Post(url string, contentType string, body io.Reader) (resp \*http.Response, err -error) : 指定的URL发出Post请求，调用方在完成读取后应关闭相应的主体。
​	func http.PostForm(url string, data url.Values) (resp \*http.Response, err error) : PostForm向指定的URL发出POST，并将数据的键和值URL编码为请求正文，Content-Type header设置为application/x-www-form-urlencoded。
​	func http.NewRequest(method string, url string, body io.Reader) (\*http.Request, error) : NewRequest使用后台上下文包装NewRequestWithContext.
​	func (\*http.Client).Do(req *http.Request) (*http.Response, error) : 按照客户端上配置的策略（如重定向、cookie、身份验证），发送到HTTP请求并返回到HTTP响应。
​	func url.Parse(rawurl string) (\*url.URL, error): Parse将rawurl解析为URL结构。
​	func (url.Values).Encode() string : Encode将值编码为按键排序的“URL编码”形式（“bar=baz&foo=qux”）。

```go
# Get 请求

resp, err := http.Get("http://example.com/")
...

# Post 请求

resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
...

# Form 表单请求

resp, err := http.PostForm("http://example.com/form",
	url.Values{"key": {"Value"}, "id": {"123"}})
...
if err != nil {
	// handle error
}
...

# 程序在使用完response后必须关闭回复的主体。

defer resp.Body.Close()

# 读取响应的数据

body, err := ioutil.ReadAll(resp.Body)
```


Tips : GET请求的参数需要使用Go语言内置的net/url标准库来处理。



简单示例.

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getMethod(url string) {
	// 1.Get请求指定地址
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	// 2.程序完毕时关闭回复的主体.
	defer resp.Body.Close()

	// 3.读取响应的网页源代码
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	html := string(body)
	
	// 4.网站请求响应数据
	fmt.Printf("resp.StatusCode : %v,\nresp.Status: %v,\nresp.Request: %#v,\nresp.Header: %#v\nresp.Cookies: %#v,\nresp.TLS: %#v\n",
		resp.StatusCode,
		resp.Status,
		resp.Request,
		resp.Header,
		resp.Cookies(),
		resp.TLS)
	
	fmt.Println("网站响应长度: ", len(html))

}

func main() {
	getMethod("https://www.weiyigeek.top")
}
```


执行结果:

```sh
dresp.StatusCode : 200,

resp.Status: 200 OK,

resp.Request: &http.Request{Method:"GET", URL:(*url.URL)(0xc000176750), Proto:"", ProtoMajor:0, ProtoMinor:0, Header:http.Header{"Referer":[]string{"https://www.weiyigeek.top"}}, Body:io.ReadCloser(nil), GetBody:(func() (io.ReadCloser, error))(nil), ContentLength:0, TransferEncoding:[]string(nil), Close:false, Host:"", Form:url.Values(nil), PostForm:url.Values(nil), MultipartForm:(*multipart.Form)(nil), Trailer:http.Header(nil), RemoteAddr:"", RequestURI:"", TLS:(*tls.ConnectionState)(nil), Cancel:(<-chan struct {})(nil), Response:(*http.Response)(0xc000176630), ctx:(*context.emptyCtx)(0xc000134010)},

resp.Header: http.Header{"Access-Control-Allow-Origin":[]string{"*"}, "Age":[]string{"340"}, "Alt-Svc":[]string{"h3=\":443\"; ma=86400, h3-29=\":443\"; ma=86400, h3-28=\":443\"; ma=86400, h3-27=\":443\"; ma=86400"}, "Cache-Control":[]string{"max-age=600"}, "Cf-Cache-Status":[]string{"DYNAMIC"}, "Cf-Ray":[]string{"6b222396a8d10d24-LAX"}, "Content-Type":[]string{"text/html; charset=utf-8"}, "Date":[]string{"Mon, 22 Nov 2021 12:25:11 GMT"}, "Expect-Ct":[]string{"max-age=604800, report-uri=\"https://report-uri.cloudflare.com/cdn-cgi/beacon/expect-ct\""}, "Expires":[]string{"Mon, 22 Nov 2021 12:27:16 GMT"}, "Last-Modified":[]string{"Fri, 08 Oct 2021 11:30:52 GMT"}, "Nel":[]string{"{\"success_fraction\":0,\"report_to\":\"cf-nel\",\"max_age\":604800}"}, "Report-To":[]string{"{\"endpoints\":[{\"url\":\"https:\\/\\/a.nel.cloudflare.com\\/report\\/v3?s=6dx48uJJynNOFy1TpNXdJU0V%2FfBjazc15i3SLtiT4xoXr8vXl0MZTRuVE11vx33KNqsU05DfmNKKuPrzDTBstP32mSG%2B%2FNWUtG%2B4vlX5Ro3hApIow24bII5D0uT8FtSh\"}],\"group\":\"cf-nel\",\"max_age\":604800}"}, "Server":[]string{"cloudflare"}, "Vary":[]string{"Accept-Encoding"}, "Via":[]string{"1.1 varnish"}, "X-Cache":[]string{"HIT"}, "X-Cache-Hits":[]string{"1"}, "X-Fastly-Request-Id":[]string{"de2ef3b07dec51677274e7fe8eb68162dc77c39d"}, "X-Github-Request-Id":[]string{"C002:9402:1A6C52:23CF10:619B8A4C"}, "X-Proxy-Cache":[]string{"MISS"}, "X-Served-By":[]string{"cache-bur17563-BUR"}, "X-Timer":[]string{"S1637583911.481935,VS0,VE1"}}
resp.Cookies: []*http.Cookie{},

resp.TLS: &tls.ConnectionState{Version:0x304, HandshakeComplete:true, DidResume:false, CipherSuite:0x1301, NegotiatedProtocol:"h2", NegotiatedProtocolIsMutual:true, ServerName:"weiyigeek.top", PeerCertificates:[]*x509.Certificate{(*x509.Certificate)(0xc000324000), (*x509.Certificate)(0xc000324580)}, VerifiedChains:[][]*x509.Certificate{[]*x509.Certificate{(*x509.Certificate)(0xc000324000), (*x509.Certificate)(0xc000324580), (*x509.Certificate)(0xc000267180)}}, SignedCertificateTimestamps:[][]uint8(nil), OCSPResponse:[]uint8{0x30, 0x82, 0x1, 0x12, 0xa, 0x1, 0x0, 0xa0, 0x82, 0x1, 0xb, 0x30, 0x82, 0x1, 0x7, 0x6, 0x9, 0x2b, 0x6, 0x1, 0x5, 0x5, 0x7, 0x30, 0x1, 0x1, 0x4, 0x81, 0xf9, 0x30, 0x81, 0xf6, 0x30, 0x81, 0x9e, 0xa2, 0x16, 0x4, 0x14, 0xa5, 0xce, 0x37, 0xea, 0xeb, 0xb0, 0x75, 0xe, 0x94, 0x67, 0x88, 0xb4, 0x45, 0xfa, 0xd9, 0x24, 0x10, 0x87, 0x96, 0x1f, 0x18, 0xf, 0x32, 0x30, 0x32, 0x31, 0x31, 0x31, 0x31, 0x38, 0x32, 0x30, 0x34, 0x32, 0x33, 0x38, 0x5a, 0x30, 0x73, 0x30, 0x71, 0x30, 0x49, 0x30, 0x9, 0x6, 0x5, 0x2b, 0xe, 0x3, 0x2, 0x1a, 0x5, 0x0, 0x4, 0x14, 0x12, 0xd7, 0x8b, 0x40, 0x2c, 0x35, 0x62, 0x6, 0xfa, 0x82, 0x7f, 0x8e, 0xd8, 0x92, 0x24, 0x11, 0xb4, 0xac, 0xf5, 0x4, 0x4, 0x14, 0xa5, 0xce, 0x37, 0xea, 0xeb, 0xb0, 0x75, 0xe, 0x94, 0x67, 0x88, 0xb4, 0x45, 0xfa, 0xd9, 0x24, 0x10, 0x87, 0x96, 0x1f, 0x2, 0x10, 0x9, 0x98, 0xa5, 0x9a, 0x26, 0x72, 0xc7, 0x24, 0x4a, 0x4d, 0xc5, 0x92, 0x80, 0xfb, 0x65, 0x5a, 0x80, 0x0, 0x18, 0xf, 0x32, 0x30, 0x32, 0x31, 0x31, 0x31, 0x31, 0x38, 0x32, 0x30, 0x32, 0x37, 0x30, 0x32, 0x5a, 0xa0, 0x11, 0x18, 0xf, 0x32, 0x30, 0x32, 0x31, 0x31, 0x31, 0x32, 0x35, 0x31, 0x39, 0x34, 0x32, 0x30, 0x32, 0x5a, 0x30, 0xa, 0x6, 0x8, 0x2a, 0x86, 0x48, 0xce, 0x3d, 0x4, 0x3, 0x2, 0x3, 0x47, 0x0, 0x30, 0x44, 0x2, 0x20, 0x5f, 0x88, 0xf2, 0xc1, 0x99, 0xcf, 0x99, 0x2b, 0x57, 0xd6, 0xd0, 0x38, 0x2e, 0x7, 0x72, 0xc7, 0x7d, 0x48, 0x34, 0x57, 0x60, 0x19, 0xe, 0x42, 0xd1, 0x32, 0x6b, 0xea, 0x5f, 0xbd, 0xfa, 0x36, 0x2, 0x20, 0x67, 0xc7, 0xc1, 0x3, 0xd4, 0xed, 0x1e, 0x32, 0xb3, 0x5f, 0x7e, 0xb3, 0xc8, 0x10, 0xb4, 0xdf, 0x88, 0x47, 0x1c, 0xf3, 0xee, 0xab, 0x3b, 0x86, 0xc7, 0xe4, 0xbc, 0xcf, 0x5c, 0x1d, 0x69, 0x48}, TLSUnique:[]uint8(nil), ekm:(func(string, []uint8, int) ([]uint8, error))(0x679b80)}

网站响应长度:  16853
```

自定义 Client&Transport

> 描述: 要管理HTTP客户端的头域、重定向策略和其他设置，创建一个Client：

```go
// 重定向策略
client := &http.Client{
  // CheckRedirect指定处理重定向的策略。如果CheckRedirect不是nil，则客户端在执行HTTP重定向之前调用它。
	CheckRedirect: redirectPolicyFunc,
  // 超时指定此客户端发出的请求的时间限制。超时包括连接时间、任何重定向和读取响应正文。Get、Head、Post或Do返回后，计时器将保持运行，并将中断Response.Body的读取。
  Timeout： 30 * time.Second,
}
resp, err := client.Get("http://example.com")

// 请求对象设置与&自定义请求头
req, err := http.NewRequest("GET", "http://example.com", nil)
req.Header.Add("If-None-Match", `W/"wyzzy"`)
resp, err := client.Do(req)
// ...
```

> 描述: 要管理代理、TLS配置、keep-alive、压缩和其他设置，创建一个Transport：

```go
tr := &http.Transport{
  // TLS 配置
	TLSClientConfig:    &tls.Config{RootCAs: pool},
  // 是否禁用压缩
	DisableCompression: true,
  // 是否保持连接(长连接、短连接)
  DisableKeepAlives:  true,
}
client := &http.Client{Transport: tr}
resp, err := client.Get("https://example.com")
```

Tips: Client和Transport类型都可以安全的被多个goroutine同时使用, 出于效率考虑，应该一次建立、尽量重用。

Tips: 如果取数据比较频繁的场景建议使用长连接,否则使用短连接即可。



综合实践

3.1 Get 请求示例

示例代码文件 get_client.go

```go
package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)
func getMethod(urlstr string) {
	// 1.URL格式校验解析
	urlParse, err := url.Parse(urlstr)
	if err != nil {
		fmt.Printf("Url %v Format Error!\nerr: %v\n", urlParse, err)
		return
	}
	// 2.URL参数设置与编码
	data := url.Values{}
	data.Set("id", "1024")
	data.Set("name", "唯一极客")
	// 处理URL中包含的中文参数,此处采用encode进行编码.
	queryStr := data.Encode()
	// URL参数设置并输出处理过后的请求字符串
	urlParse.RawQuery = queryStr // encoded query values, without '?'
	fmt.Println("QueryStr => ", queryStr)

	// 3.NewRequest使用后台上下文包装NewRequestWithContext,返回请求对象
	req, err := http.NewRequest("Get", urlParse.String(), nil)
	if err != nil {
		fmt.Printf("NewRequest %v faile!\n[error]: %v\n", urlstr, err)
		return
	}
	
	// 5.DefaultClient是默认客户端，由Get、Head和Post请求使用,此时传入上面处过的req请求对象
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Request %v faile!\n[error]: %v\n", urlstr, err)
		return
	}
	
	// 6.程序完毕时关闭回复的主体(非常重要).
	defer resp.Body.Close()
	
	// 7.从resp中把服务端返回的数据读出来
	// 方式1
	// var data []byte
	// response.Body.Read(data)
	// 方式2
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	
	// 8.网站请求响应数据集合
	fmt.Printf("resp.StatusCode : %v,\nresp.Status: %v,\nresp.Request: %#v,\nresp.Header: %#v\nresp.Cookies: %#v,\nresp.TLS: %#v\n",
		resp.StatusCode,
		resp.Status,
		resp.Request,
		resp.Header,
		resp.Cookies(),
		resp.TLS)
	fmt.Println("网站响应长度: ", len(body))

}


func main() {
	getMethod("http://10.20.172.108:8080/get")
}
```


执行结果:

```sh
QueryStr =>  id=1024&name=%E5%94%AF%E4%B8%80%E6%9E%81%E5%AE%A2 //url 编码后
resp.StatusCode : 200,
resp.Status: 200 OK,
resp.Request: &http.Request{Method:"Get", URL:(*url.URL)(0xc0000fe090), Proto:"HTTP/1.1", ProtoMajor:1, ProtoMinor:1, Header:http.Header{}, Body:io.ReadCloser(nil), GetBody:(func() (io.ReadCloser, error))(nil), ContentLength:0, TransferEncoding:[]string(nil), Close:false, Host:"10.20.172.108:8080", Form:url.Values(nil), PostForm:url.Values(nil), MultipartForm:(*multipart.Form)(nil), Trailer:http.Header(nil), RemoteAddr:"", RequestURI:"", TLS:(*tls.ConnectionState)(nil), Cancel:(<-chan struct {})(nil), Response:(*http.Response)(nil), ctx:(*context.emptyCtx)(0xc0000b0010)},
resp.Header: http.Header{"Content-Length":[]string{"68"}, "Content-Type":[]string{"application/json;charset=UTF-8"}, "Cookies":[]string{"id=1024;name=唯一极客"}, "Date":[]string{"Tue, 23 Nov 2021 05:25:41 GMT"}, "Requestmethod":[]string{"Get"}}
resp.Cookies: []*http.Cookie{},
resp.TLS: (*tls.ConnectionState)(nil)
网站响应长度:  68
```


3.2 Post 请求示例

示例代码文件 post_client.go

```sh
// POST 示例
func postMethos(urlstr string) {
	// (1) 定义Post请求上传的参数( 表单数据/json数据)
	//contentType := "application/x-www-form-urlencoded"
	//data := "name=小王子&age=18"
	contentType := "application/json"
	data := `{"id":128,"name":"Weiyi"}`

	// (2) 进行 Post 请求传入请求url,contentType以及Post上传的data数据
	// resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
	resp, err := http.Post(urlstr, contentType, strings.NewReader(data))
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}
	
	// (3) 程序结束则关闭resp资源
	defer resp.Body.Close()
	
	// (4) 读取POST请求返回的数据包
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Println("Body : ", string(body))
	
	// (5) 响应头输出
	fmt.Println("resp.Header : ", resp.Header)

}

func main() {
	postMethos("http://10.20.172.108:8080/post")
}
```


执行结果:

```sh
Body :  {"method":"POST","status":"ok","data":{"id":128,"name":"Weiyi"}}
resp.Header :  map[Content-Length:[64] Content-Type:[application/json;charset=UTF-8] Cookies:[id=128;name=Weiyi] Date:[Tue, 23 Nov 2021 05:25:41 GMT] Requestmethod:[Post]]
```

3.3 PostForm 请求示例

示例代码文件 postForm_client.go

```go
func postFormMethod(urlstr string) {
	// (1) 方式1定义Post请求上传的参数( 表单数据)
	//contentType := "application/x-www-form-urlencoded"
	//data := "id=256&name=唯一极客"

	// (2) 方式2定义Post请求上传的参数( 表单数据)
	data := url.Values{}
	data.Set("id", "256")
	data.Set("name", "WeiyiGeek-唯一极客")
	
	// (3) 进行表单上传请求
	// resp, err := http.PostForm("http://example.com/form", url.Values{"key": {"Value"}, "id": {"123"}})
	resp, err := http.PostForm(urlstr, data)
	if err != nil {
		fmt.Printf("postForm failed, err:%v\n", err)
		return
	}
	
	// (4) 程序结束则关闭resp资源
	defer resp.Body.Close()
	
	// (5) 读取POST请求返回的数据包
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("get resp failed, err:%v\n", err)
		return
	}
	fmt.Println("Body : ", string(body))
	
	// (6) 响应头输出
	fmt.Println("resp.Header : ", resp.Header)

}

func main() {
	postFormMethod("http://10.20.172.108:8080/postform")
}
```


执行结果:

```sh
Body :  {"method":"POSTFORM","status":"ok","data":{"id":256,"name":"WeiyiGeek"}}
resp.Header :  map[Content-Length:[72] Content-Type:[application/json;charset=UTF-8] Cookies:[method=form;id=256;name=WeiyiGeek-唯一极客] Date:[Tue, 23 Nov 2021 05:25:41 GMT] Requestmethod:[PostForm]]
```

3.4 Http 服务响应示例

描述: 下述http_serve.go代码包含上面三种请求示例的响应，以及请求参数获取以及输出。

服务端示例文件 http_serve.go

```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// 4.自定义Handler结构体
type myHandler struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// 5.自定义myHandler结构体的方法注册到Handler处理程序,【非常注意】、【非常注意】方法名必须为 ServeHTTP
func (handler myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 6.现在handler基于URL的路径部分（req.URL.Path）来决定执行什么逻辑。
	switch req.URL.Path {
	case "/get":
		// 7.对于客户端Get请求参数获取
		fmt.Printf("Method : %v, URL : %v \n", req.Method, req.URL)

		// 8.自动识别请求URL中的参数,参数利用Map类型变量进行存储key-value
		fmt.Println("QueryParam : ", req.URL.Query())
		queryParam := req.URL.Query()
		id := queryParam.Get("id")
		name := queryParam.Get("name")
	
		// 9.打印输出queryParam存储的value
		fmt.Printf("id = %v,name = %v\n", id, name)
	
		// 10.服务端响应头header自定义
		w.Header().Add("RequestMethod", "Get")                            // 此处将响应的RequestMethod header字段设置为Get
		w.Header().Add("Content-Type", "application/json;charset=UTF-8")  // 此处将响应的类型设置为JSON
		w.Header().Add("Cookies", fmt.Sprintf("id=%v;name=%v", id, name)) // 此处将响应的cookies设置为请求传入的参数
	
		// 11.返回给客户端的JSON数据组装
		uid, err := strconv.Atoi(id) // 将get到的id字段的值转换为整型
		if err != nil {
			errMsg := fmt.Sprintf("uid convert err! %v\n", err)
			fmt.Println(errMsg)
			w.Write([]byte(errMsg))
			return
		}
		reply := fmt.Sprintf("{\"method\":\"%v\",\"status\":\"ok\",\"data\":{\"id\":%v,\"name\":\"%v\"}}", "GET", uid, handler.Name)
		fmt.Printf("reply => %v\n\n", reply)
	
		// 12.返回响应数据给客户端
		w.Write([]byte(reply))
	
	case "/post":
		// 7.对于客户端Request请求信息获取
		fmt.Printf("Method : %v, URL : %v \n", req.Method, req.URL)
	
		// 8.服务端打印客户端发来的请求,当请求类型是application/json时才能从req.Body读取数据
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println("req Body Read Error,", err)
			return
		}
		fmt.Println("req.Body => ", string(body))
	
		//9.JSON反序列化
		json.Unmarshal(body, &handler)
	
		// 10.打印输出queryParam存储的value
		fmt.Printf("id = %v,name = %v\n", handler.Id, handler.Name)
	
		// 11.服务端响应头header自定义
		w.Header().Add("RequestMethod", "Post")                                           // 此处将响应的RequestMethod header字段设置为POST
		w.Header().Add("Content-Type", "application/json;charset=UTF-8")                  // 此处将响应的类型设置为JSON
		w.Header().Add("Cookies", fmt.Sprintf("id=%v;name=%v", handler.Id, handler.Name)) // 此处将响应的cookies设置为请求传入的参数
	
		// 12.返回给客户端的JSON数据组装
		reply := fmt.Sprintf("{\"method\":\"%v\",\"status\":\"ok\",\"data\":{\"id\":%v,\"name\":\"%v\"}}", "POST", handler.Id, handler.Name)
		fmt.Printf("reply =>  %v\n\n", reply)
	
		// 13.返回响应数据给客户端
		w.Write([]byte(reply))
	
	case "/postform":
	
		// 7.对于客户端Request请求信息获取
		fmt.Printf("Method : %v, URL : %v \n", req.Method, req.URL)
	
		// 8. 请求类型是application/x-www-form-urlencoded时解析form数据并打印
		req.ParseForm()
		fmt.Println(req.PostForm)
	
		// 9.获取postform表单中指定的字段值.
		id := req.PostForm.Get("id")
		name := req.PostForm.Get("name")
		fmt.Printf("id = %v, name = %v\n", id, name)
	
		// 10.服务端打印客户端发来的请求Body此处为[],因为但客户端请求类型是application/json时才能从req.Body读取数据
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println("req Body Read Error,", err)
			return
		}
		fmt.Println("req.Body => ", string(body))
	
		// 11.服务端响应头header自定义
		w.Header().Add("RequestMethod", "PostForm")                                   // 此处将响应的RequestMethod header字段设置为PostForm
		w.Header().Add("Content-Type", "application/json;charset=UTF-8")              // 此处将响应的类型设置为JSON
		w.Header().Add("Cookies", fmt.Sprintf("method=form;id=%v;name=%v", id, name)) // 此处将响应的cookies设置为请求传入的参数
	
		// 12.返回给客户端的JSON数据组装
		uid, err := strconv.Atoi(id) // 将get到的id字段的值转换为整型
		if err != nil {
			errMsg := fmt.Sprintf("uid convert err! %v\n", err)
			fmt.Println(errMsg)
			w.Write([]byte(errMsg))
			return
		}
		reply := fmt.Sprintf("{\"method\":\"%v\",\"status\":\"ok\",\"data\":{\"id\":%v,\"name\":\"%v\"}}", "POSTFORM", uid, handler.Name)
		fmt.Println("reply => ", reply)
		w.Write([]byte(reply))
	default:
		// 如果这个handler不能识别这个路径，它会通过调用返回客户端一个HTTP404错误,并响应给客户端表明请求的路径不存在.
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such page: %s\n", req.URL)
	}

}


func main() {
	// 1.实例化一个处理所有请求的Handler接口
	handler := myHandler{Name: "WeiyiGeek"}

	// 2.创建一个自定义的Server,注意如果Handler为nil则采用http.DefaultServeMux进行处理响应,否则需要我们自己实现结构体的ServeHTTP方法.
	server := &http.Server{
		Addr:           ":8080",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	
	// 3.启动并监听HTTP服务端
	fmt.Printf("[%v] Http Server Start.....\n", time.Now().Format("2006-01-02 15:04:05"))
	log.Fatal(server.ListenAndServe())
	defer fmt.Printf("[%v] Http Server Close.....\n", time.Now().Format("2006-01-02 15:04:05"))
}
```
执行结果:
```sh
$ demo4 go build && ./demo4
[2021-11-23 13:25:30] Http Server Start.....
Method : Get, URL : /get?id=1024&name=%E5%94%AF%E4%B8%80%E6%9E%81%E5%AE%A2
QueryParam :  map[id:[1024] name:[唯一极客]]
id = 1024,name = 唯一极客
reply => {"method":"GET","status":"ok","data":{"id":1024,"name":"WeiyiGeek"}}

Method : POST, URL : /post
req.Body =>  {"id":128,"name":"Weiyi"}
id = 128,name = Weiyi
reply =>  {"method":"POST","status":"ok","data":{"id":128,"name":"Weiyi"}}

Method : POST, URL : /postform
map[id:[256] name:[WeiyiGeek-唯一极客]]
id = 256, name = WeiyiGeek-唯一极客
req.Body =>
reply =>  {"method":"POSTFORM","status":"ok","data":{"id":256,"name":"WeiyiGeek"}}
```
至此，Go语言中HTTP Server与HTTP Client 编程学习完毕！ 