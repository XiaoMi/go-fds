package fds

import "os"

// Headers
const (
	XiaomiPrefix                = "x-xiaomi-"
	XiaomiMetaPrefix            = "x-xiaomi-meta-"
	HTTPHeaderGalaxyAccessKeyID = "GalaxyAccessKeyId"
	HTTPHeaderSignature         = "Signature"
	HTTPHeaderExpires           = "Expires"

	HTTPHeaderCacheControl          = "cache-control"
	HTTPHeaderContentLength         = "content-length"
	HTTPHeaderContentEncoding       = "content-encoding"
	HTTPHeaderLastModified          = "last-modified"
	HTTPHeaderContentMD5            = "content-md5"
	HTTPHeaderContentType           = "content-type"
	HTTPHeaderLastChecked           = "last-checked"
	HTTPHeaderUploadTime            = "upload-time"
	HTTPHeaderDate                  = "date"
	HTTPHeaderAuthorization         = "authorization"
	HTTPHeaderRange                 = "range"
	HTTPHeaderContentRange          = "content-range"
	HTTPHeaderContentMetadataLength = XiaomiMetaPrefix + HTTPHeaderContentLength
	HTTPHeaderServerSideEncryption  = XiaomiMetaPrefix + "server-side-encryption"
	HTTPHeaderStorageClass          = XiaomiMetaPrefix + "storage-class"
	HTTPHeaderOngoingRestore        = XiaomiMetaPrefix + "ongoing-restore"
	HTTPHeaderRestoreExpireDate     = XiaomiMetaPrefix + "restore-expiry"
	HTTPHeaderCRC64ECMA             = XiaomiMetaPrefix + "hash-crc64ecma"
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
	MinPartSize = 5 * 1024 * 1024        // Min part size, 5MB

	FilePermMode = os.FileMode(0664) // Default file permission

	TempFilePrefix = "fds-go-temp-" // Temp file prefix

	DefaultListObjectsMaxKeys = 1000

	URLComSuffix = ".fds.api.xiaomi.com"
	URLNetSuffix = "-fds.api.xiaomi.net"
	URLCDNSuffix = ".fds.api.mi-img.com"

	Version = "0.9.0" // Go SDK version
)

// Storage Class
type StorageClass string

const (
	Standard                 StorageClass = "STANDARD"
	StandardInfrequentAccess StorageClass = "STANDARD_IA"
	Archive                  StorageClass = "Archive"
)
