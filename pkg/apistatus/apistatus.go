package apistatus

import (
	"errors"
	"fmt"
	"net/http"
)

// Language represents a language code for status messages.
type Language string

const (
	LanguageAR = "ar"
	LanguageEN = "en"
)

// Predefined errors.
var ErrNoContent = errors.New("no db rows affected")

// Status provides methods for building and interrogating HTTP statuses.
type Status interface {
	// Configuration
	SetLanguage(language Language) Status
	SetCode(appCode int) Status
	SetMessageCode(messageCode interface{}) Status

	// Retrieval
	GetLanguage() Language
	GetCode() int
	GetMessageCode() interface{}
	GetError() error
	GetMessage() string
	GetStatus() int
	AddDetails(details string) Status
	GetDetails() string

	// 2xx Success Statuses
	OK() Status
	Created() Status
	Accepted() Status
	NoContent() Status

	// 3xx Redirection Statuses
	MovedPermanently() Status
	Found() Status
	SeeOther() Status
	TemporaryRedirect() Status

	// 4xx Client Error Statuses
	BadRequest() Status
	Unauthorized() Status
	Forbidden() Status
	NotFound() Status
	UnprocessableEntity() Status

	// 5xx Server Error Statuses
	InternalServerError() Status

	// Checks if the status represents an error (HTTP 4xx or 5xx).
	IsError() bool
}

// status is the concrete implementation of the Status interface.
type status struct {
	language    Language
	baseMessage string      // The original message provided at creation.
	message     string      // The decorated message (may include prefixes).
	messageCode interface{} // An optional message code.
	details     string      // Additional details.
	httpStatus  int         // The HTTP status code.
	appCode     int         // An optional application-specific code.
}

// New creates a new Status instance. If arguments are provided, the first one is interpreted
// as either a string or an error. If the first argument is a string, the remaining arguments
// are used with fmt.Sprintf.
func New(args ...interface{}) Status {
	s := &status{}
	if len(args) == 0 {
		return s
	}

	switch arg := args[0].(type) {
	case string:
		s.baseMessage = fmt.Sprintf(arg, args[1:]...)
	case error:
		s.baseMessage = arg.Error()
	default:
		s.baseMessage = "unknown error"
	}
	// Initialize the decorated message to the base message.
	s.message = s.baseMessage

	return s
}

// SetLanguage sets the language for the status message.
func (s *status) SetLanguage(language Language) Status {
	s.language = language
	return s
}

// SetCode sets the application-specific code.
func (s *status) SetCode(appCode int) Status {
	s.appCode = appCode
	return s
}

// SetMessageCode sets an optional message code.
func (s *status) SetMessageCode(messageCode interface{}) Status {
	s.messageCode = messageCode
	return s
}

// GetLanguage retrieves the language.
func (s *status) GetLanguage() Language {
	return s.language
}

// GetCode retrieves the application-specific code.
func (s *status) GetCode() int {
	return s.appCode
}

// GetMessageCode retrieves the optional message code.
func (s *status) GetMessageCode() interface{} {
	return s.messageCode
}

// GetError converts the current status message into an error.
func (s *status) GetError() error {
	return errors.New(s.message)
}

// GetMessage retrieves the decorated message.
func (s *status) GetMessage() string {
	return s.message
}

// GetStatus retrieves the current HTTP status code.
func (s *status) GetStatus() int {
	return s.httpStatus
}

// AddDetails adds additional details to the status.
func (s *status) AddDetails(details string) Status {
	s.details = details
	return s
}

// GetDetails retrieves any additional details.
func (s *status) GetDetails() string {
	return s.details
}

// helper function to update status code and message.
func (s *status) update(prefix string, code int) Status {
	s.httpStatus = code
	s.message = fmt.Sprintf("%s: %s", prefix, s.baseMessage)
	return s
}

// 2xx Success Statuses

// OK sets the status to 200 OK.
func (s *status) OK() Status {
	return s.update("ok", http.StatusOK)
}

// Created sets the status to 201 Created.
func (s *status) Created() Status {
	return s.update("created", http.StatusCreated)
}

// Accepted sets the status to 202 Accepted.
func (s *status) Accepted() Status {
	return s.update("accepted", http.StatusAccepted)
}

// NoContent sets the status to 204 No Content.
func (s *status) NoContent() Status {
	return s.update("no content", http.StatusNoContent)
}

// 3xx Redirection Statuses

// MovedPermanently sets the status to 301 Moved Permanently.
func (s *status) MovedPermanently() Status {
	return s.update("moved permanently", http.StatusMovedPermanently)
}

// Found sets the status to 302 Found.
func (s *status) Found() Status {
	return s.update("found", http.StatusFound)
}

// SeeOther sets the status to 303 See Other.
func (s *status) SeeOther() Status {
	return s.update("see other", http.StatusSeeOther)
}

// TemporaryRedirect sets the status to 307 Temporary Redirect.
func (s *status) TemporaryRedirect() Status {
	return s.update("temporary redirect", http.StatusTemporaryRedirect)
}

// 4xx Client Error Statuses

// BadRequest sets the status to 400 Bad Request.
func (s *status) BadRequest() Status {
	return s.update("bad request", http.StatusBadRequest)
}

// Unauthorized sets the status to 401 Unauthorized.
func (s *status) Unauthorized() Status {
	return s.update("unauthorized", http.StatusUnauthorized)
}

// Forbidden sets the status to 403 Forbidden.
func (s *status) Forbidden() Status {
	return s.update("forbidden", http.StatusForbidden)
}

// NotFound sets the status to 404 Not Found.
func (s *status) NotFound() Status {
	return s.update("not found", http.StatusNotFound)
}

// UnprocessableEntity sets the status to 422 Unprocessable Entity.
func (s *status) UnprocessableEntity() Status {
	return s.update("unprocessable entity", http.StatusUnprocessableEntity)
}

// 5xx Server Error Statuses

// InternalServerError sets the status to 500 Internal Server Error.
func (s *status) InternalServerError() Status {
	return s.update("internal server error", http.StatusInternalServerError)
}

// IsError returns true if the HTTP status code indicates an error (i.e. 400 or above).
func (s *status) IsError() bool {
	return s.httpStatus >= http.StatusBadRequest
}
