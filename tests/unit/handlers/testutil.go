package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRequest helps build HTTP requests for testing
type TestRequest struct {
	Method string
	Path   string
	Body   interface{}
	Headers map[string]string
}

// BuildRequest creates an http.Request for testing
func BuildRequest(t *testing.T, tr TestRequest) *http.Request {
	t.Helper()

	var body io.Reader
	if tr.Body != nil {
		jsonBody, err := json.Marshal(tr.Body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		body = bytes.NewReader(jsonBody)
	}

	req := httptest.NewRequest(tr.Method, tr.Path, body)
	
	// Set default headers
	if tr.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	// Set custom headers
	for key, value := range tr.Headers {
		req.Header.Set(key, value)
	}

	return req
}

// ExecuteRequest executes a test request and returns the response
func ExecuteRequest(t *testing.T, handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	
	return rr
}

// AssertStatusCode checks if the response has the expected status code
func AssertStatusCode(t *testing.T, rr *httptest.ResponseRecorder, expected int) {
	t.Helper()

	if rr.Code != expected {
		t.Errorf("Expected status code %d, got %d. Body: %s", expected, rr.Code, rr.Body.String())
	}
}

// AssertJSONResponse decodes and returns the JSON response
func AssertJSONResponse(t *testing.T, rr *httptest.ResponseRecorder, v interface{}) {
	t.Helper()

	if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
		t.Fatalf("Failed to decode JSON response: %v. Body: %s", err, rr.Body.String())
	}
}

// AssertErrorResponse checks if the response contains an error message
func AssertErrorResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedMessage string) {
	t.Helper()

	var resp struct {
		Error string `json:"error"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if resp.Error != expectedMessage {
		t.Errorf("Expected error message %q, got %q", expectedMessage, resp.Error)
	}
}

// AssertContentType checks if the response has the expected content type
func AssertContentType(t *testing.T, rr *httptest.ResponseRecorder, expected string) {
	t.Helper()

	contentType := rr.Header().Get("Content-Type")
	if contentType != expected {
		t.Errorf("Expected Content-Type %q, got %q", expected, contentType)
	}
}

// MakeAuthHeader creates an Authorization header with Bearer token
func MakeAuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// MakeAPIKeyHeader creates an Authorization header with API key
func MakeAPIKeyHeader(apiKey string) map[string]string {
	return map[string]string{
		"Authorization": "ApiKey " + apiKey,
	}
}
