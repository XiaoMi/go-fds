# galaxy-fds-sdk-go
小米文件存储FDS(File Storage Service) Go SDK.

# 简介
因为接口已经固定，[旧SDK](https://github.com/XiaoMi/galaxy-fds-sdk-golang)已经不容易改动，但是存在一些问题，是重构无法解决的，比如：

1. 内存不友好，旧SDK中，对object读写都是通过byte数组的形式，如果object较大，并且用户并发的写入，则会占用较大内存，并不能像Java SDK中使用InputStreamming那样友好。在Go语言中，提供了类似的机制，io.Reader io.Writer，这是输入输出的一种使用共识。
2. 添加新的API较为困难，旧SDK每添加一个新的API都需要考虑过多的内在逻辑，写过多的冗余代码，新SDK则只需要几行代码就可以完成添加功能，大部分重复逻辑已经简化。
3. 交互使用略显困难，比如在PutObject时，可以指定metadata或者Header，但是旧SDK这方面不易操作。新SDK通过反射的方式，解决了可选输入的问题。
4. 接口、包命名不符规范，在Go语言中，有较为严格的命名规范，比如包命名必须是小写，函数命名不能带_，如果带了是勉强可以通过运行的，但是这不符合规范。

另外，除了上面几个新SDK具有的优点之外，新版SDK参照S3和OSS的sdk设计，我也为这个SDK提供了并发的upload和download。

最后，auth.go里面的signature方法是直接使用了旧SDK中的代码。

在开发新的sdk的过程中，需要使用go version1.11+，其中使用了go语言新概念，go modules。由此，咱们的代码可以通过git submodule的方式放入galaxy工程内部，而并不是非得放到gopath里面，方便管理。

## 安装
`go get -u github.com/v2tool/galaxy-fds-sdk-go`

## 用法
```go
package main

import (
	"github.com/v2tool/galaxy-fds-sdk-go/fds"
)

func main() {
	conf := fds.NewClientConfiguration("cnbj3-staging-fds.api.xiaomi.net")
	conf.EnableHTTPS = false
	client := fds.New("xxxxxxxxxxxxx", "xxxxxxxxxxxxxxxxxx", conf)

	_ = client.CreateBucket("zenmehuishiwoshilaohu")

}
```

## License
