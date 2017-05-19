package gcache

import (
	"time"

	"github.com/blueoceans/go-common/tokenbucket"
)

var (
	tokenBucketGDriveAPI chan tokenbucket.Token
)

func init() {
	tokenBucketGDriveAPI = tokenbucket.NewTokenBucket(time.Duration(100)*time.Second, 1000) // 1000/100sec (user limit)
}
