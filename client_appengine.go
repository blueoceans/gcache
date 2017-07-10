// +build appengine

package gcache

import (
	"net/http"
	"time"

	"appengine"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	newappengine "google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

const (
	deadline = time.Duration(60) * time.Second
)

func createGDriveClient(
	r *http.Request,
) *http.Client {
	c := appengine.Timeout(appengine.NewContext(r), deadline)
	ctx := newappengine.NewContext(r)

	if oauth2TokenSource == nil {
		oauth2TokenSource = google.AppEngineTokenSource(ctx, drive.DriveFileScope)
	}

	return &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2TokenSource,
			Base: &urlfetch.Transport{
				Context: ctx,
			},
		},
	}
}
