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
	if err != nil {
		if _, ok := err.(DriveFileDoesNotExistError); !ok {
			return nil, err
		}
	}

	if existFile == nil {
		file.MimeType = MimeGSuiteDoc
		if file.Parents == nil {
			folderID, err := getRootFolderID(r)
			if err != nil {
				return nil, err
			}
			file.Parents = []string{folderID}
		}
	}

	contentType := googleapi.ContentType(mimeTxt)
	var (
		newFile      *drive.File
		refreshToken bool
	)

	n := 1
retry:
	payloadReader := bytes.NewReader(*payload)
	if existFile == nil {
		<-tokenBucketGDriveAPI
		newFile, err = service.Files.Create(file).Media(payloadReader, contentType).Do()
	} else {
		<-tokenBucketGDriveAPI
		newFile, err = service.Files.Update(existFile.Id, existFile).Media(payloadReader, contentType).Do()
	}

	if err == nil {
		return newFile, nil
	}
	refreshToken, n, err = Triable(n, err)
	if err != nil {
		return nil, err
	}
	if refreshToken {
		service, err = GetGDriveService(r)
		if err != nil {
			return nil, err
		}
	}
	goto retry
}
