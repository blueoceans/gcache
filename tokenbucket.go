package gcache

import (
	"github.com/blueoceans/go-common/tokenbucket"
)

var (
	tokenBucketGDriveAPI chan struct{}
)

func init() {
	// https://developers.google.com/drive/v3/web/handle-errors#403_user_rate_limit_exceeded
	tokenBucketGDriveAPI = tokenbucket.NewTokenBucket(msec100, 1) // 1000/100sec (userRateLimitExceeded)
}

// GetTokenBucketGDriveAPI returns a token-bucket for calling the Google Drive API.
func GetTokenBucketGDriveAPI() chan struct{} {
	return tokenBucketGDriveAPI
}
