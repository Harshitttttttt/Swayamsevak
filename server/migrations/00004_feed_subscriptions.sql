-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE feed_subscriptions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  custom_title TEXT,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  UNIQUE(user_id, feed_id)
);

CREATE INDEX feed_subscriptions_user_id_idx
ON feed_subscriptions (user_id);


-- +goose Down
DROP INDEX IF EXISTS feed_subscriptions_user_id_idx;
DROP TABLE feed_subscriptions;
