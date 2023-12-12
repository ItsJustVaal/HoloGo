package models

import "github.com/ItsJustVaal/HoloGo/internal/database"

type ApiConfig struct {
	DB *database.Queries
}

type JsonErrorResponse struct {
	Error string `json:"error"`
}
type JsonResponse struct {
	Status string `json:"status"`
}
