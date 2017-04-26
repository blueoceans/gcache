package gcache

import (
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

const (
	// MinimumField is minimim search fields on Google Drive API.
	MinimumField googleapi.Field = "files/id"

	defaultField googleapi.Field = "files(" +
		"name," +
		"properties," +
		"createdTime," +
		"modifiedTime," +
		"md5Checksum," +
		"size," +
		"id)"
)

// SetFolder sets name and permission to a top folder on Google Drive.
func SetFolder(
	name string,
	permission *drive.Permission,
) {
	setFolder(name, permission)
}

// GetFileName returns a file name on Google Drive.
func GetFileName(
	requestURI string,
) (
	string,
	error,
) {
	return getFileName(requestURI)
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
	return getGDriveFile(r, name, field)
}

// StoreGDrive stores a file to Google Drive.
func StoreGDrive(
	r *http.Request,
	name string,
	payload *[]byte,
) error {
	return storeGDrive(r, name, payload)
}
