package chirp

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/jbdoumenjou/mygoserver/internal/api"
	"github.com/jbdoumenjou/mygoserver/internal/db"
	"net/http"
	"strconv"
	"strings"
)

type ChirpStorer interface {
	CreateChirp(body string) (db.Chirp, error)
	ListChirps() ([]db.Chirp, error)
	GetChirp(id int) (*db.Chirp, error)
}

type Handler struct {
	db ChirpStorer
}

// NewHandler returns a new handler.
func NewHandler(db ChirpStorer) *Handler {
	return &Handler{db: db}
}

type ChirpParameters struct {
	Body string `json:"body"`
}

// Create creates a new chirp.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := ChirpParameters{}

	err := decoder.Decode(&params)
	if err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(params.Body) > 140 {
		api.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := cleanChirp(params.Body)
	chirp, err := h.db.CreateChirp(cleanedChirp)
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

	chirps, err := h.db.GetChirp(id)
	if err != nil {
		// the only error we can get is not found
		api.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	api.RespondWithJSON(w, http.StatusOK, chirps)
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
