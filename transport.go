// +build !appengine

package gcache

import (
	"golang.org/x/oauth2"
)

func createTransport(
	_ interface{},
	ts *oauth2.TokenSource,
) *oauth2.Transport {
	return &oauth2.Transport{
		Source: *ts,
	}
}
