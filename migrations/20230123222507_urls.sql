-- +goose Up
CREATE TABLE urls (
                       original_url varchar,
                       short_url varchar
);

-- +goose Down
DROP TABLE urls
