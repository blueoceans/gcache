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

// StoreGDrive stores a file that is a given file ID or file name on Google Drive.
func StoreGDrive(
	r *http.Request,
	file *drive.File,
	payload *[]byte,
) (
	*drive.File,
	error,
) {

	if file.Id == "" && file.Name == "" {
		return nil, errors.New("`file.Id` or `file.Name` must be enough")
	}

	var (
		err       error
		existFile *drive.File
		service   *drive.Service
	)
	if file.Id == "" {
		existFile, service, err = getGDriveFile(r, file.Name, "")
		if err != nil {
			if _, ok := err.(*DriveFileDoesNotExistError); !ok {
				return nil, err
			}
		}
	}
	if service == nil {
		service, err = GetGDriveService(r)
		if err != nil {
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
		newFile    *drive.File
		clearToken bool
	)

	n := 1
retry:
	payloadReader := bytes.NewReader(*payload)
	if existFile == nil {
		<-tokenBucketGDriveAPI
		newFile, err = service.Files.Create(file).Media(payloadReader, contentType).Do()
	} else {
		<-tokenBucketGDriveAPI
		newFile, err = service.Files.Update(existFile.Id, file).Media(payloadReader, contentType).Do()
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
