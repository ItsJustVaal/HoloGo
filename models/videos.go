package models

import (
	"context"
	"log"

	"github.com/ItsJustVaal/HoloGo/internal/database"
)

// Setup the model for videos being passed
// into a template

type VideoTPL struct {
	VideoID        string
	Playlist       string
	Title          string
	Description    string
	Thumbnail      string
	ScheduledStart string
	Schedule       string
}

// Make funcs to get all videos
// Probably split into different funcs:
// One will get all videos
// One will get region specific videos
// One will get gen specific videos
// One will get individuals videos
// One will get oshi videos
// One will get prio videos
// It will need to get input from the frontend and return the specific videos
// All of them will just return a slice of VideoTPL objects that the tpl will render

// &i.ID,
// &i.CreatedAt,
// &i.UpdatedAt,
// &i.Videoid,
// &i.Playlistid,
// &i.Title,
// &i.Description,
// &i.Thumbnail,
// &i.PublishedAt,
// &i.ScheduledStartTime,
// &i.ActualStartTime,
// &i.ActualEndTime,

func MainPageRender(db database.Queries) []VideoTPL {
	tplVideos := make([]VideoTPL, 0)
	dbVideos, err := db.GetAllVideos(context.Background())
	if err != nil {
		log.Println("Error getting all videos in main render")
		return nil
	}

	for _, video := range dbVideos {
		if video.ActualEndTime.String != "" {
			tplVideos = append(tplVideos, VideoTPL{
				VideoID:     video.Videoid,
				Playlist:    video.Playlistid,
				Title:       video.Title,
				Description: video.Description,
				Thumbnail:   video.Thumbnail,
				Schedule:    "Finished",
			})
		} else {
			tplVideos = append(tplVideos, VideoTPL{
				VideoID:        video.Videoid,
				Playlist:       video.Playlistid,
				Title:          video.Title,
				Description:    video.Description,
				Thumbnail:      video.Thumbnail,
				ScheduledStart: video.ScheduledStartTime.String,
				Schedule:       "Upcoming",
			})
		}
	}
	return tplVideos
}
