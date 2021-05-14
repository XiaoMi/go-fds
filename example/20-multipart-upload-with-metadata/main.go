package main

import (
	"bufio"
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

	bigFileName := "/tmp/test.img"

	bigFile, err := os.Open(bigFileName)
	if err != nil {
		log.Fatal(err)
	}

	bigFileReader := bufio.NewReader(bigFile)

	bucketName := "hellotest"
	objectName := "test.img"
	initRequest := &fds.InitMultipartUploadRequest{
		BucketName: bucketName,
		ObjectName: objectName,
	}

	initResponse, err := fdsClient.InitMultipartUpload(initRequest)
	if err != nil {
		log.Fatal(err)
	}

	part := make([]byte, 1024*1024*10)
	uploadPartResultList := make([]fds.UploadPartResponse, 0)

	partNumber := 1
	for {
		buffer := bytes.NewBuffer(make([]byte, 0))
		count := 0
		if count, err = bigFileReader.Read(part); err != nil {
			break
		}
		buffer.Write(part[:count])
		uploadRequest := &fds.UploadPartRequest{
			BucketName: bucketName,
			ObjectName: objectName,
			UploadID:   initResponse.UploadID,
			PartNumber: partNumber,
			Data:       buffer,
		}
		uploadResponse, err := fdsClient.UploadPart(uploadRequest)
		if err != nil {
			log.Fatal(err)
		}
		uploadPartResultList = append(uploadPartResultList, *uploadResponse)

		partNumber++
	}

	objectMetaData := fds.NewObjectMetadata()
	objectMetaData.SetContentType("image/png")
	objectMetaData.Set(fds.XiaomiMetaPrefix+"uid", "1000")
	completeMultipartUploadRequest := &fds.CompleteMultipartUploadRequest{
		BucketName:  bucketName,
		ObjectName:  objectName,
		UploadID:    initResponse.UploadID,
		UploadParts: &fds.UploadPartList{UploadPartResultList: uploadPartResultList},
		Metadata:    objectMetaData,
	}

	resp, err := fdsClient.CompleteMultipartUploadV2(completeMultipartUploadRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
}
