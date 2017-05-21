package gcache

import (
	"errors"
	"testing"
)

var vtests = []struct {
	message string
	ok      bool
}{
	{"googleapi: Error 403: User Rate Limit Exceeded, userRateLimitExceeded",
		true},
}

func TestIsRateLimit(t *testing.T) {
	for _, vt := range vtests {
		if ok := IsRateLimit(errors.New(vt.message)); ok != vt.ok {
			t.Errorf("%q, want %q", ok, vt.ok)
		}
	}
}
