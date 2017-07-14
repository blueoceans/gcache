// +build appengine

package gcache

import (
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	newappengine "google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

func createTransport(
	r *http.Request,
) *oauth2.Transport {
	ctx := newappengine.NewContext(r)

	if oauth2TokenSource == nil {
		oauth2TokenSource = google.AppEngineTokenSource(ctx, drive.DriveFileScope)
	}

	return &oauth2.Transport{
		Source: oauth2TokenSource,
		Base: &urlfetch.Transport{
			Context: ctx,
		},
	}
}
