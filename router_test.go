package main

import (
	"bytes"
	"encoding/json"
	"github.com/jbdoumenjou/mygoserver/handler"
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
	router := NewRouter()
	if router == nil {
		t.Error("Expected router to not be nil")
	}

	// GET http://localhost:8080/admin/metrics
	parameters := handler.ValidateChirpParameters{
		Body: "I had something interesting for breakfast",
	}
	body, err := json.Marshal(parameters)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", bytes.NewReader(body))
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rw.Code)
	}

	want := handler.ValidateChirpResp{Valid: true}
	wantResp, err := json.Marshal(want)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}

	if rw.Body.String() != string(wantResp) {
		t.Errorf("Expected body to be %s, got %s", string(wantResp), rw.Body.String())
	}
}
