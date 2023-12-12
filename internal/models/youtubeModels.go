package models

import "time"

type Playlist struct {
	Items []struct {
		ContentDetails struct {
			RelatedPlaylists struct {
				Uploads string `json:"uploads"`
			} `json:"relatedPlaylists"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type VideoList struct {
	Items []struct {
		ContentDetails struct {
			VideoID          string    `json:"videoId"`
			VideoPublishedAt time.Time `json:"videoPublishedAt"`
		} `json:"contentDetails"`
	} `json:"items"`
}

type Video struct {
	Items []struct {
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
				Standard struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"standard"`
				Maxres struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"maxres"`
			} `json:"thumbnails"`
			ChannelTitle         string   `json:"channelTitle"`
			Tags                 []string `json:"tags"`
			CategoryID           string   `json:"categoryId"`
			LiveBroadcastContent string   `json:"liveBroadcastContent"`
			Localized            struct {
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"localized"`
			DefaultAudioLanguage string `json:"defaultAudioLanguage"`
		} `json:"snippet"`
		LiveStreamingDetails struct {
			ActualStartTime    time.Time `json:"actualStartTime"`
			ActualEndTime      time.Time `json:"actualEndTime"`
			ScheduledStartTime time.Time `json:"scheduledStartTime"`
		} `json:"liveStreamingDetails"`
	} `json:"items"`
}
