-- +goose up
CREATE TABLE feed_follows(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
  feed_id UUID NOT NULL REFERENCES feeds ON DELETE CASCADE,
  UNIQUE(user_id, feed_id),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_feed FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose down
DROP TABLE feed_follows;
