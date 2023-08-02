package manager

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/XiaoMi/go-fds/fds"
	"github.com/XiaoMi/go-fds/fds/httpparser"
	"github.com/sirupsen/logrus"
)

// Downloader is a FDS client for file concurrency download
type Downloader struct {
	logger *logrus.Logger
	client *fds.Client

	PartSize    int64
	Concurrency int
	Breakpoint  bool
}

// NewDownloader new a downloader
func NewDownloader(client *fds.Client, partSize int64, concurrency int, breakpoint bool) (*Downloader, error) {
	if partSize < 1 {
		return nil, ErrorPartSizeSmallerThanOne
	}

	if concurrency < 1 {
		return nil, ErrorConcurrencySmallerThanOne
	}

	downloader := &Downloader{
		PartSize:    partSize,
		Concurrency: concurrency,
		Breakpoint:  breakpoint,

		client: client,
	}
	downloader.logger = logrus.New()
	downloader.logger.SetLevel(logrus.WarnLevel)

	return downloader, nil
}

// DownloadRequest is the input of Download
type DownloadRequest struct {
	fds.GetObjectRequest
	FilePath string

	// private
	breakpointFilePath string
}

// Download performs the downloading action
func (downloader *Downloader) Download(request *DownloadRequest) error {
	return downloader.DownloadWithContext(context.Background(), request)
}

// DownloadWithContext performs the downloading action with context controlling
func (downloader *Downloader) DownloadWithContext(ctx context.Context, request *DownloadRequest) error {
	if downloader.Breakpoint && request.breakpointFilePath != "" {
		request.breakpointFilePath = fmt.Sprintf("%s.download.bp", request.FilePath)
	}

	var parts []part
	var err error

	metadata, err := downloader.client.GetObjectMetadataWithContext(ctx, request.BucketName, request.ObjectName)
	if err != nil {
		return err
	}

	contentLength, err := strconv.ParseInt(metadata.Get(fds.HTTPHeaderContentMetadataLength), 10, 0)
	if err != nil {
		return err
	}

	ranges, err := httpparser.Range(request.Range)
	if err != nil {
		return err
	}

	if len(ranges) == 0 {
		ranges = append(ranges, httpparser.HTTPRange{End: contentLength - 1})
	}

	if len(ranges) > 1 {
		return ErrorRnageFormat
	}

	start := ranges[0].Start
	end := ranges[0].End + 1
	if ranges[0].Start < 0 || ranges[0].Start >= contentLength || ranges[0].End > contentLength || ranges[0].Start > ranges[0].End {
		start = 0
		end = contentLength
	}
	r := httpparser.HTTPRange{
		Start: start,
		End:   end,
	}

	bp := breakpointInfo{
		downloader: downloader,
	}
	if downloader.Breakpoint {
		// load breakpoint info
		err = bp.Load(request.breakpointFilePath)
		if err != nil {
			bp.Destroy()
		}

		// validate breakpoint info
		err = bp.Validate(request.BucketName, request.ObjectName, r)
		if err != nil {
			downloader.logger.Debug(err)
			downloader.logger.Debug("breakpoint info is invalid")
			bp.Initilize(downloader, request.BucketName, request.ObjectName, request.breakpointFilePath, r, metadata)
			bp.Destroy()
		}

		// get parts from breakpoint info
		parts = bp.UnfinishParts()
	} else {
		parts, err = downloader.splitDownloadParts(r)
		if err != nil {
			return err
		}
	}

	jobs := make(chan part, len(parts))
	results := make(chan part, len(parts))
	failed := make(chan error)
	finished := make(chan bool)

	tmpFilePath := request.FilePath + ".tmp"
	for i := 1; i < downloader.Concurrency; i++ {
		go downloader.downloaderTaskConsumer(ctx, i, request, tmpFilePath, jobs, results, failed, finished)
	}

	go downloader.downloaderTaskProducer(jobs, parts)

	completed := 0
	for completed < len(parts) {
		select {
		case p := <-results:
			completed++
			if downloader.Breakpoint {
				bp.PartStat[p.Index] = true
				bp.Dump()
			}
		case err := <-failed:
			close(finished)
			return err
		}
	}

	if downloader.Breakpoint {
		os.Remove(request.breakpointFilePath)
	}
	return os.Rename(tmpFilePath, request.FilePath)
}

