package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hujianxin/galaxy-fds-sdk-go/fds"
)

func main() {
	fdsConf, err := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	fdsClient := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), fdsConf)

	resp, err := fdsClient.ListAuthorizedBuckets()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
