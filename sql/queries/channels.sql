-- name: CreateChannel :one
INSERT INTO channels (id, created_at, updated_at, Channel, ChannelID, Region, Prio, Oshi, Gen, Tags, Company)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: GetChannel :one
SELECT * FROM channels
WHERE ChannelID = $1;
