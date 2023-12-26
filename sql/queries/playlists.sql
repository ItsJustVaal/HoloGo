-- name: CreatePlaylist :one
INSERT INTO playlists (id, created_at, updated_at, channelID, playlistID)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;