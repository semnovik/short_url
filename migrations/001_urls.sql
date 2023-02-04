-- +goose Up
CREATE TABLE urls (
                      original_url text NOT NULL UNIQUE,
                      short_url text NOT NULL UNIQUE,
                      user_uuid text
);

-- +goose Down
DROP TABLE urls
