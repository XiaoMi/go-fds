package fds

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetObjectRequest is the input of GetObject method
type GetObjectRequest struct {
	BucketName string `param:"-" header:"-"`
	ObjectName string `param:"-" header:"-"`
	Range      string `param:"-" header:"Range,omitempty"`
}

// GetObject will get full content of object
func (client *Client) GetObject(request *GetObjectRequest) (io.ReadCloser, error) {
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		QueryHeaderOptions: request,
		Method:             HTTPGet,
	}

	resp, err := client.do(req)
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
	Expect             string `header:"Expect,omitempty" param:"-"`
	Expires            string `header:"Expires,omitempty" param:"-"`
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

// PutObject will create a object
func (client *Client) PutObject(request *PutObjectRequest) (*PutObjectResponse, error) {
	result := &PutObjectResponse{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Data:               request.Data,
		QueryHeaderOptions: request,
		Method:             HTTPPut,
		Result:             result,
	}

	resp, err := client.do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	return result, nil
}

// DoesObjectExist judge wether a object exists
func (client *Client) DoesObjectExist(bucketName, objectName string) (bool, error) {
	req := &clientRequest{
		BucketName: bucketName,
		ObjectName: objectName,
		Method:     HTTPHead,
	}

	resp, err := client.do(req)
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

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err
}

type renameObjectOption struct {
	RenameTo string `param:"renameTo" header:"-"`
}

