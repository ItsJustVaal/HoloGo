-- name: CreateChannel :one
INSERT INTO channels (id, created_at, updated_at, channel, channelid, region, prio, oshi, gen, tags, company)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetChannel :one
SELECT * FROM channels
WHERE channelid = $1;

-- name: GetChannelIDs :many
SELECT channelid FROM channels;