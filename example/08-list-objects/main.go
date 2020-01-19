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

	request := &fds.ListObjectsRequest{
		BucketName: "cnbj1-nginx-log",
		Prefix:     "2018/",
		Delimiter:  "/",
		MaxKeys:    10,
	}

	listing, err := fdsClient.ListObjects(request)
	if err != nil {
		log.Fatal(err)
	}

	for {
		fmt.Printf("%+v\n", listing)
		if !listing.Truncated {
			break
		}

		listing, err = fdsClient.ListObjectsNextBatch(listing)
		fmt.Println(listing.CommonPrefixes)
		if err != nil {
			log.Fatal(err)
		}
	}

}
