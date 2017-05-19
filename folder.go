package gcache

import (
	"errors"
	"fmt"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

var (
	rootFolderID string

	folderPermission *drive.Permission
	rootFolderName   string
)

// SetRootFolder sets name and permission to a top folder on Google Drive.
func SetRootFolder(
	name string,
	permission *drive.Permission,
) error {
	if name == "" {
		return errors.New("`name` must be enough")
	}
	rootFolderName = name
	folderPermission = permission
	return nil
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
	if rootFolderName == "" {
		return "", errors.New("must SetRootFolder")
	}
	fileList, err := getFolder(r, fmt.Sprintf("name='%s'", rootFolderName), "")
	if err != nil {
		return "", err
	}

	if len(fileList.Files) == 1 {
		return fileList.Files[0].Id, nil
	}

	return createRootFolder(r)
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
		q = fmt.Sprintf("mimeType='%s' and (%s)", MimeGSuiteFolder, q)
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
		refreshToken, n, err = Triable(n, err)
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

// GetGDiveFolderIDs returns a map of folder IDs on Google Drive.
func GetGDiveFolderIDs(
	r *http.Request,
	q string,
) (
	*map[string]string,
	error,
) {

	const field googleapi.Field = "files(" +
		"name," +
		"id)"

	fileList, err := getFolder(r, q, field)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string, len(fileList.Files))
	for _, v := range fileList.Files {
		result[v.Name] = v.Id
	}
	return &result, nil
}

// CreateFolder returns a ID of new Google Drive Folder.
func CreateFolder(
	r *http.Request,
	file *drive.File,
) (
	string,
	error,
) {

	if file.Name == "" {
		return "", errors.New("`file.Name` must be enough")
	}

	file.MimeType = MimeGSuiteFolder
	if file.Parents == nil {
		if file.Name != rootFolderName {
			folderID, err := getRootFolderID(r)
			if err != nil {
				return "", err
			}
			file.Parents = []string{folderID}
		}
	}

	var refreshToken bool
	n := 1
refresh:
	service, err := GetGDriveService(r)
	if err != nil {
		return "", err
	}

retryFiles:
	file, err = service.Files.Create(file).Do()

	if err != nil {
		refreshToken, n, err = Triable(n, err)
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
		refreshToken, n, err = Triable(n, err)
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

func createRootFolder(
	r *http.Request,
) (
	string,
	error,
) {
	return CreateFolder(r,
		&drive.File{
			Name: rootFolderName,
		},
	)
}
