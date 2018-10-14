package manager

import (
	"fmt"
	"os"

	"github.com/hujianxin/galaxy-fds-sdk-go/fds"
	"github.com/sirupsen/logrus"
)

// Uploader manages your file uploading in concurrent way
type Uploader struct {
	client *fds.Client
	logger *logrus.Logger

	PartSize    int64
	Concurrency int
	Breakpoint  bool
}

// NewUploader new a Uploader
func NewUploader(client *fds.Client, partSize int64, concurrency int, breakpoint bool) (*Uploader, error) {
	if partSize < 1 {
		return nil, ErrorPartSizeSmallerThanOne
	}

	if concurrency < 1 {
		return nil, ErrorConcurrencySmallerThanOne
	}

	uploader := &Uploader{
		client:      client,
		PartSize:    partSize,
		Concurrency: concurrency,
		Breakpoint:  breakpoint,
	}

	uploader.logger = logrus.New()
	uploader.logger.SetLevel(logrus.DebugLevel)

	return uploader, nil
}

// UploadRequest is input of Upload
type UploadRequest struct {
	*fds.PutObjectRequest
	FilePath string

	breakpointFilePath string
}

// Upload file
func (uploader *Uploader) Upload(request *UploadRequest) error {
	if _, err := os.Stat(request.FilePath); os.IsNotExist(err) {
		return ErrorFileNotFound
	}

	if uploader.Breakpoint {
		request.breakpointFilePath = fmt.Sprintf("%s.upload.bp", request.FilePath)
	}

	initMultipartUploadRequest := &fds.InitMultipartUploadRequest{
		BucketName: request.BucketName,
		ObjectName: request.ObjectName,
	}
	var parts []part

	parts =

		uploader.client.InitMultipartUpload(initMultipartUploadRequest)

}

func (uploader *Uploader) splitUploadParts(request *UploadRequest) ([]part, error) {
	var parts []part

	file, err := os.Open(request.FilePath)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	splitCount := stat.Size() / uploader.PartSize
	if splitCount >= 10000 {
		return nil, ErrorTooManyUploadParts
	}

	for i := int64(0); i < splitCount; i++ {
		parts = append(parts, part{
			Index:  int(i) + 1,
			Start:  i * uploader.PartSize,
			End:    (i + 1) * uploader.PartSize,
			Offset: 0,
		})
	}

	if stat.Size()%splitCount != 0 {
		parts = append(parts, part{
			Index: len(parts) + 1,
			Start: len()
		})
	}
}
