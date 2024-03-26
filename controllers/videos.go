package controllers

import (
	"net/http"

	"github.com/ItsJustVaal/HoloGo/internal/database"
	"github.com/ItsJustVaal/HoloGo/models"
)

type Vids struct {
	Templates struct {
		Main Template
	}
	DB database.Queries
}

func (v Vids) PopulateHome(w http.ResponseWriter, r *http.Request) {
	videos := models.MainPageRender(v.DB)
	v.Templates.Main.Execute(w, r, videos)
}
