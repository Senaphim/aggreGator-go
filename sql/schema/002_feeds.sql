-- +goose up
CREATE TABLE feeds(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL UNIQUE,
  url TEXT NOT NULL UNIQUE,
  user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose down
DROP TABLE feeds;
