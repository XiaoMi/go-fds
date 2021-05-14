package main

import (
	//	"fmt"
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

	objectMetaData := fds.NewObjectMetadata()
	objectMetaData.SetContentType("image/png")
	objectMetaData.Set(fds.XiaomiMetaPrefix+"uid", "1000")
	setObjectMetadataRequest := &fds.SetObjectMetadataRequest{
		BucketName: "bucketname",
		ObjectName: "test.txt",
		Metadata:   objectMetaData,
	}

	err = fdsClient.SetObjectMetadata(setObjectMetadataRequest)
	if err != nil {
		log.Fatal(err)
	}
}
