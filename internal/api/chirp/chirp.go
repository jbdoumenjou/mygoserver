package chirp

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jbdoumenjou/mygoserver/internal/api"
	"github.com/jbdoumenjou/mygoserver/internal/api/token"
	"github.com/jbdoumenjou/mygoserver/internal/db"
)

type ChirpStorer interface {
	CreateChirp(body string, authorID int) (db.Chirp, error)
	ListChirps() ([]db.Chirp, error)
	GetChirp(id int) (*db.Chirp, error)
}

type Handler struct {
	db           ChirpStorer
	tokenManager *token.Manager
}

// NewHandler returns a new handler.
func NewHandler(db ChirpStorer, tokenManager *token.Manager) *Handler {
	return &Handler{db: db, tokenManager: tokenManager}
}

type ChirpParameters struct {
	Body string `json:"body"`
}

// Create creates a new chirp.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	// first check the access token. If it's not valid, we can't create a chirp.
	token, err := h.tokenManager.GetAccessToken(r.Header)
	if err != nil {
		api.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := h.tokenManager.GetUserID(token)
	if err != nil {
		api.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := ChirpParameters{}

	if err := decoder.Decode(&params); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(params.Body) > 140 {
		api.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := cleanChirp(params.Body)
	chirp, err := h.db.CreateChirp(cleanedChirp, userID)
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusCreated, chirp)
}

// List returns all chirps in the database
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	chirps, err := h.db.ListChirps()
	if err != nil {
		api.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, chirps)
}

// Get returns a single chirp.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := h.db.GetChirp(id)
	if err != nil {
		// the only error we can get is not found
		api.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, chirp)
}

func cleanChirp(body string) string {
	splitBody := strings.Split(body, " ")
	for i, word := range splitBody {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			splitBody[i] = "****"
		}
	}

	return strings.Join(splitBody, " ")
}
