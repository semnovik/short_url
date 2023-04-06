-- +goose Up
CREATE TABLE urls (
                      original_url TEXT NOT NULL UNIQUE,
                      short_url TEXT NOT NULL UNIQUE,
                      user_uuid TEXT,
                      is_deleted BOOLEAN DEFAULT FALSE
);

-- +goose Down
DROP TABLE urls
