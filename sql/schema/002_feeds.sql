-- +goose up
CREATE TABLE feeds (
    id uuid PRIMARY KEY,
    name varchar(255) NOT NULL,
    url varchar(255) NOT NULL UNIQUE,
    created_at timestamp,
    updated_at timestamp,
    user_id uuid REFERENCES users (id) ON DELETE CASCADE NOT NULL
);

-- +goose down
DROP TABLE feeds;

