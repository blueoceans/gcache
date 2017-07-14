// +build !appengine

package gcache

import (
	"golang.org/x/oauth2"
)

func createTransport(
	_ interface{},
) *oauth2.Transport {
	return &oauth2.Transport{
		Source: oauth2TokenSource,
	}
}
