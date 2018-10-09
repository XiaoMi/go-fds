package fds

import "os"

// Headers
const (
	XiaomiPrefix                    = "x-xiaomi-"
	XiaomiMetaPrefix                = "x-xiaomi-meta-"
	HTTPHeaderGalaxyAccessKeyID     = "GalaxyAccessKeyId"
	HTTPHeaderSignature             = "Signature"
	HTTPHeaderCacheControl          = "Cache-Control"
	HTTPHeaderContentLength         = "Content-Length"
	HTTPHeaderContentEncoding       = "Content-Encoding"
	HTTPHeaderLastModified          = "Last-Modified"
	HTTPHeaderContentMD5            = "Content-MD5"
	HTTPHeaderContentType           = "Content-Type"
	HTTPHeaderLastChecked           = "Last-Checked"
	HTTPHeaderUploadTime            = "Upload-Time"
	HTTPHeaderContentMetadataLength = XiaomiMetaPrefix + HTTPHeaderContentLength
	HTTPHeaderExpires               = "Expires"
	HTTPHeaderDate                  = "Date"
	HTTPHeaderAuthorization         = "Authorization"
	HTTPHeaderRange                 = "Range"
)

// HTTPMethod HTTP request method
type HTTPMethod string

const (
	// HTTPGet HTTP GET
	HTTPGet HTTPMethod = "GET"

	// HTTPPut HTTP PUT
	HTTPPut HTTPMethod = "PUT"

	// HTTPHead HTTP HEAD
	HTTPHead HTTPMethod = "HEAD"

	// HTTPPost HTTP POST
	HTTPPost HTTPMethod = "POST"

	// HTTPDelete HTTP DELETE
	HTTPDelete HTTPMethod = "DELETE"
)

// Other constants
const (
	MaxPartSize = 5 * 1024 * 1024 * 1024 // Max part size, 5GB
	MinPartSize = 100 * 1024             // Min part size, 100KB

	FilePermMode = os.FileMode(0664) // Default file permission

	TempFilePrefix = "fds-go-temp-" // Temp file prefix

	DefaultListObjectsMaxKeys = 1000

	URLComSuffix = ".fds.api.xiaomi.com"
	URLNetSuffix = "-fds.api.xiaomi.net"
	URLCDNSuffix = ".fds.api.mi-img.com"

	Version = "0.9.0" // Go SDK version
)
