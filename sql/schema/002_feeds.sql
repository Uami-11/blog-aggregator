-- +goose up
CREATE TABLE feeds (
    id uuid PRIMARY KEY,
    name varchar(255),
    url varchar(255) UNIQUE,
    created_at timestamp,
    updated_at timestamp,
    user_id uuid REFERENCES users (id) ON DELETE CASCADE
);

-- +goose down
DROP TABLE feeds;

