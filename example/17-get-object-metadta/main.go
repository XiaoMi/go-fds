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

	objectMetadata, err := fdsClient.GetObjectMetadata("bucketname", "test.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", objectMetadata)
	fmt.Printf("%v\n", objectMetadata.Get("x-xiaomi-meta-content-length"))
	fmt.Printf("%s\n", objectMetadata.GetContentType())

	fmt.Printf("%v\n", objectMetadata.GetRawMetadata())
}
