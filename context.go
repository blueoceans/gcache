// +build !appengine

package gcache

import (
	"golang.org/x/net/context"
)

func newContext(
	_ interface{},
) *context.Context {
	ctx := context.Background()
	return &ctx
}
