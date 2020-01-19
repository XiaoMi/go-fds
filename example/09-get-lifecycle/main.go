package main

import (
	"fmt"
	"log"
	"os"

	"github.com/XiaoMi/go-fds/fds"
)

func main() {
	fdsConf, err := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	fdsConf.EnableHTTPS = false
	if err != nil {
		log.Fatal(err)
	}

	fdsClient := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), fdsConf)

	request := &fds.GetLifecycleConfigRequest{
		BucketName: "hellonihao",
	}

	resp, err := fdsClient.GetLifecycleConfig(request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", resp)
}
