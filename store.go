package gcache

import (
	"bytes"
	"errors"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// StoreGDrive stores a file that is a given filename on Google Drive.
func StoreGDrive(
	r *http.Request,
	filename string,
	file *drive.File,
	payload *[]byte,
) (
	*drive.File,
	error,
) {

	if filename == "" {
		return nil, errors.New("`filename` must be enough")
	}

	existFile, service, err := getGDriveFile(r, filename, "")
	switch err.(type) {
	case nil:
	case *DriveFileDoesNotExistError:
		file.Name = filename
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
