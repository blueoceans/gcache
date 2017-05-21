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

	sec1 = time.Duration(1) * time.Second
)
