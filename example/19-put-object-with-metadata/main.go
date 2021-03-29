package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"net/http"

	"github.com/XiaoMi/go-fds/fds"
)

func main() {
	fdsConf, err := fds.NewClientConfiguration(os.Getenv("GO_FDS_TEST_ENDPOINT"))
	if err != nil {
		log.Fatal(err)
	}

	fdsClient := fds.New(os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID"), os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET"), fdsConf)

        // user defined metadata
        objectMetaData := fds.NewObjectMetadata()
        objectMetaData.Set(fds.XiaomiMetaPrefix + "uid", "1000")

	request := &fds.PutObjectRequest{
		BucketName: "bucketname",
		ObjectName: "test.txt",
		ContentType: "text/plain",
		ContentMd5: "5eb63bbbe01eeed093cb22bb8f5acdc3",
		Metadata: objectMetaData,
		Data:       bytes.NewReader([]byte("hello world")),
	}

	resp, err := fdsClient.PutObject(request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