func (downloader *Downloader) downloaderTaskConsumer(ctx context.Context, id int,
	request *DownloadRequest, tmpFilePath string, jobs <-chan part, results chan<- part, failed chan<- error, finished <-chan bool) {
	for p := range jobs {
		req := &fds.GetObjectRequest{
			BucketName: request.BucketName,
			ObjectName: request.ObjectName,
			Range:      fmt.Sprintf("bytes=%v-%v", p.Start, p.End),
		}

		data, err := downloader.client.GetObjectWithContext(ctx, req)
		if err != nil {
			downloader.logger.Debug(err.Error())
			failed <- err
			break
		}
		defer data.Close()

		select {
		case <-finished:
			return
		default:
		}

		fd, err := os.OpenFile(tmpFilePath, os.O_WRONLY|os.O_CREATE, os.FileMode(0664))
		if err != nil {
			failed <- err
			break
		}

		_, err = fd.Seek(p.Start-p.Offset, io.SeekStart)
		if err != nil {
			fd.Close()
			failed <- err
			break
		}

		_, err = io.Copy(fd, data)
		if err != nil {
			fd.Close()
			failed <- err
			break
		}

		fd.Close()
		results <- p
	}
}

func (downloader *Downloader) downloaderTaskProducer(jobs chan part, parts []part) {
	for _, p := range parts {
		jobs <- p
	}
	close(jobs)
}

type part struct {
	Index  int
	Start  int64
	End    int64
	Offset int64
}

func (downloader Downloader) splitDownloadParts(r httpparser.HTTPRange) ([]part, error) {
	var parts []part

	i := 0
	for offset := r.Start; offset < r.End; offset += downloader.PartSize {
		p := part{
			Index:  i,
			Start:  offset,
			End:    getEnd(offset, r.End, downloader.PartSize),
			Offset: r.Start,
		}
		i++
		parts = append(parts, p)
	}

	return parts, nil
}

func getEnd(begin int64, total int64, per int64) int64 {
	if begin+per > total {
		return total - 1
	}
	return begin + per - 1
}

type breakpointInfo struct {
	FilePath   string
	BucketName string
	ObjectName string
	ObjectStat objectStat
	Parts      []part
	PartStat   []bool
	Start      int64
	End        int64
	MD5        string

	downloader *Downloader
}

type objectStat struct {
	Size         int64  // Object size
	LastModified string // Last modified time
}

func (bp *breakpointInfo) Load(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, bp)
}

func (bp *breakpointInfo) Dump() error {
	bpi := *bp

	bpi.MD5 = ""
	data, err := json.Marshal(bpi)
	if err != nil {
		return err
	}

	sum := md5.Sum(data)
	b64 := base64.StdEncoding.EncodeToString(sum[:])
	bpi.MD5 = b64

	data, err = json.Marshal(bpi)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(bpi.FilePath, data, os.FileMode(0664))
}

func (bp *breakpointInfo) Validate(bucketName, objectName string, r httpparser.HTTPRange) error {
	if bucketName != bp.BucketName || objectName != bp.ObjectName {
		return ErrorBucketOrObjectNotMatching
	}

	bpi := *bp
	bpi.MD5 = ""
	data, err := json.Marshal(bpi)
	if err != nil {
		return err
	}
	sum := md5.Sum(data)
	b64 := base64.StdEncoding.EncodeToString(sum[:])
	if b64 != bp.MD5 {
		return ErrorMD5NotMatching
	}

	c := bp.downloader.client
	metadata, err := c.GetObjectMetadata(bucketName, objectName)
	if err != nil {
		return err
	}

	length, err := metadata.GetContentLength()
	if err != nil {
		return err
	}
	if bp.ObjectStat.Size != length || bp.ObjectStat.LastModified != metadata.Get(fds.HTTPHeaderLastModified) {
		return ErrorObjectStateNotMatching
	}

	if bp.Start != r.Start || bp.End != r.End {
		return ErrorRangeNotMatching
	}

	return nil
}

func (bp *breakpointInfo) UnfinishParts() []part {
	var result []part

	for i, s := range bp.PartStat {
		if !s {
			result = append(result, bp.Parts[i])
		}
	}

	return result
}

func (bp *breakpointInfo) Initilize(downloader *Downloader,
	bucketName, objectName, filePath string, r httpparser.HTTPRange, md *fds.ObjectMetadata) error {
	bp.MD5 = ""
	bp.BucketName = bucketName
	bp.ObjectName = objectName
	bp.FilePath = filePath
	bp.Start = r.Start
	bp.End = r.End
	bp.downloader = downloader

	contentLength, err := md.GetContentLength()
	if err != nil {
		return err
	}

	parts, err := downloader.splitDownloadParts(r)
	if err != nil {
		return err
	}
	bp.Parts = parts

	bp.PartStat = make([]bool, len(bp.Parts))

	bp.ObjectStat = objectStat{
		Size:         contentLength,
		LastModified: md.Get(fds.HTTPHeaderLastModified),
	}

	return nil
}

func (bp *breakpointInfo) Destroy() {
	os.Remove(bp.FilePath)
}
