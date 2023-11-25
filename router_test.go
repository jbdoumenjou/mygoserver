package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdminMetricsRoute(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Error("Expected router to not be nil")
	}

	// GET http://localhost:8080/admin/metrics
	req := httptest.NewRequest(http.MethodGet, "/admin/metrics", nil)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rw.Code)
	}

	if !strings.Contains(rw.Body.String(), "Welcome, Chirpy Admin") {
		t.Errorf("Expected body to contain 'Welcome, Chirpy Admin', got %s", rw.Body.String())
	}

	if !strings.Contains(rw.Body.String(), "Chirpy has been visited 0 times!") {
		t.Errorf("Expected body to contain 'Chirpy has been visited 0 times!', got %s", rw.Body.String())
	}

	// GET http://localhost:8080/app to generate a hit
	req = httptest.NewRequest(http.MethodGet, "/app", nil)
	rw = httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rw.Code)
	}

	if !strings.Contains(rw.Body.String(), "Welcome to Chirpy") {
		t.Errorf("Expected body to contain 'Welcome to Chirpy', got %s", rw.Body.String())
	}

	// GET http://localhost:8080/admin/metrics
	req = httptest.NewRequest(http.MethodGet, "/admin/metrics", nil)
	rw = httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rw.Code)
	}

	if !strings.Contains(rw.Body.String(), "Chirpy has been visited 1 times!") {
		t.Errorf("Expected body to contain 'Chirpy has been visited 1 times!', got %s", rw.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/app", nil)
	rw = httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	req = httptest.NewRequest(http.MethodGet, "/app", nil)
	rw = httptest.NewRecorder()
	router.ServeHTTP(rw, req)
	req = httptest.NewRequest(http.MethodGet, "/app", nil)
	rw = httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	req = httptest.NewRequest(http.MethodGet, "/admin/metrics", nil)
	rw = httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rw.Code)
	}

	if !strings.Contains(rw.Body.String(), "Chirpy has been visited 4 times!") {
		t.Errorf("Expected body to contain 'Chirpy has been visited 4 times!', got %s", rw.Body.String())
	}

}

func TestValidateChirpRoute(t *testing.T) {

	tests := []struct {
		name           string
		body           map[string]string
		wantResp       string
		wantStatusCode int
	}{
		{
			name: "Valid chirp",
			body: map[string]string{
				"body": "I had something interesting for breakfast",
			},
			wantResp:       `{"cleaned_body":"I had something interesting for breakfast"}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Chirp too long",
			body: map[string]string{
				"body": "lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
			},
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "Valid with extra params",
			body: map[string]string{
				"body":  "I had something interesting for breakfast",
				"extra": "should be ignored",
			},
			wantResp:       `{"cleaned_body":"I had something interesting for breakfast"}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name: "Unclean chirp",
			body: map[string]string{
				"body": "I really need a kerfuffle to go to bed sooner, Fornax !",
			},
			wantResp:       `{"cleaned_body":"I really need a **** to go to bed sooner, **** !"}`,
			wantStatusCode: http.StatusOK,
		},
	}

	router := NewRouter()
	if router == nil {
		t.Error("Expected router to not be nil")
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			body, err := json.Marshal(test.body)
			if err != nil {
				t.Errorf("Expected no error, got %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", bytes.NewReader(body))
			rw := httptest.NewRecorder()
			router.ServeHTTP(rw, req)

			if rw.Code != test.wantStatusCode {
				t.Errorf("Expected status %d, got %d", test.wantStatusCode, rw.Code)
			}

			if test.wantResp != "" && rw.Body.String() != test.wantResp {
				t.Errorf("Expected body to be %s, got %s", test.wantResp, rw.Body.String())
			}
		})
	}
}
