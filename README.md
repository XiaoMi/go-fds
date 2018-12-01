# galaxy-fds-sdk-go
FDS Go SDK.

[![Build Status](https://travis-ci.org/hujianxin/galaxy-fds-sdk-go.svg?branch=master)](https://travis-ci.org/hujianxin/galaxy-fds-sdk-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/hujianxin/galaxy-fds-sdk-go)](https://goreportcard.com/report/github.com/hujianxin/galaxy-fds-sdk-go)

[The formal Go SDK of FDS](https://github.com/XiaoMi/galaxy-fds-sdk-golang) is not well designed, but constrained by the fixed interface, I can't reconstruct it in large scale.

So, I start up this project for a good sdk design. 

:sparkles: :sparkles: :sparkles: **We got context support working, which make your concurrent program more fluent**

## Install
`go get -u github.com/hujianxin/galaxy-fds-sdk-go`

## Usage
```go
package main

import (
	"log"
	"os"

	"github.com/hujianxin/galaxy-fds-sdk-go/fds"
	"github.com/hujianxin/galaxy-fds-sdk-go/fds/manager"
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

## Development
To develop galaxy-fds-sdk-go, you'd better to upgrade your go version to 1.11+，because there is a `go modules` concept from go1.11, which can make it convenient.

## License
MIT
