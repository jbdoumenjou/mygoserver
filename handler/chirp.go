package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ValidateChirpParameters struct {
	Body string `json:"body"`
}

type ValidateChirpResp struct {
	CleanedBody string `json:"cleaned_body,omitempty"`
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := ValidateChirpParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	resp := ValidateChirpResp{CleanedBody: cleanChirp(params.Body)}
	respondWithJSON(w, http.StatusOK, resp)
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
