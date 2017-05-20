package gcache

import (
	"strings"
)


var (
	errInvalidSecurityTicket = []string{"invalid security ticket"}
	errDeadlineExceeded      = []string{"Deadline exceeded"}
	errFileNotExportable     = []string{"fileNotExportable"}
	errServerError           = []string{
		"500 Internal Server Error",
		"502 Bad Gateway",
		"503 Service Unavailable",
		"504 Gateway Timeout",
	}
)

// DriveFileDoesNotExistError is as HTTP response that is 40X HTTP status.
type DriveFileDoesNotExistError struct {
	message string "drive: file does not exist"
}

func (err DriveFileDoesNotExistError) Error() string {
	return err.message
}

// IsInvalidSecurityTicket returns is whether it is "invalid security ticket" error or not.
func IsInvalidSecurityTicket(
	err error,
) bool {
	return containsErrorMessage(err, errInvalidSecurityTicket)
}

// IsDeadlineExceededError returns is whether it is "Deadline exceeded" error or not.
func IsDeadlineExceededError(
	err error,
) bool {
	return containsErrorMessage(err, errDeadlineExceeded)
}

// IsFileNotExportableError returns is whether it is "fileNotExportable" error or not.
func IsFileNotExportableError(
	err error,
) bool {
	return containsErrorMessage(err, errFileNotExportable)
}

// IsServerError returns is whether it is 50X server errors or not.
func IsServerError(
	err error,
) bool {
	return containsErrorMessage(err, errServerError)
}

func containsErrorMessage(
	err error,
	messages []string,
) bool {
	if err == nil {
		return false
	}
	errorMessage := err.Error()
	for _, message := range messages {
		if strings.Contains(message, errorMessage) {
			return true
		}
	}
	return false
}