// RenameObject renames sourceObjectName in bucketName to targetObjectName
func (client *Client) RenameObject(bucketName, sourceObjectName, targetObjectName string) error {
	req := &clientRequest{
		BucketName:         bucketName,
		ObjectName:         sourceObjectName,
		QueryHeaderOptions: renameObjectOption{targetObjectName},
		Method:             HTTPPut,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err
}

// DeleteObject deletes objectName in bucketName
func (client *Client) DeleteObject(bucketName, objectName string) error {
	req := &clientRequest{
		BucketName: bucketName,
		ObjectName: objectName,
		Method:     HTTPDelete,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err
}

type deleteObjectsOption struct {
	EnableTrash   bool   `param:"enableTrash" header:"-"`
	DeleteObjects string `param:"deleteObjects" header:"-"`
}

// DeleteObjects will delete all objects in objectNames
func (client *Client) DeleteObjects(bucketName string, objectNames []string, put2trash bool) error {
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

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err
}

// DeleteObjectsWithPrefix will delete all objects with prefix of prefix
func (client *Client) DeleteObjectsWithPrefix(bucketName, prefix string, put2stash bool) error {
	listObjectsRequest := &ListObjectsRequest{
		BucketName: bucketName,
		Prefix:     prefix,
		Delimiter:  "",
		MaxKeys:    DefaultListObjectsMaxKeys,
	}

	objectListing, err := client.ListObjects(listObjectsRequest)
	if err != nil {
		return err
	}

	for {
		names := make([]string, 0)
		for _, o := range objectListing.ObjectSummaries {
			names = append(names, o.ObjectName)
		}

		err = client.DeleteObjects(bucketName, names, put2stash)
		if err != nil {
			return err
		}

		if !objectListing.Truncated {
			break
		}

		objectListing, err = client.ListObjectsNextBatch(objectListing)
		if err != nil {
			return err
		}
	}

	return nil
}

// ObjectMetadata is metadata of object
type ObjectMetadata struct {
	h http.Header
}

// NewObjectMetadata create a default ObjectMetadata
func NewObjectMetadata() *ObjectMetadata {
	return &ObjectMetadata{http.Header{}}
}

// Get method of ObjectMetadata
func (metadata *ObjectMetadata) Get(k string) string {
	return metadata.h.Get(k)
}

// Set method of ObjectMetadata
func (metadata *ObjectMetadata) Set(k, v string) {
	metadata.h.Set(k, v)
}

// GetContentLength gets ContentLength of object metadata
func (metadata *ObjectMetadata) GetContentLength() (int64, error) {
	return strconv.ParseInt(metadata.Get(HTTPHeaderContentMetadataLength), 10, 64)
}

// SetContentLength sets ContentLength of object metadata
func (metadata *ObjectMetadata) SetContentLength(length int64) {
	metadata.Set(HTTPHeaderContentMetadataLength, strconv.FormatInt(length, 10))
}

func (metadata *ObjectMetadata) serialize() ([]byte, error) {
	x := make(map[string]string)
	for k := range metadata.h {
		x[k] = metadata.Get(k)
	}
	data := make(map[string]map[string]string)

	data["rawMeta"] = x
	result, e := json.Marshal(data)
	if e != nil {
		return nil, e
	}

	return result, nil
}

type getObjectMetadataOption struct {
	Metadata string `param:"metadata" header:"-"`
}

// GetObjectMetadata gets metadata of objectName in bucketName
func (client *Client) GetObjectMetadata(bucketName, objectName string) (*ObjectMetadata, error) {
	result := &ObjectMetadata{}
	req := &clientRequest{
		BucketName:         bucketName,
		ObjectName:         objectName,
		Method:             HTTPGet,
		QueryHeaderOptions: getObjectMetadataOption{},
	}

	resp, err := client.do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	result.h = resp.Header

	return result, nil
}

type setObjectMetadataOption struct {
	SetMetadata string `param:"setMetaData" header:""`
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

	resp, err := client.do(req)
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
	result := &ObjectListing{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		Method:             HTTPGet,
		QueryHeaderOptions: request,
		Result:             result,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return result, err
}

// ListObjectsNextBatch list next batch of ListObjects
func (client *Client) ListObjectsNextBatch(previous *ObjectListing) (*ObjectListing, error) {
	result := &ObjectListing{}
	req := &clientRequest{
		BucketName:         previous.BucketName,
		Method:             HTTPGet,
		QueryHeaderOptions: previous,
		Result:             result,
	}

	resp, err := client.do(req)
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
	Expires            string `header:"Expires,omitempty" param:"-"`
}

// InitMultipartUploadResponse is result of InitMultipartUpload
type InitMultipartUploadResponse struct {
	BucketName string `json:"bucketName" param:"-" header:"-"`
	ObjectName string `json:"objectName" param:"-" header:"-"`
	UploadID   string `json:"uploadId" param:"uploadId" header:"-"`
}

// InitMultipartUpload starts a progress of multipart uploading
func (client *Client) InitMultipartUpload(request *InitMultipartUploadRequest) (*InitMultipartUploadResponse, error) {
	result := &InitMultipartUploadResponse{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPPut,
		QueryHeaderOptions: request,
		Result:             result,
	}

	resp, err := client.do(req)
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
	result := &UploadPartResponse{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPPut,
		Data:               request.Data,
		QueryHeaderOptions: request,
		Result:             result,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return result, err
}

// UploadPartList is required by CompleteMultipartUpload
type UploadPartList struct {
	UploadPartResultList []UploadPartResponse `json:"uploadPartResultList"`
}

// CompleteMultipartUpload completes the progress of multipart uploading
func (client *Client) CompleteMultipartUpload(request *InitMultipartUploadResponse, list *UploadPartList) (*PutObjectResponse, error) {
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

	resp, err := client.do(req)
	defer resp.Body.Close()

	return result, err
}

// AbortMultipartUpload aborts the progress of multipart uploading
func (client *Client) AbortMultipartUpload(request *InitMultipartUploadResponse) error {
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPDelete,
		QueryHeaderOptions: request,
	}

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err
}

type restoreObjectOption struct {
	Restore string `param:"restore" header:"-"`
}

// RestoreObject restore object which is deleted if this object is avaliable
func (client *Client) RestoreObject(bucketName, objectName string) error {
	req := &clientRequest{
		BucketName:         bucketName,
		ObjectName:         objectName,
		Method:             HTTPPut,
		QueryHeaderOptions: restoreObjectOption{},
	}

	resp, err := client.do(req)
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

	sig, e := signature(client.AccessSecret, request.Method, baseURL.String(), request.Metadata.h)
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
	result := &AccessControlList{}
	req := &clientRequest{
		BucketName:         request.BucketName,
		ObjectName:         request.ObjectName,
		Method:             HTTPGet,
		QueryHeaderOptions: request,
		Result:             result,
	}

	resp, err := client.do(req)
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

	resp, err := client.do(req)
	defer resp.Body.Close()

	return err

}

// SetObjectPublic is a shortcut of setting object public
func (client *Client) SetObjectPublic(bucketName, objectName string) error {
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

	return client.SetObjectACL(aclRequest)
}
