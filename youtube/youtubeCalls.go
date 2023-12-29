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
	"github.com/ItsJustVaal/HoloGo/internal/models"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func StartYoutubeCalls(db database.Queries, key string, cache models.VideoCache, interval time.Duration) {
	ticker := time.NewTicker(interval)

	for ; ; <-ticker.C {
		playlistIDs, err := db.GetPlaylistIDs(context.Background())
		if err != nil {
			log.Printf("failed to get playlist ids: %v", err.Error())
		}

		wg := &sync.WaitGroup{}
		for _, playlist := range playlistIDs {
			wg.Add(1)
			go getVideoDetails(db, key, cache, wg, playlist)
		}
		wg.Wait()
	}
}

// This function wont be run unless there are new channels
// since playlist IDs dont change
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

func getVideoIDs(db database.Queries, key string, playlist string) ([]string, error) {
	log.Println("Getting Video IDs")
	videoMap := []string{}
	service, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}

	call := service.PlaylistItems.List([]string{"contentDetails"}).MaxResults(15)
	resp, err := call.PlaylistId(playlist).Do()
	if err != nil {
		return videoMap, err
	}

	for _, item := range resp.Items {
		videoMap = append(videoMap, item.ContentDetails.VideoId)
	}
	return videoMap, nil
}

func getVideoDetails(db database.Queries, key string, cache models.VideoCache, wg *sync.WaitGroup, playlist string) error {
	defer wg.Done()
	log.Println("Starting Video Details Call")
	videoIDs, err := getVideoIDs(db, key, playlist)
	if err != nil {
		return fmt.Errorf("failed to get video IDs: %v", err.Error())
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return err
	}
	log.Printf("Playlist: %s", playlist)
	for _, video := range videoIDs {
		log.Printf("Video ID: %s\n", video)
		if cache.LastVideo[playlist] == video {
			log.Println("video found in cache, breaking")
			break
		}
		call := service.Videos.List([]string{"snippet, liveStreamingDetails"})
		resp, err := call.Id(video).Do()
		if err != nil {
			log.Fatalf("failed to get video details: %v\n", err.Error())
		}
		items := resp.Items[0]
		if items.Snippet.Thumbnails.Standard == nil {
			items.Snippet.Thumbnails.Standard = items.Snippet.Thumbnails.High
		}

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
				ScheduledStartTime: sql.NullString{String: "None", Valid: true},
				ActualStartTime:    sql.NullString{String: "None", Valid: true},
				ActualEndTime:      sql.NullString{String: "None", Valid: true},
			})
			if err != nil {
				log.Printf("failed to create video: %v\n", err.Error())
			}
			cache.LastVideo[playlist] = video
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
				ScheduledStartTime: sql.NullString{String: items.LiveStreamingDetails.ScheduledStartTime, Valid: true},
				ActualStartTime:    sql.NullString{String: items.LiveStreamingDetails.ActualStartTime, Valid: true},
				ActualEndTime:      sql.NullString{String: items.LiveStreamingDetails.ActualEndTime, Valid: true},
			})
			if err != nil {
				log.Printf("failed to create video: %v\n", err.Error())
			}
			cache.LastVideo[playlist] = video
		}
	}
	log.Println("Completed")
	return nil
}
