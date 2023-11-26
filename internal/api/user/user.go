package user

import (
	"encoding/json"
	"github.com/jbdoumenjou/mygoserver/internal/api"
	"github.com/jbdoumenjou/mygoserver/internal/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserStorer interface {
	CreateUser(username, password string) (db.User, error)
	GetUSer(email string) (*db.User, error)
}

type Handler struct {
	db UserStorer
}

func NewHandler(db UserStorer) *Handler {
	return &Handler{db: db}
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
		ID:    user.ID,
		Email: user.Email,
	}

	api.RespondWithJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.db.GetUSer(params.Email)
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

	resp := UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	api.RespondWithJSON(w, http.StatusOK, resp)
}
