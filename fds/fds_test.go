package fds_test

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v2tool/galaxy-fds-sdk-go/fds"
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

var (
	testObjectName    = "testobjectname"
	testObjectContent = "Hello world"
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

func TestClient_CreateBucket(t *testing.T) {
	exist, _ := client.DoesBucketExist(testBucket)
	if exist {
		e := client.DeleteBucket(testBucket)
		assert.Nil(t, e)
	}

	createBucketRequest := &fds.CreateBucketRequest{
		BucketName: testBucket,
	}

	err := client.CreateBucket(createBucketRequest)
	assert.Nil(t, err)

	err = client.CreateBucket(createBucketRequest)
	assert.NotNil(t, err)
}

func TestClient_DoesBucketExist(t *testing.T) {
	client.DeleteBucket(testBucket)

	b, e := client.DoesBucketExist(testBucket)
	assert.Nil(t, e)
	assert.False(t, b)
}

func TestClient_DeleteBucket(t *testing.T) {
	exist, _ := client.DoesBucketExist(testBucket)
	if !exist {
		createBucketRequest := &fds.CreateBucketRequest{
			BucketName: testBucket,
		}
		client.CreateBucket(createBucketRequest)
	}

	e := client.DeleteBucket(testBucket)
	assert.Nil(t, e)
	e = client.DeleteBucket(testBucket)
	assert.NotNil(t, e)
}

func TestClient_GetBucketInfo(t *testing.T) {
	exist, _ := client.DoesBucketExist(testBucket)
	if !exist {
		createBucketRequest := &fds.CreateBucketRequest{
			BucketName: testBucket,
		}
		client.CreateBucket(createBucketRequest)
	}

	response, e := client.GetBucketInfo(testBucket)
	assert.Nil(t, e)
	assert.NotNil(t, response)
}

func TestClient_ListBuckets(t *testing.T) {
	exist, _ := client.DoesBucketExist(testBucket)
	if !exist {
		createBucketRequest := &fds.CreateBucketRequest{
			BucketName: testBucket,
		}
		client.CreateBucket(createBucketRequest)
	}

	response, e := client.ListBuckets()
	assert.Nil(t, e)
	assert.NotNil(t, response)
}

func TestClient_GetObject(t *testing.T) {
	putObjectRequest := &fds.PutObjectRequest{
		BucketName: testBucket,
		ObjectName: testObjectName,
		Data:       strings.NewReader(testObjectContent),
	}

	_, e := client.PutObject(putObjectRequest)
	assert.Nil(t, e)

	getObjectRequest := &fds.GetObjectRequest{
		BucketName: testBucket,
		ObjectName: testObjectName,
	}

	rc, e := client.GetObject(getObjectRequest)
	assert.Nil(t, e)
	defer rc.Close()

	b, e := ioutil.ReadAll(rc)
	assert.Equal(t, string(b), testObjectContent)
}

func TestClient_PutObject(t *testing.T) {
	putObjectRequest := &fds.PutObjectRequest{
		BucketName: testBucket,
		ObjectName: testObjectName,
		Data:       strings.NewReader(testObjectContent),
	}

	putObjectResponse, e := client.PutObject(putObjectRequest)
	assert.Nil(t, e)

	assert.Equal(t, putObjectResponse.BucketName, testBucket)
	assert.Equal(t, putObjectResponse.ObjectName, testObjectName)
}
