package main

import (
	"net/http"

	"github.com/ItsJustVaal/HoloGo/internal/models"
)

func handleGetReadiness(w http.ResponseWriter, r *http.Request) {
	models.RespondWithJSON(w, http.StatusOK, models.JsonResponse{
		Status: "ok",
	})
}

func handleGetErr(w http.ResponseWriter, r *http.Request) {
	models.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
