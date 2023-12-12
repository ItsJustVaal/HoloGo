-- +goose Up
CREATE TABLE channels (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    Channel TEXT NOT NULL,
    ChannelID TEXT NOT NULL,
    Region TEXT NOT NULL,
    Prio BOOLEAN,
    Oshi BOOLEAN,
    Gen INTEGER,
    Tags TEXT,
    Company TEXT
);

-- +goose Down
DROP TABLE channels;