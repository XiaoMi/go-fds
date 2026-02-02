package main

import (
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

	objectBasicInfo, err := fdsClient.GetObjectBasicInfo("first-bucket-test-4", "test.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(objectBasicInfo)
	if objectBasicInfo.AllowOutsideAccess != nil {
		fmt.Println(*objectBasicInfo.AllowOutsideAccess)
	}
}
