package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hujianxin/galaxy-fds-sdk-go/fds"
)

func main() {
	fdsConf, err := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	fdsClient := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), fdsConf)

	// resp, err := fdsClient.DoesBucketExist("bucketname")
	// you can use context that you defined
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	resp, err := fdsClient.DoesBucketExitsWithContext(ctx, "bucketname")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp)
}
