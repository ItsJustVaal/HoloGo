package models

import (
	"github.com/ItsJustVaal/HoloGo/internal/database"
	"github.com/ItsJustVaal/HoloGo/youtube"
)

type ApiConfig struct {
	DB    *database.Queries
	Cache youtube.VideoCache
}
