package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/ItsJustVaal/HoloGo/internal/database"
	"github.com/ItsJustVaal/HoloGo/internal/models"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DBURL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln(err.Error())
	}

	queries := database.New(db)

	cfg := models.ApiConfig{
		DB: queries,
	}

	mainRouter := chi.NewRouter()
	mainRouter.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handleGetReadiness)
	v1Router.Get("/err", handleGetErr)

	mainRouter.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mainRouter,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
