-- name: CreateVideo :one
INSERT INTO videos (id, created_at, updated_at, videoID, playlistID, title, description, thumbnail, scheduled_start_time, actual_start_time, actual_end_time)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;
