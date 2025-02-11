package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthMiddleware(t *testing.T) {
	// Assume BasicAuthMiddleware returns a middleware function.
	middleware := BasicAuthMiddleware("user", "pass")

	// Create a simple handler that returns 200 OK.
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the middleware.
	handler := middleware(testHandler)

	// Test with correct credentials.
	req, _ := http.NewRequest("GET", "/", nil)
	req.SetBasicAuth("user", "pass")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200 with correct credentials, got %d", rr.Code)
	}

	// Test with incorrect credentials.
	req2, _ := http.NewRequest("GET", "/", nil)
	req2.SetBasicAuth("user", "wrongpass")
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)
	if rr2.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401 with wrong credentials, got %d", rr2.Code)
	}
}
