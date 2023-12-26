package youtube

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ItsJustVaal/HoloGo/internal/database"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetPlaylists(db database.Queries, key string) error {

	// For each ID call youtube API
	// extract string from struct
	// call next function using string
	// uses a Go routine for each ID call

	// Grab all channel IDs from DB
	channelIds, err := db.GetChannelIDs(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get Channels")
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	for _, id := range channelIds {
		call := service.Channels.List([]string{"contentDetails"})
		resp, err := call.Id(id).Do()
		if err != nil {
			log.Fatalln(err.Error())
		}

		items := resp.Items[0]
		_, err = db.CreatePlaylist(context.Background(), database.CreatePlaylistParams{
			ID:         uuid.New(),
			CreatedAt:  time.Now().UTC(),
			UpdatedAt:  sql.NullTime{Time: time.Now().UTC(), Valid: true},
			Channelid:  id,
			Playlistid: items.ContentDetails.RelatedPlaylists.Uploads,
		})
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
	return nil

}

// First is channels to get playlist IDs, contentDetails
// Then playlistitems to get video IDs
// Then videos to get video info
// Add sql field for translated video title (can make a toggle on homepage for fun)
// Need to build cache checking so I am not calling the API a ton
