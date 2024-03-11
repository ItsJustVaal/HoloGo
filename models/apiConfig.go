package models

import (
	"context"
	"fmt"
	"log"

	"github.com/ItsJustVaal/HoloGo/internal/database"
)

// Video cache map that will set the playlist to
// the most recent video ID in the DB
type VideoCache struct {
	LastVideo map[string]string
}

// API Cache
type ApiConfig struct {
	DB    *database.Queries
	Cache VideoCache
}

// Calls the DB for the most recent Video ID for each channel
// Maps the playlist to the Video ID, then sets that map in the
// API struct to be used as a server cache. This will limit the
// number of calls I need to make to the Youtube API.
func (v VideoCache) SetCache(db database.Queries, cache VideoCache) error {

	if len(cache.LastVideo) <= 0 {
		playlists, err := db.GetPlaylistIDs(context.Background())
		if err != nil {
			return fmt.Errorf("failed to get Playlist IDs for Cache: %v", err.Error())
		}

		for _, item := range playlists {
			dbCheck, err := db.GetMostRecentVideo(context.Background(), item)
			if err != nil {
				log.Printf("No recent video found")
				cache.LastVideo[item] = ""
			} else {
				cache.LastVideo[item] = dbCheck
				log.Printf("Set cache for playlist: %s", item)
			}
		}
		log.Println("Cache Init Complete")
	}
	return nil
}
