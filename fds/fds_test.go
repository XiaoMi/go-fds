package fds_test

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/v2tool/galaxy-fds-sdk-go/fds"
)

type GalaxyFDSTestSuite struct {
	suite.Suite

	client *fds.Client
	conf   *fds.ClientConfiguration

	Endpoint       string
	AccessID       string
	AccessKey      string
	TestBucketName string
}

func (suite *GalaxyFDSTestSuite) SetupSuite() {
	suite.Endpoint = os.Getenv("GO_FDS_TEST_ENDPOINT")
	suite.AccessID = os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID")
	suite.AccessKey = os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET")
	suite.TestBucketName = "galaxy-fds-sdk-go-testing-bucketname-ut"

	conf, err := fds.NewClientConfiguration(suite.Endpoint)
	if err != nil {
		log.Fatalln(err)
	}
	suite.conf = conf

	client := fds.New(suite.AccessID, suite.AccessKey, conf)
	suite.client = client
}

func (suite *GalaxyFDSTestSuite) BeforeTest(suiteName, testName string) {
	err := suite.client.DeleteObjectsWithPrefix(suite.TestBucketName, "", false)
	suite.Nil(err)
	suite.client.DeleteBucket(suite.TestBucketName)
	req := &fds.CreateBucketRequest{
		BucketName: suite.TestBucketName,
	}
	err = suite.client.CreateBucket(req)
	suite.Nil(err)
}

func (suite *GalaxyFDSTestSuite) GetRandomObjectName() string {
	pc, _, _, _ := runtime.Caller(1)
	return "golang-test-" + runtime.FuncForPC(pc).Name() + "-" + time.Now().Format(time.RFC3339)
}

// Already test it in BeforeTest
func (suite *GalaxyFDSTestSuite) TestCreateBucket() {}

func (suite *GalaxyFDSTestSuite) TestDoesBucketExist() {
	b, e := suite.client.DoesBucketExist(suite.TestBucketName)
	suite.Nil(e)
	suite.True(b)

	e = suite.client.DeleteBucket(suite.TestBucketName)
	suite.Nil(e)

	b, e = suite.client.DoesBucketExist(suite.TestBucketName)
	suite.Nil(e)
	suite.False(b)
}

// Already test it in TestDoesBucketExist
func (suite *GalaxyFDSTestSuite) TestDeleteBucket() {}

func (suite *GalaxyFDSTestSuite) TestGetBucketInfo() {
	response, e := suite.client.GetBucketInfo(suite.TestBucketName)
	suite.Nil(e)
	suite.NotNil(response)
	suite.Equal(suite.TestBucketName, response.BucketName)
}

func (suite *GalaxyFDSTestSuite) TestListBuckets() {
	response, e := suite.client.ListBuckets()
	suite.Nil(e)
	suite.NotNil(response)
}

func (suite *GalaxyFDSTestSuite) TestGetObject() {
	testObjectName := suite.GetRandomObjectName()
	testObjectContent := "Hello World"
	putObjectRequest := &fds.PutObjectRequest{
		BucketName: suite.TestBucketName,
		ObjectName: testObjectName,
		Data:       strings.NewReader(testObjectContent),
	}

	response, e := suite.client.PutObject(putObjectRequest)
	suite.Nil(e)
	suite.Equal(response.ObjectName, testObjectName)

	getObjectRequest := &fds.GetObjectRequest{
		BucketName: suite.TestBucketName,
		ObjectName: testObjectName,
	}

	rc, e := suite.client.GetObject(getObjectRequest)
	suite.Nil(e)
	defer rc.Close()

	b, e := ioutil.ReadAll(rc)
	suite.Equal(string(b), testObjectContent)
}

// Already test it in TestGetObject
func (suite *GalaxyFDSTestSuite) TestPutObject() {}

func TestGalaxyFDSuite(t *testing.T) {
	suite.Run(t, new(GalaxyFDSTestSuite))
}
