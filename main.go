package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	_ "github.com/lib/pq"

	"github.com/ItsJustVaal/HoloGo/controllers"
	"github.com/ItsJustVaal/HoloGo/internal/database"
	"github.com/ItsJustVaal/HoloGo/models"
	"github.com/ItsJustVaal/HoloGo/templates"
	"github.com/ItsJustVaal/HoloGo/views"
	"github.com/ItsJustVaal/HoloGo/youtube"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// env variable init
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	// DB init
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()
	queries := database.New(db)

	// Cache init
	cache := models.VideoCache{
		LastVideo: make(map[string]string),
	}

	// Sets API Config Struct
	cfg := models.ApiConfig{
		DB:         queries,
		Cache:      cache,
		CSRFKey:    os.Getenv("CSRF_KEY"),
		CSRFSecure: false,
	}

	vidService := controllers.Vids{
		DB: *queries,
	}

	// Sets the Server Cache with most recent videoID
	// from each channel, if no ID exists, uses zero value
	log.Println("Setting Cache")
	err = cache.SetCache(*queries, cache)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Sets main router and cors settings
	// This will be secured once im not using localhost
	csrfMW := csrf.Protect([]byte(cfg.CSRFKey), csrf.Secure(cfg.CSRFSecure))
	mainRouter := chi.NewRouter()
	mainRouter.Use(csrfMW, cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// V1 Router to check server status
	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", handleGetReadiness)
	v1Router.Get("/err", handleGetErr)

	mainRouter.Mount("/v1", v1Router)
	vidService.Templates.Main = views.Must(views.ParseFS(templates.FS, "home.gohtml", "layout.gohtml"))

	mainRouter.Get("/", vidService.PopulateHome)

	// Sets Server Struct
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mainRouter,
	}

	// Youtube calls on an interval to update the Server Cache and DB for API
	// Uses a go wait group to make a seperate call for each channel
	const interval = time.Minute * 10
	go youtube.StartYoutubeCalls(*cfg.DB, apiKey, cfg.Cache, interval)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
