package fds_test

import (
	"github.com/v2tool/galaxy-fds-sdk-go/fds"
	"os"
)

var (
	// Endpoint/ID/Key
	endpoint  = os.Getenv("FDS_TEST_ENDPOINT")
	accessID  = os.Getenv("FDS_TEST_ACCESS_KEY_ID")
	accessKey = os.Getenv("FDS_TEST_ACCESS_KEY_SECRET")

	testBucket = "galaxy-fds-sdk-go-testing-bucketname-ut"

	conf   *fds.ClientConfiguration
	client *fds.Client
)

func createTestBucket() {
	exist, _ := client.DoesBucketExist(testBucket)
	if exist {
		client.DeleteBucket(testBucket)
	}

	createBucketRequest := &fds.CreateBucketRequest{
		BucketName: testBucket,
	}

	client.CreateBucket(createBucketRequest)
}

func deleteTestBucket() {
	exist, _ := client.DoesBucketExist(testBucket)
	if exist {
		client.DeleteBucket(testBucket)
	}
}
