package gcache

import (
	"bytes"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var (
	random *rand.Rand
)

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// StoreGDrive stores a file to Google Drive.
func StoreGDrive(
	r *http.Request,
	file *drive.File,
	payload *[]byte,
) (
	*drive.File,
	error,
) {

	if file.Name == "" {
		return nil, errors.New("`file.Name` must be enough")
	}
	file.MimeType = mimeGSuiteDoc

	folderID, err := getRootFolderID(r)
	if err != nil {
		return nil, err
	}
	file.Parents = append(file.Parents, folderID)

	n := 1
refresh:
	service, err := GetGDriveService(r)
	if err != nil {
		return nil, err
	}
retry:
	file, err = service.Files.Create(file).Media(bytes.NewReader(*payload), googleapi.ContentType(mimeTxt)).Do()

	if err == nil {
		return file, nil
	}
	refreshToken, n, err := triable(n, err)
	if err != nil {
		return nil, err
	}
	if refreshToken {
		goto refresh
	}
	goto retry
}
