package gcache

import (
	"fmt"
	"net/http"

	"google.golang.org/api/drive/v3"
)

var (
	rootFolderID string

	folderParams = &drive.File{
		Name:     rootFolderName,
		MimeType: mimeGSuiteFolder,
	}

	rootFolderName   string
	folderPermission *drive.Permission
)

// SetRootFolder sets name and permission to a top folder on Google Drive.
func SetRootFolder(
	name string,
	permission *drive.Permission,
) {
	rootFolderName = name
	folderPermission = permission
}

func getRootFolderID(
	r *http.Request,
) (
	string,
	error,
) {
	if rootFolderID != "" {
		return rootFolderID, nil
	}
	rootFolderID, err := getDriveFolder(r)
	if err != nil {
		return "", err
	}
	return rootFolderID, nil
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
		fmt.Sprintf("name='%s' and mimeType='%s'", rootFolderName, mimeGSuiteFolder),
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
