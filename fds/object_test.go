package fds_test

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/v2tool/galaxy-fds-sdk-go/fds"
)

var (
	testObjectName    = "testobjectname"
	testObjectContent = "Hello world"
)

func init() {
	// Endpoint/ID/Key
	endpoint = os.Getenv("FDS_TEST_ENDPOINT")
	accessID = os.Getenv("FDS_TEST_ACCESS_KEY_ID")
	accessKey = os.Getenv("FDS_TEST_ACCESS_KEY_SECRET")
	testBucket = "galaxy-fds-sdk-go-testing-bucketname-ut"

	conf, e := fds.NewClientConfiguration(endpoint)
	if e != nil {
		log.Fatalln(e)
	}
	conf.EnableHTTPS = false
	client = fds.New(accessID, accessKey, conf)
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
