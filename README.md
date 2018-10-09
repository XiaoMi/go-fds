# galaxy-fds-sdk-go
FDS Go SDK.

[The formal Go sdk of FDS](https://github.com/XiaoMi/galaxy-fds-sdk-golang) is not well designed, but constrained by the fixed interface, I can't reconstruct it in large scale.

So, I start up this project for a good sdk design. 

## Install
In order to use this FDS Client, you should upgrade your goland to 1.11+.

`go get -u github.com/v2tool/galaxy-fds-sdk-go`

## Usage
```go
package main

import (
	"github.com/v2tool/galaxy-fds-sdk-go/fds"
)

func main() {
	conf := fds.NewClientConfiguration("cnbj3-staging-fds.api.xiaomi.net​")
	conf.EnableHTTPS = false
	client := fds.New("xxxxxxxxxxxxxxxxxx", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", conf)

	_ = client.CreateBucket("zenmehuishiwoshilaohu")

}
```

For more sample, please look into `sample` package.

## License
