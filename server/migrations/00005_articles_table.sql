-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS articles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  guid TEXT NOT NULL,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  author TEXT,
  content TEXT,
  summary TEXT,
  published_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

  -- full-text search helper
  search_vector tsvector GENERATED ALWAYS AS (
    to_tsvector('simple', coalesce(title, '') || ' ' || coalesce(content, '') || ' ' || coalesce(summary, ''))
  ) STORED,

  UNIQUE(feed_id, guid)
);

CREATE INDEX IF NOT EXISTS articles_feed_id_idx ON articles (feed_id);
CREATE INDEX IF NOT EXISTS articles_feed_published_idx ON articles (feed_id, published_at DESC);
CREATE INDEX IF NOT EXISTS articles_search_gin_idx ON articles USING GIN (search_vector);

-- +goose Down
DROP INDEX IF EXISTS articles_search_gin_idx;
DROP INDEX IF EXISTS articles_feed_published_idx;
DROP INDEX IF EXISTS articles_feed_id_idx;
DROP TABLE IF EXISTS articles;
