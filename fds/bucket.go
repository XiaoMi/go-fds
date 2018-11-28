package fds

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// CreateBucketRequest if request of creating bucket
// OrgID is option, if setted, bucket will be created under orgnization of orgid
type CreateBucketRequest struct {
	BucketName string `param:"-" header:"-"`
	OrgID      string `param:"orgId,omitempty" header:"-"`
}

// CreateBucket creates new bucket
func (client *Client) CreateBucket(request *CreateBucketRequest) error {
	buf := new(bytes.Buffer)

	req := &clientRequest{
		BucketName:         request.BucketName,
		Method:             HTTPPut,
		Data:               buf,
		QueryHeaderOptions: request,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err
}

// CreateBucketWithContext creates new bucket with context controlling
func (client *Client) CreateBucketWithContext(ctx *context.Context, request *CreateBucketRequest) error {
	return nil
}

// DoesBucketExist judge whether a bucket exist
func (client *Client) DoesBucketExist(bucketName string) (bool, error) {

	req := &clientRequest{
		BucketName: bucketName,
		Method:     HTTPHead,
	}

	resp, err := client.do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// DoesBucketExits judge whether bucket exitst with context controlling
func (client *Client) DoesBucketExits(ctx *context.Context, bucketName string) (bool, error) {
	return false, nil
}

// DeleteBucket delete a bucket
func (client *Client) DeleteBucket(bucketName string) error {
	req := &clientRequest{
		BucketName: bucketName,
		Method:     HTTPDelete,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err
}

// GetBucketInfoResponse is result of GetBucketInfo
type GetBucketInfoResponse struct {
	AllowOutsideAccess bool   `json:"allowOutsideAccess"`
	CreationTime       int64  `json:"creationTime"`
	BucketName         string `json:"name"`
	ObjectNum          int64  `json:"numObjects"`
	UsedSpace          int64  `json:"usedSpace"`
}

// GetBucketInfo get information of a bucket
func (client *Client) GetBucketInfo(bucketName string) (*GetBucketInfoResponse, error) {
	result := &GetBucketInfoResponse{}
	req := &clientRequest{
		BucketName: bucketName,
		Method:     HTTPGet,
		Result:     result,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return result, err
}

// Owner is owner of bucket or object
type Owner struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

// ListBucketsResponse is result of ListBuckets
type ListBucketsResponse struct {
	Owner   Owner                   `json:"owner"`
	Buckets []GetBucketInfoResponse `json:"buckets"`
}

// ListBuckets list all buckets
func (client *Client) ListBuckets() (*ListBucketsResponse, error) {
	result := &ListBucketsResponse{}
	req := &clientRequest{
		Method: HTTPGet,
		Result: result,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return result, err
}

type listAuthorizedBucketsOption struct {
	AuthorizedBuckets string `param:"authorizedBuckets" header:"-"`
}

// ListAuthorizedBuckets will return all buckets you could access
func (client *Client) ListAuthorizedBuckets() (*ListBucketsResponse, error) {
	result := &ListBucketsResponse{}
	req := &clientRequest{
		Method:             HTTPGet,
		QueryHeaderOptions: listAuthorizedBucketsOption{""},
		Result:             result,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return result, err
}

type migrateBucketOption struct {
	Migrate string `param:"migrate" header:"-"`
}

// MigrateBucketRequest is request of migrating bucket to other org and team
type MigrateBucketRequest struct {
	migrateBucketOption
	BucketName string `param:"-" header:"-"`
	OrgID      string `param:"orgId" header:"-"`
	TeamID     string `param:"teamId" header:"-"`
}

// MigrateBucket will change bucket's orgId and teamId
func (client *Client) MigrateBucket(request *MigrateBucketRequest) error {
	buf := new(bytes.Buffer)

	req := &clientRequest{
		BucketName:         request.BucketName,
		Method:             HTTPPut,
		QueryHeaderOptions: request,
		Data:               buf,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()
	return err
}

// GetBucketACL will return AccessControlList of bucket
func (client *Client) GetBucketACL(bucketName string) (*AccessControlList, error) {
	result := &AccessControlList{}
	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPGet,
		QueryHeaderOptions: aclOption{""},
		Result:             result,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return result, err
}

// SetBucketACL sets AccessControlList for bucket
func (client *Client) SetBucketACL(bucketName string, acl *AccessControlList) error {
	aclBytes, e := json.Marshal(acl)
	if e != nil {
		return errors.New("fds client: can't marshal acl")
	}

	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPPut,
		QueryHeaderOptions: aclOption{""},
		Data:               bytes.NewReader(aclBytes),
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err

}

// LifecycleConfig is ttl config of bucket
type LifecycleConfig struct {
}

// GetLifecycleConfig returns LifecycleConfig of bucket
func (client *Client) GetLifecycleConfig(bucketName string) (*LifecycleConfig, error) {
	return &LifecycleConfig{}, nil
}

// SetLifecycleConfig sets LifecycleConfig of bucket
func (client *Client) SetLifecycleConfig(bucketName string, config *LifecycleConfig) error {
	return nil
}

type accessLogOption struct {
	AccessLog string `param:"accessLog" header:"-"`
}

// AccessLog is accesslog
type AccessLog struct {
	BucketName    string `json:"bucketName"`
	Enabled       bool   `json:"enabled"`
	LogBucketName string `json:"logBucketName"`
	LogPrefix     string `json:"logPrefix"`
}

// GetAccessLog gets access log
func (client *Client) GetAccessLog(bucketName string) (*AccessLog, error) {
	result := &AccessLog{}
	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPGet,
		QueryHeaderOptions: accessLogOption{},
		Result:             result,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return result, err
}

// SetAccessLog sets acccess log
func (client *Client) SetAccessLog(bucketName string, accessLog *AccessLog) error {
	data, err := json.Marshal(accessLog)
	if err != nil {
		return err
	}

	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPPut,
		QueryHeaderOptions: accessLogOption{},
		Data:               bytes.NewReader(data),
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err
}
