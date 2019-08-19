package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hujianxin/go-fds/fds"
	"io/ioutil"
)

func main() {
	fdsConf, err := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	fdsClient := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), fdsConf)

	request := &fds.GetObjectRequest{
		BucketName: "hello",
		ObjectName: "build.log",
	}

	body, err := fdsClient.GetObject(request)
	if err != nil {
		log.Fatal(err)
	}
	defer body.Close()

	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", b)
}
