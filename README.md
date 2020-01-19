# go-fds
FDS Go SDK.

[![Build Status](https://travis-ci.org/XiaoMi/go-fds.svg?branch=master)](https://travis-ci.org/XiaoMi/go-fds)
[![Go Report Card](https://goreportcard.com/badge/github.com/XiaoMi/go-fds)](https://goreportcard.com/report/github.com/XiaoMi/go-fds)

[The formal Go SDK of FDS](https://github.com/XiaoMi/go-fdslang) is not well designed, but constrained by the fixed interface, I can't reconstruct it in large scale.

So, I start up this project for a good sdk design. 

:sparkles: :sparkles: :sparkles: **We got context support working, which make your concurrent program more fluent**

## Install
`go get -u github.com/XiaoMi/go-fds`

## Usage
```go
package main

import (
	"log"
	"os"

	"github.com/XiaoMi/go-fds/fds"
	"github.com/XiaoMi/go-fds/fds/manager"
)

func main() {
	conf, _ := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	client := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), conf)

	downloader, _ := manager.NewDownloader(client, 1024*1024, 10, true)

	request := &manager.DownloadRequest{
		GetObjectRequest: fds.GetObjectRequest{
			BucketName: "hellodf",
			ObjectName: "build.log",
		},
		FilePath: "/home/XiaoMi/tmp/build.log",
	}
	err := downloader.Download(request)
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Done")
	}
}
```

For more sample, please look into `example` package

## Development
To develop go-fds, you'd better to upgrade your go version to 1.13+.
