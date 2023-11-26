package user

import (
	"encoding/json"
	"github.com/jbdoumenjou/mygoserver/internal/api"
	"github.com/jbdoumenjou/mygoserver/internal/db"
	"net/http"
)

type UserStorer interface {
	CreateUser(username string) (db.User, error)
}

type Handler struct {
	db UserStorer
}

func NewHandler(db UserStorer) *Handler {
	return &Handler{db: db}
}

type Parameters struct {
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

	user, err := h.db.CreateUser(params.Email)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, user)
}
