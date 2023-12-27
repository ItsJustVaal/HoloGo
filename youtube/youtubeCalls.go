package youtube

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/ItsJustVaal/HoloGo/internal/database"
	"github.com/google/uuid"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type VideoCache struct {
	LastVideo map[string]string
}

// This function wont be run unless there are new channels
// since playlist IDs dont change
func GetPlaylists(db database.Queries, key string, cache VideoCache) error {
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
	go getVideoDetails(db, key, cache)
	log.Printf("Errors in Playlists: %v\n", errSlice)
	return nil
}

func getVideoIDs(db database.Queries, key string, cache VideoCache) (map[string][]string, error) {
	log.Println("Getting Video IDs")
	videoMap := make(map[string][]string)

	playlistIDs, err := db.GetPlaylistIDs(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist ids: %v", err.Error())
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return nil, err
	}

	for _, id := range playlistIDs {
		_, ok := videoMap[id]
		if !ok {
			videoMap[id] = []string{}
		}

		call := service.PlaylistItems.List([]string{"contentDetails"}).MaxResults(15)
		resp, err := call.PlaylistId(id).Do()
		if err != nil {
			return videoMap, err
		}
		for _, item := range resp.Items {
			videoMap[id] = append(videoMap[id], item.ContentDetails.VideoId)
		}

	}
	for _, items := range videoMap {
		slices.Reverse(items)
	}
	return videoMap, nil
}

func getVideoDetails(db database.Queries, key string, cache VideoCache) error {
	log.Println("Setting Cache")
	err := setCache(db, cache)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("Starting Video Details Call")
	videoIDs, err := getVideoIDs(db, key, cache)
	if err != nil {
		return fmt.Errorf("failed to get video IDs: %v", err.Error())
	}

	service, err := youtube.NewService(context.Background(), option.WithAPIKey(key))
	if err != nil {
		return err
	}

	for channel, videos := range videoIDs {
		log.Printf("Playlist ID: %s\n", channel)
		for _, video := range videos {
			if cache.LastVideo[channel] == video {
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
					Playlistid:         channel,
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
				cache.LastVideo[channel] = video
			} else {
				_, err = db.CreateVideo(context.Background(), database.CreateVideoParams{
					ID:                 uuid.New(),
					CreatedAt:          time.Now().UTC(),
					UpdatedAt:          sql.NullTime{Time: time.Now().UTC(), Valid: true},
					Videoid:            video,
					Playlistid:         channel,
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
				cache.LastVideo[channel] = video
			}
		}
	}
	log.Println("Completed")
	log.Printf("Cache Complete: %v", cache.LastVideo)
	return nil
}

func NewCache() VideoCache {
	c := VideoCache{
		LastVideo: make(map[string]string),
	}
	return c
}

func setCache(db database.Queries, cache VideoCache) error {
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
				cache.LastVideo[item] = dbCheck.Videoid
				log.Printf("Set cache for playlist: %s", item)
			}
		}
		log.Println("Cache Init Complete")
	}
	return nil
}

// Need to build cache checking so I am not calling the API a ton
