package gcache

import (
	"strings"
)

const (
	// https://developers.google.com/drive/v3/web/handle-errors
	// 400
	reasonBadRequest            = "badRequest"
	reasonInvalidSharingRequest = "invalidSharingRequest"
	// 401
	reasonAuthError = "authError"
	// 403
	reasonDailyLimitExceeded          = "dailyLimitExceeded"
	reasonUserRateLimitExceeded       = "userRateLimitExceeded"
	reasonRateLimitExceeded           = "rateLimitExceeded"
	reasonSharingRateLimitExceeded    = "sharingRateLimitExceeded"
	reasonAppNotAuthorizedToFile      = "appNotAuthorizedToFile"
	reasonInsufficientFilePermissions = "insufficientFilePermissions"
	reasonDomainPolicy                = "domainPolicy"
	// 404
	reasonNotFound = "notFound"
	// 500
	reasonBackendError = "backendError"
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
	errRateLimit = []string{
		reasonUserRateLimitExceeded,
		reasonRateLimitExceeded,
	}
)

// DriveFileDoesNotExistError is as HTTP response that is 40X HTTP status.
type DriveFileDoesNotExistError struct {
	message string
}

func (err DriveFileDoesNotExistError) Error() string {
	return err.message
}

// NewDriveFileDoesNotExistError returns a DriveFileDoesNotExistError.
func NewDriveFileDoesNotExistError() *DriveFileDoesNotExistError {
	return &DriveFileDoesNotExistError{message: "drive: file does not exist"}
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

// IsRateLimit returns is whether it is "userRateLimitExceeded" or "rateLimitExceeded" server errors or not.
func IsRateLimit(
	err error,
) bool {
	return containsErrorMessage(err, errRateLimit)
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
		if strings.Contains(errorMessage, message) {
			return true
		}
	}
	return false
}
