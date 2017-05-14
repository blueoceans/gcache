package gcache

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

const ()

var (
	oauth2TokenSource oauth2.TokenSource // The token is valid for 30 minutes.
)

func init() {
}

// GetFileName returns a file name on Google Drive.
func GetFileName(
	requestURI string,
) (
	string,
	error,
) {
	uri, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return "", err
	}

	splits := strings.SplitN(uri.Host, ":", 2) // hostname:port

	// remove the standard port number
	switch len(splits) {
	case 2:
		switch {
		case
			uri.Scheme == "http" && splits[1] == "80",
			uri.Scheme == "https" && splits[1] == "443":
			uri.Host = splits[0]
		}
	}

	return uri.String(), nil
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
