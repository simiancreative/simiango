package errors

import (
	"net/http"
)

const (
	ValidationFailed = iota
	IntegrationNotFound
	StatusForbidden
)

var DefaultReasons = Reasons{
	ValidationFailed: {
		Status:      http.StatusUnprocessableEntity,
		Key:         "validation_failed",
		Description: "Validation failed",
	},

	IntegrationNotFound: {
		Status:      http.StatusNotImplemented,
		Key:         "integration_not_found",
		Description: "Integration not found",
	},

	StatusForbidden: {
		Status:      http.StatusForbidden,
		Key:         "forbidden_request",
		Description: "Forbidden Request",
	},
}
