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
	StorageClassType StorageClass `param:"storageClass,omitempty" header:"-"`
	OrgID      string `param:"orgId,omitempty" header:"-"`
}

// CreateBucket creates new bucket
func (client *Client) CreateBucket(request *CreateBucketRequest) error {
	return client.CreateBucketWithContext(context.Background(), request)
}

// CreateBucketWithContext creates new bucket with context controlling
func (client *Client) CreateBucketWithContext(ctx context.Context, request *CreateBucketRequest) error {
	buf := new(bytes.Buffer)

	req := &clientRequest{
		BucketName:         request.BucketName,
		Method:             HTTPPut,
		Data:               buf,
		QueryHeaderOptions: request,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

// DoesBucketExist judge whether a bucket exist
func (client *Client) DoesBucketExist(bucketName string) (bool, error) {
	return client.DoesBucketExitsWithContext(context.Background(), bucketName)
}

// DoesBucketExitsWithContext judge whether bucket exitst with context controlling
func (client *Client) DoesBucketExitsWithContext(ctx context.Context, bucketName string) (bool, error) {
	req := &clientRequest{
		BucketName: bucketName,
		Method:     HTTPHead,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

// DeleteBucket delete a bucket
func (client *Client) DeleteBucket(bucketName string) error {
	return client.DeleteBucketWithContext(context.Background(), bucketName)
}

// DeleteBucketWithContext delete bucket with context controlling
func (client *Client) DeleteBucketWithContext(ctx context.Context, bucketName string) error {
	req := &clientRequest{
		BucketName: bucketName,
		Method:     HTTPDelete,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
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
	return client.GetBucketInfoWithContext(context.Background(), bucketName)
}

// GetBucketInfoWithContext get information of bucket with context controlling
func (client *Client) GetBucketInfoWithContext(ctx context.Context, bucketName string) (*GetBucketInfoResponse, error) {
	result := &GetBucketInfoResponse{}
	req := &clientRequest{
		BucketName: bucketName,
		Method:     HTTPGet,
		Result:     result,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}
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
	return client.ListBucketsWithContext(context.Background())
}

// ListBucketsWithContext list all buckets with context controlling
func (client *Client) ListBucketsWithContext(ctx context.Context) (*ListBucketsResponse, error) {
	result := &ListBucketsResponse{}
	req := &clientRequest{
		Method: HTTPGet,
		Result: result,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return result, err
}

type listAuthorizedBucketsOption struct {
	AuthorizedBuckets string `param:"authorizedBuckets" header:"-"`
}

// ListAuthorizedBuckets will return all buckets user could access
func (client *Client) ListAuthorizedBuckets() (*ListBucketsResponse, error) {
	return client.ListAuthorizedBucketsWithContext(context.Background())
}

// ListAuthorizedBucketsWithContext will return all buckets users could access with context controlling
func (client *Client) ListAuthorizedBucketsWithContext(ctx context.Context) (*ListBucketsResponse, error) {
	result := &ListBucketsResponse{}
	req := &clientRequest{
		Method:             HTTPGet,
		QueryHeaderOptions: listAuthorizedBucketsOption{""},
		Result:             result,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}
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
	return client.MigrateBucketWithContext(context.Background(), request)
}

// MigrateBucketWithContext will change bucket's orgId and teamId with context controlling
func (client *Client) MigrateBucketWithContext(ctx context.Context, request *MigrateBucketRequest) error {
	buf := new(bytes.Buffer)

	req := &clientRequest{
		BucketName:         request.BucketName,
		Method:             HTTPPut,
		QueryHeaderOptions: request,
		Data:               buf,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return err
}

// GetBucketACL will return AccessControlList of bucket
func (client *Client) GetBucketACL(bucketName string) (*AccessControlList, error) {
	return client.GetBucketACLWithContext(context.Background(), bucketName)
}

// GetBucketACLWithContext will return AccessControlList of bucket with context controlling
func (client *Client) GetBucketACLWithContext(ctx context.Context, bucketName string) (*AccessControlList, error) {
	result := &AccessControlList{}
	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPGet,
		QueryHeaderOptions: aclOption{""},
		Result:             result,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return result, err
}

// SetBucketACL sets AccessControlList for bucket
func (client *Client) SetBucketACL(bucketName string, acl *AccessControlList) error {
	return client.SetBucketACLWithContext(context.Background(), bucketName, acl)
}

// SetBucketACLWithContext sets AccessControlList for bucket with context controlling
func (client *Client) SetBucketACLWithContext(ctx context.Context, bucketName string, acl *AccessControlList) error {
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

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

// LifecycleBaseItem replaces days in action
type LifecycleBaseItem struct {
	Days float64 `json:"days"`
}

// LifecycleActionType is enum type of action
type LifecycleActionType string

// Three LifecycleAction Types
const (
	Expiration                     LifecycleActionType = "expiration"
	NonCurrentVersionExpiration    LifecycleActionType = "nonCurrentVersionExpiration"
	AbortIncompleteMultipartUpload LifecycleActionType = "abortIncompleteMultipartUpload "
)

// LifecycleAction is action in LifecycleRule
type LifecycleAction map[LifecycleActionType]LifecycleBaseItem

// LifecycleRule is rule of LifecycleConfig
type LifecycleRule struct {
	ID      string          `json:"id"`
	Prefix  string          `json:"prefix"`
	Enabled bool            `json:"enabled"`
	Action  LifecycleAction `json:"actions"`
}

// NewLifecycleRuleFromJSON is a shortcut for translating json string to LifecycleRule.
// Becuase, constructing a LifecycleRule is too hard
func NewLifecycleRuleFromJSON(content []byte) (*LifecycleRule, error) {
	result := &LifecycleRule{}
	err := json.Unmarshal(content, result)
	return result, err
}

// LifecycleConfig is ttl config of bucket
type LifecycleConfig struct {
	Rules []LifecycleRule `json:"rules"`
}

// NewLifecycleConfigFromJSON is a shortcut for translating json string to LifecycleConfig.
// Becuase, constructing a LifecycleConfig is too hard
func NewLifecycleConfigFromJSON(content []byte) (*LifecycleConfig, error) {
	result := &LifecycleConfig{}
	err := json.Unmarshal(content, result)
	return result, err
}

// GetLifecycleConfig returns LifecycleConfig of bucket
func (client *Client) GetLifecycleConfig(request *GetLifecycleConfigRequest) (*LifecycleConfig, error) {
	return client.GetLifecycleConfigWithContext(context.Background(), request)
}

// GetLifecycleConfigRequest is request for getting lifecycle config
type GetLifecycleConfigRequest struct {
	BucketName string `param:"-" header:"-"`
	RuleID     string `param:"lifecycle" header:"-"`
}

// GetLifecycleConfigWithContext returns LifecycleConfig of bucket with context controlling
func (client *Client) GetLifecycleConfigWithContext(ctx context.Context, request *GetLifecycleConfigRequest) (*LifecycleConfig, error) {
	result := &LifecycleConfig{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		Method:             HTTPGet,
		QueryHeaderOptions: request,
		Result:             result,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return result, err
}

// SetLifecycleConfig sets LifecycleConfig of bucket
func (client *Client) SetLifecycleConfig(bucketName string, config *LifecycleConfig) error {
	return client.SetLifecycleConfigWithContext(context.Background(), bucketName, config)
}

type lifecycleOption struct {
	Lifecycle string `param:"lifecycle" header:"-"`
}

// SetLifecycleConfigWithContext sets LifecycleConfig of bucket with context controlling
func (client *Client) SetLifecycleConfigWithContext(ctx context.Context, bucketName string, config *LifecycleConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPPut,
		QueryHeaderOptions: lifecycleOption{},
		Data:               bytes.NewReader(data),
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

// SetLifecycleRule sets LifecycleRule of bucket
func (client *Client) SetLifecycleRule(bucketName string, rule *LifecycleRule) error {
	return client.SetLifecycleRuleWithContext(context.Background(), bucketName, rule)
}

// SetLifecycleRuleWithContext sets LifecycleRule of bucket with context controlling
func (client *Client) SetLifecycleRuleWithContext(ctx context.Context, bucketName string, rule *LifecycleRule) error {
	data, err := json.Marshal(rule)
	if err != nil {
		return err
	}

	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPPut,
		QueryHeaderOptions: lifecycleOption{"rule"},
		Data:               bytes.NewReader(data),
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

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
	return client.GetAccessLogWithContext(context.Background(), bucketName)
}

// GetAccessLogWithContext gets access log with context controlling
func (client *Client) GetAccessLogWithContext(ctx context.Context, bucketName string) (*AccessLog, error) {
	result := &AccessLog{}
	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPGet,
		QueryHeaderOptions: accessLogOption{},
		Result:             result,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return result, err
}

// SetAccessLog sets acccess log
func (client *Client) SetAccessLog(bucketName string, accessLog *AccessLog) error {
	return client.SetAccessLogWithContext(context.Background(), bucketName, accessLog)
}

// SetAccessLogWithContext sets acccess log with context controlling
func (client *Client) SetAccessLogWithContext(ctx context.Context, bucketName string, accessLog *AccessLog) error {
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

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}
