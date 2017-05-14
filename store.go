package gcache

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

const ()

var (
	parentFolderID string
	random         *rand.Rand

	folderParams = &drive.File{
		Name:     folderName,
		MimeType: mimeGSuiteFolder,
	}

	folderName       string
	folderPermission *drive.Permission
)

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// SetRootFolder sets name and permission to a top folder on Google Drive.
func SetRootFolder(
	name string,
	permission *drive.Permission,
) {
	folderName = name
	folderPermission = permission
}

func getParentFolderID(
	r *http.Request,
) (
	string,
	error,
) {
	if parentFolderID != "" {
		return parentFolderID, nil
	}
	parentFolderID, err := getDriveFolder(r)
	if err != nil {
		return "", err
	}
	return parentFolderID, nil
}

func getDriveFolder(
	r *http.Request,
) (
	string,
	error,
) {

retry:
	service, err := drive.New(createGDriveClient(r))
	if err != nil {
		if IsInvalidSecurityTicket(err) {
			oauth2TokenSource = nil
			goto retry
		}
		return "", err
	}

	fileList, err := service.Files.List().PageSize(1).Spaces("drive").Q(
		fmt.Sprintf("name='%s' and mimeType='%s'", folderName, mimeGSuiteFolder),
	).Fields(MinimumField).Do()
	if err != nil {
		return "", err
	}
	if len(fileList.Files) == 1 {
		return fileList.Files[0].Id, nil
	}

	return createDriveFolder(r)
}

func createDriveFolder(
	r *http.Request,
) (
	string,
	error,
) {

retry:
	service, err := drive.New(createGDriveClient(r))
	if err != nil {
		if IsInvalidSecurityTicket(err) {
			oauth2TokenSource = nil
			goto retry
		}
		return "", err
	}

	file, err := service.Files.Create(folderParams).Do()
	if err != nil {
		return "", err
	}
	_, err = service.Permissions.Create(file.Id, folderPermission).Do()
	if err != nil {
		return "", err
	}
	return file.Id, nil
}

// StoreGDrive stores a file to Google Drive.
func StoreGDrive(
	r *http.Request,
	name string,
	payload *[]byte,
) (
	*drive.File,
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

	folderID, err := getParentFolderID(r)
	if err != nil {
		return nil, err
	}

	file, err := service.Files.Create(&drive.File{
		Name:     name,
		MimeType: mimeGSuiteDoc,
		Parents:  []string{folderID},
	}).Media(bytes.NewReader(*payload), googleapi.ContentType(mimeTxt)).Do()
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
