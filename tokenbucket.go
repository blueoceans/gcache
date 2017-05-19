package gcache

import (
	"time"

	"github.com/blueoceans/go-common/tokenbucket"
)

var (
	tokenBucketGDriveAPI chan byte
)

func init() {
	tokenBucketGDriveAPI = tokenbucket.NewTokenBucket(time.Duration(100)*time.Second, 1000) // 1000/100sec (user limit)
}

// GetTokenBucketGDriveAPI returns a token-bucket for calling the Google Drive API.
func GetTokenBucketGDriveAPI() chan byte {
	return tokenBucketGDriveAPI
}
