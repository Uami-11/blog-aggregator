-- +goose Up
CREATE TABLE users (
    id uuid PRIMARY KEY,
    created_at timestamp,
    updated_at timestamp,
    name varchar(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;

