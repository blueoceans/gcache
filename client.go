package gcache

import (
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

func createGDriveClient(
	r *http.Request,
) *http.Client {
	ctx := newContext(r)

	if oauth2TokenSource == nil {
		oauth2TokenSource = (&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			Scopes:       []string{drive.DriveFileScope},
		}).TokenSource(
			*ctx,
			&oauth2.Token{
				RefreshToken: refreshToken,
			},
		)
	}

	return &http.Client{
		Transport: createTransport(ctx, &oauth2TokenSource),
	}
}
