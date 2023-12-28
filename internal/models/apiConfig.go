package models

import (
	"github.com/ItsJustVaal/HoloGo/internal/database"
)

type VideoCache struct {
	LastVideo map[string]string
}

type ApiConfig struct {
	DB    *database.Queries
	Cache VideoCache
}
