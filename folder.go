package gcache

import (
	"fmt"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var (
	rootFolderID string

	folderParams = &drive.File{
		Name:     rootFolderName,
		MimeType: mimeGSuiteFolder,
	}

	folderPermission *drive.Permission
	rootFolderName   string
)

// SetRootFolder sets name and permission to a top folder on Google Drive.
func SetRootFolder(
	name string,
	permission *drive.Permission,
) {
	folderPermission = permission
	rootFolderName = name
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
	fileList, err := getFolder(r, fmt.Sprintf("name='%s'", rootFolderName), "")
	if err != nil {
		return "", err
	}

	if len(fileList.Files) == 1 {
		return fileList.Files[0].Id, nil
	}

	return createFolder(r)
}

func getFolder(
	r *http.Request,
	q string,
	field googleapi.Field,
) (
	*drive.FileList,
	error,
) {
	if q != "" {
		q = fmt.Sprintf("mimeType='%s' and (%s)", mimeGSuiteFolder, q)
	}
	if field == "" {
		field = MinimumField
	}

	var refreshToken bool
	n := 1
refresh:
	service, err := GetGDriveService(r)
	if err != nil {
		return nil, err
	}

retry:
	fileList, err := service.Files.List().PageSize(1).Spaces("drive").Q(q).Fields(field).Do()

	if err != nil {
		refreshToken, n, err = triable(n, err)
		if err != nil {
			return nil, err
		}
		if refreshToken {
			goto refresh
		}
		goto retry
	}

	return fileList, nil
}

func createFolder(
	r *http.Request,
) (
	string,
	error,
) {

	var refreshToken bool
	n := 1
refresh:
	service, err := GetGDriveService(r)
	if err != nil {
		return "", err
	}

retryFiles:
	file, err := service.Files.Create(folderParams).Do()

	if err != nil {
		refreshToken, n, err = triable(n, err)
		if err != nil {
			return "", err
		}
		if refreshToken {
			goto refresh
		}
		goto retryFiles
	}

retryPermissions:
	_, err = service.Permissions.Create(file.Id, folderPermission).Do()

	if err != nil {
		refreshToken, n, err = triable(n, err)
		if err != nil {
			return "", err
		}
		if refreshToken {
			service, err = GetGDriveService(r)
			if err != nil {
				return "", err
			}
		}
		goto retryPermissions
	}

	return file.Id, nil
}
