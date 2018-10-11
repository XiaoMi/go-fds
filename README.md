# galaxy-fds-sdk-go
FDS Go SDK.

[![Build Status](https://travis-ci.org/v2tool/galaxy-fds-sdk-go.svg?branch=master)](https://travis-ci.org/v2tool/galaxy-fds-sdk-go)

[The formal Go sdk of FDS](https://github.com/XiaoMi/galaxy-fds-sdk-golang) is not well designed, but constrained by the fixed interface, I can't reconstruct it in large scale.

So, I start up this project for a good sdk design. 

## Install
In order to use this FDS Client, you should upgrade your goland to 1.11+.

`go get -u github.com/v2tool/galaxy-fds-sdk-go`

## Usage
```go
package main

import (
	"log"
	"os"

	"github.com/v2tool/galaxy-fds-sdk-go/fds"
	"github.com/v2tool/galaxy-fds-sdk-go/fds/manager"
)

func main() {
	conf, _ := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	client := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), conf)

	downloader := manager.NewDownloader(client, 1024*1024, 10, true)

	request := &manager.DownloadRequest{
		GetObjectRequest: fds.GetObjectRequest{
			BucketName: "hellodf",
			ObjectName: "build.log",
		},
		FilePath: "/home/hujianxin/tmp/build.log",
	}
	err := downloader.Download(request)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Done")
	}
}
```

For more sample, please look into `fds/fds_test.go` file.

## License
