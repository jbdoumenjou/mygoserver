package token

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	issuerRefresh = "chirpy-refresh"
	issuerAccess  = "chirpy-access"
)

type Manager struct {
	jwtSecret string
}

func NewManager(jwtSecret string) *Manager {
	return &Manager{jwtSecret: jwtSecret}
}

func (t *Manager) GetAccessToken(header http.Header) (*jwt.Token, error) {
	return t.getToken(header, issuerAccess)
}

func (t *Manager) GetRefreshToken(header http.Header) (*jwt.Token, error) {
	return t.getToken(header, issuerRefresh)
}

func (t *Manager) GetUserID(token *jwt.Token) (int, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, errors.New("invalid token")
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid token")
	}

	return userID, nil
}

func (t *Manager) getToken(header http.Header, expectedIssuer string) (*jwt.Token, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("missing Authorization header")
	}

	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		return nil, errors.New("invalid Authorization header")
	}

	tokenString := split[1]
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.jwtSecret), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if issuer != expectedIssuer {
		return nil, errors.New("invalid token issuer")
	}

	return token, nil
}

func (t *Manager) CreateAccessToken(userID int) (string, error) {
	return t.createToken(userID, issuerAccess, time.Hour)
}

func (t *Manager) CreateRefreshToken(userID int) (string, error) {
	return t.createToken(userID, issuerRefresh, 60*24*time.Hour)
}

func (t *Manager) createToken(userID int, issuer string, expiresAt time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := jwt.RegisteredClaims{
		Issuer:    issuer,
		Subject:   strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresAt)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(t.jwtSecret))
}
