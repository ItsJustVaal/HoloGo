// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: videos.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createVideo = `-- name: CreateVideo :one
INSERT INTO videos (id, created_at, updated_at, videoID, playlistID, title, description, thumbnail, scheduled_start_time, actual_start_time, actual_end_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, created_at, updated_at, videoid, playlistid, title, description, thumbnail, scheduled_start_time, actual_start_time, actual_end_time
`

type CreateVideoParams struct {
	ID                 uuid.UUID
	CreatedAt          time.Time
	UpdatedAt          sql.NullTime
	Videoid            string
	Playlistid         string
	Title              string
	Description        string
	Thumbnail          string
	ScheduledStartTime sql.NullString
	ActualStartTime    sql.NullString
	ActualEndTime      sql.NullString
}

func (q *Queries) CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error) {
	row := q.db.QueryRowContext(ctx, createVideo,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Videoid,
		arg.Playlistid,
		arg.Title,
		arg.Description,
		arg.Thumbnail,
		arg.ScheduledStartTime,
		arg.ActualStartTime,
		arg.ActualEndTime,
	)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Videoid,
		&i.Playlistid,
		&i.Title,
		&i.Description,
		&i.Thumbnail,
		&i.ScheduledStartTime,
		&i.ActualStartTime,
		&i.ActualEndTime,
	)
	return i, err
}

const getMostRecentVideo = `-- name: GetMostRecentVideo :one
SELECT videoID, playlistID FROM videos
WHERE playlistID = $1
ORDER BY created_at DESC
LIMIT 1
`

type GetMostRecentVideoRow struct {
	Videoid    string
	Playlistid string
}

func (q *Queries) GetMostRecentVideo(ctx context.Context, playlistid string) (GetMostRecentVideoRow, error) {
	row := q.db.QueryRowContext(ctx, getMostRecentVideo, playlistid)
	var i GetMostRecentVideoRow
	err := row.Scan(&i.Videoid, &i.Playlistid)
	return i, err
}
