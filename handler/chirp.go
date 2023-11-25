package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type ValidateChirpParameters struct {
	Body string `json:"body"`
}

type ValidateChirpResp struct {
	Valid bool   `json:"valid,omitempty"`
	Error string `json:"error,omitempty"`
}

func ValidateChirp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(r.Body)
	params := ValidateChirpParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		// an error will be thrown if the JSON is invalid or has the wrong types
		// any missing fields will simply have their values in the struct set to their zero value
		log.Printf("Something went wrong: %s", err)
		content, err := json.Marshal(ValidateChirpResp{Error: err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(content)

		return
	}

	if len(params.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		content, err := json.Marshal(ValidateChirpResp{Error: "Chirp is too long"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(content)

		return
	}

	w.WriteHeader(http.StatusOK)
	content, err := json.Marshal(ValidateChirpResp{Valid: true})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(content)
}
