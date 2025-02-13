package api

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"messaging-app/config"
)

func TestRateLimitingIntegration(t *testing.T) {
	// Create a dummy service.
	ds := &dummyService{}
	// Create a dummy configuration with auth and rate limit settings.
	testConfig := &config.Config{
		AuthUsername: "red",
		AuthPassword: "abc123",
		RateLimit:    100,
		HTTPPort:     "3000",
	}

	// Create the API handler using the dummy service.
	handler := NewHandler(ds, testConfig)

	// Create the router using your actual NewRouter function.
	router := NewRouter(handler, testConfig)

	// Create an HTTP test server with the router.
	ts := httptest.NewServer(router)
	defer ts.Close()

	client := ts.Client()

	// Prepare the BasicAuth header (the BasicAuth middleware in your router expects "red" / "abc123").
	authStr := "red:abc123"
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authStr))

	// Use a constant IP to simulate requests from the same client.
	const testIP = "1.2.3.4"

	// We'll test the GET /chats/1/messages endpoint.
	url := ts.URL + "/chats/1/messages"

	// Send 100 requests that should succeed.
	for i := 0; i < 100; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			t.Fatalf("failed to create request %d: %v", i+1, err)
		}
		req.Header.Set("Authorization", "Basic "+encodedAuth)
		req.Header.Set("X-Forwarded-For", testIP)
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("failed to perform request %d: %v", i+1, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 on request %d, got %d", i+1, resp.StatusCode)
		}
		resp.Body.Close()
	}

	// The 101st request should be rate-limited (HTTP 429).
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("failed to create request 101: %v", err)
	}
	req.Header.Set("Authorization", "Basic "+encodedAuth)
	req.Header.Set("X-Forwarded-For", testIP)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to perform request 101: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 on request 101, got %d", resp.StatusCode)
	}
}
