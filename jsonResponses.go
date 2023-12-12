package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ItsJustVaal/HoloGo/internal/models"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, models.JsonErrorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "applcation/json")
	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marhsalling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}
