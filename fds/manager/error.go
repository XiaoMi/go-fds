package manager

import "errors"

// Errors
var (
	ErrorPartSizeSmallerThanOne    = errors.New("PartSize can not be smaller than 1")
	ErrorConcurrencySmallerThanOne = errors.New("Concurrency can not be smaller than 1")
	ErrorRnageFormat               = errors.New("Does not support (bytes=i-j,m-n) format, only support (bytes=i-j)")
	ErrorBucketOrObjectNotMatching = errors.New("BucketName or ObjectName is not matching")
	ErrorMD5NotMatching            = errors.New("MD5 is not matching")
	ErrorObjectStateNotMatching    = errors.New("Object state is not matching")
	ErrorRangeNotMatching          = errors.New("Range is not matching")
	ErrorFileNotFound              = errors.New("File is not found")
	ErrorTooManyUploadParts        = errors.New("Too many upload parts, increase PartSize please")
)
