-- +goose up
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    title VARCHAR(225) NOT NULL,
    url VARCHAR(225) NOT NULL UNIQUE,
    description VARCHAR(225) NOT NULL,
    published_at TIMESTAMP,
    feed_id UUID REFERENCES feeds(id) ON DELETE CASCADE NOT NULL
);

-- +goose down
DROP TABLE posts;
