// +build !appengine

package gcache

import (
	"net/http"

	"golang.org/x/oauth2"
)

const ()

func createGDriveClient(
	_ interface{},
) *http.Client {
	return &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2TokenSource,
		},
	}
}
