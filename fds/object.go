package fds

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"strings"
)

// GetObjectRequest is the input of GetObject method
type GetObjectRequest struct {
	BucketName string `param:"-" header:"-"`
	ObjectName string `param:"-" header:"-"`
	Range      string `param:"-" header:"Range,omitempty"`
}

// GetObject will get full content of object
func (client *Client) GetObject(request *GetObjectRequest) (io.ReadCloser, error) {
	return client.GetObjectWithContext(context.Background(), request)
}

// GetObjectWithContext will get full content of object with context controlling
func (client *Client) GetObjectWithContext(ctx context.Context, request *GetObjectRequest) (io.ReadCloser, error) {
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		QueryHeaderOptions: request,
		Method:             HTTPGet,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

// PutObjectRequest is the input of PutObject method
type PutObjectRequest struct {
	BucketName string    `param:"-" header:"-"`
	ObjectName string    `param:"-" header:"-"`
	Data       io.Reader `param:"-" header:"-"`

	CacheControl       string `header:"Cache-Control,omitempty" param:"-"`
	ContentDisposition string `header:"Content-Disposition,omitempty" param:"-"`
	ContentEncoding    string `header:"Content-Encoding,omitempty" param:"-"`
	ContentType        string `header:"Content-Type,omitempty" param:"-"`
	ContentLength      int    `header:"Content-Length,omitempty" param:"-"`
	ContentMd5         string `header:"Content-Md5,omitempty" param:"-"`
	Expect             string `header:"Expect,omitempty" param:"-"`
	Metadata           *ObjectMetadata `header:"-" param:"-"`
}

// PutObjectResponse is the result of PutObject method
type PutObjectResponse struct {
	BucketName        string `json:"bucketName"`
	ObjectName        string `json:"objectName"`
	AccessKeyID       string `json:"accessKeyId"`
	Signature         string `json:"signature"`
	Expires           int64  `json:"expires"`
	PreviousVersionID string `json:"previousVersionId"`
	OutsideAccess     bool   `json:"outsideAccess"`
}

// PutObject will create object
func (client *Client) PutObject(request *PutObjectRequest) (*PutObjectResponse, error) {
	return client.PutObjectWithContext(context.Background(), request)
}

// PutObjectWithContext will create object with context controlling
func (client *Client) PutObjectWithContext(ctx context.Context, request *PutObjectRequest) (*PutObjectResponse, error) {
	result := &PutObjectResponse{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Data:               request.Data,
		QueryHeaderOptions: request,
		Metadata:           request.Metadata,
		Method:             HTTPPut,
		Result:             result,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	return result, nil
}

// DoesObjectExist judge wether object exists
func (client *Client) DoesObjectExist(bucketName, objectName string) (bool, error) {
	return client.DoesObjectExistWithContext(context.Background(), bucketName, objectName)
}

// DoesObjectExistWithContext judge wether object exists with context controlling
func (client *Client) DoesObjectExistWithContext(ctx context.Context, bucketName, objectName string) (bool, error) {
	req := &clientRequest{
		BucketName: bucketName,
		ObjectName: objectName,
		Method:     HTTPHead,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, nil
}

type copyObjectOption struct {
	Copy string `param:"cp" header:"-"`
}

// CopyObjectRequest is the input of CopyObject method
type CopyObjectRequest struct {
	copyObjectOption
	SourceBucketName string `param:"-" header:"-"`
	SourceObjectName string `param:"-" header:"-"`
	TargetBucketName string `param:"-" header:"-"`
	TargetObjectName string `param:"-" header:"-"`
}

// CopyObject copy object from a bucket to other bucket
func (client *Client) CopyObject(request *CopyObjectRequest) error {
	return client.CopyObjectWithContext(context.Background(), request)
}

// CopyObjectWithContext copy object from a bucket to other bucket with context controlling
func (client *Client) CopyObjectWithContext(ctx context.Context, request *CopyObjectRequest) error {
	dataString := map[string]string{
		"srcBucketName": request.SourceBucketName,
		"srcObjectName": request.SourceObjectName,
	}

	data, e := json.Marshal(dataString)
	if e != nil {
		return e
	}

	req := &clientRequest{
		BucketName:         request.TargetBucketName,
		ObjectName:         request.TargetObjectName,
		QueryHeaderOptions: request,
		Method:             HTTPPut,
		Data:               bytes.NewReader(data),
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type renameObjectOption struct {
	RenameTo string `param:"renameTo" header:"-"`
}

// RenameObject renames sourceObjectName in bucketName to targetObjectName
func (client *Client) RenameObject(bucketName, sourceObjectName, targetObjectName string) error {
	return client.RenameObjectWithContext(context.Background(), bucketName, sourceObjectName, targetObjectName)
}

// RenameObjectWithContext renames sourceObjectName in bucketName to targetObjectName with context controlling
func (client *Client) RenameObjectWithContext(ctx context.Context, bucketName, sourceObjectName, targetObjectName string) error {
	req := &clientRequest{
		BucketName:         bucketName,
		ObjectName:         sourceObjectName,
		QueryHeaderOptions: renameObjectOption{targetObjectName},
		Method:             HTTPPut,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteObject deletes objectName in bucketName
func (client *Client) DeleteObject(bucketName, objectName string) error {
	return client.DeleteObjectWithContext(context.Background(), bucketName, objectName)
}

// DeleteObjectWithContext deletes object in bucket with context controlling
func (client *Client) DeleteObjectWithContext(ctx context.Context, bucketName, objectName string) error {
	req := &clientRequest{
		BucketName: bucketName,
		ObjectName: objectName,
		Method:     HTTPDelete,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type deleteObjectsOption struct {
	EnableTrash   bool   `param:"enableTrash" header:"-"`
	DeleteObjects string `param:"deleteObjects" header:"-"`
}

// DeleteObjects will delete all objects in objectNames
func (client *Client) DeleteObjects(bucketName string, objectNames []string, put2trash bool) error {
	return client.DeleteObjectsWithContext(context.Background(), bucketName, objectNames, put2trash)
}

// DeleteObjectsWithContext will delete all objects in bucket with context controlling
func (client *Client) DeleteObjectsWithContext(ctx context.Context, bucketName string, objectNames []string, put2trash bool) error {
	data, err := json.Marshal(objectNames)
	if err != nil {
		return err
	}

	req := &clientRequest{
		BucketName:         bucketName,
		Method:             HTTPPut,
		QueryHeaderOptions: deleteObjectsOption{EnableTrash: put2trash},
		Data:               bytes.NewReader(data),
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// DeleteObjectsWithPrefix will delete all objects with prefix of prefix
func (client *Client) DeleteObjectsWithPrefix(bucketName, prefix string, put2stash bool) error {
	return client.DeleteObjectsWithPrefixWithContext(context.Background(), bucketName, prefix, put2stash)
}

// DeleteObjectsWithPrefixWithContext will delete all objects with prefix of prefix with context controlling
func (client *Client) DeleteObjectsWithPrefixWithContext(ctx context.Context, bucketName, prefix string, put2stash bool) error {
	listObjectsRequest := &ListObjectsRequest{
		BucketName: bucketName,
		Prefix:     prefix,
		Delimiter:  "",
		MaxKeys:    DefaultListObjectsMaxKeys,
	}

	objectListing, err := client.ListObjectsWithContext(ctx, listObjectsRequest)
	if err != nil {
		return err
	}

	for {
		names := make([]string, 0)
		for _, o := range objectListing.ObjectSummaries {
			names = append(names, o.ObjectName)
		}

		err = client.DeleteObjectsWithContext(ctx, bucketName, names, put2stash)
		if err != nil {
			return err
		}

		if !objectListing.Truncated {
			break
		}

		objectListing, err = client.ListObjectsNextBatchWithContext(ctx, objectListing)
		if err != nil {
			return err
		}
	}

	return nil
}

// ObjectMetadata is metadata of object
type ObjectMetadata struct {
	metadata map[string]string
}


var predefinedMetadata = map[string]string{
        HTTPHeaderCacheControl          : "",
        HTTPHeaderContentLength         : "",
        HTTPHeaderContentEncoding       : "",
        HTTPHeaderLastModified          : "",
        HTTPHeaderContentMD5            : "",
        HTTPHeaderContentType           : "",
        HTTPHeaderLastChecked           : "",
        HTTPHeaderUploadTime            : "",
        HTTPHeaderDate                  : "",
        HTTPHeaderAuthorization         : "",
        HTTPHeaderRange                 : "",
        HTTPHeaderContentRange          : "",
        HTTPHeaderContentMetadataLength : "",
        HTTPHeaderServerSideEncryption  : "",
        HTTPHeaderStorageClass          : "",
        HTTPHeaderOngoingRestore        : "",
        HTTPHeaderRestoreExpireDate     : "",
        HTTPHeaderCRC64ECMA             : "",
}

// NewObjectMetadata create a default ObjectMetadata
func NewObjectMetadata() *ObjectMetadata {
	return &ObjectMetadata{map[string]string{}}
}

// Get method of ObjectMetadata
func (metadata *ObjectMetadata) Get(k string) string {
	key := strings.ToLower(k)
	return metadata.metadata[key]
}

// Set method of ObjectMetadata
func (metadata *ObjectMetadata) Set(k, v string) error {
	key := strings.ToLower(k)
	_, ok := predefinedMetadata[key]
	if ok || strings.HasPrefix(key, XiaomiMetaPrefix) {
		metadata.metadata[key] = v
		return nil
	} else {
		return errors.New("Invalid metadata: " + k)
	}
}

func (metadata *ObjectMetadata) GetRawMetadata() map[string]string {
	data := make(map[string]string)
	for k,v := range metadata.metadata {
		data[k] = v
	}
	return data
}

// GetContentLength gets ContentLength of object metadata
func (metadata *ObjectMetadata) GetContentLength() (int64, error) {
	return strconv.ParseInt(metadata.Get(HTTPHeaderContentMetadataLength), 10, 64)
}

// SetContentLength sets ContentLength of object metadata
func (metadata *ObjectMetadata) SetContentLength(length int64) {
	metadata.Set(HTTPHeaderContentMetadataLength, strconv.FormatInt(length, 10))
}

func (metadata *ObjectMetadata) GetContentType() string {
	return metadata.Get(HTTPHeaderContentType)
}

func (metadata *ObjectMetadata) SetContentType(contentType string) {
	metadata.Set(HTTPHeaderContentType, contentType)
}

func (metadata *ObjectMetadata) serialize() ([]byte, error) {
	data := make(map[string]map[string]string)

	data["rawMeta"] = metadata.metadata
	result, e := json.Marshal(data)
	if e != nil {
		return nil, e
	}

	return result, nil
}

func parseObjectMetadataFromHeader(header http.Header) *ObjectMetadata {
	objectMetadata := NewObjectMetadata()
	for k := range header {
		key := strings.ToLower(k)
		_, ok := predefinedMetadata[key]
		if ok || strings.HasPrefix(key, XiaomiMetaPrefix) {
			objectMetadata.Set(key, header.Get(k))
		}
	}
	return objectMetadata
}

type getObjectMetadataOption struct {
	Metadata string `param:"metadata" header:"-"`
}

// GetObjectMetadata gets metadata of objectName in bucketName
func (client *Client) GetObjectMetadata(bucketName, objectName string) (*ObjectMetadata, error) {
	return client.GetObjectMetadataWithContext(context.Background(), bucketName, objectName)
}

// GetObjectMetadataWithContext gets metadata of objectName in bucketName with context controlling
func (client *Client) GetObjectMetadataWithContext(ctx context.Context, bucketName, objectName string) (*ObjectMetadata, error) {
	req := &clientRequest{
		BucketName:         bucketName,
		ObjectName:         objectName,
		Method:             HTTPGet,
		QueryHeaderOptions: getObjectMetadataOption{},
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return &ObjectMetadata{}, err
	}
	defer resp.Body.Close()
	result := parseObjectMetadataFromHeader(resp.Header)

	return result, nil
}

type setObjectMetadataOption struct {
	SetMetadata string `param:"setMetaData" header:"-"`
}

// SetObjectMetadataRequest is input of SetObjectMetadata
type SetObjectMetadataRequest struct {
	setObjectMetadataOption
	BucketName string          `param:"-" header:"-"`
	ObjectName string          `param:"-" header:"-"`
	Metadata   *ObjectMetadata `param:"-" header:"-"`
}

// SetObjectMetadata sets metadata of object
func (client *Client) SetObjectMetadata(request *SetObjectMetadataRequest) error {
	return client.SetObjectMetadataWithContext(context.Background(), request)
}

// SetObjectMetadataWithContext sets metadata of object with context controlling
func (client *Client) SetObjectMetadataWithContext(ctx context.Context, request *SetObjectMetadataRequest) error {
	data, e := request.Metadata.serialize()
	if e != nil {
		return e
	}

	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPPut,
		QueryHeaderOptions: request,
		Data:               bytes.NewReader(data),
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

// ObjectSummary beans
type ObjectSummary struct {
	ETag         string    `json:"etag"`
	ObjectName   string    `json:"name"`
	Owner        Owner     `json:"owner"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"lastModified"`
	UploadTime   int64     `json:"uploadTime"`
}

// ListObjectsRequest is input of ListObjectsRequest
type ListObjectsRequest struct {
	BucketName string `param:"-" header:"-"`
	Prefix     string `param:"prefix" header:"-"`
	Delimiter  string `param:"delimiter" header:"-"`
	MaxKeys    int    `param:"maxKeys" header:"-"`
}

// ObjectListing bean
type ObjectListing struct {
	BucketName      string          `json:"name" param:"-" header:"-"`
	Prefix          string          `json:"prefix"  param:"prefix" header:"-"`
	MaxKeys         int             `json:"maxKeys" param:"maxKeys" header:"-"`
	Marker          string          `json:"marker" param:"-" header:"-"`
	Truncated       bool            `json:"truncated" param:"-" header:"-"`
	NextMarker      string          `json:"nextMarker" param:"marker" header:"-"`
	Delimiter       string          `json:"delimiter" param:"delimiter" header:"-"`
	ObjectSummaries []ObjectSummary `json:"objects" param:"-" header:"-"`
	CommonPrefixes  []string        `json:"commonPrefixes" param:"-" header:"-"`
}

// ListObjects list all objects with Prefix and Delimiter
func (client *Client) ListObjects(request *ListObjectsRequest) (*ObjectListing, error) {
	return client.ListObjectsWithContext(context.Background(), request)
}

// ListObjectsWithContext list all objects with Prefix and Delimiter with context controlling
func (client *Client) ListObjectsWithContext(ctx context.Context, request *ListObjectsRequest) (*ObjectListing, error) {
	result := &ObjectListing{}
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

// ListObjectsNextBatch list next batch of ListObjects
func (client *Client) ListObjectsNextBatch(previous *ObjectListing) (*ObjectListing, error) {
	return client.ListObjectsNextBatchWithContext(context.Background(), previous)
}

// ListObjectsNextBatchWithContext list next batch of ListObjects with context controlling
func (client *Client) ListObjectsNextBatchWithContext(ctx context.Context, previous *ObjectListing) (*ObjectListing, error) {
	result := &ObjectListing{}
	req := &clientRequest{
		BucketName:         previous.BucketName,
		Method:             HTTPGet,
		QueryHeaderOptions: previous,
		Result:             result,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return result, err
}

type initMultipartUploadOption struct {
	Uploads string `param:"uploads" header:"-"`
}

// InitMultipartUploadRequest is input of InitMultipartUpload
type InitMultipartUploadRequest struct {
	initMultipartUploadOption
	BucketName string `param:"-" header:"-"`
	ObjectName string `param:"-" header:"-"`

	CacheControl       string `header:"Cache-Control,omitempty" param:"-"`
	ContentDisposition string `header:"Content-Disposition,omitempty" param:"-"`
	ContentEncoding    string `header:"Content-Encoding,omitempty" param:"-"`
	ContentType        string `header:"Content-Type,omitempty" param:"-"`
	ContentLength      int    `header:"Content-Length,omitempty" param:"-"`
	Expect             string `header:"Expect,omitempty" param:"-"`
}

// InitMultipartUploadResponse is result of InitMultipartUpload
type InitMultipartUploadResponse struct {
	BucketName string `json:"bucketName" param:"-" header:"-"`
	ObjectName string `json:"objectName" param:"-" header:"-"`
	UploadID   string `json:"uploadId" param:"uploadId" header:"-"`
}

// InitMultipartUpload starts a progress of multipart uploading
func (client *Client) InitMultipartUpload(request *InitMultipartUploadRequest) (*InitMultipartUploadResponse, error) {
	return client.InitMultipartUploadWithContext(context.Background(), request)
}

// InitMultipartUploadWithContext starts a progress of multipart uploading with context controlling
func (client *Client) InitMultipartUploadWithContext(ctx context.Context, request *InitMultipartUploadRequest) (*InitMultipartUploadResponse, error) {
	result := &InitMultipartUploadResponse{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPPut,
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

// UploadPartRequest is input of UploadPart
type UploadPartRequest struct {
	BucketName string    `param:"-" header:"-"`
	ObjectName string    `param:"-" header:"-"`
	UploadID   string    `param:"uploadId" header:"-"`
	PartNumber int       `param:"partNumber" header:"-"`
	Data       io.Reader `param:"-" header:"-"`
}

// UploadPartResponse is result of UploadPart
type UploadPartResponse struct {
	PartNumber int    `json:"partNumber"`
	ETag       string `json:"etag"`
	PartSize   int64  `json:"partSize"`
}

// UploadPart upload part of multipart uploading
func (client *Client) UploadPart(request *UploadPartRequest) (*UploadPartResponse, error) {
	return client.UploadPartWithContext(context.Background(), request)
}

// UploadPartWithContext upload part of multipart uploading with context controlling
func (client *Client) UploadPartWithContext(ctx context.Context, request *UploadPartRequest) (*UploadPartResponse, error) {
	result := &UploadPartResponse{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPPut,
		Data:               request.Data,
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

// UploadPartList is required by CompleteMultipartUpload
type UploadPartList struct {
	UploadPartResultList []UploadPartResponse `json:"uploadPartResultList"`
}

// CompleteMultipartUpload completes the progress of multipart uploading
func (client *Client) CompleteMultipartUpload(request *InitMultipartUploadResponse, list *UploadPartList) (*PutObjectResponse, error) {
	return client.CompleteMultipartUploadWithContext(context.Background(), request, list)
}

// CompleteMultipartUploadWithContext completes the progress of multipart uploading with context controlling
func (client *Client) CompleteMultipartUploadWithContext(ctx context.Context, request *InitMultipartUploadResponse, list *UploadPartList) (*PutObjectResponse, error) {
	result := &PutObjectResponse{}
	data, e := json.Marshal(*list)
	if e != nil {
		return result, e
	}

	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPPut,
		Data:               bytes.NewReader(data),
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

// AbortMultipartUpload aborts the progress of multipart uploading
func (client *Client) AbortMultipartUpload(request *InitMultipartUploadResponse) error {
	return client.AbortMultipartUploadWithContext(context.Background(), request)
}

// AbortMultipartUploadWithContext aborts the progress of multipart uploading with context controlling
func (client *Client) AbortMultipartUploadWithContext(ctx context.Context, request *InitMultipartUploadResponse) error {
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPDelete,
		QueryHeaderOptions: request,
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

type restoreObjectOption struct {
	Restore string `param:"restore" header:"-"`
}

// RestoreObject restore object which is deleted if this object is avaliable
func (client *Client) RestoreObject(bucketName, objectName string) error {
	return client.RestoreObjectWithContext(context.Background(), bucketName, objectName)
}

// RestoreObjectWithContext restore object which is deleted if this object is avaliable with context controlling
func (client *Client) RestoreObjectWithContext(ctx context.Context, bucketName, objectName string) error {
	req := &clientRequest{
		BucketName:         bucketName,
		ObjectName:         objectName,
		Method:             HTTPPut,
		QueryHeaderOptions: restoreObjectOption{},
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

// GenerateAbsoluteObjectURL generates a absoluted url
func (client *Client) GenerateAbsoluteObjectURL(bucketName, objectName string) *url.URL {
	return client.buildRequestURL(bucketName, objectName, "", false)
}

// GeneratePresignedURLRequest is input of GeneratePresignedURL
type GeneratePresignedURLRequest struct {
	CDN        bool
	BucketName string
	ObjectName string
	Method     HTTPMethod
	Expiration time.Time
	Metadata   *ObjectMetadata
}

// GeneratePresignedURL generates presigned url
func (client *Client) GeneratePresignedURL(request *GeneratePresignedURLRequest) (*url.URL, error) {
	baseURL := client.buildRequestURL(request.BucketName, request.ObjectName, "", request.CDN)

	params := url.Values{}
	if request.Method == HTTPHead {
		params.Add("metadata", "")
	}

	params.Add(HTTPHeaderGalaxyAccessKeyID, client.AccessID)
	params.Add(HTTPHeaderExpires, fmt.Sprintf("%d", request.Expiration.UnixNano()/int64(time.Millisecond)))
	baseURL.RawQuery = params.Encode()

	header := http.Header{}
	for k,v := range request.Metadata.metadata {
		header.Set(k, v)
	}

	sig, e := signature(client.AccessSecret, request.Method, baseURL.String(), header)
	if e != nil {
		return nil, e
	}

	return url.Parse(baseURL.String() + "&" + HTTPHeaderSignature + "=" + sig)
}

// GetObjectACLRequest is input of GetObjectACL
type GetObjectACLRequest struct {
	aclOption
	BucketName string `param:"-" header:"-"`
	ObjectName string `param:"-" header:"-"`
	VersionID  string `param:"versionId,omitempty" header:"-"`
}

// GetObjectACL will return AccessControlList of object
func (client *Client) GetObjectACL(request *GetObjectACLRequest) (*AccessControlList, error) {
	return client.GetObjectACLWithContext(context.Background(), request)
}

// GetObjectACLWithContext will return AccessControlList of object with context controlling
func (client *Client) GetObjectACLWithContext(ctx context.Context, request *GetObjectACLRequest) (*AccessControlList, error) {
	result := &AccessControlList{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
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

// SetObjectACLRequest is input of SetObjectACL
type SetObjectACLRequest struct {
	aclOption
	BucketName string             `param:"-" header:"-"`
	ObjectName string             `param:"-" header:"-"`
	VersionID  string             `param:"versionId,omitempty" header:"-"`
	ACL        *AccessControlList `param:"-" header:"-"`
}

// SetObjectACL sets AccessControlList for object
func (client *Client) SetObjectACL(request *SetObjectACLRequest) error {
	return client.SetObjectACLWithContext(context.Background(), request)
}

// SetObjectACLWithContext sets AccessControlList for object with context controlling
func (client *Client) SetObjectACLWithContext(ctx context.Context, request *SetObjectACLRequest) error {
	aclBytes, e := json.Marshal(request.ACL)
	if e != nil {
		return errors.New("fds client: can't marshal acl")
	}

	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPPut,
		QueryHeaderOptions: request,
		Data:               bytes.NewReader(aclBytes),
	}

	resp, err := client.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return err
}

// SetObjectPublic is a shortcut of setting object public
func (client *Client) SetObjectPublic(bucketName, objectName string) error {
	return client.SetObjectPublicWithContext(context.Background(), bucketName, objectName)
}

// SetObjectPublicWithContext is a shortcut of setting object public with context controlling
func (client *Client) SetObjectPublicWithContext(ctx context.Context, bucketName, objectName string) error {
	grant := Grant{
		Grantee: GrantKey{
			ID: "ALL_USERS",
		},
		Permission: GrantPermissionRead,
		Type:       GrantTypeGroup,
	}

	controlList := &AccessControlList{}
	controlList.AddGrant(grant)

	aclRequest := &SetObjectACLRequest{
		BucketName: bucketName,
		ObjectName: objectName,
		ACL:        controlList,
	}

	return client.SetObjectACLWithContext(ctx, aclRequest)
}
