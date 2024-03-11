package youtube

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/ItsJustVaal/HoloGo/internal/database"
	"github.com/ItsJustVaal/HoloGo/models"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Sets up the interval and wait group then starts a
// a routine for each channel call
func StartYoutubeCalls(db database.Queries, key string, cache models.VideoCache, interval time.Duration) {
	ticker := time.NewTicker(interval)
	log.Println("Starting Youtube Calls")
	for range ticker.C {
		// playlistIDs, err := db.GetPlaylistIDs(context.Background())
		// if err != nil {
		// 	log.Printf("failed to get playlist ids: %v", err.Error())
		// }
		// log.Printf(playlistIDs[0])

		wg := &sync.WaitGroup{}
		for playlist := range cache.LastVideo {
			log.Println("PLAYLIST ID " + playlist)
			wg.Add(1)
			go getVideoDetails(db, key, cache, wg, playlist)
		}
		wg.Wait()
		log.Println("Completed")
	}
}

// Gets the most recent video IDs for each channel
// This uses the least amount of Youtube API call resources
// as possible
func getVideoIDs(key string, playlist string) ([]string, error) {
	videoMap := []string{}
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}

	resp, err := service.PlaylistItems.List([]string{"contentDetails"}).PlaylistId(playlist).MaxResults(15).Do()
	if err != nil {
		return videoMap, err
	}

	for _, item := range resp.Items {
		videoMap = append(videoMap, item.ContentDetails.VideoId)
	}
	return videoMap, nil
}

func getVideoDetails(db database.Queries, key string, cache models.VideoCache, wg *sync.WaitGroup, playlist string) {
	defer wg.Done()

	// First Call to get most recent video IDs for each channel
	videoIDs, err := getVideoIDs(key, playlist)
	if err != nil {
		log.Printf("Failed to get video IDs for playlist: %s %v", playlist, err.Error())
	}

	// Youtube Client init
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		log.Fatalf("Failed to make youtube service: %v", err.Error())
	}

	for _, video := range videoIDs {
		// Checks if the most recent video is in the cache
		// if it is, it ends the go routine
		if cache.LastVideo[playlist] == video {
			cache.LastVideo[playlist] = videoIDs[0]
			log.Println("video found in cache, breaking")
			return
		}

		// Youtube calls start here if the cache != the most recent video
		resp, err := service.Videos.List([]string{"snippet, liveStreamingDetails"}).Id(video).Do()
		if err != nil {
			log.Fatalf("failed to get video details: %v\n", err.Error())
		}

		// Checks for nil thumbnail value
		items := resp.Items[0]
		if items.Snippet.Thumbnails.Standard == nil {
			items.Snippet.Thumbnails.Standard = items.Snippet.Thumbnails.High
		}

		// Converts the published at time
		timeCnv, err := time.Parse(time.RFC3339, items.Snippet.PublishedAt)
		if err != nil {
			log.Println(err.Error())
		}

		// Runs if the video is not a live stream
		if items.LiveStreamingDetails == nil {
			_, err = db.CreateVideo(context.Background(), database.CreateVideoParams{
				ID:                 uuid.New(),
				CreatedAt:          time.Now().UTC(),
				UpdatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
				Videoid:            video,
				Playlistid:         playlist,
				Title:              items.Snippet.Title,
				Description:        items.Snippet.Description,
				Thumbnail:          items.Snippet.Thumbnails.High.Url,
				PublishedAt:        sql.NullTime{Time: timeCnv, Valid: true},
				ScheduledStartTime: sql.NullString{String: "None", Valid: true},
				ActualStartTime:    sql.NullString{String: "None", Valid: true},
				ActualEndTime:      sql.NullString{String: "None", Valid: true},
			})
			if err != nil {
				log.Printf("failed to create video with ID: %s with error: %v\n", video, err.Error())
			}
			// Runs if the video is/was a live stream
		} else {
			_, err = db.CreateVideo(context.Background(), database.CreateVideoParams{
				ID:                 uuid.New(),
				CreatedAt:          time.Now().UTC(),
				UpdatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
				Videoid:            video,
				Playlistid:         playlist,
				Title:              items.Snippet.Title,
				Description:        items.Snippet.Description,
				Thumbnail:          items.Snippet.Thumbnails.High.Url,
				PublishedAt:        sql.NullTime{Time: timeCnv, Valid: true},
				ScheduledStartTime: sql.NullString{String: items.LiveStreamingDetails.ScheduledStartTime, Valid: true},
				ActualStartTime:    sql.NullString{String: items.LiveStreamingDetails.ActualStartTime, Valid: true},
				ActualEndTime:      sql.NullString{String: items.LiveStreamingDetails.ActualEndTime, Valid: true},
			})
			if err != nil {
				log.Printf("failed to create video with ID: %s with error: %v\n", video, err.Error())
			}
		}
	}
}

// This function wont be run unless there are new channels
// since playlist IDs dont change, will change this into
// a proper admin command or a CLI tool
func GetPlaylists(db database.Queries, key string, cache models.VideoCache) error {
	// Grab all channel IDs from DB
	errSlice := make([]string, 0)
	channelIds, err := db.GetChannelIDs(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get Channels: %v", err.Error())
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return err
	}

	for _, id := range channelIds {
		call := service.Channels.List([]string{"contentDetails"})
		resp, err := call.Id(id).Do()
		if err != nil {
			return err
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
			if !slices.Contains(errSlice, err.Error()) {
				errSlice = append(errSlice, err.Error())
			}
		}
	}
	log.Printf("Errors in Playlists: %v\n", errSlice)
	return nil
}
