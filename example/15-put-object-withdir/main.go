package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/XiaoMi/go-fds/fds"
)

func main() {
	fdsConf, err := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	fdsClient := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), fdsConf)

	request := &fds.PutObjectRequest{
		BucketName: "hello",
		ObjectName: "dir3/dir4/test.txt",
		Data:       bytes.NewReader([]byte("hello world")),
	}

	resp, err := fdsClient.PutObject(request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
