package gcache

import (
	"net/http"
)

func createGDriveClient(
	r *http.Request,
) *http.Client {
	return &http.Client{
		Transport: createTransport(r),
	}
}
