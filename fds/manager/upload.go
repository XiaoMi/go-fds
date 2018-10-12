package manager

import (
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

// UploadRequest is input of Upload
type UploadRequest struct {
	fds.PutObjectRequest
	FilePath           string
	BreakpointFilePath string
}

// Upload file
func (uploader *Uploader) Upload(request *UploadRequest) error {
	return nil
}
