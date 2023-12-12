package main

import (
	"net/http"

	"github.com/ItsJustVaal/HoloGo/internal/models"
)

func handleGetReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, models.JsonResponse{
		Status: "ok",
	})
}

func handleGetErr(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
