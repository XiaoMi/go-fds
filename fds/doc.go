/*
Package fds 是小米云对象存储服务FDS的go语言客户端。

sdk通过封装原始的JSON Restful API，来实现方便地调用。

更多的使用方法，可以通过源码中的fds_test.go文件学到。

Usage:

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

		downloader := manager.NewDownloader(client, 1024*1024, 10, true)

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

*/
package fds
