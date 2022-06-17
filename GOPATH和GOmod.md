GOROOT 是环境变量

GOPATH 是Golang 1.5版本之前一个重要的环境变量配置，是存放 Golang 项目代码的文件路径。

## 如何查看 GOPATH 路径

```sh
go env GOPATH
```

或者输入:

```sh
go env | grep GOPATH
```

进入GOPATH目录，查看该目录下的所有文件。

```
go
├── bin
├── pkg
└── src
    ├── github.com
    ├── golang.org
    ├── google.golang.org
    ....
```

可以看到有三个文件夹。

•bin  存放编译生成的二进制文件。比如 执行命令 go get github.com/google/gops，bin目录会生成 gops 的二进制文件。

• pkg  其中pkg目录下有三个文件夹。

• xx_amd64:其中 xx 是目标操作系统，如 mac 对应的是darwin_amd64, linux 系统对应的是 linux_amd64，存放的是.a结尾的文件。

•mod: 当开启go Modules 模式下，go get命令缓存下依赖包存放的位置

•sumdb: go get命令缓存下载的checksum数据存放的位置

•src 存放golang项目代码的位置



![图片](https://mmbiz.qpic.cn/mmbiz_jpg/krB1icImoA9vfbXeBKibkoUID1O89HdTULuzAH0nxxynFwO1uezYmss9pGhuk8oW94nlIvcttR6eDoaImfbQ5uDw/640?wx_fmt=jpeg&wxfrom=5&wx_lazy=1&wx_co=1)

因此在使用 GOPATH 模式下，我们需要将应用代码存放在固定的`$GOPATH/src`目录下，并且如果执行`go get`来拉取外部依赖会自动下载并安装到`$GOPATH`目录下。

简单来说，GOPATH模式下，项目代码不能想放哪里就放哪里。



## GOPATH 缺点

•go get 命令的时候，无法指定获取的版本•引用第三方项目的时候，无法处理v1、v2、v3等不同版本的引用问题，因为在GOPATH 模式下项目路径都是 github.com/foo/project•无法同步一致第三方版本号，在运行 Go 应用程序的时候，无法保证其它人与所期望依赖的第三方库是相同的版本。

## 为什么需要Go Modules

在go 1.11 官方出手了推出了 Go Modules， 通过设置环境变量 GO111MODULE 进行开启或者关闭 go mod 模式。

•auto 自动模式，当项目根目录有 go.mod 文件，启用 Go modules•off 关闭 go mod 模式•on 开启go mod 模式

开启 go mod 模式后，你的项目代码想放哪里就放哪里，你想引用哪个版本就用哪个版本。

## GOPROXY

环境变量 GOPROXY 就是设置 Go 模块代理的，其作用直接通过镜像站点来快速拉取所需项目代码。

常见代理配置

•阿里云：https://mirrors.aliyun.com/goproxy/• 七牛云：https://goproxy.cn,direct

执行命令：

```sh
go env -w GOPROXY="https://goproxy.cn,direct" 
```

## 初始化Modules

基于 go1.17.3 版本

新创建一个空目录test_mod，进入该目录，执行命令

```sh
//test_mod 为项目名称
go mod int test_mod
```

会在根目录生成一个 go.mod 文件，内容如下：

```go
module test_mod

go 1.17
```

如果想引入第三方网络包，在该项目目录执行 go get 仓库地址。比如引入定时任务：

```go
go get github.com/robfig/cron/v3
```

go.mod 会变成为, indirect 代表是间接依赖，因为当前项目是空的，所以并没有发现这个模块的明确引用。

```go
module test_mod

go 1.17

require github.com/robfig/cron/v3 v3.0.1 // indirect
```

并且也会新增一个go.sum文件, 它的作用是保证项目所依赖的模块版本，不会被篡改。

```
github.com/robfig/cron/v3 v3.0.1 h1:WdRxkvbJztn8LMz/QEvLN5sBU+xKpSqwwUO1Pjr4qDs=github.com/robfig/cron/v3 v3.0.1/go.mod h1:eQICP3HwyT7UooqI/z+Ov+PtYAWygg1TEWWzGIFLtro=
```

注意此时，我们的项目是没有任何go代码文件的，现在只有 go.mod 和 go.sum 两个文件。

## go mod tidy

如果我们 go.mod 导入了第三方包，但项目代码中我不用，就是玩。领导发现后，不小心一个 go mod tidy 命令，直接把你回到解放前。

观察 go.mod 会发现已经没有了这串神秘代码

```
require github.com/robfig/cron/v3 v3.0.1 // indirect
```

go mod tidy 就是去掉go.mod文件中项目不需要的依赖。

## go mod edit

## 方法一

执行命令：

```sh
go mod edit -replace [old git package]@[version]=[new git package]@[version]
```

例如：

```sh
go mod edit -replace github.com/bndr/gojenkins=github.com/Bpazy/gojenkins@latest
```

执行后 ，会发现 go.mod 文件最后有一串神秘代码

```
replace github.com/bndr/gojenkins => github.com/Bpazy/gojenkins v1.0.2-0.20200708084040-3655c428bba9
```

## 方法二 

简单粗暴，直接修改go.mod文件，在go.mod文件最后添加以下神秘代码

```
replace github.com/bndr/gojenkins => github.com/Bpazy/gojenkins v1.0.2-0.20200708084040-3655c428bba9
```

即可完美解决此问题，replace 还有一个隐藏的秘密，那就是可引入本地项目代码

```
replace github.com/bndr/gojenkins => ../gojenkins
```