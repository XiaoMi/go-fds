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
