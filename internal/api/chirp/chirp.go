package chirp

import (
	"encoding/json"
	"github.com/jbdoumenjou/mygoserver/internal/db"
	"net/http"
	"strings"
)

type ChirpStorer interface {
	CreateChirp(body string) (db.Chirp, error)
	ListChirps() ([]db.Chirp, error)
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
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedChirp := cleanChirp(params.Body)
	chirp, err := h.db.CreateChirp(cleanedChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

// List returns all chirps in the database
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	chirps, err := h.db.ListChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirps)
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	content, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(content)
}
