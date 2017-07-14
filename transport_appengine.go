// +build appengine

package gcache

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/urlfetch"
)

func createTransport(
	ctx *context.Context,
	ts *oauth2.TokenSource,
) *oauth2.Transport {
	return &oauth2.Transport{
		Source: *ts,
		Base: &urlfetch.Transport{
			Context: *ctx,
		},
	}
}
