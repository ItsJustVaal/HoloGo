package models

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	RespondWithJSON(w, code, JsonErrorResponse{
		Error: msg,
	})
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, response interface{}) {
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
