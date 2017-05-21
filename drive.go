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
	if err == nil {
		return service, nil
	}
	_, n, err = Triable(n, err)
	if err != nil {
		return nil, err
	}
	goto retry
}

func getGDriveFile(
	r *http.Request,
	name string,
	field googleapi.Field,
) (
	*drive.File,
	*drive.Service,
	error,
) {
	if field == "" {
		field = MinimumField
	}

	var refreshToken bool
	n := 1
refresh:
	service, err := GetGDriveService(r)
	if err != nil {
		return nil, nil, err
	}

retry:
	<-tokenBucketGDriveAPI
	fileList, err := service.Files.List().PageSize(1).Spaces("drive").Q(fmt.Sprintf("name='%s'", name)).Fields(field).Do()

	if err != nil {
		refreshToken, n, err = Triable(n, err)
		if err != nil {
			return nil, nil, err
		}
		if refreshToken {
			goto refresh
		}
		goto retry
	}

	if len(fileList.Files) <= 0 {
		return nil, nil, &DriveFileDoesNotExistError{}
	}
	return fileList.Files[0], service, nil
}

// GetGDriveFile returns a file on Google Drive.
func GetGDriveFile(
	r *http.Request,
	name string,
	field googleapi.Field,
) (
	*drive.File,
	error,
) {
	if field == "" {
		field = defaultField
	}

	file, _, err := getGDriveFile(r, name, field)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetGDriveFileContent returns a file with content on Google Drive.
func GetGDriveFileContent(
	r *http.Request,
	name string,
	field googleapi.Field,
) (
	*drive.File,
	[]byte, // payload
	error,
) {
	if field == "" {
		field = defaultField
	}

	file, service, err := getGDriveFile(r, name, field)
	if err != nil {
		return nil, nil, err
	}
	fileID := file.Id

	var payload []byte

	var refreshToken bool
	n := 1
retry:
	httpResponse, err := service.Files.Export(fileID, mimeTxt).Download()
	if IsFileNotExportableError(err) {
		<-tokenBucketGDriveAPI
		err = service.Files.Delete(fileID).Do()
		if err != nil {
			// pass
		}
		return nil, nil, &DriveFileDoesNotExistError{}
	}
	if err != nil {
		refreshToken, n, err = Triable(n, err)
		if err != nil {
			return nil, nil, err
		}
		if refreshToken {
			service, err = GetGDriveService(r)
			if err != nil {
				return nil, nil, err
			}
		}
		goto retry
	}

	defer httpResponse.Body.Close()
	payload, err = ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, nil, err
	}
	return file, payload, err
}

// Triable returns whether it can retry or not.
func Triable(
	retries int,
	err error,
) (
	bool, //refreshToken
	int, //retries
	error,
) {
	if err == nil {
		return false, retries, nil
	}
	switch {
	case IsInvalidSecurityTicket(err):
		oauth2TokenSource = nil
		return true, retries, nil
	case IsRateLimit(err):
		time.Sleep(msec100)
		return false, retries, nil
	case IsServerError(err):
		retries, err = sleeping(retries)
		if err != nil {
			return false, retries, err
		}
		return false, retries, nil
	}
	return false, retries, err
}

func sleeping(
	n int,
) (
	int,
	error,
) {
	switch {
	case n > 16:
		return 0, errors.New("Sleeping Timeout")
	case n < 1:
		n = 1
	}
	time.Sleep(time.Duration(n)*time.Second + time.Duration(random.Intn(1000))*time.Millisecond)
	return n * 2, nil
}
