// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Channel   string
	Channelid string
	Region    string
	Prio      sql.NullBool
	Oshi      sql.NullBool
	Gen       sql.NullInt32
	Tags      sql.NullString
	Company   sql.NullString
}

type Playlist struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  sql.NullTime
	Channelid  string
	Playlistid string
}

type Video struct {
	ID                 uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          sql.NullTime
	Videoid            string
	Playlistid         string
	Title              string
	Description        string
	Thumbnail          string
	PublishedAt        sql.NullTime
	ScheduledStartTime sql.NullString
	ActualStartTime    sql.NullString
	ActualEndTime      sql.NullString
}
