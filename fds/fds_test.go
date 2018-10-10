package fds_test

import (
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"
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

func (suite *GalaxyFDSTestSuite) TestListAuthorizedBuckets() {
	resp, e := suite.client.ListAuthorizedBuckets()
	suite.Nil(e)
	suite.NotNil(resp)
}

func (suite *GalaxyFDSTestSuite) TestGetBucketACL() {
	acl, err := suite.client.GetBucketACL(suite.TestBucketName)
	suite.Nil(err)
	suite.Equal(len(acl.Grants), 1)
	suite.Equal(acl.Grants[0].Permission, fds.GrantPermissionFullControl)
	suite.Equal(acl.Grants[0].Type, fds.GrantTypeUser)
}

func (suite *GalaxyFDSTestSuite) TestSetBucketACL() {
	grant := fds.Grant{
		Grantee: fds.GrantKey{
			ID: "ALL_USERS",
		},
		Permission: fds.GrantPermissionRead,
		Type:       fds.GrantTypeGroup,
	}

	controlList := &fds.AccessControlList{}
	controlList.AddGrant(grant)
	e := suite.client.SetBucketACL(suite.TestBucketName, controlList)
	suite.Nil(e)

	acl, e := suite.client.GetBucketACL(suite.TestBucketName)
	suite.Nil(e)
	suite.Equal(len(acl.Grants), 2)
}

func (suite *GalaxyFDSTestSuite) TestGetAccessLog() {
	accessLog, err := suite.client.GetAccessLog(suite.TestBucketName)
	suite.Nil(err)
	suite.NotNil(accessLog)
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

func (suite *GalaxyFDSTestSuite) TestDoesObjectExist() {
	objectName := suite.GetRandomObjectName()
	exist, e := suite.client.DoesObjectExist(suite.TestBucketName, objectName)
	suite.Nil(e)
	suite.False(exist)

	testObjectContent := "Hello World"
	putObjectRequest := &fds.PutObjectRequest{
		BucketName: suite.TestBucketName,
		ObjectName: objectName,
		Data:       strings.NewReader(testObjectContent),
	}

	_, e = suite.client.PutObject(putObjectRequest)
	suite.Nil(e)

	exist, e = suite.client.DoesObjectExist(suite.TestBucketName, objectName)
	suite.Nil(e)
	suite.True(exist)
}

func (suite *GalaxyFDSTestSuite) TestCopyObject() {
	objectName := suite.GetRandomObjectName()
	testCopyObjectBucketName := suite.TestBucketName + "cp"
	suite.client.DeleteObjectsWithPrefix(testCopyObjectBucketName, "", false)
	suite.client.DeleteBucket(testCopyObjectBucketName)
	createBucketRequest := &fds.CreateBucketRequest{BucketName: testCopyObjectBucketName}
	e := suite.client.CreateBucket(createBucketRequest)
	suite.Nil(e)

	testCopyObjectObjectName := objectName + "cp"
	testObjectContent := "Hello World"
	putObjectRequest := &fds.PutObjectRequest{
		BucketName: suite.TestBucketName,
		ObjectName: objectName,
		Data:       strings.NewReader(testObjectContent),
	}

	_, e = suite.client.PutObject(putObjectRequest)
	suite.Nil(e)

	copyObjectRequest := &fds.CopyObjectRequest{
		SourceBucketName: suite.TestBucketName,
		SourceObjectName: objectName,
		TargetBucketName: testCopyObjectBucketName,
		TargetObjectName: testCopyObjectObjectName,
	}
	e = suite.client.CopyObject(copyObjectRequest)
	suite.Nil(e)

	getObjectRequest := &fds.GetObjectRequest{
		BucketName: testCopyObjectBucketName,
		ObjectName: testCopyObjectObjectName,
	}

	rc, e := suite.client.GetObject(getObjectRequest)
	suite.Nil(e)
	defer rc.Close()

	b, e := ioutil.ReadAll(rc)
	suite.Equal(string(b), testObjectContent)

	suite.client.DeleteBucket(testCopyObjectBucketName)
}

func (suite *GalaxyFDSTestSuite) TestRenameObject() {
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

	testRenameObjectName := testObjectName + "rename"

	e = suite.client.RenameObject(suite.TestBucketName, testObjectName, testRenameObjectName)
	suite.Nil(e)

	b, e := suite.client.DoesObjectExist(suite.TestBucketName, testObjectName)
	suite.Nil(e)
	suite.False(b)

	b, e = suite.client.DoesObjectExist(suite.TestBucketName, testRenameObjectName)
	suite.Nil(e)
	suite.True(b)
}

func (suite *GalaxyFDSTestSuite) TestDeleteObject() {
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

	e = suite.client.DeleteObject(suite.TestBucketName, testObjectName)
	suite.Nil(e)

	b, e := suite.client.DoesObjectExist(suite.TestBucketName, testObjectName)
	suite.Nil(e)
	suite.False(b)
}

func (suite *GalaxyFDSTestSuite) TestDeleteObjects() {
	testObjectContent := "Hello World"
	var names []string

	for i := 0; i < 10; i++ {
		objectName := strconv.FormatInt(int64(i), 10)
		putObjectRequest := &fds.PutObjectRequest{
			BucketName: suite.TestBucketName,
			ObjectName: objectName,
			Data:       strings.NewReader(testObjectContent),
		}

		response, e := suite.client.PutObject(putObjectRequest)
		suite.Nil(e)
		suite.Equal(response.ObjectName, objectName)
	}

	for i := 0; i < 10; i++ {
		objectName := strconv.FormatInt(int64(i), 10)
		names = append(names, objectName)
		b, e := suite.client.DoesObjectExist(suite.TestBucketName, objectName)
		suite.Nil(e)
		suite.True(b)
	}

	e := suite.client.DeleteObjects(suite.TestBucketName, names, false)
	suite.Nil(e)

	for i := 0; i < 10; i++ {
		objectName := strconv.FormatInt(int64(i), 10)
		b, e := suite.client.DoesObjectExist(suite.TestBucketName, objectName)
		suite.Nil(e)
		suite.False(b)
	}
}

func (suite *GalaxyFDSTestSuite) TestDeleteObjectsWithPrefix() {
	testObjectContent := "Hello World"
	var names []string

	for i := 0; i < 10; i++ {
		objectName := "prefix/" + strconv.FormatInt(int64(i), 10)
		putObjectRequest := &fds.PutObjectRequest{
			BucketName: suite.TestBucketName,
			ObjectName: objectName,
			Data:       strings.NewReader(testObjectContent),
		}

		response, e := suite.client.PutObject(putObjectRequest)
		suite.Nil(e)
		suite.Equal(response.ObjectName, objectName)
	}

	for i := 0; i < 10; i++ {
		objectName := "prefix/" + strconv.FormatInt(int64(i), 10)
		names = append(names, objectName)
		b, e := suite.client.DoesObjectExist(suite.TestBucketName, objectName)
		suite.Nil(e)
		suite.True(b)
	}

	e := suite.client.DeleteObjectsWithPrefix(suite.TestBucketName, "prefix/", false)
	suite.Nil(e)

	for i := 0; i < 10; i++ {
		objectName := "prefix/" + strconv.FormatInt(int64(i), 10)
		b, e := suite.client.DoesObjectExist(suite.TestBucketName, objectName)
		suite.Nil(e)
		suite.False(b)
	}
}

func TestGalaxyFDSuite(t *testing.T) {
	suite.Run(t, new(GalaxyFDSTestSuite))
}
