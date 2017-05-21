package gcache

import (
	"time"

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

	msec100 = time.Duration(100) * time.Millisecond
)
