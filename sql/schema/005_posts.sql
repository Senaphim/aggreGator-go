-- +goose up
CREATE TABLE posts(
  id UUID NOT NULL,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  title TEXT,
  url TEXT NOT NULL UNIQUE,
  description TEXT,
  published_at TEXT,
  feed_id UUID NOT NULL REFERENCES feeds ON DELETE CASCADE,
  CONSTRAINT fk_feeds FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose down
DROP TABLE posts;

