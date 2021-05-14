package fds

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/XiaoMi/go-fds/fds/httpparser"
	"github.com/sirupsen/logrus"
)

// Client supplies an interface for interaction with FDS
type Client struct {
	logger     *logrus.Logger
	httpClient *http.Client

	Configuration *ClientConfiguration
	AccessID      string
	AccessSecret  string
}

// New a FDSClient
func New(accessID, accessSecret string, conf *ClientConfiguration) *Client {
	client := &Client{}
	client.Configuration = conf
	client.AccessID = accessID
	client.AccessSecret = accessSecret
	client.httpClient = &http.Client{}
	client.logger = logrus.New()

	client.logger.SetLevel(logrus.WarnLevel)

	return client
}

type clientRequest struct {
	BucketName         string
	ObjectName         string
	Method             HTTPMethod
	QueryHeaderOptions interface{}
	Metadata           *ObjectMetadata
	Data               io.Reader
	Result             interface{}
}

// make request
func (client *Client) do(ctx context.Context, request *clientRequest) (*http.Response, error) {
	// parse http url query string
	queryString, e := httpparser.QueryString(request.QueryHeaderOptions)
	if e != nil {
		return nil, e
	}

	query := queryString.Encode()

	u := client.buildRequestURL(request.BucketName, request.ObjectName, query, false)

	// parse http header
	header, e := httpparser.Header(request.QueryHeaderOptions)
	if e != nil {
		return nil, e
	}
	if request.Metadata != nil {
		for k, v := range request.Metadata.metadata {
			header.Add(k, v)
		}
	}

	return client.doRequest(ctx, request.Method, u, header, request.Data, request.Result)
}

func (client *Client) doRequest(ctx context.Context, method HTTPMethod, url *url.URL, header http.Header,
	data io.Reader, result interface{}) (*http.Response, error) {
	methodString := strings.ToUpper(string(method))
	req := &http.Request{
		Method:     methodString,
		URL:        url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       url.Host,
	}

	// inject context
	req = req.WithContext(ctx)

	dataFile := client.doHandleRequestBody(req, data)
	if dataFile != nil {
		defer func() {
			dataFile.Close()
			os.Remove(dataFile.Name())
		}()
	}

	if header != nil {
		for k := range header {
			req.Header.Set(k, header.Get(k))
		}
	}

	data = dataFile

	//req.Header.Add(HTTPHeaderContentMD5, "")
	req.Header.Set(HTTPHeaderDate, time.Now().Format(time.RFC1123))

	signature, err := signature(client.AccessSecret, method, url.String(), req.Header)
	if err != nil {
		return nil, err
	}
	req.Header.Set(HTTPHeaderAuthorization, fmt.Sprintf("Galaxy-V2 %s:%s", client.AccessID, signature))

	for k, v := range req.Header {
		client.logger.Debug(fmt.Sprintf(" >>> HTTP Header: k=%s, v=%s", k, v))
	}
	client.logger.Debug(fmt.Sprintf(" >>> HTTP URL: %s", req.URL.String()))

	response, err := client.httpClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		return nil, err
	}

	// check http status
	statusNeed2Check := []int{http.StatusOK}
	if method == HTTPHead {
		statusNeed2Check = append(statusNeed2Check, http.StatusNotFound)
	}
	err = checkResponseStatus(response, statusNeed2Check)
	if err != nil {
		return response, err
	}

	// unmarshal response body into result
	if result != nil {
		if w, ok := result.(io.Writer); ok {
			io.Copy(w, response.Body)
		} else {
			err = client.jsonResponseUnmarshal(response.Body, result)
		}
	}

	return response, err
}

func checkResponseStatus(response *http.Response, allowed []int) error {
	for _, v := range allowed {
		if response.StatusCode == v {
			return nil
		}
	}

	statusCode := response.StatusCode

	var err error
	var respBody []byte
	if statusCode >= 400 && statusCode <= 505 {
		respBody, err = readResponseBody(response)
		if err != nil {
			return err
		}

		err = newServerError(string(respBody), response.StatusCode)
		response.Body = ioutil.NopCloser(bytes.NewReader(respBody))
	} else if statusCode >= 300 && statusCode <= 307 {
		err = newServerError(fmt.Sprintf("fds: service returned %s", response.Status), response.StatusCode)
	}
	return err
}

func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	if err == io.EOF {
		err = nil
	}
	return out, err
}

func (client *Client) doHandleRequestBody(req *http.Request, body io.Reader) *os.File {
	var file *os.File
	switch v := body.(type) {
	case *bytes.Buffer:
		req.ContentLength = int64(v.Len())
	case *bytes.Reader:
		req.ContentLength = int64(v.Len())
	case *strings.Reader:
		req.ContentLength = int64(v.Len())
	case *os.File:
		fileInfo, _ := v.Stat()
		req.ContentLength = fileInfo.Size()
	case *io.LimitedReader:
		req.ContentLength = int64(v.N)
	}

	req.Header.Set(HTTPHeaderContentLength, strconv.FormatInt(req.ContentLength, 10))

	if body != nil && client.Configuration.EnableMd5Calculate && req.Header.Get(HTTPHeaderContentMD5) == "" {
		md5hash := ""
		md5hash, file, _ = calMD5(body)
		req.Header.Set(HTTPHeaderContentMD5, md5hash)
	}

	bodyCloser, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		bodyCloser = ioutil.NopCloser(body)
	}
	req.Body = bodyCloser

	return file
}

func (client *Client) buildRequestURL(bucketName string, objectName string, params string, cdn bool) *url.URL {
	var buf bytes.Buffer
	basicURL := client.basicURL(cdn)
	buf.WriteString(basicURL)

	objectName = url.QueryEscape(objectName)
	objectName = strings.Replace(objectName, "+", "%20", -1)

	buf.WriteByte('/')
	if bucketName != "" {
		buf.WriteString(bucketName)
	}
	if objectName != "" {
		buf.WriteByte('/')
		buf.WriteString(objectName)
	}

	if params != "" {
		buf.WriteByte('?')
		buf.WriteString(params)
	}

	u, _ := url.ParseRequestURI(buf.String())
	return u
}

func (client *Client) basicURL(cdn bool) string {
	var buf bytes.Buffer
	httpSchema := client.httpSchema()
	buf.WriteString(fmt.Sprintf("%s://", httpSchema))
	if cdn {
		buf.WriteString(client.Configuration.cdnEndpoint)
	} else {
		buf.WriteString(client.Configuration.Endpoint)
	}

	return buf.String()
}

func (client *Client) httpSchema() string {
	var httpSchema string
	if client.Configuration.EnableHTTPS {
		httpSchema = "https"
	} else {
		httpSchema = "http"
	}
	return httpSchema
}

func (client *Client) jsonResponseUnmarshal(body io.Reader, v interface{}) error {
	data, e := ioutil.ReadAll(body)
	if e != nil {
		return e
	}

	if len(data) == 0 {
		return nil
	}

	e = json.Unmarshal(data, v)
	client.logger.Debug(fmt.Sprintf(" <<< client Response: %s", v))
	return e
}

func calMD5(data io.Reader) (string, *os.File, error) {
	tempFile, err := ioutil.TempFile(os.TempDir(), TempFilePrefix)
	if tempFile != nil {
		md5hash := md5.New()
		writer := io.MultiWriter(md5hash, tempFile)
		io.Copy(writer, data)
		sum := md5hash.Sum(nil)

		return hex.EncodeToString(sum), tempFile, nil
	}
	return "", tempFile, err
}
