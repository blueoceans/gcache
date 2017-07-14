// +build appengine

package gcache

import (
	"net/http"

	"golang.org/x/net/context"
	newappengine "google.golang.org/appengine"
)

func newContext(
	r *http.Request,
) *context.Context {
	ctx := newappengine.NewContext(r)
	return &ctx
}
