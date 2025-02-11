package apistatus

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestNew_WithString(t *testing.T) {
	msg := "this is a test message"
	s := New(msg)
	if s.GetMessage() != msg {
		t.Errorf("expected message %q, got %q", msg, s.GetMessage())
	}
}

func TestNew_WithError(t *testing.T) {
	err := errors.New("error occurred")
	s := New(err)
	if s.GetMessage() != err.Error() {
		t.Errorf("expected message %q, got %q", err.Error(), s.GetMessage())
	}
}

func TestNew_WithNoArgs(t *testing.T) {
	s := New()
	// With no arguments, no base message is set so we expect an empty string.
	if s.GetMessage() != "" {
		t.Errorf("expected empty message, got %q", s.GetMessage())
	}
}

func TestSetLanguage(t *testing.T) {
	s := New("test")
	s.SetLanguage(LanguageEN)
	if s.GetLanguage() != LanguageEN {
		t.Errorf("expected language %q, got %q", LanguageEN, s.GetLanguage())
	}
}

func TestSetCode(t *testing.T) {
	s := New("test")
	s.SetCode(123)
	if s.GetCode() != 123 {
		t.Errorf("expected app code 123, got %d", s.GetCode())
	}
}

func TestSetMessageCode(t *testing.T) {
	s := New("test")
	expected := "CODE123"
	s.SetMessageCode(expected)
	if s.GetMessageCode() != expected {
		t.Errorf("expected message code %q, got %v", expected, s.GetMessageCode())
	}
}

func TestAddAndGetDetails(t *testing.T) {
	s := New("test")
	details := "extra details"
	s.AddDetails(details)
	if s.GetDetails() != details {
		t.Errorf("expected details %q, got %q", details, s.GetDetails())
	}
}

func TestHTTPStatusMethods(t *testing.T) {
	// Table-driven test for each status helper method.
	tests := []struct {
		name           string
		apply          func(Status) Status
		expectedStatus int
		expectedPrefix string
	}{
		{"OK", func(s Status) Status { return s.OK() }, http.StatusOK, "ok"},
		{"Created", func(s Status) Status { return s.Created() }, http.StatusCreated, "created"},
		{"Accepted", func(s Status) Status { return s.Accepted() }, http.StatusAccepted, "accepted"},
		{"NoContent", func(s Status) Status { return s.NoContent() }, http.StatusNoContent, "no content"},
		{"MovedPermanently", func(s Status) Status { return s.MovedPermanently() }, http.StatusMovedPermanently, "moved permanently"},
		{"Found", func(s Status) Status { return s.Found() }, http.StatusFound, "found"},
		{"SeeOther", func(s Status) Status { return s.SeeOther() }, http.StatusSeeOther, "see other"},
		{"TemporaryRedirect", func(s Status) Status { return s.TemporaryRedirect() }, http.StatusTemporaryRedirect, "temporary redirect"},
		{"BadRequest", func(s Status) Status { return s.BadRequest() }, http.StatusBadRequest, "bad request"},
		{"Unauthorized", func(s Status) Status { return s.Unauthorized() }, http.StatusUnauthorized, "unauthorized"},
		{"Forbidden", func(s Status) Status { return s.Forbidden() }, http.StatusForbidden, "forbidden"},
		{"NotFound", func(s Status) Status { return s.NotFound() }, http.StatusNotFound, "not found"},
		{"UnprocessableEntity", func(s Status) Status { return s.UnprocessableEntity() }, http.StatusUnprocessableEntity, "unprocessable entity"},
		{"InternalServerError", func(s Status) Status { return s.InternalServerError() }, http.StatusInternalServerError, "internal server error"},
	}

	baseMsg := "test message"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(baseMsg)
			s = tt.apply(s)
			if s.GetStatus() != tt.expectedStatus {
				t.Errorf("expected HTTP status %d, got %d", tt.expectedStatus, s.GetStatus())
			}
			expectedMsg := fmt.Sprintf("%s: %s", tt.expectedPrefix, baseMsg)
			if s.GetMessage() != expectedMsg {
				t.Errorf("expected message %q, got %q", expectedMsg, s.GetMessage())
			}
		})
	}
}

func TestIsError(t *testing.T) {
	// A status with a 200 OK should not be considered an error.
	s := New("test").OK()
	if s.IsError() {
		t.Error("expected IsError to return false for 200 OK")
	}

	// A status with a 400 Bad Request should be considered an error.
	s = New("test").BadRequest()
	if !s.IsError() {
		t.Error("expected IsError to return true for 400 Bad Request")
	}
}
