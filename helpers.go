package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"time"

	"github.com/ItsJustVaal/HoloGo/internal/database"
	"github.com/ItsJustVaal/HoloGo/models"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

// v1 Readiness Handler
func handleGetReadiness(w http.ResponseWriter, r *http.Request) {
	models.RespondWithJSON(w, http.StatusOK, models.JsonResponse{
		Status: "ok",
	})
}

// v1 Error Handler
func handleGetErr(w http.ResponseWriter, r *http.Request) {
	models.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

// This is just a helper function for me to add all the chs to DB from a csv, setup function
// This will also be changed into a admin or CLI action later
func AddChannelsToDB(db database.Queries) {
	channelsFile, err := os.OpenFile("holodash.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer channelsFile.Close()

	channels := []*models.CSVChannel{}

	if err := gocsv.UnmarshalFile(channelsFile, &channels); err != nil {
		fmt.Println(err.Error())
	}
	errSlice := make([]string, 0)
	uniqueCheck := 0
	for _, channel := range channels {
		_, err := db.CreateChannel(context.Background(), database.CreateChannelParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
			Channel:   channel.Channel,
			Channelid: channel.Channelid,
			Region:    channel.Region,
			Prio:      sql.NullBool{Bool: channel.Prio, Valid: true},
			Oshi:      sql.NullBool{Bool: channel.Oshi, Valid: true},
			Gen:       sql.NullInt32{Int32: channel.Gen, Valid: true},
			Tags:      sql.NullString{String: channel.Tags, Valid: true},
			Company:   sql.NullString{String: channel.Company, Valid: true},
		})
		if err != nil {
			uniqueCheck++
			if !slices.Contains(errSlice, err.Error()) {
				errSlice = append(errSlice, err.Error())
			}
		}
	}
	log.Printf("Num of Channels Added: %d, Num of Duplicates: %d\n", len(channels)-uniqueCheck, uniqueCheck)
	if len(errSlice) != 0 {
		log.Printf("List of errors: %v\n", errSlice)
	}
}

// Move this to a new package so it can be run from the cli
