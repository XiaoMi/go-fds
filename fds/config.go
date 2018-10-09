package fds

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

// HTTPTimeout defines HTTP timeout.
type HTTPTimeout struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	HeaderTimeout    time.Duration
	LongTimeout      time.Duration
	IdleConnTimeout  time.Duration
}

// ClientConfiguration required by FDSClient initialization
type ClientConfiguration struct {
	regionName  string
	cdnEndpoint string

	Endpoint               string
	EnableHTTPS            bool
	EnableCDNForUpload     bool
	EnableCDNForDownload   bool
	EnableMd5Calculate     bool
	Timeout                uint
	HTTPTimeout            HTTPTimeout
	MaxConnection          uint
	BatchDeleteSize        uint
	RetryCount             uint
	RetryInterval          uint
	PartSize               uint
	DownloadBandwidth      uint64
	UploadBandwidth        uint64
	HTTPKeepAliveTimeoutMs uint64
}

// NewClientConfiguration create a usable ClientConfiguration
func NewClientConfiguration(endpoint string) (*ClientConfiguration, error) {
	conf := defaultFDSClientConfiguration()
	conf.Endpoint = endpoint

	// parsing regionName name for configuration
	var urlSuffix string
	host := strings.Split(conf.Endpoint, ":")[0]
	if strings.HasSuffix(host, URLNetSuffix) {
		urlSuffix = URLNetSuffix
	} else if strings.HasSuffix(host, URLComSuffix) {
		urlSuffix = URLComSuffix
	} else {
		return conf, fmt.Errorf("[client] error endpoint")
	}
	conf.regionName = host[0 : len(host)-len(urlSuffix)]

	// parsing cdn endpoint for configuration
	var buf bytes.Buffer
	buf.WriteString("cdn.")
	buf.WriteString(conf.regionName)
	buf.WriteString(URLCDNSuffix)
	conf.cdnEndpoint = buf.String()

	return conf, nil
}

// CDNEndpoint is endpoint of cdn
func (conf *ClientConfiguration) CDNEndpoint() string {
	return conf.cdnEndpoint
}

// RegionName get region name
func (conf *ClientConfiguration) RegionName() string {
	return conf.regionName
}

func defaultFDSClientConfiguration() *ClientConfiguration {
	config := ClientConfiguration{}
	config.Endpoint = ""
	config.EnableHTTPS = true
	config.EnableCDNForUpload = false
	config.EnableCDNForDownload = true
	config.EnableMd5Calculate = false
	config.Timeout = 50
	config.HTTPTimeout.ConnectTimeout = time.Second * 30   // 30s
	config.HTTPTimeout.ReadWriteTimeout = time.Second * 60 // 60s
	config.HTTPTimeout.HeaderTimeout = time.Second * 60    // 60s
	config.HTTPTimeout.LongTimeout = time.Second * 300     // 300s
	config.HTTPTimeout.IdleConnTimeout = time.Second * 50  // 50s
	config.MaxConnection = 20
	config.BatchDeleteSize = 1000
	config.RetryCount = 3
	config.RetryInterval = 500 // ms
	config.PartSize = 10 * 1024 * 1024
	config.DownloadBandwidth = 10 * 1024 * 1024
	config.UploadBandwidth = 10 * 1024 * 1024

	return &config
}
