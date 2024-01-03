-- +goose Up
CREATE TABLE videos (
    id UUID UNIQUE PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    videoID TEXT NOT NULL UNIQUE,
    playlistID TEXT NOT NULL REFERENCES playlists(playlistID),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    thumbnail TEXT NOT NULL,
    published_at TIMESTAMP,
    scheduled_start_time TEXT,
    actual_start_time TEXT,
    actual_end_time TEXT
);

-- +goose Down
DROP TABLE videos;