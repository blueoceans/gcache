package gcache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var (
	oauth2TokenSource oauth2.TokenSource // The token is valid for 30 minutes.
)

func sleeping(
	n int,
) (
	int,
	error,
) {
	if n > 16 {
		return 0, errors.New("Sleeping Timeout")
	}
	time.Sleep(time.Duration(n)*time.Second + time.Duration(random.Intn(1000))*time.Millisecond)
	return n * 2, nil
}

// GetGDriveService returns the API service of Google Drive.
func GetGDriveService(
	r *http.Request,
) (
	*drive.Service,
	error,
) {
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
	return service, nil
}

// GetGDriveFile returns a file on Google Drive.
func GetGDriveFile(
	r *http.Request,
	name string,
	field googleapi.Field,
) (
	*drive.File,
	[]byte, // payload
	error,
) {
	var payload []byte

retry:
	service, err := drive.New(createGDriveClient(r))
	if err != nil {
		if IsInvalidSecurityTicket(err) {
			oauth2TokenSource = nil
			goto retry
		}
		return nil, nil, err
	}

	fileList, err := service.Files.List().PageSize(1).Spaces("drive").Q(fmt.Sprintf("name='%s'", name)).Fields(field).Do()
	if err != nil {
		return nil, nil, err
	}
	if len(fileList.Files) <= 0 {
		return nil, nil, &DriveFileDoesNotExistError{}
	}
	file := fileList.Files[0]
	fileID := file.Id
	httpResponse, err := service.Files.Export(fileID, mimeTxt).Download()
	if IsFileNotExportableError(err) {
		err = service.Files.Delete(fileID).Do()
		if err != nil {
			// pass
		}
		return nil, nil, &DriveFileDoesNotExistError{}
	}
	if err != nil {
		return nil, nil, err
	}

	defer httpResponse.Body.Close()
	payload, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, nil, err
	}
	return file, payload, err
}
