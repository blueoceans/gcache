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

	n := 1
retry:
	service, err := drive.New(createGDriveClient(r))
	if err != nil {
		if IsInvalidSecurityTicket(err) {
			oauth2TokenSource = nil
			goto retry
		} else if IsServerError(err) {
			n, err = sleeping(n)
			if err == nil {
				goto retry
			}
		}
		return nil, err
	}

	folderID, err := getRootFolderID(r)
	if err != nil {
		return nil, err
	}

	file.Parents = append(file.Parents, folderID)
	file.MimeType = mimeGSuiteDoc

	file, err = service.Files.Create(file).Media(bytes.NewReader(*payload), googleapi.ContentType(mimeTxt)).Do()
	if err != nil {
		if IsInvalidSecurityTicket(err) {
			oauth2TokenSource = nil
			goto retry
		} else if IsServerError(err) {
			n, err = sleeping(n)
			if err == nil {
				goto retry
			}
		}
		return nil, err
	}
	return file, nil
}
