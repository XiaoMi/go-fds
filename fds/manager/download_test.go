package manager

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/XiaoMi/go-fds/fds"
	"golang.org/x/time/rate"
)

func TestDownloader_DownloaderTaskConsumer(t *testing.T) {
	// prepare
	// 1. client
	endpoint := os.Getenv("GO_FDS_TEST_ENDPOINT")
	accessID := os.Getenv("GO_FDS_TEST_ACCESS_KEY_ID")
	accessSecret := os.Getenv("GO_FDS_TEST_ACCESS_KEY_SECRET")

	conf, err := fds.NewClientConfiguration(endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	client := fds.New(accessID, accessSecret, conf)
	// 2. downloader
	downloader, err := NewDownloader(client, 16*1024*1024, 20, true)
	if err != nil {
		log.Fatalln(err)
	}
	// 3. limiter
	// intervel 200ms means 1s download 5 part = 80M
	// download 1024M should cost 12.8s
	limiter := rate.NewLimiter(rate.Every(200*time.Millisecond), 1)
	downloader.SetLimiter(limiter)

	// 4. request
	cwd, _ := os.Getwd()
	output := filepath.Join(cwd, "output")
	request := &DownloadRequest{
		GetObjectRequest: fds.GetObjectRequest{
			BucketName: "log-test",
			ObjectName: "fds-test-up",
		},
		FilePath: output,
	}

	err = downloader.Download(request)
	if err != nil {
		log.Fatal(err)
	}

}
