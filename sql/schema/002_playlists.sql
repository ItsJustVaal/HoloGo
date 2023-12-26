-- +goose Up
CREATE TABLE playlists (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    channelID TEXT NOT NULL UNIQUE,
    playlistID TEXT NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE playlists;