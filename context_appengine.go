// +build appengine

package gcache

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	newappengine "google.golang.org/appengine"
)

const (
	deadline = time.Duration(60) * time.Second
)

func newContext(
	r *http.Request,
) *context.Context {
	ctx, _ := context.WithTimeout(newappengine.NewContext(r), deadline)
	return &ctx
}
