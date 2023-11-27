package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jbdoumenjou/mygoserver/internal/api/token"

	"github.com/jbdoumenjou/mygoserver/internal/db"
)

func TestAdminMetricsRoute(t *testing.T) {
	mockDB := NewMockDB()
	tokenManager := token.NewManager("mysecret")
	router := NewRouter(mockDB, tokenManager)
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

type MockDB struct {
	Chirps []db.Chirp
}

func (m *MockDB) CreateUser(username, password string) (db.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDB) UpdateUser(id int, email, password string) (db.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDB) GetUserByEmail(email string) (*db.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDB) GetUser(id int) (*db.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockDB) RevokeToken(token string) string {
	//TODO implement me
	panic("implement me")
}

func (m *MockDB) IsTokenRevoked(token string) bool {
	//TODO implement me
	panic("implement me")
}

func NewMockDB() *MockDB {
	return &MockDB{Chirps: []db.Chirp{}}
}

func (m *MockDB) CreateChirp(body string, authorID int) (db.Chirp, error) {
	chirp := db.Chirp{ID: 1, AuthorID: 1, Body: body}
	m.Chirps = append(m.Chirps, chirp)
	return chirp, nil
}

func (m *MockDB) ListChirps() ([]db.Chirp, error) {
	return m.Chirps, nil
}

func (m *MockDB) GetChirp(id int) (*db.Chirp, error) {
	for _, chirp := range m.Chirps {
		if chirp.ID == id {
			return &chirp, nil
		}
	}
	return nil, errors.New("not found")
}

func TestCreateChirpRoute(t *testing.T) {
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
			wantResp:       `{"id":1,"author_id":1,"body":"I had something interesting for breakfast"}`,
			wantStatusCode: http.StatusCreated,
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
			wantResp:       `{"id":1,"author_id":1,"body":"I had something interesting for breakfast"}`,
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "Unclean chirp",
			body: map[string]string{
				"body": "I really need a kerfuffle to go to bed sooner, Fornax !",
			},
			wantResp:       `{"id":1,"author_id":1,"body":"I really need a **** to go to bed sooner, **** !"}`,
			wantStatusCode: http.StatusCreated,
		},
	}
	mockDB := NewMockDB()
	tokenManager := token.NewManager("mysecret")
	accessToken, err := tokenManager.CreateAccessToken(1)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
	router := NewRouter(mockDB, tokenManager)

	if router == nil {
		t.Error("Expected router to not be nil")
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.body)
			if err != nil {
				t.Errorf("Expected no error, got %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, "/api/chirps", bytes.NewReader(body))
			req.Header.Set("Authorization", "Bearer "+accessToken)
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

func TestGetChirp(t *testing.T) {
	mockDB := NewMockDB()
	tokenManager := token.NewManager("mysecret")
	chirp := db.Chirp{ID: 1, Body: "I had something interesting for breakfast"}
	mockDB.Chirps = []db.Chirp{chirp}
	router := NewRouter(mockDB, tokenManager)
	if router == nil {
		t.Error("Expected router to not be nil")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/chirps/1", http.NoBody)
	rw := httptest.NewRecorder()
	router.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Errorf("Expected StatusOk, got %d", rw.Code)
	}
	want, err := json.Marshal(chirp)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
	if rw.Body.String() != string(want) {
		t.Errorf("Expected body to be %s, got %s", string(want), rw.Body.String())
	}
}
