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

	existFile, service, err := getGDriveFile(r, file.Name, "")
	switch err.(type) {
	case nil:
	case *DriveFileDoesNotExistError:
		file.MimeType = MimeGSuiteDoc
		if file.Parents == nil {
			folderID, err := getRootFolderID(r)
			if err != nil {
				return nil, err
			}
			file.Parents = []string{folderID}
		}
	default:
		return nil, err
	}

	contentType := googleapi.ContentType(mimeTxt)
	var (
		newFile    *drive.File
		clearToken bool
	)

	n := 1
retry:
	payloadReader := bytes.NewReader(*payload)
	<-tokenBucketGDriveAPI
	if existFile == nil {
		newFile, err = service.Files.Create(file).Media(payloadReader, contentType).Do()
	} else {
		newFile, err = service.Files.Update(existFile.Id, existFile).Media(payloadReader, contentType).Do()
	}

	if err == nil {
		return newFile, nil
	}
	clearToken, n, err = Triable(n, err)
	if err != nil {
		return nil, err
	}
	if clearToken {
		service, err = GetGDriveService(r)
		if err != nil {
			return nil, err
		}
	}
	goto retry
}
