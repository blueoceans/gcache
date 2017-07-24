package gcache

import (
	"bytes"
	"errors"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// StoreGDriveByID stores a file that is a given id on Google Drive.
func StoreGDriveByID(
	r *http.Request,
	id string,
	file *drive.File,
	payload *[]byte,
) (
	*drive.File,
	error,
) {

	if id == "" {
		return nil, errors.New("`id` must be enough")
	}

	contentType := googleapi.ContentType(mimeTxt)
	var (
		newFile    *drive.File
		clearToken bool
	)

	n := 1
refresh:
	service, err := GetGDriveService(r)
	if err != nil {
		return nil, err
	}

retry:
	payloadReader := bytes.NewReader(*payload)
	<-tokenBucketGDriveAPI
	newFile, err = service.Files.Update(id, file).Media(payloadReader, contentType).Do()

	if err == nil {
		return newFile, nil
	}
	clearToken, n, err = Triable(n, err)
	if err != nil {
		return nil, err
	}
	if clearToken {
		goto refresh
	}
	goto retry
}
