package user

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jbdoumenjou/mygoserver/internal/api"
	"github.com/jbdoumenjou/mygoserver/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type UserStorer interface {
	CreateUser(username, password string) (db.User, error)
	UpdateUser(id int, email, password string) (db.User, error)
	GetUserByEmail(email string) (*db.User, error)
	GetUser(id int) (*db.User, error)
}

type Handler struct {
	db        UserStorer
	jwtSecret string
}

func NewHandler(db UserStorer, jwtSecret string) *Handler {
	return &Handler{db: db, jwtSecret: jwtSecret}
}

type Parameters struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.db.CreateUser(params.Email, string(bcryptPassword))
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := UserResponse{
		Email: user.Email,
		ID:    user.ID,
	}

	api.RespondWithJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		api.RespondWithError(w, http.StatusUnauthorized, "missing Authorization header")
		return
	}
	split := strings.Split(authHeader, " ")
	if len(split) != 2 {
		api.RespondWithError(w, http.StatusUnauthorized, "invalid Authorization header")
		return
	}
	tokenString := split[1]
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})
	if err != nil {
		api.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	subject, err := token.Claims.GetSubject()
	if err != nil {
		api.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := strconv.Atoi(subject)
	if err != nil {
		api.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	updatedUser, err := h.db.UpdateUser(userID, params.Email, string(bcryptPassword))
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := UserResponse{
		Email: updatedUser.Email,
		ID:    updatedUser.ID,
	}

	api.RespondWithJSON(w, http.StatusOK, resp)
}

type UserLoginParameters struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

type UserLoginResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := UserLoginParameters{}

	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.db.GetUserByEmail(params.Email)
	if err != nil {
		// the only error we expect is "not found"
		api.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		api.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// default expiration is 24 hours
	expiresInSeconds := 60 * 60 * 24 * time.Second
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds < 60*60*24 {
		expiresInSeconds = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	// ok we can create a jwt token
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresInSeconds)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := UserLoginResponse{
		ID:    user.ID,
		Email: user.Email,
		Token: tokenString,
	}

	api.RespondWithJSON(w, http.StatusOK, resp)
}
