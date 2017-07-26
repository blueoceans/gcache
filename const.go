package gcache

import (
	"time"

	"google.golang.org/api/googleapi"
)

const (
	// MinimumFilesField is minimim search fields on Google Drive API.
	MinimumFilesField googleapi.Field = "files/id"
	minimumField      googleapi.Field = "id"
	defaultFilesField googleapi.Field = "files(" +
		"appProperties," +
		"modifiedTime," +
		"name," +
		"id)"
	defaultField googleapi.Field = "" +
		"appProperties," +
		"modifiedTime," +
		"name," +
		"id"

	msec100 = time.Duration(100) * time.Millisecond
)
